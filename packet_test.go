package radigo

import (
	"fmt"
	"reflect"
	"sync"
	"testing"
)

func TestPacketDecode(t *testing.T) {
	// sample packet taken out of RFC2865 -Section 7.2.
	encdPkt := []byte{
		0x01, 0x01, 0x00, 0x47, 0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
		0x0d, 0x33, 0x67, 0xa2, 0x01, 0x08, 0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79, 0x03, 0x13, 0x16, 0xe9,
		0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75, 0x04,
		0x06, 0xc0, 0xa8, 0x01, 0x10, 0x05, 0x06, 0x00, 0x00, 0x00, 0x14, 0x06, 0x06, 0x00, 0x00, 0x00,
		0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01, 0x1a, 0x13, 0x00, 0x00, 0x00, 0x09,
		0x17, 0x0d, 0x43, 0x47, 0x52, 0x61, 0x74, 0x65, 0x53, 0x2e, 0x6f, 0x72, 0x67,
	}
	ePkt := &Packet{
		Code:       AccessRequest,
		Identifier: 1,
		Authenticator: [16]byte{0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
			0x0d, 0x33, 0x67, 0xa2},
		AVPs: []*AVP{
			&AVP{
				Number:   uint8(1),                                   // User-Name
				RawValue: []byte{0x66, 0x6c, 0x6f, 0x70, 0x73, 0x79}, // flopsy
			},
			&AVP{
				Number: uint8(3), // CHAPPassword
				RawValue: []byte{0x16, 0xe9,
					0x75, 0x57, 0xc3, 0x16, 0x18, 0x58, 0x95, 0xf2, 0x93, 0xff, 0x63, 0x44, 0x07, 0x72, 0x75}, // 3
			},
			&AVP{
				Number:   uint8(4),                       // NASIPAddress
				RawValue: []byte{0xc0, 0xa8, 0x01, 0x10}, // 192.168.1.16
			},
			&AVP{
				Number:   uint8(5),                       // NASPort
				RawValue: []byte{0x00, 0x00, 0x00, 0x14}, // 20
			},
			&AVP{
				Number:   uint8(6),                       // ServiceType
				RawValue: []byte{0x00, 0x00, 0x00, 0x02}, // 2
			},
			&AVP{
				Number:   uint8(7),                       // FramedProtocol
				RawValue: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
			&AVP{
				Number: VendorSpecificNumber, // VSA
				RawValue: []byte{0x00, 0x00, 0x00, 0x09,
					0x17, 0x0d, 0x43, 0x47, 0x52, 0x61, 0x74, 0x65, 0x53, 0x2e, 0x6f, 0x72, 0x67}, // VendorID: 9(Cisco), VSA-type: 23(Remote-Gateway-ID), VSA-Data: CGRateS.org
			},
		},
	}
	pkt := new(Packet)
	if err := pkt.Decode(encdPkt); err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(ePkt, pkt) {
		t.Errorf("Expecting: %+v, received: %+v", ePkt, pkt)
	}
}

func TestPacketEncode(t *testing.T) {
	pkt := &Packet{
		Code:       AccessAccept,
		Identifier: 1,
		Authenticator: [16]byte{0x2a, 0xee, 0x86, 0xf0, 0x8d, 0x0d, 0x55, 0x96, 0x9c, 0xa5, 0x97, 0x8e,
			0x0d, 0x33, 0x67, 0xa2}, // Authenticator out of origin request
		AVPs: []*AVP{
			&AVP{
				Number:   6,                              // ServiceType
				RawValue: []byte{0x00, 0x00, 0x00, 0x02}, // 2
			},
			&AVP{
				Number:   7,                              // FramedProtocol
				RawValue: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
			&AVP{
				Number:   8,                              // FramedIPAddress
				RawValue: []byte{0xff, 0xff, 0xff, 0xfe}, // 255.255.255.254
			},
			&AVP{
				Number:   10,                             // FramedRouting
				RawValue: []byte{0x00, 0x00, 0x00, 0x02}, // 0
			},
			&AVP{
				Number:   13,                             // FramedCompression
				RawValue: []byte{0x00, 0x00, 0x00, 0x01}, // 1
			},
			&AVP{
				Number:   12,                             // FramedMTU
				RawValue: []byte{0x00, 0x00, 0x05, 0xdc}, // 1500
			},
			&AVP{
				Number: 26, // VSA
				RawValue: []byte{0x00, 0x00, 0x00, 0x09,
					0x17, 0x0d, 0x43, 0x47, 0x52, 0x61, 0x74, 0x65, 0x53, 0x2e, 0x6f, 0x72, 0x67}, // VendorID: 9(Cisco), VSA-type: 23(Remote-Gateway-ID), VSA-Data: CGRateS.org
			},
		},
	}
	ePktEncd := []byte{
		0x02, 0x01, 0x00, 0x4b, 0x0c, 0x51, 0xfd, 0x77, 0xec, 0xb6, 0x5a, 0xac, 0x43, 0x8b, 0x79, 0x99,
		0xe4, 0x12, 0x55, 0x18, 0x06, 0x06, 0x00, 0x00, 0x00, 0x02, 0x07, 0x06, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0xff, 0xff, 0xff, 0xfe, 0x0a, 0x06, 0x00, 0x00, 0x00, 0x02, 0x0d, 0x06, 0x00, 0x00,
		0x00, 0x01, 0x0c, 0x06, 0x00, 0x00, 0x05, 0xdc, 0x1a, 0x13, 0x00, 0x00, 0x00, 0x09,
		0x17, 0x0d, 0x43, 0x47, 0x52, 0x61, 0x74, 0x65, 0x53, 0x2e, 0x6f, 0x72, 0x67,
	}
	var buf [4096]byte
	n, err := pkt.Encode(buf[:])
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(ePktEncd, buf[:n]) {
		t.Errorf("Expecting: % x, received: % x", ePktEncd, buf[:n])
	}

}

func TestPacketStringer(t *testing.T) {
	p := AccessRequest
	exp := "AccessRequest"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = AccessAccept
	exp = "AccessAccept"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = AccessReject
	exp = "AccessReject"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = AccountingRequest
	exp = "AccountingRequest"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = AccountingResponse
	exp = "AccountingResponse"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = AccessChallenge
	exp = "AccessChallenge"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = StatusServer
	exp = "StatusServer"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = StatusClient
	exp = "StatusClient"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = Reserved
	exp = "Reserved"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}

	p = PacketCode(60)
	exp = "unknown packet code"
	if rcv := p.String(); rcv != exp {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}
}

func TestPacketHas(t *testing.T) {
	p := &Packet{
		AVPs: []*AVP{
			{
				Number: 1,
			},
			{
				Number: 25,
			},
			{
				Number: 5,
			},
		},
	}
	attrNr := uint8(5)

	rcv := p.Has(attrNr)

	if rcv != true {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", true, rcv)
	}
}

func TestPacketHasNot(t *testing.T) {
	p := &Packet{
		AVPs: []*AVP{
			{
				Number: 1,
			},
			{
				Number: 25,
			},
			{
				Number: 5,
			},
		},
	}
	attrNr := uint8(6)

	rcv := p.Has(attrNr)

	if rcv != false {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", false, rcv)
	}
}

func TestPacketNewPacket(t *testing.T) {
	exp := &Packet{
		Code:       Reserved,
		Identifier: uint8(5),
		dict:       &Dictionary{},
		coder:      Coder{},
		secret:     "testString",
	}

	rcv := NewPacket(Reserved, uint8(5), &Dictionary{}, Coder{}, "testString")

	if !reflect.DeepEqual(rcv, exp) {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", exp, rcv)
	}
}

func TestPacketEncodeNilRawValue(t *testing.T) {
	b := make([]byte, 100)
	p := &Packet{
		RWMutex:       sync.RWMutex{},
		Code:          Reserved,
		Identifier:    uint8(5),
		Authenticator: [16]byte{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5, 6},
		AVPs: []*AVP{
			{Number: 0},
			{Number: 1},
		},
	}

	experr := fmt.Sprintf("avp: %+v, no value", p.AVPs[0])
	n, err := p.Encode(b)
	if err == nil || err.Error() != experr {
		t.Fatalf("\nExpected: <%+v>, \nReceived: <%+v>", experr, err)
	}

	if n != 0 {
		t.Errorf("\nExpected: <%+v>, \nReceived: <%+v>", 0, n)
	}
}
