# rap generator service

## Reminder : 

* for the google cloud bucket connection, do not forget to `export GOOGLE_APPLICATION_CREDENTIALS="/<PATH_TO_JSON>/keen-dispatch.json"` and `export GOOGLE_SOUND_BUCKET="rapgo-bucket-1"` and `export DATA_FOLDER="/data/sounds/"` in the Kubernetes pod deployment.


## TODO : 

* voice chunks should be the same length or almost the same length (use the phase vocoder or split again)
* Increase the volume of the voice on the initial beat
