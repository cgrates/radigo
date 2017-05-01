package codecs

import (
	"errors"
	"net"
)

// AddressCodec is a codec for address values
type AddressCodec struct{}

// Decode is part of AVPCoder interface
func (cdc AddressCodec) Decode(b []byte) (v interface{}, s string, err error) {
	ip := net.IP(b)
	return ip, ip.String(), nil

}

// Encode is part of AVPCoder interface
func (cdc AddressCodec) Encode(v interface{}) (b []byte, err error) {
	ipVal, ok := v.(net.IP)
	if !ok {
		err = errors.New("cannot cast to net.IP")
		return
	}
	ipVal = ipVal.To4()
	if ipVal == nil {
		err = errors.New("cannot enforce IPv4")
		return
	}
	return []byte(ipVal), nil
}
