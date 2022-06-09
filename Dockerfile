# taken from: https://docs.docker.com/language/golang/build-images/
FROM golang:1.18-alpine

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