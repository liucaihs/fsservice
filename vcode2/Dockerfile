FROM golang:1.8
COPY . /go/src/vcode2/
RUN go install vcode2

WORKDIR ../app

EXPOSE 8088
ENTRYPOINT ["vcode2"]
