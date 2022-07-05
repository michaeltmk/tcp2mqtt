package main

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleMsg(t *testing.T) {
	ip := "127.0.0.1"
	port := "7777"
	schema := []byte(`{"{{.netData.IMEI | printf "%.f"}}":[{"addr":"{{.netData.addr}}","cmd":"{{.netData.cmd}}","count":{{.netData.count}},"ipFrom":"{{.ip}}","port":"{{.port}}"}]}`)
	netData := (`{"addr":"1","cmd":"upload","count":[0,0,0,0,0,0,0],"IMEI":861714058319892}`)
	result, err := handleMsg(ip, port, netData, schema)
	assert.NoError(t, err)
	assert.Equal(t, (`{"861714058319892":[{"addr":"1","cmd":"upload","count":[0 0 0 0 0 0 0],"ipFrom":"127.0.0.1","port":"7777"}]}`), string(result))
	_, err = json.Marshal(result)
	assert.NoError(t, err)
}
