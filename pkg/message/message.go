package message

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"tcp2mqtt/pkg/template"
)

type HandledMessage struct {
	Msg      string
	Username string
	Password string
}

// use template to hangle username
func handleUsername(bytesRepresentation []byte, schema []byte) ([]byte, error) {
	result, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

type csvError struct {
	s string
}

func (e *csvError) Error() string {
	return e.s
}

// csv2map converts csv to map
func csv2map(csvStr string) (map[string]any, error) {
	elementMap := make(map[string]any)
	if csvStr[0] == '{' {
		err := &csvError{"csv2map: csvStr is not csv"}
		return map[string]any{}, err
	}
	s := strings.Split(csvStr, ",")
	for i := 0; i < len(s); i += 1 {
		num, err := strconv.ParseFloat(s[i], 64)
		if err != nil {
			elementMap[fmt.Sprintf("key%d", i+1)] = s[i]
		} else {
			elementMap[fmt.Sprintf("key%d", i+1)] = &num
		}
	}
	return elementMap, nil
}

// use template to hangle password
func handlePassword(bytesRepresentation []byte, schema []byte) ([]byte, error) {
	result, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// use template to hangle msg
func handleMsg(ip string, port string, bytesRepresentation []byte, schema []byte) ([]byte, error) {
	handledNetData, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	msgData := make(map[string]interface{})
	err = json.Unmarshal(handledNetData, &msgData)
	if err != nil {
		return nil, err
	}
	msgData["ipFrom"] = ip
	msgData["port"] = port
	return json.Marshal(msgData)
}

func HandleWholeMsg(ip string, port string, netData string, msgSchema []byte, usernameSchema []byte, passwordSchema []byte, messagetypeSchema []byte) (HandledMessage, error) {
	netDataJS := make(map[string]any)
	var err error
	nilHandledMessage := HandledMessage{}
	if string(messagetypeSchema) == "csv" {
		netDataJS, err = csv2map(netData)
		if err != nil {
			return nilHandledMessage, err
		}
	} else {
		err := json.Unmarshal([]byte(strings.TrimSpace(netData)), &netDataJS)
		if err != nil {
			return nilHandledMessage, err
		}
	}
	bytesRepresentation, err := json.Marshal(netDataJS)
	if err != nil {
		return nilHandledMessage, err
	}
	msg, err := handleMsg(ip, port, bytesRepresentation, msgSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	username, err := handleUsername(bytesRepresentation, usernameSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	password, err := handlePassword(bytesRepresentation, passwordSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	return HandledMessage{string(msg), string(username), string(password)}, nil

}
