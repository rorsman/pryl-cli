package timeutil

import "time"

type Unit string

const (
	Seconds      Unit = "seconds"
	Milliseconds Unit = "milliseconds"
)

// EpochToISO converts an explicitly unit-qualified Unix timestamp to UTC ISO 8601.
func EpochToISO(value int64, unit Unit) string {
	if unit == Milliseconds {
		seconds := value / 1000
		nanos := (value % 1000) * int64(time.Millisecond)
		return time.Unix(seconds, nanos).UTC().Format(time.RFC3339Nano)
	}

	return time.Unix(value, 0).UTC().Format(time.RFC3339)
}
