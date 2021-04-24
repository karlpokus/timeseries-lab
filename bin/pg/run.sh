#!/bin/bash

docker run -d -p 5432:5432 \
  --name pg \
  --restart unless-stopped \
  -v pg:/var/lib/postgresql/data \
  -e POSTGRES_PASSWORD=secret \
  --network telemetry \
  postgres
