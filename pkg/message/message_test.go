package message

import (
	"encoding/json"
	"strings"
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
	messagetypeSchema := []byte(`json`)
	result, err := HandleWholeMsg(ip, port, netData, msgSchema, usernameSchema, passwordSchema, messagetypeSchema)
	assert.NoError(t, err)
	assert.Equal(t, (`{"IMEI":861714058319892,"addr":"1","cmd":"upload","ipFrom":"127.0.0.1","port":"7777"}`), string(result.Msg))
	assert.Equal(t, (`861714058319892`), string(result.Username))
	assert.Equal(t, (``), string(result.Password))
	_, err = json.Marshal(result)
	assert.NoError(t, err)

	// csv format
	msgSchema = []byte(`{{.| fjson}}`)
	usernameSchema = []byte(`{{.key5 | printf "%.f"}}`)
	passwordSchema = []byte(``)
	netData = (`1,2,3,a,4`)
	messagetypeSchema = []byte(`csv`)
	result, err = HandleWholeMsg(ip, port, netData, msgSchema, usernameSchema, passwordSchema, messagetypeSchema)
	assert.NoError(t, err)
	assert.Equal(t, (`{"ipFrom":"127.0.0.1","key1":1,"key2":2,"key3":3,"key4":"a","key5":4,"port":"7777"}`), string(result.Msg))
	assert.Equal(t, (`4`), string(result.Username))
	assert.Equal(t, (``), string(result.Password))
	_, err = json.Marshal(result)
	assert.NoError(t, err)
}

func TestHandleMsg(t *testing.T) {
	ip := "127.0.0.1"
	port := "7777"
	msgSchema := []byte(`{{. | fjson}}`)
	netData := (`{"addr":"1","cmd":"upload","IMEI":861714058319892}`)
	netDataJS := make(map[string]any)
	err := json.Unmarshal([]byte(strings.TrimSpace(netData)), &netDataJS)
	assert.NoError(t, err)
	result, err := handleMsg(ip, port, netDataJS, msgSchema)
	assert.NoError(t, err)
	assert.Equal(t, (`{"IMEI":861714058319892,"addr":"1","cmd":"upload","ipFrom":"127.0.0.1","port":"7777"}`), string(result))
	_, err = json.Marshal(result)
	assert.NoError(t, err)
}

func Test_csv2map(t *testing.T) {
	csvStr := `1,2,3,4,5`
	elementMap, err := csv2map(csvStr)
	assert.NoError(t, err)
	assert.Equal(t, (map[string]interface{}{"key1": 1., "key2": 2., "key3": 3., "key4": 4., "key5": 5.}), elementMap)

	// wrong msg format
	csvStr = "{1,2,3,4,5,6"
	elementMap, err = csv2map(csvStr)
	assert.Equal(t, (map[string]any{}), elementMap)
	assert.EqualErrorf(t, err, "csv2map: csvStr is not csv", "Error should be: %v, got: %v", "csv2map: csvStr is not csv", err)
}
