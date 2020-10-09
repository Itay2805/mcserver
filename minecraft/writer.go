package minecraft

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/google/uuid"
	"log"
)

type Writer struct {
	bytes.Buffer
}

func (writer *Writer) Bytes() []byte {
	return writer.Buffer.Bytes()[:writer.Buffer.Len()]
}

func (writer *Writer) WriteBytes(data []byte) {
	writer.Buffer.Write(data)
}

func (writer *Writer) WriteBoolean(val bool) {
	if val {
		writer.WriteByte(0x01)
	} else {
		writer.WriteByte(0x00)
	}
}

func (writer *Writer) WriteByte(val byte) {
	writer.Buffer.WriteByte(val)
}

func (writer *Writer) WriteShort(val int16) {
	writer.Grow(2)
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteUShort(val uint16) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteInt(val int32) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteUInt(val uint32) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteLong(val int64) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteULong(val uint64) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteFloat(val float32) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteDouble(val float64) {
	_ = binary.Write(&writer.Buffer, binary.BigEndian, val)
}

func (writer *Writer) WriteString(val string) {
	writer.WriteVarint(int32(len(val)))
	writer.Buffer.WriteString(val)
}

func (writer *Writer) WriteChat(val Chat) {
	b := val.ToJSON()
	writer.WriteVarint(int32(len(b)))
	writer.WriteBytes(b)
}

func (writer *Writer) WriteJson(val interface{}) {
	b, err := json.Marshal(val)
	if err != nil {
		log.Panicln(err)
	}
	writer.WriteVarint(int32(len(b)))
	writer.WriteBytes(b)
}

// TODO: writeIdentifier

func (writer *Writer) WriteVarint(val int32) {
	raw := uint32(val)
	for {
		temp := byte(raw & 0b01111111)
		raw >>= 7
		if raw != 0 {
			temp |= 0b10000000
		}
		writer.WriteByte(temp)
		if raw == 0 {
			break
		}
	}
}

func (writer *Writer) WriteVarlong(val int64) {
	raw := uint64(val)
	for {
		temp := byte(raw & 0b01111111)
		raw >>= 7
		if raw != 0 {
			temp |= 0b10000000
		}
		writer.WriteByte(temp)
		if raw == 0 {
			break
		}
	}
}

// TODO: entity metaData

// TODO: SLot

func (writer *Writer) WritePosition(pos Position) {
	writer.WriteULong(pos.Pack())
}

func (writer *Writer) WriteAngle(angle Angle) {
	writer.WriteByte(byte(angle))
}

func (writer *Writer) WriteUUID(val uuid.UUID) {
	b, _ := val.MarshalBinary()
	writer.WriteBytes(b)
}

func (writer *Writer) WriteUUIDAsString(val uuid.UUID) {
	b, _ := val.MarshalText()
	writer.WriteVarint(int32(len(b)))
	writer.WriteBytes(b)
}

func (writer *Writer) StartNBT() NbtWriter {
	return NbtWriter{
		w:                   writer,
		hierarchy:           []uint8{},
		listSizeStack:       []int{},
		listSizeOffsetStack: []int{},
	}
}

func (writer *Writer) StartEntityMetadata() *EntityMetadataWriter {
	return &EntityMetadataWriter{
		writer,
	}
}

type EntityMetadataWriter struct {
	w *Writer
}

type Pose int32

const (
	PoseStanding = Pose(0)
	PoseFallFlying = Pose(1)
	PoseSleeping = Pose(2)
	PoseSwimming = Pose(3)
	PoseSpinAttack = Pose(4)
	PoseSneaking = Pose(5)
	PoseDying = Pose(6)
)

func (writer *EntityMetadataWriter) WriteByte(index byte, val byte) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(0)
	writer.w.WriteByte(val)
}

func (writer *EntityMetadataWriter) WriteVarint(index byte, val int32) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(1)
	writer.w.WriteVarint(val)
}
func (writer *EntityMetadataWriter) WriteFloat(index byte, val float32) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(2)
	writer.w.WriteFloat(val)
}
func (writer *EntityMetadataWriter) WriteString(index byte, val string) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(3)
	writer.w.WriteString(val)
}
func (writer *EntityMetadataWriter) WriteBoolean(index byte, val bool) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(6)
	writer.w.WriteBoolean(val)
}

func (writer *EntityMetadataWriter) WritePosition(index byte, val Position) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(8)
	writer.w.WritePosition(val)
}

func (writer *EntityMetadataWriter) StartNBT(index byte) NbtWriter {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(13)
	return NbtWriter{
		w:                   writer.w,
		hierarchy:           []uint8{},
		listSizeStack:       []int{},
		listSizeOffsetStack: []int{},
	}
}

func (writer *EntityMetadataWriter) WritePose(index byte, val Pose) {
	writer.w.WriteByte(index)
	writer.w.WriteVarint(18)
	writer.w.WriteVarint(int32(val))
}

func (writer *EntityMetadataWriter) Done() {
	writer.w.WriteByte(0xFF)
}
