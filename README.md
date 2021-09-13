# TimeWalker

[![Go Reference](https://pkg.go.dev/badge/github.com/daneroo/timewalker.svg)](https://pkg.go.dev/github.com/daneroo/timewalker)
![Github Build Status](https://github.com/daneroo/timewalker/actions/workflows/test.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/daneroo/timewalker)](https://goreportcard.com/report/github.com/daneroo/timewalker)

go package to generate sequential time intervals
for days,months, and years, accounting for Timezones.

Truncate for D,M,Y: <http://play.golang.org/p/PUNNHq9sh6>

## TODO

- Update Codeship 
- Move to GitHub actions
- Clean up package examples
  - Find DST boundary days
  - Walk back days
- Wrapper/Interface type for Time
- Consider `*time.time` in Interval, or `*Interval` in walker
- Separate benchmarks

## Testing

We have setup continuous testing on Codeship.

### Benchmarking

    export TZ=America/Montreal
    go test --bench .
    go test --bench Round
    go test --bench Construct

### Test coverage

    go test -coverprofile cover.out ; \
    go tool cover -html=cover.out

## Documentation (godoc.org html output locally)

    godoc -http=:8080
    open http://`hostname`:8080/pkg/github.com/daneroo/timewalker/    
