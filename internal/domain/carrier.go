package domain

// CarrierInfo is a mode-agnostic value object describing whoever operates
// the shipment. The same fields cover all transport modes:
//
//	TRUCK → OperatorName=driver,   OperatorPhone=driver phone, UnitIdentifier=plate number
//	AIR   → OperatorName=pilot,    OperatorPhone=ops line,     UnitIdentifier=flight number
//	SEA   → OperatorName=captain,  OperatorPhone=ops line,     UnitIdentifier=vessel IMO
//	RAIL  → OperatorName=engineer, OperatorPhone=ops line,     UnitIdentifier=train number
type CarrierInfo struct {
	OperatorName   string
	OperatorPhone  string
	UnitIdentifier string
}

func (c CarrierInfo) Validate() error {
	if c.OperatorName == "" || c.UnitIdentifier == "" {
		return ErrInvalidCarrier
	}
	return nil
}
