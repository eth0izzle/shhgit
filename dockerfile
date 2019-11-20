#!/usr/bin/env -S docker build --compress -t eth0izzle/shhgit -f

FROM debian as build

RUN apt update
RUN apt install -y gcc git curl

RUN curl -skL https://dl.google.com/go/go1.13.linux-amd64.tar.gz \
	| tar --strip-components=0 -xzC /usr/local

ENV PATH "$PATH:/usr/local/go/bin"

WORKDIR /root/go/src/github.com/eth0izzle/shhgit
COPY ./ ./
RUN echo get build test install | xargs -n1 | xargs -n1 -I% -- go % .
RUN sed -i.orig "s:^  - '':  - \${GITHUB_TOKEN}:" ./config.yaml

FROM debian
RUN apt update
RUN apt install -y ca-certificates
WORKDIR /data
COPY --from=build \
	/root/go/bin/shhgit \
	/usr/local/bin/shhgit
COPY --from=build \
	/root/go/src/github.com/eth0izzle/shhgit/config.yaml \
	./config.yaml
ENTRYPOINT [ "/usr/local/bin/shhgit" ]
CMD        [ ]
