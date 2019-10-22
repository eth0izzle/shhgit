FROM golang:alpine AS builder

WORKDIR /go/src/
ADD . .

RUN go install && go build

FROM alpine:latest AS runtime
WORKDIR /app
VOLUME /tmp/shhgit

COPY --from=builder /go/src/shhgit .
COPY --from=builder /go/src/config.yaml .

ENTRYPOINT [ "/app/shhgit" ]
