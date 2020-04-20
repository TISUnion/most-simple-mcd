package utils

import "encoding/json"

type TellrowMessage map[string]string

func NewTellrowMessage() TellrowMessage {
	return make(TellrowMessage)
}

func (t TellrowMessage) SetText(text string) TellrowMessage {
	t["text"] = text
	return t
}

func (t TellrowMessage) SetTextColor(text, color string) TellrowMessage {
	t["text"] = text
	t["color"] = color
	return t
}

func (t TellrowMessage) SetParam(param, val string) TellrowMessage {
	t[param] = val
	return t
}

func (t TellrowMessage) JSON() string {
	rawmsgbyte, _ := json.Marshal(t)
	return string(rawmsgbyte)
}
