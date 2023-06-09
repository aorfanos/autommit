FROM golang:1.20-alpine AS build-env
ENV CGO_ENABLED=0
LABEL name="Autommit" version="0.0.12"
COPY . /build
WORKDIR /build
RUN go build -a -installsuffix cgo -ldflags "-w -s" -o autommit cmd/autommit/main.go

FROM debian:sid-slim
WORKDIR /app
RUN apt-get update && \
    apt-get install --no-install-recommends -y git && \
    rm -rf /var/lib/apt/lists/* && \
    apt-get clean

COPY --from=build-env /build/autommit /autommit
ENTRYPOINT ["/autommit"]
