package station

type Station struct {
	StationId  uint32
	ChargerIds []uint32
	Downtime   [][2]uint32
	Min        uint32
	Max        uint32
}
