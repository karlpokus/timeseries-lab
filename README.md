# timeseries-lab
Let's collect some timeseries data from my laptop and shove them into postgres. If we're feeling jazzy we'll show them in grafana too.

telemetry types
- [ ] boot/sleep
- [x] cpu heat
- [x] battery charge %
- [x] cpuhogs

# usage
This will most likely only work on linux on dell xps hardware.

````bash
# run postgres and grafana
$ make pg && make grafana
# run collectors
$ go run cmd/main.go
````

# deploy
The agent contains the telemetry collectors. The api is just a test to try out the grafana JSON data source to build a REST data service.

````bash
$ make deploy-agent
$ make deploy-api
````

# todos
- [ ] optimize postgres index, ingestion, retention
- [x] run postgres under docker
- [x] run collectors under systemd
- [x] grafana
- [ ] pg replication log https://wiki.postgresql.org/wiki/Streaming_Replication
- [ ] dump records inserted every x
- [x] use db connection pool w reconnects
- [x] api service
- [ ] bonus: run go binary in telemetry network namespace without docker
- [ ] sql query builer https://github.com/Masterminds/squirrel
- [x] store interface

# license
MIT
