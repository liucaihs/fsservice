FROM golang:1.9
ADD . /go/src/userpwd/
RUN go install userpwd
WORKDIR ../app
EXPOSE 10001
ENTRYPOINT ["userpwd"]
