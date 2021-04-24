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
$ ./bin/pg/run.sh
$ ./bin/grafana/run.sh
# run collectors
$ go run cmd/main.go
````

# todos
- [ ] optimize postgres index, ingestion
- [x] run postgres under docker
- [ ] run collectors under systemd
- [x] grafana
- [ ] pg replication log https://wiki.postgresql.org/wiki/Streaming_Replication

# license
MIT
