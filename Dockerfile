FROM alpine
RUN apk add --no-cache ca-certificates
COPY bin/telemetry-api /src/api
ENTRYPOINT ["/src/api"]
