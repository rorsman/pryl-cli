package timeutil

import "testing"

func TestEpochToISO(t *testing.T) {
	for _, test := range []struct {
		name, want string
		value      int64
		unit       Unit
	}{
		{name: "seconds", value: 0, unit: Seconds, want: "1970-01-01T00:00:00Z"},
		{name: "milliseconds", value: 1712345678000, unit: Milliseconds, want: "2024-04-05T19:34:38Z"},
		{name: "negative milliseconds", value: -1, unit: Milliseconds, want: "1969-12-31T23:59:59.999Z"},
	} {
		t.Run(test.name, func(t *testing.T) {
			got := EpochToISO(test.value, test.unit)
			if got != test.want {
				t.Fatalf("EpochToISO(%d, %q) = %q; want %q", test.value, test.unit, got, test.want)
			}
		})
	}
}
