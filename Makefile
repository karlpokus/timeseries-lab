VERSION := $(shell git describe --always --dirty --tags)
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: grafana pg ping-pg deploy

pg:
	@docker run -d -p 5432:5432 \
	  --name pg \
	  --restart unless-stopped \
	  -v pg:/var/lib/postgresql/data \
	  -e POSTGRES_PASSWORD=secret \
	  --network telemetry \
	  postgres

ping-pg:
	@pg_isready -h localhost

grafana:
	@docker run -d -p 3000:3000 \
	  --name grafana \
	  --restart unless-stopped \
	  -v grafana:/var/lib/grafana \
	  --network telemetry \
	  grafana/grafana

build:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/telemetry cmd/main.go

deploy: build
	ln -svi ~/dev/timeseries-lab/telemetry.service ~/.config/systemd/user
	systemctl --user enable telemetry
	systemctl --user start telemetry
