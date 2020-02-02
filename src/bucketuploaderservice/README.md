# Bucket uploader service

The title is quite self-descriptive. Indeed, this service is meant to upload data to bucket while listening to the Kafka topic `toBucket`. It also returns some data in the `toHeartbeat` Kafka topic with the key `<UUID>` corresponding of the UUId of the input voice. This will be used by the heartbeat service 
to give an insight of the stage of the process on the frontend part. The heartbeat service will listen to this topic and send data through webSocket protocol to the frontend. 