package message

import (
	"github.com/dist-ribut-us/crypto"
	"github.com/dist-ribut-us/log"
	"github.com/dist-ribut-us/rnet"
	"github.com/dist-ribut-us/serial"
	"github.com/golang/protobuf/proto"
)

// BitFlag is used to set flags on the header flag field
type BitFlag uint32

// Flag masks Flags field
const (
	QueryFlag = BitFlag(1 << (iota))
	ResponseFlag
	FromNet
	ToNet
)

// Type for message header
type Type uint32

// Types for messages
const (
	Undefined = Type(iota)
	GetPort
	GetIP
	RegisterService
	Test
	Ping
	NetSend
	NetReceive
	AddBeacon
	GetPubKey
	Die
	SessionData
	StaticKey
	RandomKey
)

// NewHeader takes a type and a body. See SetBody for valid body types
func NewHeader(t Type, body interface{}) *Header {
	h := &Header{
		Type32: uint32(t),
		Id:     crypto.RandUint32(),
	}
	return h.SetBody(body)
}

// Unmarshal a buffer to a header. If there is an error it is logged and a nil
// value is returned.
func Unmarshal(buf []byte) *Header {
	h := &Header{}
	err := proto.Unmarshal(buf, h)
	if log.Error(err) {
		return nil
	}
	return h
}

// Marshal a header to a buffer. If there is an error it is logged and nil is
// returned.
func (h *Header) Marshal() []byte {
	buf, err := proto.Marshal(h)
	if log.Error(err) {
		return nil
	}
	return buf
}

// SetBody takes various typesand tries to set the Body field. If body is nil,
// that will set the body to nil. If body is a []byte, that will be set on the
// body. If it is a string, it will be case to []byte. If the body is a
// proto.Message it will be marshalled and the serialized data will be used. If
// it is a uint32, dist.ribut.us/serail will be used to marshal the value into
// body.
func (h *Header) SetBody(body interface{}) *Header {
	if body == nil {
		h.Body = nil
	} else if ui, ok := body.(uint32); ok {
		h.Body = serial.MarshalUint32(ui, nil)
	} else if bts, ok := body.([]byte); ok {
		h.Body = bts
	} else if str, ok := body.(string); ok {
		h.Body = []byte(str)
	} else if msg, ok := body.(proto.Message); ok {
		buf, err := proto.Marshal(msg)
		if !log.Error(err) {
			h.Body = buf
		}
	}
	return h
}

// GetType returns the Type
func (h *Header) GetType() Type {
	return Type(h.Type32)
}

// SetType sets the Type
func (h *Header) SetType(t Type) {
	h.Type32 = uint32(t)
}

// CheckFlag checks if a bit flag is set
func (h *Header) CheckFlag(flag BitFlag) bool {
	return h.Flags&uint32(flag) == uint32(flag)
}

// IsQuery checks if the underlying type is a query
func (h *Header) IsQuery() bool {
	return h.CheckFlag(QueryFlag)
}

// IsResponse checks if the underlying type is a query
func (h *Header) IsResponse() bool {
	return h.CheckFlag(ResponseFlag)
}

// IsFromNet checks if the FromNet flag is set
func (h *Header) IsFromNet() bool {
	return h.CheckFlag(FromNet)
}

// IsToNet checks if the FromNet flag is set
func (h *Header) IsToNet() bool {
	return h.CheckFlag(ToNet)
}

// SetFlag sets a bit in the flag field
func (h *Header) SetFlag(flag BitFlag) {
	h.Flags |= uint32(flag)
}

// UnsetFlag removes a flag
func (h *Header) UnsetFlag(flag BitFlag) {
	h.Flags ^= uint32(flag)
}

// BodyToUint32 uses dist.ribut.us/serial to Unmarshal the body.
func (h *Header) BodyToUint32() uint32 {
	if len(h.Body) < 4 {
		return 0
	}
	return serial.UnmarshalUint32(h.Body)
}

// BodyString returns the body as a string
func (h *Header) BodyString() string {
	return string(h.Body)
}

// Unmarshal the body of the header
func (h *Header) Unmarshal(pb proto.Message) error {
	return proto.Unmarshal(h.Body, pb)
}

// SetAddr sets the underlying Addrpb struct in the Header
func (h *Header) SetAddr(addr *rnet.Addr) *Header {
	h.Addrpb = FromAddr(addr)
	return h
}

// FromAddr creates and Addrpb from an rnet.Addr
func FromAddr(addr *rnet.Addr) *Addrpb {
	return &Addrpb{
		Ip:   addr.IP,
		Port: uint32(addr.UDPAddr.Port),
		Zone: addr.Zone,
	}
}

// GetAddr returns the Addrpb as an rnet.Addr
func (h *Header) GetAddr() *rnet.Addr {
	if h.Addrpb == nil {
		return nil
	}
	return h.Addrpb.GetAddr()
}

// UnmarshalAddrpb unmarshals a buffer to an Addrpb. If there is an error it is
// logged and a nil value is returned.
func UnmarshalAddrpb(buf []byte) *Addrpb {
	a := &Addrpb{}
	err := proto.Unmarshal(buf, a)
	if log.Error(err) {
		return nil
	}
	return a
}

// GetAddr returns the Addrpb as an rnet.Addr
func (a *Addrpb) GetAddr() *rnet.Addr {
	return rnet.NewAddr(a.Ip, int(a.Port), a.Zone)
}

// Marshal a header to a buffer. If there is an error it is logged and nil is
// returned.
func (a *Addrpb) Marshal() []byte {
	buf, err := proto.Marshal(a)
	if log.Error(err) {
		return nil
	}
	return buf
}
