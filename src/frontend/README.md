# frontend service

The frontend part let the user generating a rap voice with his voice using a nice UI.

## Test with Docker
just execute the following commands :

* `docker build . -t frontend`
* `docker run -it -p 80:80 frontend`

(The Docker container has been optimized using a production-ready image using multistage builds.) 