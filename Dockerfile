ARG DOCKER_REG=docker.io
FROM ${DOCKER_REG}/qnib/alplain-golang AS build

WORKDIR /usr/local/src/github.com/qnib/doxy
COPY main.go ./main.go
COPY proxy ./proxy
COPY vendor/ vendor/
RUN govendor install

## Build final image
FROM alpine:3.5

COPY --from=build /usr/local/bin/doxy /usr/local/bin/
CMD ["/usr/local/bin/doxy"]
