import sys
import os
import time
import glob
import pickle
import random
import fnmatch
import math
from os.path import isfile, join
from aubio import onset, source, tempo
import numpy as np
import pydub
from pydub.silence import split_on_silence
from pydub.silence import detect_nonsilent
import lib.alignment

class MusicAssembler:

    def __init__(self,voicefile):
        #voicefile is under the form : input_<UUID>.mp3
        #we can find the metadata innside TMP_FOLDER/metadata_<UUID>
        self.voicefile = voicefile
        self.voiceUUID = self.voicefile.split("_")[1].split(".")[0]
        self.TMP_FOLDER = os.environ.get("TMP_FOLDER", "/data/tmp/")
        self.METADATA_FOLDER = self.TMP_FOLDER+os.environ.get("METADATA_FOLDER_PREFIX","metadata_")+self.voiceUUID+"/"

        self.beatsDistribution = [] #same dim as self.beatsIntensity
        self.beatsIntensity = []
        self.verse_interval = []

        self.beatfileDuration = 0.0  #duration of beat file in seconds
        self.bpm = 0.0 #average bpm of the beatfile
        self.randomPadding = 0 # integer representing how many beat we should skip for the next chunk
        self.chunks = []

        #OUTPUT
        self.outputfilePrefix    = "output_"
        self.outputformat  = "mp3"

    def voice_splitter(self):
        """ 
        Generate .mp3 files in which you can find the chunks of voice
        without the silence parts.
        """
        print("[VOICE_SPLITTER] Starting voice splitting...")
        voice = pydub.AudioSegment.from_file(self.TMP_FOLDER+str(self.voicefile).split("/")[-1].split("'")[0])
        print("[VOICE_SPLITTER] voice dBFS : "+str(voice.dBFS))
        dBFS = voice.dBFS
        # Split track where the silence is 1 second or more and get chunks using the imported function
        # for the further C implementation, use ffmpeg C API 'silencedetect' filter (https://stackoverflow.com/questions/44799312/ffprobe-ffmpg-silence-detection-command)
        # which will be way faster !
        chunks = split_on_silence (
            voice,
            min_silence_len = 500,
            # Consider a chunk silent if it's quieter than -16 dBFS.
            # (We may want to adjust this parameter.)
            silence_thresh = dBFS-16,
        )
        print("[VOICE_SPLITTER] "+str(len(chunks))+" chunks generated.")
        for chunk in chunks:
            self.chunks.append(chunk)
        print("[VOICE_SPLITTER] voice splitting has ended.")

    def __match_target_amplitude(self,aChunk, target_dBFS):
        ''' Normalize given audio chunk '''
        change_in_dBFS = target_dBFS - aChunk.dBFS
        return aChunk.apply_gain(change_in_dBFS)

    def __read(self,f, normalized=False):
        """MP3 to numpy array ==> useful for the input of __phase_vocoder function"""
        a = pydub.AudioSegment.from_mp3(f)
        y = np.array(a.get_array_of_samples())
        if a.channels == 2:
            y = y.reshape((-1, 2))
        if normalized:
            return a.frame_rate, np.float32(y) / 2**15
        else:
            return a.frame_rate, y

    def __write(self,f, sr, x, normalized=False):
        """numpy array to MP3"""
        channels = 2 if (x.ndim == 2 and x.shape[1] == 2) else 1
        if normalized:  # normalized array - each item should be a float in [-1, 1)
            y = np.int16(x * 2 ** 15)
        else:
            y = np.int16(x)
        song = pydub.AudioSegment(y.tobytes(), frame_rate=sr, sample_width=2, channels=channels)
        song.export(f, format="mp3", bitrate="320k")

    def __get_sound_paddings(self):
        """
        calculate the average sound padding for the chunks during the active beat phases of the song
        """
        avg_sound_paddings = []
        for i in range(1, len(self.verse_interval)):
            startBeatPhaseIdx = self.verse_interval[i-1][1]
            endBeatPhaseIdx = self.verse_interval[i][0]

            avgRythm = 0.0
            for t in range(startBeatPhaseIdx,endBeatPhaseIdx+1):
                if t != endBeatPhaseIdx:
                    avgRythm += self.beatsDistribution[t+1]-self.beatsDistribution[t]
            avgRythm /= endBeatPhaseIdx-startBeatPhaseIdx
            avg_sound_paddings.append(avgRythm)
        return avg_sound_paddings

    def __exportMergedSound(self,mergedResult):
        mergedResult.export(self.TMP_FOLDER+self.outputfilePrefix+self.voiceUUID,format=self.outputformat)

    def merger(self):
        """
        merge `musicData` which contains the beats distribution
        calculated from getBeats function and saved into pickle object
        and the intervals telling us when we are in the verse part.
        For now, we will put the voices only on the "hooks" and not on the "verses" 
        `voiceData` contains liste of voice chunks name with their context, 
         i.e : is this chunk in the chorus part ? 
        """
        
        ### ONLY FOR TESTING ###
        #self.beatsDistribution = pickle.load(open("testPickle/beatDistribution","rb"))
        #self.verse_interval = pickle.load(open("testPickle/verseInterval","rb"))
        #self.chunks = pickle.load(open("testPickle/chunks","rb"))
        ########################
        soundPaddings = self.__get_sound_paddings()
        print("[VOICE MERGER] the sound padding are : "+str(soundPaddings))
        print("[VOICE MERGER] verse distribution : "+str(self.verse_interval))
        mergedResult = self.__loadBeat() #the final result (for now, contains only the beatfile but the chunks will be merged progressively)
        
        chunkNumber = 0
        beatCounted = 0
        while(beatCounted != len(self.beatsDistribution)):
            try:
                c = self.chunks[chunkNumber]
            except IndexError:
                print("chunkNumber overflow. Break.")
                break
            #Alignment
            endT = self.beatsDistribution[beatCounted]+c.duration_seconds
            #alignChunk, newBeatsCounted = alignment.align(c,endT,self.beatsDistribution,beatCounted) # after the MVP : use the self.beatsIntensity
            newBeatsCounted = lib.alignment.align(c,endT,self.beatsDistribution,beatCounted)
            startBeatCounted = beatCounted
            beatCounted += newBeatsCounted-beatCounted
            # leftChunkPart is a slice of the original chunk 
            # rightChunkPart is the "stretched"(slower or faster) complementary slice of the original chunk. 
            #mergedResult = alignment.overlay(mergedResult,alignChunk,startBeatCounted,beatCounted,self.beatsDistribution) # add the chunk to the beatsound
            mergedResult = lib.alignment.overlay(mergedResult, c, startBeatCounted, beatCounted, self.beatsDistribution) # add the chunk to the beatsound
            #Padding (padding means moving forward beatCounted )
            # for now, we'll choose a constant padding of 250ms
            beatCounted += random.randint(2,5)
            if beatCounted >= len(self.beatsDistribution):
                print("EARLY break")
                break
            print("beats counted end loop : "+str(beatCounted))
            chunkNumber += 1

        # final step : exportation of the result under the defined output file format
        self.__exportMergedSound(mergedResult)
        print("Merger finished ")
        return True

    def run(self):

        from os import listdir
        from os.path import isfile, join
        metadata = [f for f in listdir(self.METADATA_FOLDER) if isfile(join(self.METADATA_FOLDER, f))]

        for f in metadata:
            prefix = f.split("_")[0]
            if prefix == "duration":
                with open(self.METADATA_FOLDER+f, "rb") as f_in:
                    self.beatfileDuration = pickle.load(f_in)
            elif prefix == "bpm":
                with open(self.METADATA_FOLDER+f, "rb") as f_in:
                    self.bpm = pickle.load(f_in)
                self.randomPadding = math.ceil((self.bpm/60)/(random.randint(1,3)))
            elif prefix == "sound":
                print("voicefile detected !")
            elif prefix == "tempDist":
                with open(self.METADATA_FOLDER+f, "rb") as f_in:
                    self.beatsDistribution = pickle.load(f_in)
            elif prefix == "tempInt":
                with open(self.METADATA_FOLDER+f, "rb") as f_in:
                    self.beatsIntensity = pickle.load(f_in)
            elif prefix == "verseInterval":
                with open(self.METADATA_FOLDER+f, "rb") as f_in:
                    self.verse_interval = pickle.load(f_in)
            else:
                print("the prefix found does not match")

        #the attributes are loaded, we can begin the voice spliter
        self.voice_splitter()
        # merge voice chunks with beats
        res = self.merger()
        return res
    
    # obselete
    #def __assembler(self):
    #    """
    #    Reassemble the chunks together to form the raw beat pydub.AudioSegment object
    #    """
    #    chunkList = glob.glob(self.METADATA_FOLDER+"beatchunk*")
    #    uuid = chunkList[0].split("_")[1]
    #    chunkList.sort()
    #    res = pydub.AudioSegment.empty()
    #    for c in chunkList:
    #        res += pickle.load(open(c,"rb"))
    #    return res
    def __loadBeat(self):
        beatfile = glob.glob(self.METADATA_FOLDER+"beat_*")[0].split("/")[-1]
        res = pydub.AudioSegment.from_mp3(self.METADATA_FOLDER+beatfile)
        return res


#BEAT_PATH  = "testBeat/"
#VOICE_PATH = "testVoice/"

#ma = MusicAssembler(BEAT_PATH+"song_548ba4d4-64a3-4db1-ba2e-f9037c535932.mp3",VOICE_PATH+"greatSpeech.mp3")
#ma.getBeats() ==> will be used by the ingestion engine
#ma.verse_detector() ==> will be used by the ingestion engine
#ma.voice_splitter()
#ma.merger()