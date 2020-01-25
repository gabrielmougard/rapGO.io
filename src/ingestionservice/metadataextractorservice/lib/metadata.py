"""
metadata.py : each function compute an important variable
that will be useful for the rap generation.
"""

import time
import pickle
import numpy as np
from aubio import onset, source, tempo

def get_beats(beatfile):
    uuid = beatfile.split("_")[1].split(".")[0]
    start_time = time.time()
    print("[BEATS RECOGNIZER] Starting beats recognition...")
    win_s = 512
    hop_s = win_s // 2
    samplerate = 44000
    s = source(beatfile, samplerate, hop_s)
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

    desc_times = [float(t) * hop_s / samplerate for t in range(len(desc))]
    desc_max = max(desc) if max(desc) != 0 else 1.
    desc_plot = [d / desc_max for d in desc]
    CONSIDERATION_THRESHOLD = 0.7
    desc_times_filtered = []
    desc_plot_filtered = []
    for x, y in zip(desc_times, desc_plot):
        if y > CONSIDERATION_THRESHOLD:
            desc_times_filtered.append(x)
            desc_plot_filtered.append(y)
    allsamples_max = (allsamples_max > 0) * allsamples_max
    allsamples_max_times = [float(t)*hop_s/downsample/samplerate for t in range(len(allsamples_max))]
    duration = allsamples_max_times[-1]

    #update object
    beatfile_duration = duration
    beats_distribution = desc_times_filtered
    beats_intensity = desc_plot_filtered
    #write the pickle binaries
    with open("duration_"+uuid, "wb") as f:
        pickle.dump(beatfile_duration, f)
    with open("tempDist_"+uuid, "wb") as f:
        pickle.dump(beats_distribution, f)
    with open("tempInt_"+uuid, "wb") as f:
        pickle.dump(beats_intensity, f)
    print("[BEATS RECOGNIZER] Beats recognition ended.")
    end_time = time.time()
    print("[BEATS RECOGNIZER] Time elasped (in secs): "+str(end_time-start_time))
    return beats_distribution, beatfile_duration

def get_verse_intervals(beatfile, beats_distribution, beatfile_duration):
    """
    Return an array of tuple in which you have
    the starting time of a verse and the ending time
    of the verse.
    """
    uuid = beatfile.split("_")[1].split(".")[0]
    start_time = time.time()
    print("[VERSE_DETECTOR] Starting verse detector of voice file...")
    verse_intervals = []
    MINIMUM_VERSE_LENGTH = 5.0 # if the interval is greater or equal than 5s, it's a verse.
    if beats_distribution[0] > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((-1, 0))
    for i in range(0, len(beats_distribution)-1):
        if abs(beats_distribution[i]-beats_distribution[i+1]) >= MINIMUM_VERSE_LENGTH:
            verse_intervals.append((i, i+1))
    if abs(beatfile_duration-beats_distribution[-1]) > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((len(beats_distribution)-1, -1))

    with open("verseInterval_"+uuid, "wb") as f:
        pickle.dump(verse_intervals, f)
    print("[VERSE_DETECTOR] verse detector ended.")
    end_time = time.time()
    print("[VERSE_DETECTOR] Time elasped (in secs): "+str(end_time-start_time))

def get_bpm(beatfile):
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
    with open("bpm_"+uuid, "wb") as f:
        pickle.dump(b, f)
