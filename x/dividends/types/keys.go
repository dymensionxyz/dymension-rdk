package types

import "cosmossdk.io/collections"

const (
	// ModuleName defines the module name.
	ModuleName = "dividends"

	// StoreKey defines the primary module store key.
	StoreKey = ModuleName

	// RouterKey is the message route for slashing.
	RouterKey = ModuleName
)

const (
	ParamsByte    = iota // Module params: Params
	LastGaugeByte        // GaugeID sequence
	GaugesByte           // Gauges: GaugeID -> Gauge
)

var (
	ParamsKey    = collections.NewPrefix(ParamsByte)
	LastGaugeKey = collections.NewPrefix(LastGaugeByte)
	GaugesKey    = collections.NewPrefix(GaugesByte)
)
