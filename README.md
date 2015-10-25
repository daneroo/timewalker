# TimeWalker
go package to generate sequential time intervals
for days,months, and years, accounting for Timezones.

Truncate for D,M,Y: http://play.golang.org/p/PUNNHq9sh6

## TODO

* Clean up package examples
    * Find DST boundary days
    * Walk back days
* Wrapper/Interface type for Time
* Consider *time.time in Interval, or *Interval in walker
* Seperate benchmarks

## Testing

### Benchmarking

    go test --bench .
    go test --bench Round
    go test --bench Construct

### Test coverage

    go test -coverprofile cover.out
    go tool cover -html=cover.out

## Documentation (godoc.org html output locally)

    godoc -http=:8080
    open http://`hostname`:8080/pkg/github.com/daneroo/timewalker/    