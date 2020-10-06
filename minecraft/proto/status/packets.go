package status

import "github.com/itay2805/mcserver/minecraft"

type ServerListVersion struct {
	Name string `json:"name"`
	Protocol int32 `json:"protocol"`
}

type ServerListPlayerSample struct {
	Name string `json:"name"`
	Id string `json:"id"`
}

type ServerListPlayers struct {
	Max int `json:"max"`
	Online int `json:"online"`
	Sample []ServerListPlayerSample `json:"sample"`
}

type ServerListResponse struct {
	Version ServerListVersion `json:"version"`
	Players ServerListPlayers `json:"players"`
	Description minecraft.Chat `json:"description"`
	Favicon string `json:"favicon,omitempty"`
}

type Response struct {
	Response ServerListResponse
}

func (r Response) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x00)
	writer.WriteJson(r.Response)
}

type Pong struct {
	Payload int64
}

func (r Pong) Encode(writer *minecraft.Writer) {
	writer.WriteVarint(0x01)
	writer.WriteLong(r.Payload)
}