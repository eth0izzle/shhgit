FROM golang:alpine AS builder

WORKDIR /go/src
ADD . .

RUN export CGO_ENABLED=0 && go install && go build -o /

FROM scratch AS runtime
VOLUME /tmp/shhgit

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /shhgit .

ENTRYPOINT [ "/shhgit" ]
