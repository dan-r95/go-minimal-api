# taken from: https://docs.docker.com/language/golang/build-images/
FROM ubuntu:22.04

#TODO: add database
RUN yes | apt-get update
RUN yes | apt-get upgrade

# otherwise, prompts are shown while installing
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow
RUN yes | apt-get install curl
RUN yes | apt-get install postgresql postgresql-contrib


# install golang and db
RUN curl -s https://storage.googleapis.com/golang/go1.18.3.linux-amd64.tar.gz| tar -v -C /usr/local -xz
ENV PATH="/usr/local/go/bin:${PATH}"


WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /little-api

EXPOSE 8090

CMD [ "/little-api" ]