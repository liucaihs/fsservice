FROM golang:1.9
COPY . /go/src/vbs/
RUN go install vbs
WORKDIR ../app
EXPOSE 8093
ENTRYPOINT ["vbs"]
