FROM golang:1.8
COPY . /go/src/device-frontend/
RUN go install device-frontend

WORKDIR ../app

EXPOSE 8090
ENTRYPOINT ["device-frontend"]
