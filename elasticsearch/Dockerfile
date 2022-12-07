# https://www.docker.elastic.co/
FROM docker.elastic.co/elasticsearch/elasticsearch:8.5.2

# Add your elasticsearch plugins setup here
# Example: RUN elasticsearch-plugin install analysis-icu

COPY ./healthcheck.sh /usr/share/elasticsearch/healthcheck.sh
HEALTHCHECK --interval=10s --timeout=5s --start-period=5s --retries=3 CMD [ "./healthcheck.sh" ]

COPY ./config/elasticsearch.yml /usr/share/elasticsearch/config/elasticsearch.yml
