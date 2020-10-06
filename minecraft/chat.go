package minecraft

import (
	"encoding/json"
	"log"
)

type Chat struct {
	Text			string 	`json:"text,omitempty"`
	Bold			bool	`json:"bold,omitempty"`
	Italic			bool	`json:"italic,omitempty"`
	Underlined		bool	`json:"underlined,omitempty"`
	Strikethrough	bool	`json:"strikethrough,omitempty"`
	Obfuscated		bool	`json:"obfuscated,omitempty"`
	Color			string	`json:"color,omitempty"`
	Translate		string 	`json:"translate,omitempty"`
	With			[]Chat	`json:"with,omitempty"`
	Extra			[]Chat	`json:"extra,omitempty"`
}

func NewChat(js []byte) Chat {
	m := Chat{}
	if js[0] == '"' {
		err := json.Unmarshal(js, &m.Text)
		if err != nil {
			log.Panicln("Error decoding chat object")
		}
	} else {
		err := json.Unmarshal(js, m)
		if err != nil {
			log.Panicln("Error decoding chat object")
		}
	}
	return m
}

func (m Chat) ToJSON() []byte {
	code, err := json.Marshal(m)
	if err != nil {
		log.Panicln(err)
	}
	return code
}

func Text(str string) Chat {
	return Chat{
		Text: str,
	}
}
