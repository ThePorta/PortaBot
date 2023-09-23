package types

// Code generated by github.com/tinylib/msgp DO NOT EDIT.

import (
	"github.com/tinylib/msgp/msgp"
)

// DecodeMsg implements msgp.Decodable
func (z *AccountInfo) DecodeMsg(dc *msgp.Reader) (err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, err = dc.ReadMapHeader()
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, err = dc.ReadMapKeyPtr()
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "WalletConnectSession":
			z.WalletConnectSession, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "WalletConnectSession")
				return
			}
		case "ChatId":
			z.ChatId, err = dc.ReadString()
			if err != nil {
				err = msgp.WrapError(err, "ChatId")
				return
			}
		default:
			err = dc.Skip()
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	return
}

// EncodeMsg implements msgp.Encodable
func (z AccountInfo) EncodeMsg(en *msgp.Writer) (err error) {
	// map header, size 2
	// write "WalletConnectSession"
	err = en.Append(0x82, 0xb4, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e)
	if err != nil {
		return
	}
	err = en.WriteString(z.WalletConnectSession)
	if err != nil {
		err = msgp.WrapError(err, "WalletConnectSession")
		return
	}
	// write "ChatId"
	err = en.Append(0xa6, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64)
	if err != nil {
		return
	}
	err = en.WriteString(z.ChatId)
	if err != nil {
		err = msgp.WrapError(err, "ChatId")
		return
	}
	return
}

// MarshalMsg implements msgp.Marshaler
func (z AccountInfo) MarshalMsg(b []byte) (o []byte, err error) {
	o = msgp.Require(b, z.Msgsize())
	// map header, size 2
	// string "WalletConnectSession"
	o = append(o, 0x82, 0xb4, 0x57, 0x61, 0x6c, 0x6c, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x53, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e)
	o = msgp.AppendString(o, z.WalletConnectSession)
	// string "ChatId"
	o = append(o, 0xa6, 0x43, 0x68, 0x61, 0x74, 0x49, 0x64)
	o = msgp.AppendString(o, z.ChatId)
	return
}

// UnmarshalMsg implements msgp.Unmarshaler
func (z *AccountInfo) UnmarshalMsg(bts []byte) (o []byte, err error) {
	var field []byte
	_ = field
	var zb0001 uint32
	zb0001, bts, err = msgp.ReadMapHeaderBytes(bts)
	if err != nil {
		err = msgp.WrapError(err)
		return
	}
	for zb0001 > 0 {
		zb0001--
		field, bts, err = msgp.ReadMapKeyZC(bts)
		if err != nil {
			err = msgp.WrapError(err)
			return
		}
		switch msgp.UnsafeString(field) {
		case "WalletConnectSession":
			z.WalletConnectSession, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "WalletConnectSession")
				return
			}
		case "ChatId":
			z.ChatId, bts, err = msgp.ReadStringBytes(bts)
			if err != nil {
				err = msgp.WrapError(err, "ChatId")
				return
			}
		default:
			bts, err = msgp.Skip(bts)
			if err != nil {
				err = msgp.WrapError(err)
				return
			}
		}
	}
	o = bts
	return
}

// Msgsize returns an upper bound estimate of the number of bytes occupied by the serialized message
func (z AccountInfo) Msgsize() (s int) {
	s = 1 + 21 + msgp.StringPrefixSize + len(z.WalletConnectSession) + 7 + msgp.StringPrefixSize + len(z.ChatId)
	return
}
