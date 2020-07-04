FROM golang:1.14.4-alpine3.12

LABEL maintainer=vikas@reachvikas.com

ENV REDIS_VERSION=6.0.5 \
    REDIS_DOWNLOAD_URL=http://download.redis.io/releases/redis-6.0.5.tar.gz \
    REDIS_DOWNLOAD_SHA1=9aabd4ad8a007342741933c73314a87402f6e279 \
    GLIDE_VERSION=0.13.3 \
    SUPERVISOR_VERSION=4.2.0-r0

## Install Redis
RUN addgroup -S redis && adduser -S -G redis redis && \
    apk add --no-cache 'su-exec>=0.2'
RUN set -x && \
    apk add --no-cache --virtual .build-deps \
    gcc linux-headers make musl-dev tar curl \
    && curl -o redis.tar.gz -L "${REDIS_DOWNLOAD_URL}" \
    && echo "$REDIS_DOWNLOAD_SHA1 *redis.tar.gz" | sha1sum -c - \
    && mkdir -p /usr/src/redis \
    && tar xzf redis.tar.gz -C /usr/src/redis --strip-components=1 \
    && rm redis.tar.gz \
    && make -C /usr/src/redis \
    && make -C /usr/src/redis install \
    && rm -r /usr/src/redis \
    && apk del .build-deps
RUN mkdir /data && chown redis:redis /data

## Install Glide
RUN apk add --no-cache git curl && \
    curl -o glide.tar.gz -L "https://github.com/Masterminds/glide/releases/download/v${GLIDE_VERSION}/glide-v${GLIDE_VERSION}-linux-arm64.tar.gz" \
    && tar xf glide.tar.gz \
    && mv -v linux-arm64/glide /go/bin/glide \
    && rm -rf glide.tar.gz linux-arm64 \
    && glide --version

## Install Counter with Redis Database
WORKDIR /go/src/github.com/vikas027/go-redis-counter
COPY glide.lock /go/src/github.com/vikas027/go-redis-counter/
COPY glide.yaml /go/src/github.com/vikas027/go-redis-counter/
RUN glide install
COPY . /go/src/github.com/vikas027/go-redis-counter
RUN go install

## Install Supervisord
RUN apk add --no-cache supervisor=${SUPERVISOR_VERSION} && \
    rm -rf /tmp/* /var/cache/apk/*
ADD supervisord.conf /etc/
ENTRYPOINT ["supervisord", "--configuration", "/etc/supervisord.conf"]
