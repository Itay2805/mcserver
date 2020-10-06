package minecraft

import (
	"encoding/binary"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math"
)

type Reader struct {
	Data 	[]byte
	Offset 	int
}

func (reader *Reader) Read(p []byte) (int, error) {
	n := copy(p, reader.Data[reader.Offset:])
	reader.Offset += n
	return n, nil
}

func (reader *Reader) ReadBytes(size int) []byte {
	if reader.Offset + size > len(reader.Data) {
		log.Panicln("Got to end of packet!")
	}
	Data := reader.Data[reader.Offset:reader.Offset + size]
	reader.Offset += size
	return Data
}

func (reader *Reader) ReadBoolean() bool {
	b := reader.ReadByte()
	if b == 0x01 {
		return true
	} else if b == 0x00 {
		return false
	} else {
		panic(fmt.Sprint("Invalid boolean value ", b))
	}
}

func (reader *Reader) ReadByte() byte {
	return reader.ReadBytes(1)[0]
}

func (reader *Reader) ReadShort() int16 {
	return int16(binary.BigEndian.Uint16(reader.ReadBytes(2)))
}

func (reader *Reader) ReadUShort() uint16 {
	return binary.BigEndian.Uint16(reader.ReadBytes(2))
}

func (reader *Reader) ReadInt() int32 {
	return int32(binary.BigEndian.Uint32(reader.ReadBytes(4)))
}

func (reader *Reader) ReadUInt() uint32 {
	return binary.BigEndian.Uint32(reader.ReadBytes(4))
}

func (reader *Reader) ReadLong() int64 {
	return int64(binary.BigEndian.Uint64(reader.ReadBytes(8)))
}

func (reader *Reader) ReadULong() uint64 {
	return binary.BigEndian.Uint64(reader.ReadBytes(8))
}

func (reader *Reader) ReadFloat() float32 {
	return math.Float32frombits(binary.BigEndian.Uint32(reader.ReadBytes(4)))
}

func (reader *Reader) ReadDouble() float64 {
	return math.Float64frombits(binary.BigEndian.Uint64(reader.ReadBytes(8)))
}

func (reader *Reader) ReadString(maxLen int) string {
	l := int(reader.ReadVarint())
	if l > maxLen {
		log.Panicln("String is too large!")
	}
	return string(reader.ReadBytes(l))
}

func (reader *Reader) ReadChat() Chat {
	l := int(reader.ReadVarint())
	if l > 32767 {
		log.Panicln("String is too large!")
	}
	return NewChat(reader.ReadBytes(32767))
}

func (reader *Reader) ReadIdentifier() string {
	return reader.ReadString(32767)
}

func (reader *Reader) ReadVarint() int32 {
	numRead := 0
	result := int32(0)
	for {
		read := reader.ReadByte()
		value := read & 0b01111111
		result |= int32(value) << (7 * numRead)

		numRead++
		if numRead > 5 {
			log.Panicln("Varint is too big")
		}

		if (read & 0b10000000) == 0 {
			return result
		}
	}
}

func (reader *Reader) ReadVarlong() int64 {
	numRead := 0
	result := int64(0)
	for {
		read := reader.ReadByte()
		value := read & 0b01111111
		result |= int64(value) << (7 * numRead)

		numRead++
		if numRead > 10 {
			log.Panicln("Varlong is too big")
		}

		if (read & 0b10000000) == 0 {
			return result
		}
	}
}

// TODO: entity metaData

// TODO: Slot

// TODO: NBT Tag

func (reader *Reader) ReadPosition() Position {
	return ParsePosition(reader.ReadULong())
}

func (reader *Reader) ReadAngle() Angle {
	return Angle(reader.ReadByte())
}

func (reader *Reader) ReadUUID() uuid.UUID {
	u, _ := uuid.FromBytes(reader.ReadBytes(16))
	return u
}

func (reader *Reader) ReadUUIDFromString() uuid.UUID {
	return uuid.MustParse(reader.ReadString(36))
}
