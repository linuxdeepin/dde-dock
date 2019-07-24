package network

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"time"

	dbus "pkg.deepin.io/lib/dbus1"
	"pkg.deepin.io/lib/dbusutil"
)

const (
	typeEchoReply     = 0
	typeEchoRequest   = 8
	typeNetUnrechable = 3

	codeNetUnrechable      = 0
	codeHostUnrechable     = 1
	codeProtocolUnrechable = 2
	codePortUnrechable     = 3
)

// ICMP icmp protocol header
type ICMP struct {
	Type        uint8
	Code        uint8
	CheckSum    uint16
	Identifier  uint16
	SequenceNum uint16
}

func calcCheckSum(data []byte) uint16 {
	var (
		sum    uint32
		idx    int
		length = len(data)
	)

	for length > 1 {
		sum += uint32(data[idx])<<8 + uint32(data[idx+1])
		idx += 2
		length -= 2
	}
	if length > 0 {
		sum += uint32(data[idx])
	}
	sum += (sum >> 16)

	return uint16(^sum)
}

var getSequenceNum = func() func() uint16 {
	var num uint16
	return func() uint16 {
		num++
		return num
	}
}()

func makeEchoReqHeader() *ICMP {
	var icmp = ICMP{
		Type:        typeEchoRequest,
		SequenceNum: getSequenceNum(),
	}
	var buffer bytes.Buffer
	_ = binary.Write(&buffer, binary.BigEndian, &icmp)
	icmp.CheckSum = calcCheckSum(buffer.Bytes())
	buffer.Reset()

	return &icmp
}

func sendEchoRequest(conn *net.IPConn) error {
	var buffer bytes.Buffer
	var icmp = makeEchoReqHeader()

	err := binary.Write(&buffer, binary.BigEndian, icmp)
	if err != nil {
		return err
	}

	_, err = conn.Write(buffer.Bytes())
	if err != nil {
		return err
	}
	_ = conn.SetDeadline(time.Now().Add(time.Second * 3))
	buffer.Reset()

	return nil
}

func recvEchoReply(conn *net.IPConn) (*ICMP, error) {
	var reply = make([]byte, 1024)
	_, err := conn.Read(reply)
	if err != nil {
		return nil, err
	}
	if len(reply) < 4 {
		return nil, fmt.Errorf("invalid icmp reply")
	}
	var icmp ICMP
	err = binary.Read(bytes.NewBuffer(reply), binary.BigEndian, &icmp)
	if err != nil {
		return nil, err
	}
	// get real type from ipv4
	// the header contains ipv4 header: version(0x4) and length(0x5)
	// TODO(jouyouyun): optimization
	icmp.Type = uint8(reply[0] & 0xba)
	return &icmp, nil
}

func newICMPConn(host string) (*net.IPConn, error) {
	raddr, err := net.ResolveIPAddr("ip4", host)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialIP("ip4:icmp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func handleICMPReply(icmp *ICMP) error {
	if icmp.Type == typeEchoReply {
		return nil
	}
	if icmp.Type != typeNetUnrechable {
		return fmt.Errorf("unknown error")
	}

	var msg string
	switch icmp.Code {
	case codeNetUnrechable:
		msg = "network unreachable"
	case codeHostUnrechable:
		msg = "host unreachable"
	case codeProtocolUnrechable:
		msg = "protocol unreachable"
	case codePortUnrechable:
		msg = "host port unreachable"
	}
	return fmt.Errorf(msg)
}

// Ping ping remote host, blocked operation.
func (n *Network) Ping(host string) *dbus.Error {
	conn, err := newICMPConn(host)
	if err != nil {
		return dbusutil.ToError(err)
	}
	defer conn.Close()

	err = sendEchoRequest(conn)
	if err != nil {
		return dbusutil.ToError(err)
	}

	icmp, err := recvEchoReply(conn)
	if err != nil {
		return dbusutil.ToError(err)
	}

	logger.Debugf("Reply: %#v", icmp)
	err = handleICMPReply(icmp)
	if err != nil {
		return dbusutil.ToError(err)
	}
	return nil
}
