# rap generator service

## Reminder : 

* for the google cloud bucket connection, do not forget to `export GOOGLE_APPLICATION_CREDENTIALS="/<PATH_TO_JSON>/keen-dispatch.json"` and `export GOOGLE_SOUND_BUCKET="rapgo-bucket-1"` and `export DATA_FOLDER="/data/sounds/"` in the Kubernetes pod deployment.


## TODO : 

* voice chunks should be the same length or almost the same length (use the phase vocoder or split again)
* Increase the volume of the voice on the initial beat

## After the MVP

* The rap generation task is quite CPU consuming. Thus, we would like to distibute the computation. That's why we will use the python distributed computing library [Ray](https://ray.io/) on top of Kubernetes. The idea is the following :
  * inside Kubernetes cluster, create a "Ray" namespace (or create an other Kubernetes cluster for Ray workers only) 
  * The Ray master node consume the Kafka topic "toCore" and launch a task on a healthy worker node.