VERSION := $(shell git describe --always --dirty --tags)
GOOS    := $(shell go env GOOS)
GOARCH  := $(shell go env GOARCH)

.PHONY: grafana pg ping-pg deploy-agent deploy-api

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
	@docker run -d -p 3001:3000 \
	  --name grafana \
	  --restart unless-stopped \
	  -v grafana:/var/lib/grafana \
	  --network telemetry \
	  grafana/grafana

agent:
	CGO_ENABLED=0 go build -ldflags="-s -w -X main.version=$(VERSION)" -o bin/telemetry-agent cmd/agent.go

deploy-agent: agent
	ln -svi ~/dev/timeseries-lab/telemetry.service ~/.config/systemd/user
	systemctl --user enable telemetry
	systemctl --user restart telemetry

api:
	CGO_ENABLED=0 go build -o bin/telemetry-api cmd/api.go
	docker build -t cf/telem-api .

deploy-api: api
	-docker stop telemapi 2>&1 > /dev/null
	-docker rm telemapi 2>&1 > /dev/null
	@docker run -d -p 8989:8989 \
		--name telemapi \
		--restart unless-stopped \
		--network telemetry \
		cf/telem-api
