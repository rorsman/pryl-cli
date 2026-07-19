package subnet

import "testing"

func TestCalculateIPv4(t *testing.T) {
	info, err := Calculate("192.168.1.42/24")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Prefix.String(), "192.168.1.0/24"; got != want {
		t.Fatalf("prefix = %q; want %q", got, want)
	}
	if got, want := info.Network.String(), "192.168.1.0"; got != want {
		t.Fatalf("network = %q; want %q", got, want)
	}
	if got, want := info.Last.String(), "192.168.1.255"; got != want {
		t.Fatalf("last = %q; want %q", got, want)
	}
	if got, want := info.UsableCount.String(), "254"; got != want {
		t.Fatalf("usable count = %q; want %q", got, want)
	}
}

func TestCalculateIPv4Exceptions(t *testing.T) {
	for _, test := range []struct{ cidr, first, last, count string }{
		{"192.0.2.0/31", "192.0.2.0", "192.0.2.1", "2"},
		{"192.0.2.7/32", "192.0.2.7", "192.0.2.7", "1"},
	} {
		info, err := Calculate(test.cidr)
		if err != nil {
			t.Fatal(err)
		}
		if info.UsableFirst.String() != test.first || info.UsableLast.String() != test.last || info.UsableCount.String() != test.count {
			t.Errorf("Calculate(%q) = %s-%s (%s); want %s-%s (%s)", test.cidr, info.UsableFirst, info.UsableLast, info.UsableCount, test.first, test.last, test.count)
		}
	}
}

func TestCalculateIPv6(t *testing.T) {
	info, err := Calculate("2001:db8::1234/126")
	if err != nil {
		t.Fatal(err)
	}
	if got, want := info.Network.String(), "2001:db8::1234"; got != want {
		t.Fatalf("network = %q; want %q", got, want)
	}
	if got, want := info.Last.String(), "2001:db8::1237"; got != want {
		t.Fatalf("last = %q; want %q", got, want)
	}
	if got, want := info.UsableCount.String(), "4"; got != want {
		t.Fatalf("usable count = %q; want %q", got, want)
	}
}

func TestContains(t *testing.T) {
	contained, err := Contains("10.0.0.0/8", "10.42.1.2")
	if err != nil || !contained {
		t.Fatalf("Contains returned %v, %v; want true, nil", contained, err)
	}
	contained, err = Contains("10.0.0.0/8", "192.0.2.1")
	if err != nil || contained {
		t.Fatalf("Contains returned %v, %v; want false, nil", contained, err)
	}
}
