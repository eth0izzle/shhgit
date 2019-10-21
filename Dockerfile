FROM golang:1.13-alpine

WORKDIR /go/src/app
ADD . src

WORKDIR /go/src/app/src
RUN ["go", "install"]
RUN ["go", "build"]

WORKDIR /go/src/app/
RUN cp src/shhgit .
RUN cp src/config.yaml .
RUN rm -R src

VOLUME ["/tmp/shhgit"]
VOLUME ["/go/app/config.yaml"]
CMD ["shhgit"]
