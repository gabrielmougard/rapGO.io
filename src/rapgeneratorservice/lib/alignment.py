import phase_vocoder as pvoc
import numpy as np 
import math
import random
import pydub 

def align(c,endT,beatDistList,beatCounted):
    nearestTime = binarySearch(endT,beatDistList)
    xL, xR = None, None
    if (nearestTime > endT):
        xR = nearestTime
    else:
        xL = nearestTime

    #First step : slicing
    sliceIdx = 0
    if (xR):
        sliceIdx = randomChunkSlicer(xR,beatDistList)
    else:
        sliceIdx = randomChunkSlicer(xL,beatDistList)
    
    leftChunkSlice = c[:sliceIdx] 
    rightChunkSlice = c[sliceIdx:]

    #Stretching Problem
    #Purpose of this problem is to calculate "r_stretch"
    #Second step : stretching rightChunkSlice
    r_stretch = 1.0
    if (xR):
        r_stretch = len(rightChunkSlice)/(len(rightChunkSlice)+(xR-endT))
    else:
        r_stretch = (len(rightChunkSlice)+(endT-xL))/len(rightChunkSlice)
    
    rightChunkSliceRaw = rightChunkSlice.get_array_of_samples()
    rightChunkSliceRaw = np.array(rightChunkSliceRaw).astype('float32') # convert the type to float32
    rightStftChunk = pvoc.stft(rightChunkSliceRaw)
    rightStftChunkVocoded = pvoc.phase_vocoder(rightStftChunk,r_stretch)
    rightIstftChunk = pvoc.istft(rightStftChunkVocoded).astype('int16')
    rightFinalChunkArray = np.ravel(rightIstftChunk)
    rightVocodedSound = rightChunkSlice._spawn(rightFinalChunkArray)
    leftChunkSlice.append(rightVocodedSound) # merge the original part and the stretched one

    newBeatCountedIdx = 0
    if (xR): # calculate the number of beat in the final stretched chunk
        newBeatCountedIdx = beatDistList.index(xR)
    else:
        newBeatCountedIdx = beatDistList.index(xL) 

    return leftChunkSlice, newBeatCountedIdx

def binarySearch(t,beatDistList):
    """
    Find the nearest time in beatDistList for the target 't'
    run in O(lg(n))
    """
    if (t < beatDistList[0]):
        return beatDistList[0]
    if (t > beatDistList[len(beatDistList)-1]):
        return beatDistList[len(beatDistList)-1]
    
    lo = 0
    hi = len(beatDistList)-1

    while(lo <= hi):
        mid = (hi+lo)//2

        if (t < beatDistList[mid]):
            hi = mid - 1
        elif (t > beatDistList[mid]):
            lo = mid + 1
        else:
            return beatDistList[mid]
    return  beatDistList[lo] if ((beatDistList[lo] - t) < (t - beatDistList[hi])) else beatDistList[hi]

def randomChunkSlicer(xS,beatDistList):
    """
    return `sliceIdx`
    """
    cutIdx = len(beatDistList)-beatDistList.index(xS) 
    randomMultiple = random.randint(1,len(beatDistList)//(2*cutIdx))
    sliceIdx = len(beatDistList)-(cutIdx*randomMultiple)
    return sliceIdx

def setToTargetLevel(sound, targetLevel):
    difference = targetLevel - sound.dBFS
    return sound.apply_gain(difference)

def overlay(beatObject,chunkObject,startBeatIdx,endBeatIdx,beatDistList):
    """
    Overlay two sounds

    Parameters
    ----------
    beatObject : pydub.AudioSegment 
    chunkObject : pydub.AudioSegment
    startT : float (in secs)
    endT : float (in secs)
    ----------

    Return
    ------
    res : The entire beatObject with an overlayed part with the stretched voice chunk.
    Also, the volume of the beat part has been a bit lowered compared to the coice chunk
    to hear more about the voice.
    ------
    """
    res = beatObject
    startT = math.ceil(beatDistList[startBeatIdx]*1000) #convert in millisecs and convert to integer
    endT   = math.ceil(beatDistList[endBeatIdx]*1000) 

    beatPart = beatObject[startT:endT]
    #adjust the sound volume and put more gain on the voice chunk
    voiceChunkAdjusted = setToTargetLevel(chunkObject,0.0)
    beatChunkAdjusted = setToTargetLevel(beatPart,-6.0)

    combined = beatChunkAdjusted.overlay(voiceChunkAdjusted)
    res[startT:endT] = combined
    return res