FROM golang:1.8
COPY . /go/src/rcs/
RUN go get "github.com/Shopify/sarama"
RUN go get "github.com/gorilla/websocket"
RUN go get "github.com/mattn/go-sqlite3"
RUN go install rcs

WORKDIR ../app

EXPOSE 9090 2800
ENTRYPOINT ["rcs"]
