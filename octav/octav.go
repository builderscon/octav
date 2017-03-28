//go:generate gendb -d db
//go:generate genmodel -d .
//go:generate gentransport -d .
//go:generate genoctavctl -s ../spec/v1/api.json -t spec/octavctl.json -o cmd/octavctl/octavctl.go

package octav

