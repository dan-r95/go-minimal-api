# taken from: https://docs.docker.com/language/golang/build-images/
FROM ubuntu:22.04

RUN yes | apt-get update
RUN yes | apt-get upgrade

# install curl to download go
RUN yes | apt-get install curl

# install golang
RUN curl -s https://storage.googleapis.com/golang/go1.18.3.linux-amd64.tar.gz| tar -v -C /usr/local -xz  > /dev/null
ENV PATH="/usr/local/go/bin:${PATH}"

ENV  DB_HOST=db
ENV POSTGRES_USER=postgres
ENV POSTGRES_PASSWORD=mysecretpassword
ENV POSTGRES_DB=mini-db
ENV POSTGRES_PORT=5432

# build process
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o /little-api .
EXPOSE 8090

CMD [ "/little-api" ]