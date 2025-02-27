FROM golang:1.22.1

WORKDIR /go/app/src

COPY lim .

RUN go mod download
RUN go build -o lim-rest-api internal/lim-handler/cmd/swagger-l-i-m-server/main.go

WORKDIR /app

RUN cp /go/app/src/lim-rest-api .

EXPOSE 7777
