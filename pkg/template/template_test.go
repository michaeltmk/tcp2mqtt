package template

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	templateData := []byte(`{"{{.netData.IMEI | printf "%.f"}}":[{"addr":"{{.netData.addr}}","cmd":"{{.netData.cmd}}","count":{{.netData.count}},"ipFrom":"{{.ip}}","port":"{{.port}}"}]}`)
	data := []byte(`{"netData":{"addr":"1","cmd":"upload","count":[0,0,0,0,0,0,0],"IMEI":861714058319892},"ip":"127.0.0.1","port":"7777"}`)
	result, err := ApplyTemplate(templateData, data)

	assert.NoError(t, err)
	assert.Equal(t, `{"861714058319892":[{"addr":"1","cmd":"upload","count":[0 0 0 0 0 0 0],"ipFrom":"127.0.0.1","port":"7777"}]}`, string(result))

}

func TestApplyTemplate(t *testing.T) {
	templateData := []byte("hello {{.Name}}")
	data := []byte(`{"Name": "world"}`)
	result, err := ApplyTemplate(templateData, data)
	assert.NoError(t, err)
	assert.Equal(t, []byte("hello world"), result)

}
