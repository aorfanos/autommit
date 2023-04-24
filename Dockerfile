FROM golang:1.19-alpine AS build-env
ARG GH_MAIL="foo@users.github.com" GH_NAME="GitHub User"
ENV CGO_ENABLED=0 \
    GIT_ACC_MAIL=$GH_MAIL \
    GIT_ACC_NAME=$GH_NAME
COPY . /build
WORKDIR /build
RUN go build -a -installsuffix cgo -ldflags "-w -s" -o autommit cmd/autommit/main.go

FROM debian:sid-slim
WORKDIR /app
RUN apt-get update && apt-get install -y git

RUN git config --global user.email "${GIT_ACC_MAIL}" && \
    git config --global user.name "${GIT_ACC_NAME}"
COPY --from=build-env /build/autommit /autommit
ENTRYPOINT ["/autommit"]
