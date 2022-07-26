package message

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleWholeMsg(t *testing.T) {
	ip := "127.0.0.1"
	port := "7777"
	msgSchema := []byte(`{{. | fjson}}`)
	usernameSchema := []byte(`{{.IMEI | printf "%.f"}}`)
	passwordSchema := []byte(``)
	netData := (`{"addr":"1","cmd":"upload","IMEI":861714058319892}`)
	result, err := HandleWholeMsg(ip, port, netData, msgSchema, usernameSchema, passwordSchema)
	assert.NoError(t, err)
	assert.Equal(t, (`{"IMEI":861714058319892,"addr":"1","cmd":"upload","ipFrom":"127.0.0.1","port":"7777"}`), string(result.Msg))
	assert.Equal(t, (`861714058319892`), string(result.Username))
	assert.Equal(t, (``), string(result.Password))
	_, err = json.Marshal(result)
	assert.NoError(t, err)
}
