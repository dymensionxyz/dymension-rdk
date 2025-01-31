package types

func NewGauge(
	id uint64,
	address string,
	queryCondition *QueryCondition,
	vestingCondition *VestingCondition,
) Gauge {
	return Gauge{
		Id:               id,
		Address:          address,
		QueryCondition:   queryCondition,
		VestingCondition: vestingCondition,
	}
}
