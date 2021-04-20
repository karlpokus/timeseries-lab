# timeseries-lab
Let's collect some timeseries data from my laptop and shove them into postgres. If we're feeling jazzy we'll show them in grafana too.

# telemetry
- [ ] boot/sleep
- [ ] cpu heat
- [x] battery charge state
- [ ] cpuhogs

# usage
This will most likely only work on linux on dell xps.

````bash
# run postgres
$ ./script/pg.sh
# run collectors
$ go run cmd/main.go
````

# todos
- [ ] optimize postgres index, ingestion
- [x] run postgres under docker
- [ ] run collectors under systemd
- [ ] grafana

# license
MIT
