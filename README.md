# loggingPlayground

- [Fluentd Logging Driver](https://docs.docker.com/config/containers/logging/fluentd/)

## Loki Config

### Displaying `\n` as new line for stacktraces

Loki Query: `{job="fluentbit"} | json | line_format "{{.log}}"`