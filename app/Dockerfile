FROM golang:alpine AS builder

RUN apk update && apk add --no-cache -t .build-deps bison jansson jansson-dev build-base libc-dev file-dev automake autoconf libtool git \
  && cd /tmp \
  && git clone --recursive --branch v4.0.2 https://github.com/VirusTotal/yara.git \
  && cd /tmp/yara \
  && ./bootstrap.sh && sync && ./configure \
  && make && make install

WORKDIR /tmp/shhgit
COPY ./src /tmp/shhgit
RUN export CGO_ENABLED=1 && go install ./cmd/shhgit && go build -o app ./cmd/shhgit

FROM golang:alpine AS runtime
WORKDIR /app

RUN apk update && apk add --no-cache git mercurial
COPY --from=builder /usr/local/lib /usr/local/lib
COPY --from=builder /tmp/shhgit/app /app/shhgit

ENTRYPOINT [ "/app/shhgit" ]