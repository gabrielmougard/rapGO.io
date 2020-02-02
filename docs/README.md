# Global architecture of rapGO.co

## Overview
![the global architecture](0-Global-arch.png)

The image above describes the overall operations of the platform. This diagram represents the architecture of the version `test/alpha` which has been tested using `docker-compose`

## Explaination

In the early days of the idea (december 2019), I started thinking about a way to generate a "rap" of something that one could've said or written. For now, the user has to speak in the microphone module embedded in the browser start the generation process.
* When the voice is received, a websocket connection is created with the hearbeat service listening on a Kafka cluster for any incoming message to be sent to the client.
* Using a filewatcher system and a Kafka event-driven communication scheme, all the microservices are running independently to perform all the needed tasks. Pulling the metadata from the bucket for the generation, the generation itself, the upload of the new data and the emission of the heartbeats to the in-mem RB-tree storage... Everything is executed concurrently.
* Finally, the client receive the generated output in a nice react audio player. 

## Output demo

Since running the app imply that you should have a bucket configured on GCP, I provided a [test output](output_rapgo.mp3) that you can listen. This is the rap generation of the very famous speech of Martin Luther King.

## Requirements for running the `test/alpha` version.

* create a GCP account and create a project.
* create a google bucket (see google storage website)
* export your google cloud project configuration in a `json` file. Mine is called `rapgo-storage.json` look like this :

`"type": "service_account",
  "project_id": "tidy-hold-247521",
  "private_key_id": "<PRIVATE_KEY_ID>",
  "private_key": "-----BEGIN PRIVATE KEY---- <MY_PRIVATE_KEY> -----END PRIVATE KEY----",
  "client_email" : "rapgo-bucket@tidy-hold-247521.iam.gserviceaccount.com",
  "client_id": "<CLIENT_ID>",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/rapgo-bucket%40tidy-hold-247521.iam.gserviceaccount.com"
`

* put this json file in `src/bucketuploaderservice` and in `src/rapgeneratorservice` and change the name of the ENV variable inside the associated dockerfiles. 
* run the command `cd test/alpha && docker-compose build && docker-compose up`
* go to `http://localhost` for testing the app.

## What's next ?

* The most important will be to continuously increase the quality of the generation model (which is super simple at this time). I want to integrate my phase vocoder function for modifying the speed of the voice without changing the frequencies. Also, there will be some options allowing the user to generate a song according to a "build strategy" (you want to generate an output with the same style as Gucci-Gang ? You will be able to select this generation pattern)

* I have already started writing the Kubernetes manifests to deploy the application on GKE. However, since this is not a trivial application, it could take some time. Here is the K8s overview :
![the global architecture](GKE-topology.png)

* When the application and the ingestion engine will be deployed on GKE, I plan to deploy a distributed storage cluster to store my large files (currently stored in the google bucket), like `GlusterFS` or `CephFS` which have a higher throughput.

* At the same time, it will be necessary to develop classic user account management and monthly subscription functions as well as the marketplace that will be the place where users can share their creations. For now, I plan to use Google Spanner or Google BigTable since it's easier to have a managed service at the beginning.