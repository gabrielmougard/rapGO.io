import numpy as np
import array
import pydub
from pydub.playback import play
import time

import phase_vocoder as pvoc

# test of vocoder when accelerating/deccelerating speed of chunk file
chunk = pydub.AudioSegment.from_mp3("chunks/chunk39.mp3")
chunkSecondPart = chunk[len(chunk)//2:]

chunkSecondPartRaw = chunkSecondPart.get_array_of_samples()
chunkSecondPartRaw = np.array(chunkSecondPartRaw).astype('float32')
startTime = time.time()
stftChunk = pvoc.stft(chunkSecondPartRaw)
endTime = time.time()
print("time elapsed for the STFT execution (in secs) :"+str(endTime-startTime))
print(stftChunk.shape)
startTime = time.time()
vocodedChunk = pvoc.phase_vocoder(stftChunk,2.0)
endTime = time.time()
print("time elapsed for the phase_vocoder execution (in secs) :"+str(endTime-startTime))
print(vocodedChunk.shape)
startTime = time.time()
finalChunk = pvoc.istft(vocodedChunk).astype('int16')
endTime = time.time()
print("time elapsed for the ISTFT execution (in secs) :"+str(endTime-startTime))
print(finalChunk[:20])

finalChunkArray = np.ravel(finalChunk)
vocodedSound = chunkSecondPart._spawn(finalChunkArray)

result = chunk[:len(chunk)//2].append(vocodedSound)
play(result)



