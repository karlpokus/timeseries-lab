#!/bin/bash

# admin:admin

docker run -d -p 3000:3000 \
  --name grafana \
  --restart unless-stopped \
  -v grafana:/var/lib/grafana \
  --network telemetry \
  grafana/grafana
