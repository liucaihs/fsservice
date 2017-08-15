FROM golang:1.8
COPY . /go/src/device-backend/
RUN go install device-backend

WORKDIR ../app

EXPOSE 8091
ENTRYPOINT ["device-backend"]
