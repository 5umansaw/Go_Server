FROM golang:1.17-alpine as builder

WORKDIR /app

# COPY go.mod, go.sum and download the dependencies
COPY go.* ./
RUN go mod download

# COPY All things inside the project and build
COPY . .
RUN go build -o /go-docker-ping

EXPOSE 8080

CMD [ "/go-docker-ping" ]