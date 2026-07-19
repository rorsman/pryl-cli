// Package subnet calculates information about IP network prefixes.
package subnet

import (
	"encoding/json"
	"fmt"
	"math/big"
	"net/netip"
)

// Info contains the calculated range for a network prefix.
type Info struct {
	Prefix       netip.Prefix
	Network      netip.Addr
	Last         netip.Addr
	AddressCount *big.Int
	UsableFirst  netip.Addr
	UsableLast   netip.Addr
	UsableCount  *big.Int
}

// Calculate returns information about cidr. Host bits in the input address
// are ignored, so 192.168.1.42/24 is treated as 192.168.1.0/24.
func Calculate(cidr string) (Info, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return Info{}, fmt.Errorf("invalid subnet %q: %w", cidr, err)
	}
	prefix = prefix.Masked()
	network := prefix.Addr()
	addressBits := network.BitLen()
	hostBits := addressBits - prefix.Bits()

	addressCount := new(big.Int).Lsh(big.NewInt(1), uint(hostBits))
	last := addAddress(network, new(big.Int).Sub(new(big.Int).Set(addressCount), big.NewInt(1)))
	usableFirst, usableLast := network, last
	usableCount := new(big.Int).Set(addressCount)

	// Traditional IPv4 networks reserve the network and broadcast addresses.
	// RFC 3021 /31 networks and /32 host routes are exceptions.
	if network.Is4() && prefix.Bits() <= 30 {
		usableFirst = network.Next()
		usableLast = last.Prev()
		usableCount.Sub(usableCount, big.NewInt(2))
	}

	return Info{
		Prefix:       prefix,
		Network:      network,
		Last:         last,
		AddressCount: addressCount,
		UsableFirst:  usableFirst,
		UsableLast:   usableLast,
		UsableCount:  usableCount,
	}, nil
}

func addAddress(address netip.Addr, offset *big.Int) netip.Addr {
	if address.Is4() {
		addressBytes := address.As4()
		value := new(big.Int).SetBytes(addressBytes[:])
		value.Add(value, offset)
		var bytes [4]byte
		copy(bytes[:], value.FillBytes(make([]byte, 4)))
		return netip.AddrFrom4(bytes)
	}

	addressBytes := address.As16()
	value := new(big.Int).SetBytes(addressBytes[:])
	value.Add(value, offset)
	var bytes [16]byte
	copy(bytes[:], value.FillBytes(make([]byte, 16)))
	return netip.AddrFrom16(bytes)
}

// Contains reports whether address belongs to cidr.
func Contains(cidr, address string) (bool, error) {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		return false, fmt.Errorf("invalid subnet %q: %w", cidr, err)
	}
	ip, err := netip.ParseAddr(address)
	if err != nil {
		return false, fmt.Errorf("invalid address %q: %w", address, err)
	}
	return prefix.Contains(ip), nil
}

// JSON returns a stable machine-readable representation of Info.
func (i Info) JSON() ([]byte, error) {
	value := struct {
		CIDR         string `json:"cidr"`
		Network      string `json:"network"`
		Last         string `json:"last"`
		AddressCount string `json:"address_count"`
		UsableFirst  string `json:"usable_first"`
		UsableLast   string `json:"usable_last"`
		UsableCount  string `json:"usable_count"`
	}{
		CIDR:         i.Prefix.String(),
		Network:      i.Network.String(),
		Last:         i.Last.String(),
		AddressCount: i.AddressCount.String(),
		UsableFirst:  i.UsableFirst.String(),
		UsableLast:   i.UsableLast.String(),
		UsableCount:  i.UsableCount.String(),
	}
	return json.Marshal(value)
}
