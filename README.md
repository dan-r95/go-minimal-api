# PIXXIO Backend Homecase

Instructions:

1. Clone this github repository
2. Compile the code using `go build` and execute the generated binary which starts the webserver
3. Call [localhost:8090/ping](http://localhost:8090/ping) in your browser to ensure that the running webserver is reachable
4. Modify the code according to the 'Homecase' instructions and tasks given to you by pixx.io

You are allowed to change everything in the code and use external libraries of your choice.



## Setup:
* build image: `docker build --tag little-api .`
* run image as container: `docker run little-api --name my-api`

!todo: use blob store for caching
