package login

import (
	"github.com/google/uuid"
	"github.com/itay2805/mcserver/minecraft"
)

type Disconnect struct {
	Reason minecraft.Chat
}

func (r Disconnect) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x00)
	writer.WriteChat(r.Reason)
}

type EncryptionRequest struct {
	ServerId	string
	PublicKey 	[]byte
	VerifyToken	[]byte
}

func (r EncryptionRequest) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x01)
	writer.WriteString(r.ServerId)
	writer.WriteVarint(int32(len(r.PublicKey)))
	writer.WriteBytes(r.PublicKey)
	writer.WriteVarint(int32(len(r.VerifyToken)))
	writer.WriteBytes(r.VerifyToken)
}

type LoginSuccess struct {
	Uuid 		uuid.UUID
	Username 	string
}

func (r LoginSuccess) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x02)
	writer.WriteString(r.Uuid.String())
	writer.WriteString(r.Username)
}

type SetCompression struct {
	Threshold 	int32
}

func (r SetCompression) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x03)
	writer.WriteVarint(r.Threshold)
}