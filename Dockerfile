# taken from: https://docs.docker.com/language/golang/build-images/
FROM ubuntu:22.04

COPY --from=golang:1.18-alpine /usr/local/go/ /usr/local/go/
ENV PATH="/usr/local/go/bin:${PATH}"
RUN echo $PATH
RUN ls /usr/local/go
RUN go

#TODO: add database
RUN yes | apt-get update
RUN yes | apt-get upgrade

# otherwise, prompts are shown while installing
ARG DEBIAN_FRONTEND=noninteractive
ENV TZ=Europe/Moscow
RUN yes | apt-get install postgresql postgresql-contrib


WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /little-api

EXPOSE 8090

CMD [ "/little-api" ]