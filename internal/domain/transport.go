package domain

// TransportMode identifies the physical means of carrying a shipment.
// Adding a new mode only requires a new constant here — no other domain
// logic needs to change because the status machine is mode-agnostic.
type TransportMode string

const (
	TransportModeTruck TransportMode = "TRUCK"
	TransportModeAir   TransportMode = "AIR"
	TransportModeSea   TransportMode = "SEA"
	TransportModeRail  TransportMode = "RAIL"
)

var validTransportModes = map[TransportMode]struct{}{
	TransportModeTruck: {},
	TransportModeAir:   {},
	TransportModeSea:   {},
	TransportModeRail:  {},
}

func (m TransportMode) Validate() error {
	if _, ok := validTransportModes[m]; !ok {
		return ErrInvalidTransportMode
	}
	return nil
}
