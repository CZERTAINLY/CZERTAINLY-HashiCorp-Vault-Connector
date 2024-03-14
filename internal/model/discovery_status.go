package model

import (
	"fmt"
)

// DiscoveryStatus : Status of Discovery
type DiscoveryStatus string

// List of DiscoveryStatus
const (
	IN_PROGRESS DiscoveryStatus = "inProgress"
	FAILED      DiscoveryStatus = "failed"
	COMPLETED   DiscoveryStatus = "completed"
	WARNING     DiscoveryStatus = "warning"
)

// AllowedDiscoveryStatusEnumValues is all the allowed values of DiscoveryStatus enum
var AllowedDiscoveryStatusEnumValues = []DiscoveryStatus{
	"inProgress",
	"failed",
	"completed",
	"warning",
}

// validDiscoveryStatusEnumValue provides a map of DiscoveryStatuss for fast verification of use input
var validDiscoveryStatusEnumValues = map[DiscoveryStatus]struct{}{
	"inProgress": {},
	"failed":     {},
	"completed":  {},
	"warning":    {},
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v DiscoveryStatus) IsValid() bool {
	_, ok := validDiscoveryStatusEnumValues[v]
	return ok
}

// NewDiscoveryStatusFromValue returns a pointer to a valid DiscoveryStatus
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewDiscoveryStatusFromValue(v string) (DiscoveryStatus, error) {
	ev := DiscoveryStatus(v)
	if ev.IsValid() {
		return ev, nil
	} else {
		return "", fmt.Errorf("invalid value '%v' for DiscoveryStatus: valid values are %v", v, AllowedDiscoveryStatusEnumValues)
	}
}

// AssertDiscoveryStatusRequired checks if the required fields are not zero-ed
func AssertDiscoveryStatusRequired(obj DiscoveryStatus) error {
	return nil
}

// AssertDiscoveryStatusConstraints checks if the values respects the defined constraints
func AssertDiscoveryStatusConstraints(obj DiscoveryStatus) error {
	return nil
}
