FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN mkdir output &&\
    go build -ldflags="-s -w -linkmode internal" -o output/esd main.go &&\
    bins/upx output/esd &&\
    cp -r conf output/
#
FROM alpine:3.18
COPY --from=builder /app/output /data
WORKDIR /data
ENTRYPOINT ["/data/esd"]