package dbfm

import (
	"encoding/json"
	"testing"
)

func TestGetChannels(t *testing.T) {
	GetChannels()
}

func TestParseChannels(t *testing.T) {
	var s string = `
{"channels":[{"name_en":"Personal Radio","seq_id":0,"abbr_en":"My","name":"私人兆赫","channel_id":0}]}`
	channels := new(FmChannels)
	json.Unmarshal([]byte(s), channels)
}
