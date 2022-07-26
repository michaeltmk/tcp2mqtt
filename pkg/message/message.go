package message

import (
	"encoding/json"
	"strings"
	"tcp2mqtt/pkg/template"
)

type HandledMessage struct {
	Msg      string
	Username string
	Password string
}

//use template to hangle username
func handleUsername(netDataJS map[string]any, schema []byte) ([]byte, error) {
	bytesRepresentation, err := json.Marshal(netDataJS)
	if err != nil {
		return nil, err
	}
	result, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//use template to hangle password
func handlePassword(netDataJS map[string]any, schema []byte) ([]byte, error) {
	bytesRepresentation, err := json.Marshal(netDataJS)
	if err != nil {
		return nil, err
	}
	result, err := template.ApplyTemplate(schema, bytesRepresentation)
	if err != nil {
		return nil, err
	}
	return result, nil
}

//use template to hangle msg
func handleMsg(ip string, port string, netDataJS map[string]any, schema []byte) ([]byte, error) {
	bytesRepresentation, err := json.Marshal(netDataJS)
	if err != nil {
		return nil, err
	}
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

func HandleWholeMsg(ip string, port string, netData string, msgSchema []byte, usernameSchema []byte, passwordSchema []byte) (HandledMessage, error) {
	netDataJS := make(map[string]any)
	err := json.Unmarshal([]byte(strings.TrimSpace(netData)), &netDataJS)
	nilHandledMessage := HandledMessage{}
	if err != nil {
		return nilHandledMessage, err
	}
	msg, err := handleMsg(ip, port, netDataJS, msgSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	username, err := handleUsername(netDataJS, usernameSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	password, err := handlePassword(netDataJS, passwordSchema)
	if err != nil {
		return nilHandledMessage, err
	}
	return HandledMessage{string(msg), string(username), string(password)}, nil

}
