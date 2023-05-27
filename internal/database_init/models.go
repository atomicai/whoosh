package database_init

import "time"

type Clash struct {
	HexId            int     `json:"hex_id" gorethink:"hex_id"`
	ClashesShare     float64 `json:"clashes_share" gorethink:"clashes_share"`
	ClashPowerMedian float64 `json:"clash_power_median" gorethink:"clash_power_median"`
}

type ScootersAtParkings struct {
	TsUtc             time.Time `json:"ts_utc" gorethink:"ts_utc"`
	ParkingId         int       `json:"parking_id" gorethink:"parking_id"`
	ScootersAtParking int       `json:"scooters_at_parking" gorethink:"scooters_at_parking"`
}
