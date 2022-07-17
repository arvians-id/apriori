FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o apriori-backend

EXPOSE 8080

CMD ./apriori-backend