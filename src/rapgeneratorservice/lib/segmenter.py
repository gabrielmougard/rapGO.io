import sys
from aubio import onset, source
from numpy import hstack, zeros

def getBeats(filename):
    win_s = 512
    hop_s = win_s // 2
    samplerate = 44000

    s = source("../../bin/musicParser/data/"+filename,samplerate, hop_s)
    samplerate = s.samplerate

    o = onset("default", win_s, hop_s, samplerate)

    onsets = []
    desc = []
    tdesc = []
    downsample = 2 #to plot n samples / hop_s
    allsamples_max = zeros(0,)
    #total number of frames read
    total_frames = 0
    while True:
        samples, read = s()
        if o(samples):
            onsets.append(o.get_last())
        new_maxes = (abs(samples.reshape(hop_s//downsample, downsample))).max(axis=0)
        allsamples_max = hstack([allsamples_max, new_maxes])
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
    return (desc_times_filtered,desc_plot_filtered,duration)

def verse_detector(beatsTimeline,duration):
    """
    Return an array of tuple in which you have
    the starting time of a verse and the ending time
    of the verse.
    """
    verse_intervals = []
    MINIMUM_VERSE_LENGTH = 5.0 # if the interval is greater or equal than 5s, it's a verse.
    if beatsTimeline[0] > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((0.0,beatsTimeline[0]))
    for i in range(0,len(beatsTimeline)-1):
        if abs(beatsTimeline[i]-beatsTimeline[i+1]) >= MINIMUM_VERSE_LENGTH:
            verse_intervals.append((beatsTimeline[i],beatsTimeline[i+1]))
    if abs(duration-beatsTimeline[-1]) > MINIMUM_VERSE_LENGTH:
        verse_intervals.append((beatsTimeline[-1],duration))
    
    return verse_intervals

(desc_times_filtered,desc_plot_filtered,duration) = getBeats("song_0dcf4fae-65a9-403e-88e2-829215421c32.mp3")
print(verse_detector(desc_times_filtered,duration))