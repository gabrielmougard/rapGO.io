# rapGO.co : an experimental platform for generating a rap from any input voice.

When I initially started this project, I wanted to dig deeper in the audio manipulation of sounds. 
Then I wanted to add a scalable infrastructure layer using Kubernetes and a distributed filesystem.

This project is a kind of training for deploying a complex kubernetes application with several microservices
powered by several languages like Go, C (coming soon for the rapgenration service), Python and Reactjs for the frontend services.

For now, the `src/rapgeneratorservice/lib` which is the core of the application is a bit simple and the results are very far from being audible. However, I'm improving this core library over the time to create a better user experience.

To learn more about the architecture of rapGO.io, please refer to the [docs](/docs/README.md) section.

## Demo

The demo (available [here](https://drive.google.com/file/d/1XtkA3sAqpnCYaEMn2g3RyA2EgFZA9w-q/view?usp=sharing)) ha been tested with the `test/alpha` version on docker-compose. This is very far from being perfect (especially the core algorithm), but the global infrastructure is working.

For the ingestion engine, wa can have a idea of what it will look like on this [video](https://drive.google.com/file/d/1oq8C7PLc_7nMztTs38DcD6P16b-raWFx/view?usp=sharing). I'm still working on it and also on the ingestion service backend.

## Kubernetes deployment on GKE & images repository
The docker images used for this project are :

* `rapgoio/frontend:alpha`
* `rapgoio/bucketserv:alpha`
* `rapgoio/convertserv:alpha`
* `rapgoio/heartbeatserv:alpha`
* `rapgoio/rapserv:alpha`

For the ingestion engine, the available images are :

* `rapgoio/ingest-frontend:alpha`