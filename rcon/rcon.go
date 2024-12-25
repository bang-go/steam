package rcon

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/bang-go/network/tcpx"
	"time"
)

const (
	RespAuthFailed = -1 //认证失败
)

type Rcon interface {
	Dail() error
	Auth(string) error
	ExecCommand([]byte) ([]byte, error)
	Close()
}
type Config struct {
	Addr    string
	Timeout time.Duration
}

type rconEntity struct {
	*Config
	conn         tcpx.Connect
	lastPacketID int32
	isAuthed     bool
}

func New(conf *Config) Rcon {
	return &rconEntity{Config: conf}
}

func (s *rconEntity) Dail() (err error) {
	client := tcpx.NewClient(&tcpx.ClientConfig{Addr: s.Addr, Timeout: s.Timeout})
	s.conn, err = client.Dail()
	return
}

func (s *rconEntity) Auth(password string) (err error) {
	s.lastPacketID++
	packet := NewPacket(DataTypeAuth, s.lastPacketID, []byte(password))
	err = s.writePacket(packet)
	if err != nil {
		return
	}
	respPacket, err := s.readPacket()
	if err != nil {
		return
	}
	return s.isAuthSuccess(respPacket)
}

// ExecCommand 执行命令
func (s *rconEntity) ExecCommand(command []byte) (body []byte, err error) {
	if s.isAuthed == false {
		err = errors.New("未认证")
		return
	}
	//发送命令
	s.lastPacketID++
	packet := NewPacket(DataTypeExecCommand, s.lastPacketID, command)
	err = s.writePacket(packet)
	if err != nil {
		return
	}
	//发送空包 为了循环接受响应直到空的响应为止 抓包发现不需要
	//err = s.writePacket(s.getCommandEmptyPackage())
	//if err != nil {
	//	return
	//}
	respPacket, err := s.readPacket()
	if err != nil {
		return
	}
	err = s.checkID(respPacket)
	if err != nil {
		return
	}
	body = respPacket.Body()
	return
}

func (s *rconEntity) isAuthSuccess(packet *Packet) (err error) {
	if packet.Type != DataTypeAuthResponse {
		err = errors.New(fmt.Sprintf("不匹配的type类型,type:%d", packet.Type))
		return
	}
	if packet.ID == RespAuthFailed { //认证失败，密码校验失败
		err = errors.New(fmt.Sprintf("认证失败"))
		return
	}
	if packet.ID > RespAuthFailed {
		s.isAuthed = true
		return
	}
	return
}

// 通过tcp发送packet
func (s *rconEntity) writePacket(packet *Packet) (err error) {
	data, err := s.encodePacket(packet)
	if err != nil {
		return
	}
	err = s.conn.Send(data)
	return
}

// 通过tcp连接读取packet
func (s *rconEntity) readPacket() (p *Packet, err error) {
	//input io.Reader
	reader := bufio.NewReader(s.conn.Conn())
	p = &Packet{}
	err = binary.Read(reader, binary.LittleEndian, &p.Size)
	if err != nil {
		return
	}
	if p.Size < MinPacketSize {
		err = fmt.Errorf("size小于最小值,size:%d", p.Size)
		return
	}

	err = binary.Read(reader, binary.LittleEndian, &p.ID)
	if err != nil {
		return
	}
	err = binary.Read(reader, binary.LittleEndian, &p.Type)
	if err != nil {
		return
	}

	bodyLen := p.Size - (PacketIdSize + PacketTypeSize)
	bodyBuf := make([]byte, bodyLen)
	err = binary.Read(reader, binary.LittleEndian, &bodyBuf)
	if err != nil {
		return
	}
	if len(bodyBuf) < int(PacketPaddingSize) {
		err = fmt.Errorf("异常的body大小, lne:%d", len(bodyBuf))
	}
	p.body = bodyBuf[:bodyLen-PacketPaddingSize] //删除后两位last null terminated ascii
	return
}

func (s *rconEntity) encodePacket(p *Packet) (data []byte, err error) {
	var buf bytes.Buffer
	_ = binary.Write(&buf, binary.LittleEndian, p.Size)
	_ = binary.Write(&buf, binary.LittleEndian, p.ID)
	_ = binary.Write(&buf, binary.LittleEndian, p.Type)
	_ = binary.Write(&buf, binary.LittleEndian, p.body)
	_ = binary.Write(&buf, binary.LittleEndian, [2]byte{})
	if buf.Len() > int(MaxPacketSize) {
		err = fmt.Errorf("packet长度超出最大限制,len:%d", buf.Len())
		return
	}
	data = buf.Bytes()
	return
}

//	func (s *rconEntity) getCommandEmptyPackage() *Packet {
//		return NewPacket(DataTypeResponseValue, s.lastPacketID, "")
//	}
//
// 检测请求包是否和请求包匹配
func (s *rconEntity) checkID(p *Packet) (err error) {
	if p.ID != s.lastPacketID {
		err = fmt.Errorf("packet包中ID不匹配: origin_id:%d,current_id:%d", s.lastPacketID, p.ID)
		return
	}
	return
}

func (s *rconEntity) Close() {
	s.conn.Close()
}
