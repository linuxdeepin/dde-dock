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

	ipHeaderMinLen = 20
	ipProtocolICMP = 1
)

/**
 * IP Header
 *
 *    0     4         8         16    19       32
 *    | 版本 | 头部长度 | 区分服务 |    总长度     |  --->
 *    | 标识                    | 标志 | 片偏移  |      |
 *    |     生存时间   |   协议  |   首部校验和   |       | --> 固定部分
 *    |                 源地址                  |      |
 *    |                 目标地址                |  ---->
 *    |         可选字段(长度可变)       |  填充  |
 *    |                 数据部分                |
 *
 *    头部长度：除去数据部分之外的长度(头部总字节数/4)，最小为 5,最大为 15。因为头部长度只占 4 位，最大能够表达的值为 15.
 *             通过这也可知，头部最小有 20 字节数据，最大有 60 字节数据
 *
**/

// IPHeader store ip header data
type IPHeader struct {
	Version     uint8 // the first byte[0:4]
	Length      uint8 // the first byte[4:8]
	TOS         uint8 // type of service
	TotalLen    uint16
	Identifier  uint16
	Flag        uint8  // only has 3 bits
	Offset      uint16 // has 13 bits
	TTL         uint8
	Protocol    uint8 // ICMP: 1, TCP: 6, UDP: 7。 定义在 /etc/protocols
	CheckSum    uint16
	Source      [4]uint8
	Destination [4]uint8
	Options     [40]uint8 // 可选字段，根据 Length 来判断是否存在
}

/**
 * ICMP Header:
 *
 *    0     8     16          32
 *    | 类型 | 代码 |   校验和   |
 *    |    标识符   |   序列号   |
 *    |  选项(根据类型和代码而定)  |
 *
 * icmp 实际是附加在 ip 协议上的，所有收到的响应要先剥离 IP Header
 *
 **/

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
	_ = conn.SetDeadline(time.Now().Add(time.Second * 5))
	buffer.Reset()

	return nil
}

func recvEchoReply(conn *net.IPConn) (*ICMP, error) {
	var reply = make([]byte, 1024)
	_, err := conn.Read(reply)
	if err != nil {
		return nil, err
	}

	ipHeader, err := unmarshalIPHeader(reply)
	if err != nil {
		return nil, err
	}
	if ipHeader.Protocol != ipProtocolICMP {
		return nil, fmt.Errorf("not excepted protocol: %d", ipHeader.Protocol)
	}

	var icmp ICMP
	idx := int(ipHeader.Length * 4)
	err = binary.Read(bytes.NewBuffer(reply[idx:]), binary.BigEndian, &icmp)
	if err != nil {
		return nil, err
	}
	return &icmp, nil
}

func unmarshalIPHeader(datas []byte) (*IPHeader, error) {
	if len(datas) < ipHeaderMinLen {
		return nil, fmt.Errorf("invalid ip data: %v", datas)
	}

	var (
		header IPHeader
		i8     uint8
		i16    uint16
	)
	_ = binary.Read(bytes.NewBuffer(datas[:1]), binary.BigEndian, &i8)
	header.Version = i8 & 0xf0
	header.Length = i8 & 0xf
	_ = binary.Read(bytes.NewBuffer(datas[1:2]), binary.BigEndian, &header.TOS)
	_ = binary.Read(bytes.NewBuffer(datas[2:4]), binary.BigEndian, &header.TotalLen)
	_ = binary.Read(bytes.NewBuffer(datas[4:6]), binary.BigEndian, &header.Identifier)
	_ = binary.Read(bytes.NewBuffer(datas[6:7]), binary.BigEndian, &i8)
	header.Flag = i8 & 0xd0
	_ = binary.Read(bytes.NewBuffer(datas[6:8]), binary.BigEndian, &i16)
	header.Offset = i16 & 0x1f
	_ = binary.Read(bytes.NewBuffer(datas[8:9]), binary.BigEndian, &header.TTL)
	_ = binary.Read(bytes.NewBuffer(datas[9:10]), binary.BigEndian, &header.Protocol)
	_ = binary.Read(bytes.NewBuffer(datas[10:12]), binary.BigEndian, &header.CheckSum)

	idx := 12
	for i := 0; i < len(header.Source); i++ {
		_ = binary.Read(bytes.NewBuffer(datas[idx:idx+1]), binary.BigEndian, &header.Source[i])
		idx++
	}

	for i := 0; i < len(header.Destination); i++ {
		_ = binary.Read(bytes.NewBuffer(datas[idx:idx+1]), binary.BigEndian, &header.Destination[i])
		idx++
	}

	hlen := int(header.Length * 4)
	for i := 0; idx < hlen; {
		_ = binary.Read(bytes.NewBuffer(datas[idx:idx+1]), binary.BigEndian, &header.Options[i])
		idx++
	}

	return &header, nil
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
