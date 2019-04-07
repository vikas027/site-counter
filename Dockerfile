FROM golang:1.12.2-alpine3.9

LABEL maintainer=vikas@reachvikas.com

## Install Redis
RUN addgroup -S redis && adduser -S -G redis redis && \
    apk add --no-cache 'su-exec>=0.2'
ENV REDIS_VERSION 3.0.7
ENV REDIS_DOWNLOAD_URL http://download.redis.io/releases/redis-3.0.7.tar.gz
ENV REDIS_DOWNLOAD_SHA1 e56b4b7e033ae8dbf311f9191cf6fdf3ae974d1c
RUN set -x && \
    apk add --no-cache --virtual .build-deps \
    gcc linux-headers make musl-dev tar \
    && wget -O redis.tar.gz "$REDIS_DOWNLOAD_URL" \
    && echo "$REDIS_DOWNLOAD_SHA1 *redis.tar.gz" | sha1sum -c - \
    && mkdir -p /usr/src/redis \
    && tar -xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && rm redis.tar.gz \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -r /usr/src/redis \
    && apk del .build-deps
RUN mkdir /data && chown redis:redis /data

## Install Counter with Redis Database
RUN apk add --no-cache curl git && \
    curl https://glide.sh/get | sh
WORKDIR /go/src/github.com/vikas027/go-redis-counter
COPY glide.lock /go/src/github.com/vikas027/go-redis-counter/
COPY glide.yaml /go/src/github.com/vikas027/go-redis-counter/
RUN glide install
COPY . /go/src/github.com/vikas027/go-redis-counter
RUN go install

## Install Supervisord
RUN apk add --no-cache supervisor=3.3.4-r1 && \
    rm  -rf /tmp/* /var/cache/apk/*
ADD supervisord.conf /etc/
ENTRYPOINT ["supervisord", "--configuration", "/etc/supervisord.conf"]
