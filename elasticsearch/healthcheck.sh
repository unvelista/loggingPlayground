#!/usr/bin/env bash 

elasticsearch_host="${ELASTICSEARCH_HOST:-elasticsearch}"
output=$(curl -s -D- -m 15 -w "%{http_code}" "http://${elasticsearch_host}:9200/" -u elastic:"${ELASTIC_PASSWORD}") || exit 1

if [[ "${output: -3}" -eq 200 ]]; then
    exit 0
fi
exit 1