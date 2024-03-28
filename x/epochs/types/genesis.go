package types

import (
	"errors"
	"time"
)

func NewGenesisState(epochs []EpochInfo) *GenesisState {
	return &GenesisState{Epochs: epochs}
}

var (
	MINUTE_EPOCH = "minute"
	HOUR_EPOCH   = "hour"
	DAY_EPOCH    = "day"
	WEEK_EPOCH   = "week"
	YEAR_EPOCH   = "year"
)

// DefaultGenesis returns the default Epochs genesis state
func DefaultGenesis() *GenesisState {
	epochs := []EpochInfo{
		{
			Identifier:              MINUTE_EPOCH,
			StartTime:               time.Time{},
			Duration:                time.Minute,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              HOUR_EPOCH,
			StartTime:               time.Time{},
			Duration:                time.Hour,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              DAY_EPOCH,
			StartTime:               time.Time{},
			Duration:                time.Hour * 24,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              WEEK_EPOCH,
			StartTime:               time.Time{},
			Duration:                time.Hour * 24 * 7,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
		{
			Identifier:              YEAR_EPOCH,
			StartTime:               time.Time{},
			Duration:                time.Hour * 24 * 365,
			CurrentEpoch:            0,
			CurrentEpochStartHeight: 0,
			CurrentEpochStartTime:   time.Time{},
			EpochCountingStarted:    false,
		},
	}
	return NewGenesisState(epochs)
}

// Validate performs basic genesis state validation returning an error upon any
// failure.
func (gs GenesisState) Validate() error {
	epochIdentifiers := map[string]bool{}
	for _, epoch := range gs.Epochs {
		if err := epoch.Validate(); err != nil {
			return err
		}
		if epochIdentifiers[epoch.Identifier] {
			return errors.New("epoch identifier should be unique")
		}
		epochIdentifiers[epoch.Identifier] = true
	}
	return nil
}

// Validate also validates epoch info.
func (epoch EpochInfo) Validate() error {
	if epoch.Identifier == "" {
		return errors.New("epoch identifier should NOT be empty")
	}
	if epoch.Duration == 0 {
		return errors.New("epoch duration should NOT be 0")
	}
	if epoch.CurrentEpoch < 0 {
		return errors.New("epoch CurrentEpoch must be non-negative")
	}
	if epoch.CurrentEpochStartHeight < 0 {
		return errors.New("epoch CurrentEpochStartHeight must be non-negative")
	}
	return nil
}
