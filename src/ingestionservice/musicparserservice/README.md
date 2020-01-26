# Music Parser

## API description

### Note : 
For the `POST` requests, if the sources are coming from freemusicarchive.org, and if the source already exists in the bucket, it will not be ingested. We do not want any duplicate data.
Also, you will be notified if the bucket can handle the upload (if there's not enough space available in it for example)

* **POST** : `/ingest/url`

* **POST** : `/ingest/urlfile`

* **POST** : `/ingest/rawfile`

* **POST** : `/ingest/genre/{genre}`
  * **parameters** :
    * **limit** : (int) describe the max number of file to be parsed by genre in the source website (freemusicarchive.org)
    * **random** (boolean) choose randomly the beats in the {genre} category. If it set to `False`, it will parse the found URL in order.
* **GET** : `/ingest/genre`
  * returns the list of genres that we can parse on freemusicarchive.org 
* **GET (Websocket later)** : `/bucket/data`
  * returns the total space and the available space, the number of objects specified by types, etc...

## TODO

* make a `swagger.yaml` description of this API.  