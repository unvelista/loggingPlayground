FROM fluent/fluentd:v1.14.0-1.0

USER root

RUN apk add --no-cache tzdata
ENV TZ=Europe/Berlin

RUN ["gem", "install", "fluent-plugin-elasticsearch"]

COPY ./config/fluent.conf /fluentd/etc/fluent.conf

USER fluent
