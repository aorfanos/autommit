FROM golang:1.19-alpine AS build-env
ENV CGO_ENABLED=0
COPY . /build
WORKDIR /build
RUN go build -a -installsuffix cgo -ldflags "-w -s" -o autommit cmd/autommit/main.go

FROM debian:sid-slim
WORKDIR /app
RUN apt-get update && apt-get install -y git && \
    git config --global user.email $GIT_ACC_MAIL && \
    git config --global user.name $GIT_ACC_NAME
COPY --from=build-env /build/autommit /autommit
ENTRYPOINT ["/autommit"]
