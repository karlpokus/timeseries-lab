FROM alpine:3.13.5
RUN apk add --no-cache ca-certificates
COPY bin/telemetry-api /src/api
ENTRYPOINT ["/src/api"]
