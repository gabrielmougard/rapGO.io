
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
    bucket = storage_client.get_bucket(STORAGE_BUCKET_NAME)
    blob = bucket.blob(filename_bucket)
    try:
        blob.upload_from_filename(filename_local)
        return True
    except:
        return False

def bucket_download(storage_client, source_filename, destination_file_name, voiceUUID):
    #if (source_filename.split("_")[0] == "beatchunk"):
    #    print("Starting downloading BEATCHUNKS...")
    #    #for the beatchunks, we do a search by prefix and download all the chunks
    #    exploded = source_filename.split("_")
    #    blobs = list(storage_client.list_blobs(STORAGE_BUCKET_NAME, prefix=exploded[0]+"_"+exploded[1]+"_"))
    #    blobsExploded = [blobs[i:i+10] for i in range(0, len(blobs)-9, 10)] # for concurrency purpose
    #    threadList = []
#
    #    # we will start 10 worker threads to download concurrently the beatchunks
    #    def downloadBlobsSlice(blobsSlice, voiceUUID):
    #        for b in blobsSlice:
    #            b.download_to_filename(TMP_FOLDER+METADATA_FOLDER_PREFIX+voiceUUID+"/"+b.name)
    #            print(b.name+" has been downloaded to "+TMP_FOLDER+METADATA_FOLDER_PREFIX+voiceUUID+"/"+b.name)
#
    #    startTime = time.time()
    #    for blobsSlice in blobsExploded:
    #        threadList.append(threading.Thread(target=downloadBlobsSlice, args=(blobsSlice,voiceUUID,)))
    #    for t in threadList:
    #        t.start()
    #    for t in threadList:
    #        t.join()
    #    print("temps ecoule : "+str(time.time()-startTime))
    #    
    #else:
    bucket = storage_client.get_bucket(STORAGE_BUCKET_NAME)
    blob = bucket.blob(source_filename)
    blob.download_to_filename(destination_file_name)
    print(blob.name+" has been downloaded to "+destination_file_name)

def getRandomBeatData(producer, voiceUUID, storage_client):
    '''
    Connect to bucket and retrieve a random beatFile with its associated metadata.
    Then, get the binaries of the file and save it inside TMP_FOLDER folder with the name
    `beat_<filenameUUID>.mp3`. Finally, return `beat_<filenameUUID>.mp3` as a string.
    '''
    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Fetching metadata...")

    blobs = storage_client.list_blobs(STORAGE_BUCKET_NAME, prefix="beat_")
    random_uuid = random.choice([blob.name for blob in blobs]).split("_")[1].split(".")[0]
    print("the chosen random_uuid is : "+random_uuid)
    metadata_prefixes = ["duration_", "bpm_", "beat_", "tempDist_", "tempInt_", "verseInterval_"]
    thread_list = list()
    
    for p in metadata_prefixes:
        source_filename = p+random_uuid
        if p == "beat_":
            source_filename += ".mp3"
        t = threading.Thread(target=bucket_download, args=(storage_client, source_filename, TMP_FOLDER+METADATA_FOLDER_PREFIX+voiceUUID+"/"+source_filename, voiceUUID,))
        thread_list.append(t)
        t.start()

    for t in thread_list:
        t.join()


    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Metadata fetched successfully !")
    
def consume_loop(consumer, producer, topics, storage_client):
    print("starting consumming messages...")
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
                print("message detected : "+str(msg.value()))
                if msg.topic() == KAFKA_TOCORE_TOPIC:
                    threading.Thread(target=processToCoreMsg, args=(msg,producer,storage_client,)).start()
                else:
                    sys.stderr.write('topic %s is not recognized for now.\n' % msg.topic())
    finally:
        # Close down consumer to commit final offsets.
        consumer.close()

def processToCoreMsg(message, producer, storage_client):
    """
    Consume our kafka core message
    """
    print("starting processing message "+str(message.value()))
    voiceUUID = str(message.value()).split("_")[1].split(".")[0]

    #create metadata folder
    metadataPath = TMP_FOLDER+METADATA_FOLDER_PREFIX+voiceUUID
    os.mkdir(metadataPath)

    producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Starting core processing...")
    getRandomBeatData(producer, voiceUUID, storage_client) #fetching the needed metadata and write them in TMP_FOLDER/metadata_<UUID>/ folder

    # MusicAssembler in the core processing class
    ma = MusicAssembler(str(message.value()))
    finished = ma.run() #Run the process. If successful, it return the output filename which should be 'output_<UUID>.mp3'

    if (finished):
        producer.produce(KAFKA_TOHEARTBEAT_TOPIC, key=voiceUUID, value="Rap generated successfully !")
        #to_bucket(outputfilename, outputfilename) #the bucketuploaderservice handle it automatically. The output filename should be : 'output_<UUID>.mp3'
        clean_storage(voiceUUID) #delete the metadata folder

def shutdown():
    running = False

if __name__ == "__main__":
    # env variables
    global KAFKA_TOCORE_TOPIC
    global KAFKA_TOBUCKET_TOPIC
    global KAFKA_TOHEARTBEAT_TOPIC
    global KAFKA_BROKER
    global KAFKA_GROUP_ID
    global STORAGE_BUCKET_NAME
    global TMP_FOLDER
    global METADATA_FOLDER_PREFIX

    KAFKA_TOCORE_TOPIC = os.environ.get("KAFKA_TOCORE_TOPIC")
    KAFKA_TOBUCKET_TOPIC = os.environ.get("KAFKA_TOBUCKET_TOPIC")
    KAFKA_TOHEARTBEAT_TOPIC = os.environ.get("KAFKA_TOHEARTBEAT_TOPIC")
    KAFKA_BROKER = os.environ.get("KAFKA_BROKER")
    KAFKA_GROUP_ID = os.environ.get("KAFKA_GROUP_ID")
    STORAGE_BUCKET_NAME = os.environ.get('STORAGE_BUCKET_NAME')
    TMP_FOLDER = os.environ.get('TMP_FOLDER')
    METADATA_FOLDER_PREFIX = os.environ.get("METADATA_FOLDER_PREFIX")

    #instantiate the google bucket client
    try:
        storage_client = storage.Client()
        print("google.cloud storage client is alive.")
    except:
        print("google.cloud storage could not be instantiated")
        exit(1)
    
    print("Waiting for kafka cluster to start...")
    time.sleep(50) # secured delay for the kafka cluster to setup (leader election can take some time)
    
    #setup the producer and the consumer
    confProducer = {'bootstrap.servers': KAFKA_BROKER, 'client.id': socket.gethostname()}
    confConsumer = {'bootstrap.servers': KAFKA_BROKER, 'group.id': KAFKA_GROUP_ID, 'auto.offset.reset': 'smallest'}
    
    global running
    running = True #indicate whether the consumer is listening or not

    # Instantiate kafka producer and consumer
    producer = Producer(confProducer)
    print("kafka producer is alive")
    consumer = Consumer(confConsumer)
    print("kafka consumer is alive")

    consume_loop(consumer, producer, [KAFKA_TOCORE_TOPIC], storage_client)