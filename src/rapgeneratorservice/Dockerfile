FROM python:3.7-slim-buster

WORKDIR /rapgenerator
COPY . .

#install python dependencies (in requirements.txt)
#RUN apt-get install python3-aubio python-aubio aubio-tools
RUN apt-get update \
 && apt-get install -y \
 coreutils \
 gcc \
 libogg-dev \
 libass-dev \
 libvpx-dev \
 libvorbis-dev \
 libwebp-dev \
 libtheora-dev \
 pkg-config \
 wget \
 yasm \
 ffmpeg


RUN pip3 install -r requirements.txt

#Environment variables
ENV KAFKA_TOCORE_TOPIC toCore
ENV KAFKA_TOBUCKET_TOPIC toBucket
ENV KAFKA_TOHEARTBEAT_TOPIC toHeartbeat
ENV KAFKA_BROKER kafka:9092
ENV KAFKA_GROUP_ID go-kafka-consumer

ENV STORAGE_BUCKET_NAME rapgo-bucket-2
ENV STORAGE_PROJECT_ID tidy-hold-247521
ENV GOOGLE_APPLICATION_CREDENTIALS /rapgenerator/rapgo-storage.json

ENV TMP_FOLDER /data/tmp/

ENTRYPOINT ["python3", "rapgenerator.py"]
