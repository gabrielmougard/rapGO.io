
import os
import glob
import fnmatch
import random
import socket
import time
import sys
import threading

from google.cloud import storage
from confluent_kafka import Consumer #we will use the Confluent Kafka client instead (I find it more understandable...)
from confluent_kafka import Producer 

from lib.musicAssembler import MusicAssembler #The core™ lib !!

print("Waiting for kafka cluster to start...")
time.sleep(50) # secured delay for the kafka cluster to setup (leader election can take some time)

# env variables
KAFKA_TOCORE_TOPIC = os.environ.get("KAFKA_TOCORE_TOPIC","toCore")
KAFKA_TOBUCKET_TOPIC = os.environ.get("KAFKA_TOBUCKET_TOPIC","toBucket")
KAFKA_TOHEARTBEAT_TOPIC = os.environ.get("KAFKA_TOHEARTBEAT_TOPIC","toHeartbeat")
KAFKA_BROKER = os.environ.get("KAFKA_BROKER","kafka:9092")
KAFKA_GROUP_ID = os.environ.get("KAFKA_GROUP_ID","go-kafka-consumer")
STORAGE_BUCKET_NAME = os.environ.get('STORAGE_BUCKET_NAME', 'rapgo-bucket-2')
TMP_FOLDER = os.environ.get('TMP_FOLDER', '/data/tmp/')

#setup the producer and the consumer
confProducer = {'bootstrap.servers': KAFKA_BROKER, 'client.id': socket.gethostname()}
confConsumer = {'bootstrap.servers': KAFKA_BROKER, 'group.id': KAFKA_GROUP_ID, 'auto.offset.reset': 'smallest'}
running = True #indicate whether the consumer is listening or not

# Instantiate kafka producer and consumer
producer = Producer(confProducer)
consumer = Consumer(confConsumer)

#instantiate the google bucket client
try:
    storage_client = storage.Client()
except:
    print("google.cloud storage could not be instantiated")

def clean_storage(voiceUUID):
    """
    Delete the temporary folder 'metadata_<voiceUUID>
    """
    metadata_folder_prefix = os.environ.get("METADATA_FOLDER_PREFIX","metadata_")
    for root, dirs, files in os.walk(TMP_FOLDER+metadata_folder_prefix+voiceUUID, topdown=False):
        for name in files:
            os.remove(os.path.join(root, name))
        for name in dirs:
            os.rmdir(os.path.join(root, name))

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

def getRandomBeatData(producer, voiceUUID):
    '''
    Connect to bucket and retrieve a random beatFile with its associated metadata.
    Then, get the binaries of the file and save it inside /data/sounds/ folder with the name
    `beat_<filenameUUID>.mp3`. Finally, return `beat_<filenameUUID>.mp3` as a string.
    '''
    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Fetching metadata...")

    blobs = storage_client.list_blobs(STORAGE_BUCKET_NAME, prefix="beat_")
    random_uuid = random.choice([blob.name for blob in blobs]).split("_")[1].split(".")[0]
    metadata_folder_prefix = os.environ.get("METADATA_FOLDER_PREFIX","metadata_")
    bucket = storage_client.get_bucket(STORAGE_BUCKET_NAME)
    metadata_prefixes = ["duration_", "bpm_", "sound_", "tempDist_", "tempInt_", "verseInterval_"]
    thread_list = list()
    for p in metadata_prefixes:
        if p == "sound_":
            source_filename = "beat_"+p+random_uuid+".mp3"
        else: # it's binary objects
            source_filename = "beat_"+p+random_uuid
        t = threading.Thread(target=bucket_download, args=(bucket, source_filename, TMP_FOLDER+metadata_folder_prefix+voiceUUID+"/"+source_filename,))
        thread_list.append(t)
        t.start()

    for idx, t in enumerate(thread_list):
        t.join()
        print("download #"+str(idx)+" ended.")

    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Metadata fetched successfully !")
    
def consume_loop(consumer, producer, topics):
    try:
        consumer.subscribe(topics)

        while running:
            msg = consumer.poll(timeout=1.0)
            if msg is None: continue

            if msg.error():
                if msg.error().code() == confluent_kafka.KafkaError._PARTITION_EOF:
                    # End of partition event
                    sys.stderr.write('%% %s [%d] reached end at offset %d\n' %
                                     (msg.topic(), msg.partition(), msg.offset()))
                elif msg.error():
                    raise confluent_kafka.KafkaException(msg.error())
            else:
                if msg.topic() == KAFKA_TOCORE_TOPIC:
                    threading.Thread(target=processToCoreMsg, args=(msg,producer,)).start()
                else:
                    sys.stderr.write('topic %s is not recognized for now.\n' % msg.topic())
    finally:
        # Close down consumer to commit final offsets.
        consumer.close()

def processToCoreMsg(message, producer):
    """
    Consume our kafka core message
    """
    voiceUUID = message.value().split("_")[1].split(".")[0]
    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Starting core processing...")
    getRandomBeatData(producer, voiceUUID) #fetching the needed metadata and write them in TMP_FOLDER/metadata_<UUID>/ folder

    # MusicAssembler in the core processing class
    ma = MusicAssembler(message.value(), producer)
    outputfilename = ma.run() #Run the process. If successful, it return the output filename which should be 'output_<UUID>.mp3'

    if (outputfilename):
        producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Rap generated successfully !")
        producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Uploading result to cloud...")
        to_bucket(outputfilename, outputfilename)
        producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Upload successful !")
        clean_storage(voiceUUID) #delete the metadata folder

def shutdown():
    running = False

if __name__ == "__main__":
    consume_loop(consumer, producer, [KAFKA_TOCORE_TOPIC])