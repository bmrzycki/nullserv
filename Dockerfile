# build stage
FROM golang:alpine AS build-env
WORKDIR /go/src/app
COPY . .
RUN ./mkbuildinfo.sh && \
    CGO_ENABLED=0 GOOS=linux \
    go build -a -ldflags '-w -extldflags "-static"' -o nullserv *.go

# final stage
FROM scratch
COPY --from=build-env /go/src/app/nullserv /
ENTRYPOINT ["/nullserv"]
