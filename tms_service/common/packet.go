package common

import (
	"bytes"
	"encoding/binary"
)

//默认udp包大小
const BUFSIZE = 65535

var buff = make([]byte, BUFSIZE)

type PacketHeader struct {
	Stx      byte
	Len      uint32
	Ver      uint8
	SeqId    uint32
	CmdId    uint16
	ClientIp uint32
	Resv     [10]byte
}

const PKG_HEAD_LEN = 26
const C_STX = 'H'
const C_VER = 0

// 自定义包
type Packet struct {
	header PacketHeader
	body   []byte
}

func NewPacket(seqId uint32, body []byte) *Packet {
	return &Packet{
		header: PacketHeader{
			Stx:      C_STX,
			Len:      uint32(PKG_HEAD_LEN + len(body)),
			Ver:      C_VER,
			SeqId:    seqId,
			CmdId:    0,
			ClientIp: 0,
		},
		body: body,
	}
}

func (p *Packet) GetSeqId() uint32 {
	return p.header.SeqId
}

func (p *Packet) GetLen() uint32 {
	return p.header.Len
}

func (p *Packet) GetBody() []byte {
	return p.body
}

// 打包
func (p *Packet) Pack() []byte {
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.BigEndian, p.header)
	binary.Write(buf, binary.BigEndian, p.body)
	return buf.Bytes()
}

// 解包
func Unpack(buf []byte) (*Packet, error) {
	packet := Packet{}
	buffer := bytes.NewBuffer(buf)
	if err := binary.Read(buffer, binary.BigEndian, &packet.header); err != nil {
		return nil, err
	}
	packet.body = buffer.Next(buffer.Len())
	return &packet, nil
}
