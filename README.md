# My minimal image api


## Setup:
* build image: `docker build --tag little-api .`
* run image as container: `docker run little-api --name my-api`

!todo: use blob store for caching

### Local:
* deploy postgres db with
`docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=mini-db -p 9920:5432 -d postgres`


* Postman requests can be found in `postman_collection.json` and in `*.http` files.


