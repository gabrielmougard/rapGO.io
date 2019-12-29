import sys
from os import listdir
import time
import pickle
import random
import math
from os.path import isfile, join
from aubio import onset, source, tempo
import numpy as np
import pydub
from pydub.silence import split_on_silence
from pydub.silence import detect_nonsilent
import alignment

class MusicAssembler:

    def __init__(self,beatfile,voicefile):
        self.beatfile  = beatfile
        self.voicefile = voicefile

        self.beatsDistribution = [] #same dim as self.beatsIntensity
        self.beatsIntensity = []
        self.verse_interval = []
        self.beatfileDuration = 0.0  #duration of beat file in seconds
        self.bpm = self.__getBPM() #average bpm of the beatfile
        self.randomPadding = math.ceil((self.bpm/60)/(random.randint(1,3))) # integer representing how many beat we should skip for the next chunk
        self.chunks = []
        self.CHUNKS_PATH = "chunks/"

        #OUTPUT
        self.outputfile    = "testOutput/output"
        self.outputformat  = "mp3"
    
    def getBeats(self):
        startTime = time.time()
        print("[BEATS RECOGNIZER] Starting beats recognition...")
        win_s = 512
        hop_s = win_s // 2
        samplerate = 44000

        s = source(self.beatfile,samplerate, hop_s)
        samplerate = s.samplerate

        o = onset("default", win_s, hop_s, samplerate)

        onsets = []
        desc = []
        tdesc = []
        downsample = 2 #to plot n samples / hop_s
        allsamples_max = np.zeros(0,)
        #total number of frames read
        total_frames = 0
        while True:
            samples, read = s()
            if o(samples):
                onsets.append(o.get_last())
            new_maxes = (abs(samples.reshape(hop_s//downsample, downsample))).max(axis=0)
            allsamples_max = np.hstack([allsamples_max, new_maxes])
            desc.append(o.get_descriptor())
            tdesc.append(o.get_thresholded_descriptor())
            total_frames += read
            if read < hop_s: break

        desc_times = [ float(t) * hop_s / samplerate for t in range(len(desc)) ]
        desc_max = max(desc) if max(desc) != 0 else 1.
        desc_plot = [d / desc_max for d in desc]
        CONSIDERATION_THRESHOLD = 0.7
        desc_times_filtered = []
        desc_plot_filtered = []
        for x,y in zip(desc_times,desc_plot):
            if y > CONSIDERATION_THRESHOLD:
                desc_times_filtered.append(x)
                desc_plot_filtered.append(y)

        allsamples_max = (allsamples_max > 0) * allsamples_max
        allsamples_max_times = [ float(t) * hop_s / downsample / samplerate for t in range(len(allsamples_max)) ]
        duration = allsamples_max_times[-1]

        #update object
        self.beatfileDuration  = duration
        self.beatsDistribution = desc_times_filtered
        self.beatsIntensity    = desc_plot_filtered
        with open("testPickle/beatDistribution","wb") as f:
            pickle.dump(self.beatsDistribution, f)
        with open("testPickle/beatIntensity","wb") as f:
            pickle.dump(self.beatsIntensity, f)

        print("[BEATS RECOGNIZER] Beats recognition ended.")
        endTime = time.time()
        print("[BEATS RECOGNIZER] Time elasped (in secs): "+str(endTime-startTime))

    def __getBPM(self):
        """
        return the average bpm of a song. In prod, this will be executed in the
        ingestion engine and it will be persisted in a binary file (pickle object)
        with the beatDistribution and beatAmplDistribution.
        """
        print("Getting average BPM of song...")
        samplerate, win_s, hop_s = 44100, 1024, 512

        s = source(self.beatfile, samplerate, hop_s)
        samplerate = s.samplerate
        o = tempo("specdiff", win_s, hop_s, samplerate)
        # List of beats, in samples
        beats = []
        # Total number of frames read
        total_frames = 0

        while True:
            samples, read = s()
            is_beat = o(samples)
            if is_beat:
                this_beat = o.get_last_s()
                beats.append(this_beat)
                #if o.get_confidence() > .2 and len(beats) > 2.:
                #    break
            total_frames += read
            if read < hop_s:
                break

        # Convert to periods and to bpm 
        if len(beats) > 1:
            if len(beats) < 4:
                print("few beats found in {:s}".format(self.beatfile))
            bpms = 60./np.diff(beats)
            b = np.median(bpms)
        else:
            b = 0
            print("not enough beats found in {:s}".format(self.beatfile))
        print("The song is "+str(b)+" bpm.")
        self.bpm = b
        return b

    def verse_detector(self):
        """
        Return an array of tuple in which you have
        the starting time of a verse and the ending time
        of the verse.
        """
        startTime = time.time()
        print("[VERSE_DETECTOR] Starting verse detector of voice file...")
        verse_intervals = []
        MINIMUM_VERSE_LENGTH = 5.0 # if the interval is greater or equal than 5s, it's a verse.
        if self.beatsDistribution[0] > MINIMUM_VERSE_LENGTH:
            verse_intervals.append((-1,0))
        for i in range(0,len(self.beatsDistribution)-1):
            if abs(self.beatsDistribution[i]-self.beatsDistribution[i+1]) >= MINIMUM_VERSE_LENGTH:
                verse_intervals.append((i,i+1))
        if abs(self.beatfileDuration-self.beatsDistribution[-1]) > MINIMUM_VERSE_LENGTH:
            verse_intervals.append((len(self.beatsDistribution)-1,-1))

        self.verse_interval = verse_intervals
        with open("testPickle/verseInterval","wb") as f:
            pickle.dump(self.verse_interval, f)
        print("[VERSE_DETECTOR] verse detector ended.")
        endTime = time.time()
        print("[VERSE_DETECTOR] Time elasped (in secs): "+str(endTime-startTime))

    def voice_splitter(self):
        """ 
        Generate .mp3 files in which you can find the chunks of voice
        without the silence parts.
        """
        print("[VOICE_SPLITTER] Starting voice splitting...")
        voice = pydub.AudioSegment.from_mp3(self.voicefile)
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
        print("BEFORE PICKLE DURATION : "+str(self.chunks[0].duration_seconds))
        with open("testPickle/chunks","wb") as f:
            pickle.dump(self.chunks, f)
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
        mergedResult.export(self.outputfile,format=self.outputformat)

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
        self.beatsDistribution = pickle.load(open("testPickle/beatDistribution","rb"))
        self.verse_interval = pickle.load(open("testPickle/verseInterval","rb"))
        self.chunks = pickle.load(open("testPickle/chunks","rb"))
        ########################
        soundPaddings = self.__get_sound_paddings()
        print("[VOICE MERGER] the sound padding are : "+str(soundPaddings))
        print("[VOICE MERGER] verse distribution : "+str(self.verse_interval))
        mergedResult = pydub.AudioSegment.from_mp3(self.beatfile) #the final result (for now, contains only the beatfile but the chunks will be merged progressively)
        
        chunkNumber = 0
        beatCounted = 0
        while(beatCounted != len(self.beatsDistribution)):
            c = self.chunks[chunkNumber]
            #Alignment
            endT = self.beatsDistribution[beatCounted]+c.duration_seconds
            #alignChunk, newBeatsCounted = alignment.align(c,endT,self.beatsDistribution,beatCounted) # after the MVP : use the self.beatsIntensity
            newBeatsCounted = alignment.align(c,endT,self.beatsDistribution,beatCounted)
            startBeatCounted = beatCounted
            beatCounted += newBeatsCounted-beatCounted
            # leftChunkPart is a slice of the original chunk 
            # rightChunkPart is the "stretched"(slower or faster) complementary slice of the original chunk. 
            #mergedResult = alignment.overlay(mergedResult,alignChunk,startBeatCounted,beatCounted,self.beatsDistribution) # add the chunk to the beatsound
            mergedResult = alignment.overlay(mergedResult, c, startBeatCounted, beatCounted, self.beatsDistribution) # add the chunk to the beatsound
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

    def run(self):

BEAT_PATH  = "testBeat/"
VOICE_PATH = "testVoice/"

ma = MusicAssembler(BEAT_PATH+"song_548ba4d4-64a3-4db1-ba2e-f9037c535932.mp3",VOICE_PATH+"greatSpeech.mp3")
ma.getBeats()
ma.verse_detector()
ma.voice_splitter()
ma.merger()