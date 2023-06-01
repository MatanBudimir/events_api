FROM golang:alpine

WORKDIR /app

COPY . /app

RUN go build -o api /app/cmd/api/main.go

EXPOSE $PORT

ENTRYPOINT ["./api"]