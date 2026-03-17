package domain

type Status string

const (
	StatusPending   Status = "PENDING"
	StatusAssigned  Status = "ASSIGNED"
	StatusPickedUp  Status = "PICKED_UP"
	StatusInTransit Status = "IN_TRANSIT"
	StatusDelivered Status = "DELIVERED"
	StatusFailed    Status = "FAILED"
	StatusCancelled Status = "CANCELLED"
)

var validTransitions = map[Status][]Status{
	StatusPending:   {StatusAssigned, StatusCancelled},
	StatusAssigned:  {StatusPickedUp, StatusCancelled},
	StatusPickedUp:  {StatusInTransit},
	StatusInTransit: {StatusDelivered, StatusFailed},
	StatusDelivered: {},
	StatusFailed:    {},
	StatusCancelled: {},
}

func IsTerminal(s Status) bool {
	return s == StatusDelivered || s == StatusFailed || s == StatusCancelled
}

func CanTransition(from, to Status) bool {
	allowed, ok := validTransitions[from]
	if !ok {
		return false
	}
	for _, s := range allowed {
		if s == to {
			return true
		}
	}
	return false
}
