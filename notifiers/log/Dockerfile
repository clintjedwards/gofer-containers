FROM golang:1.18 AS build-environment

LABEL org.opencontainers.image.source=https://github.com/clintjedwards/gofer-containers

ARG GOARCH
ARG SEMVER

ADD . /files
WORKDIR /files/notifiers/log

ENV GOOS=linux GOARCH=${GOARCH} CGO_ENABLED=0

RUN go build -o /build/notifier

FROM gcr.io/distroless/static
WORKDIR /app
COPY --from=build-environment /build /app
ENTRYPOINT [ "./notifier" ]
