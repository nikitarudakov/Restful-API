package timeparse

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// ParseExpirationTime converts the quantifier and unit to a time.Duration.
func ParseExpirationTime(quantifier string, unit string) (time.Duration, error) {
	quantifierInt, err := strconv.Atoi(quantifier)
	if err != nil {
		return 0, fmt.Errorf("unsupported time unit: %w", err)
	}

	// Convert unit string to time.Duration.
	var durationUnit time.Duration
	switch strings.ToLower(unit) {
	case "second", "seconds":
		durationUnit = time.Second
	case "minute", "minutes":
		durationUnit = time.Minute
	case "hour", "hours":
		durationUnit = time.Hour
	default:
		return 0, fmt.Errorf("unsupported time unit: %s", unit)
	}

	// Calculate the total duration.
	totalDuration := time.Duration(quantifierInt) * durationUnit

	return totalDuration, nil
}
