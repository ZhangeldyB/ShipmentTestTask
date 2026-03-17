package domain

import "errors"

var (
	ErrInvalidTransition = errors.New("invalid status transition")
	ErrTerminalState     = errors.New("shipment is in a terminal state")
	ErrDuplicateStatus   = errors.New("shipment is already in this status")
	ErrShipmentNotFound  = errors.New("shipment not found")
	ErrInvalidAmount     = errors.New("amount must be greater than zero")
	ErrInvalidRevenue    = errors.New("driver revenue must be greater than zero and not exceed amount")
)
