import time
import os
import glob
import pickle
import threading
from aubio import onset, source, tempo
import pydub
import numpy as np
from pydub.silence import split_on_silence
from pydub.silence import detect_nonsilent
from google.cloud import storage
from google.cloud.exceptions import GoogleCloudError
#### sound operation
def getBeats(beatfile):
    uuid = beatfile.split("_")[1].split(".")[0]
    startTime = time.time()
    print("[BEATS RECOGNIZER] Starting beats recognition...")
    win_s = 512
    hop_s = win_s // 2
    samplerate = 44000
    s = source(beatfile,samplerate, hop_s)
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
    beatfileDuration  = duration
    beatsDistribution = desc_times_filtered
    beatsIntensity    = desc_plot_filtered
    #write the pickle binaries
    with open("duration_"+uuid,"wb") as f:
        pickle.dump(beatfileDuration, f)
    with open("tempDist_"+uuid,"wb") as f:
        pickle.dump(beatsDistribution, f)
    with open("tempInt_"+uuid,"wb") as f:
        pickle.dump(beatsIntensity, f)
    print("[BEATS RECOGNIZER] Beats recognition ended.")
    endTime = time.time()
    print("[BEATS RECOGNIZER] Time elasped (in secs): "+str(endTime-startTime))

def getBPM(beatfile):
    """
    return the average bpm of a song. In prod, this will be executed in the
    ingestion engine and it will be persisted in a binary file (pickle object)
    with the beatDistribution and beatAmplDistribution.
    """
    uuid = beatfile.split("_")[1].split(".")[0]
    print("Getting average BPM of song...")
    samplerate, win_s, hop_s = 44100, 1024, 512
    s = source(beatfile, samplerate, hop_s)
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
            print("few beats found in {:s}".format(beatfile))
        bpms = 60./np.diff(beats)
        b = np.median(bpms)
    else:
        b = 0
        print("not enough beats found in {:s}".format(beatfile))
    print("The song is "+str(b)+" bpm.")
    #write the pickle binaries
    with open("bpm_"+uuid,"wb") as f:
        pickle.dump(b, f)

def verse_detector(beatfile, beatsDistribution, beatfileDuration):
    """
    Return an array of tuple in which you have
    the starting time of a verse and the ending time
    of the verse.
    """
    uuid = beatfile.split("_")[1].split(".")[0]
    startTime = time.time()
    print("[VERSE_DETECTOR] Starting verse detector of voice file...")
    verse_intervals = []
    MINIMUM_VERSE_LENGTH = 5.0 # if the interval is greater or equal than 5s, it's a verse.
    if beatsDistribution[0] > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((-1,0))
    for i in range(0,len(beatsDistribution)-1):
        if abs(beatsDistribution[i]-beatsDistribution[i+1]) >= MINIMUM_VERSE_LENGTH:
            verse_intervals.append((i,i+1))
    if abs(beatfileDuration-beatsDistribution[-1]) > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((len(beatsDistribution)-1,-1))
    
    with open("verseInterval_"+uuid,"wb") as f:
        pickle.dump(verse_intervals, f)
    print("[VERSE_DETECTOR] verse detector ended.")
    endTime = time.time()
    print("[VERSE_DETECTOR] Time elasped (in secs): "+str(endTime-startTime))

##############

#test bucket uploading
os.environ["GOOGLE_APPLICATION_CREDENTIALS"] = "rapgo-storage.json"
STORAGE_BUCKET_NAME = os.environ.get('STORAGE_BUCKET_NAME', 'rapgo-bucket-2')

try:
    storage_client = storage.Client()
except:
    print("google.cloud storage could not be instantiated")

def upload(filename_local, filename_bucket):
    """
    Send the generated data to the bucket
    """
    print("starting uploading "+filename_local+"to bucket...")
    bucket = storage_client.get_bucket(STORAGE_BUCKET_NAME)
    blob = bucket.blob(filename_bucket)
    
    try:
        with open(filename_local, "rb") as f:
            blob.upload_from_file(f)
        print(filename_local+" uploaded !")
        return True
    except GoogleCloudError as e:
        print(e.errors())
        return False

def chunker(beatfile):
    """
    This is meant to avoid the timeout when we do an upload of a big file to the bucket
    """
    print("chunker running...")
    uuid = beatfile.split("_")[1].split(".")[0]
    rawBeat = pydub.AudioSegment.from_mp3(beatfile)
    chunkRatio = 0.01 # a chunk is of length 1/100 of the length of the entire file 
    rawBeatLength = len(rawBeat)
    if rawBeatLength <= 10000:
        print("No need for chunker : the file is too small.")
        return False

    chunkSize = int(rawBeatLength*chunkRatio)
    chunks = []
    idx = 0
    while idx < rawBeatLength:
        if idx < rawBeatLength-chunkSize:
            chunks.append(rawBeat[idx:idx+chunkSize])
        else:
            chunks.append(rawBeat[idx:])
        idx += chunkSize
    print("rawBeat splitted into "+str(len(chunks))+" chunks")
    for i in range(len(chunks)):
        idx = ""
        if i < 10:
            idx = "00"+str(i)
        elif i >= 10 and i <= 99:
            idx = "0"+str(i)
        else:
            idx = str(i)

        with open("beatchunk#"+idx+"_"+uuid,"wb") as f:
            pickle.dump(chunks[i], f)

def assembler():
    """
    Reassemble the chunks together
    """
    chunkList = glob.glob("beatchunk*")
    uuid = chunkList[0].split("_")[1]
    chunkList.sort()
    res = pydub.AudioSegment.empty()
    for c in chunkList:
        res += pickle.load(open(c,"rb"))
    #exportation
    res.export("assembled_"+uuid+".mp3", format="mp3")

def run(beatfile):
    getBeats(beatfile)
    getBPM(beatfile)

    uuid = beatfile.split("_")[1].split(".")[0]
    beatfileDuration = pickle.load(open("duration_"+uuid,"rb"))
    beatsDistribution = pickle.load(open("tempDist_"+uuid,"rb"))

    verse_detector(beatfile, beatsDistribution, beatfileDuration)

    #chunk the rawBeat file since we may encounter timeout in uploading a large file.
    chunker(beatfile)

    #upload
    for file in glob.glob("*_"+uuid+"*"):
        if file != "beat_"+uuid+".mp3": #we do not upload the mp3 since we would timeout
            upload(file,file)

if __name__ == "__main__":
    # for aubio : sudo apt-get install python3-aubio python-aubio aubio-tools

    # list all the available beat in the current folder to be ingested.
    beats = glob.glob("beat_*.mp3") # for now we take only one file
    #threadList = []
    #run the ingestion 
    #for b in beats:
    #    threadList.append(threading.Thread(target=run,args=(b,)))
    #for t in threadList:
    #    t.start()
    startTime = time.time()
    run(beats[1])
    endTime = time.time()
    print("song ingested in : "+str(endTime-startTime)+" secs.")
