FROM golang:1.8

ADD . /go/src/idfa
RUN go install idfa

EXPOSE 8005
CMD ["idfa"]

