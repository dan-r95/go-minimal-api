# taken from: https://docs.docker.com/language/golang/build-images/
FROM golang:1.18-alpine

WORKDIR /app


#TODO: add database


COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./

RUN go build -o /little-api

EXPOSE 8090

CMD [ "/little-api" ]