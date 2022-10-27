FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o api-message-broker

EXPOSE 8080

CMD ./api-message-broker