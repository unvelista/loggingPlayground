FROM golang:alpine

ARG APPNAME

RUN apk add --no-cache tzdata
ENV TZ=Europe/Berlin

# Pull dependencies
WORKDIR /src/
COPY go.* .
RUN go mod download

# Build binary
COPY . .
WORKDIR /src/${APPNAME}
RUN go build -o /go/bin/main ${APPNAME}.go 

ENTRYPOINT [ "main" ]
