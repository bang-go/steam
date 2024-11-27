package rcon

const (
	PacketSizeSelfSize int32 = 4 // size字段本身的大小(size字段的值没有包含其本身的大小)
	PacketIdSize       int32 = 4 //ID大小
	PacketTypeSize     int32 = 4 //类型大小
	PacketPaddingSize  int32 = 2 //填充大小(包括body结束,包尾部)
	MinPacketSize            = PacketIdSize + PacketTypeSize + PacketPaddingSize
	MaxPacketSize      int32 = 4096
)
const (
	DataTypeAuth          = PacketType(3) // Authentication packet
	DataTypeAuthResponse  = PacketType(2) // Authentication response packet
	DataTypeExecCommand   = PacketType(2) // Command packet
	DataTypeResponseValue = PacketType(0) // Response packet
)

type PacketType int32

type Packet struct {
	Size int32
	ID   int32
	Type PacketType
	body []byte
}

func NewPacket(packageType PacketType, packetId int32, bodyStr string) *Packet {
	size := len([]byte(bodyStr)) + int(MinPacketSize)
	return &Packet{
		Size: int32(size),
		Type: packageType,
		ID:   packetId,
		body: []byte(bodyStr),
	}
}

func (s *Packet) Body() string {
	return string(s.body)
}
