import os 
import glob
import fnmatch
import random
from flask import Flask, flash, request, redirect, url_for, jsonify
from werkzeug.utils import secure_filename
from lib.musicAssembler import MusicAssembler
from google.cloud import storage
import threading

app = Flask(__name__)

try:
    storage_client = storage.Client()
except:
    print("google.cloud storage could not be instantiated")

from logger import getJSONLogger
logger = getJSONLogger('rapgeneratorservice-server')

BUCKET_NAME = os.environ.get('BUCKET_NAME', 'rapgo-bucket-1')
DATA_FOLDER = os.environ.get('SOUNDS_FOLDER', 'voiceTempStorage/')

def clean_storage(uuidBeatData, voicefilename):
    """
    Delete the temporary voicefilename and the beat data containing the uuidBeatData from the DATA_FOLDER 
    """
    try:
        os.remove(DATA_FOLDER+voicefilename)
        print("voice file "+voicefilename+" deleted")
    except:
        print("error in deleting the voice file "+voicefilename)
    files = glob.glob(DATA_FOLDER+"*")
    for filename in fnmatch.filter(files, uuidBeatData):
        try:
            os.remove(DATA_FOLDER+filename)
            print("file "+filename+" deleted.")
        except:
            print("error in deleting the file "+filename)

def to_bucket(filename_local, filename_bucket):
    """
    Send the generated data to the bucket
    """
    bucket = storage_client.get_bucket(BUCKET_NAME)
    blob = bucket.blob(filename_bucket)
    try:
        blob.upload_from_filename(filename_local)
        return True
    except:
        return False

def bucket_download(bucket, source_filename, destination_file_name):
    blob = bucket.blob(source_filename)
    blob.download_to_filename(destination_file_name)

def getRandomBeatData():
    '''
    Connect to bucket and retrieve a random beatFile with its associated metadata.
    Then, get the binaries of the file and save it inside /data/sounds/ folder with the name
    `beat_<filenameUUID>.mp3`. Finally, return `beat_<filenameUUID>.mp3` as a string.
    '''
    blobs = storage_client.list_blobs(BUCKET_NAME, prefix="beat_")
    random_uuid = random.choice([blob.name for blob in blobs]).split("_")[1].split(".")[0]

    bucket = storage_client.get_bucket(BUCKET_NAME)
    metadata_prefixes = ["duration_", "bpm_", "sound_", "tempDist_", "tempInt_", "verseInterval_"]
    thread_list = list()
    for p in metadata_prefixes:
        if p == "sound_":
            source_filename = "beat_"+p+random_uuid+".mp3"
        else: # it's binary objects
            source_filename = "beat_"+p+random_uuid
        t = threading.Thread(target=bucket_download, args=(bucket, source_filename, DATA_FOLDER+source_filename,))
        thread_list.append(t)
        t.start()

    for idx, t in enumerate(thread_list):
        t.join()
        print("download #"+str(idx)+" ended.")
    
    return random_uuid

@app.route("/<voicefilename>", methods=['POST'])
def rapgenerator(voicefilename):
    """
    The incoming request comes from the audio streaming microservice which communicates the registered filename (as a string)
    of the voice file uploaded in the bucket.
    """
    if request.method == 'POST':
        uuidBeatData = getRandomBeatData() # download beat data from the bucket
        bucket = storage_client.get_bucket(BUCKET_NAME)
        bucket_download(bucket, voicefilename, DATA_FOLDER+voicefilename) #download the voice from the bucket and store it into DATA_FOLDER/filename

        # call MusicAssembler and run the model
        ma = MusicAssembler(uuidBeatData, voicefilename)
        outputfilename = ma.run()
        if (outputfilename):
            to_bucket(outputfilename, outputfilename)
            clean_storage(uuidBeatData, voicefilename)
            return jsonify(
                statusCode=200,
                outputfilename=outputfilename 
            )
        else:
            return jsonify(
                statusCode=500
            )
    else:
        return jsonify(
                statusCode=404 
            )

if __name__ == "__main__":
    app.run() # need to run generate_cert.go here to get cert.pem and key.pem