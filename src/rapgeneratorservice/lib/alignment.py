import phase_vocoder as pvoc
import numpy as np 
import math
import random
import pydub 
from pydub.playback import play

def align(c,endT,beatDistList,beatCounted):
    nearestTime = binarySearch(endT,beatDistList)
    xL, xR = None, None
    if (nearestTime > endT):
        xR = nearestTime
    else:
        xL = nearestTime
    print("endT : "+str(endT))
    print("XR : "+str(xR))
    print("XL : "+str(xL))

    #First step : slicing
    sliceIdx = 0
    if (xR):
        # do the same as bellow but for xR
        deltaT = int((xR-endT)*1000) # deltaT > int(0)
        randomizationFactor = 0
        rightBoundary = abs((len(c)//2)//deltaT)
        if rightBoundary == 0:
            randomizationFactor = 1
        else:
            randomizationFactor = random.randint(1, rightBoundary)
        sliceIdx = len(c) - randomizationFactor*deltaT # sliceIdx < len(c) since deltaT > 0 and randomizationFactor >= 1

    else:
        deltaT = int((xL-endT)*1000) # deltaT < int(0)
        # get the randomizationFactor
        randomizationFactor = 0
        rightBoundary = abs((len(c)//2)//deltaT)
        if rightBoundary == 0:
            randomizationFactor = 1
        else:
            randomizationFactor = random.randint(1, rightBoundary)
        
        sliceIdx = len(c) + randomizationFactor*deltaT # sliceIdx < len(c) since deltaT < 0 and randomizationFactor >= 1

    print("DURATION : "+str(c.duration_seconds))
    print("IDX : "+str(sliceIdx))
    print("LEN C : "+str(len(c)))
    leftChunkSlice = c[:len(c)-sliceIdx] 
    rightChunkSlice = c[len(c)-sliceIdx:]

    #Stretching Problem
    #Purpose of this problem is to calculate "r_stretch"
    #Second step : stretching rightChunkSlice
    r_stretch = 1.0
    if (xR):
        r_stretch = len(rightChunkSlice)/(len(rightChunkSlice)+abs(int((xR-endT)*1000))) # r_strecth < 1
    else:
        r_stretch = (len(rightChunkSlice)+abs(int((xL-endT)*1000)))/len(rightChunkSlice) # r_stretch > 1

    #rightChunkSliceRaw = rightChunkSlice.get_array_of_samples()
    #rightChunkSliceRaw = np.array(rightChunkSliceRaw).astype('float32') # convert the type to float32
    #rightStftChunk = pvoc.stft(rightChunkSliceRaw)
    #rightStftChunkVocoded = pvoc.phase_vocoder(rightStftChunk, r_stretch)
    #rightIstftChunk = pvoc.istft(rightStftChunkVocoded).astype('int16') # was not working for int16 so I chose s104
    #rightFinalChunkArray = np.ravel(rightIstftChunk)
    #rightVocodedSound = rightChunkSlice._spawn(rightFinalChunkArray)
    #
    #leftChunkSlice.append(rightVocodedSound, crossfade=0) # merge the original part and the stretched one with crossfade 0 to avoid crossfade errors
    #play(leftChunkSlice)

    newBeatCountedIdx = 0
    if (xR): # calculate the number of beat in the final stretched chunk
        newBeatCountedIdx = beatDistList.index(xR)
    else:
        newBeatCountedIdx = beatDistList.index(xL) 

    #return leftChunkSlice, newBeatCountedIdx
    print("newBeatCountedIdx : "+str(newBeatCountedIdx))
    return newBeatCountedIdx

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
    voiceChunkAdjusted = setToTargetLevel(chunkObject,1.0)
    #beatChunkAdjusted = setToTargetLevel(beatPart,-6.0)

    combined = beatObject.overlay(voiceChunkAdjusted, position=startT)
    return combined