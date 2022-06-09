# minimal image api :rocket:


## Docker Image:
* build image: `docker build --tag little-api .`
* run containers via docker-compose: `docker-compose up`
* Can also run image as container: `docker run --name my-api little-api`.


### Local:
* deploy postgres db with
`docker run --name postgres -e POSTGRES_PASSWORD=mysecretpassword -e POSTGRES_DB=mini-db -p 5432:5432 -d postgres`
* Run the main file. `source user.env && go build go-minimal-api`. This assumes you have the environment variables set in `user.env` on your system
* For GoLand this is not needed. Uses [Env file plugin for variables](https://plugins.jetbrains.com/plugin/7861-envfile).
* Postman requests can be found in `postman_collection.json` and in `*.http` files.
