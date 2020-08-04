FROM debian:buster-slim AS builder

RUN apt-get update -y \
        && apt-get install --no-install-recommends -y wget git unzip lsb-release gnupg2 dpkg-dev ca-certificates \
        && echo "deb-src http://nginx.org/packages/`lsb_release -is | tr '[:upper:]' '[:lower:]'` `lsb_release -cs` nginx" | tee /etc/apt/sources.list.d/nginx.list \
        && wget http://nginx.org/keys/nginx_signing.key && apt-key add nginx_signing.key && rm nginx_signing.key \
        && cd /tmp \
        && apt-get update \
        && apt-get source nginx \
        && apt-get build-dep nginx --no-install-recommends -y \
        && git clone https://github.com/wandenberg/nginx-push-stream-module.git nginx-push-stream-module \
        && cd nginx-1* \
        && sed -i "s@--with-stream_ssl_module@--with-stream_ssl_module --add-module=/tmp/nginx-push-stream-module @g" debian/rules \
        && dpkg-buildpackage -uc -us -b \
        && cd .. \
        && mv nginx_1*~buster_amd64.deb nginx.deb

FROM debian:buster-slim AS runner

COPY --from=builder /tmp/nginx.deb /tmp

RUN apt-get update -y \
        && apt-get install --no-install-recommends -y libssl1.1 lsb-base \
        && dpkg -i /tmp/nginx.deb \
        && apt-mark hold nginx

COPY nginx.conf /etc/nginx/nginx.conf

CMD ["nginx", "-g", "daemon off;"]