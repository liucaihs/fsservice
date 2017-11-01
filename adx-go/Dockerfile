FROM golang:1.9
COPY . /go/src/rcs/
 
RUN go install rcs

WORKDIR ../app

EXPOSE 9090 2800 8091
ENTRYPOINT ["rcs"]
