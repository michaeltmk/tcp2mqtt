package template

import (
	"bytes"
	"text/template"

	"github.com/bytedance/sonic"
)

// ApplyTemplate gets template of an eventID and applies it
func ApplyTemplate(templateData, data []byte) ([]byte, error) {
	t, err := template.New("event-template").Parse(string(templateData))
	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)

	dataMap := make(map[string]any)
	if err := sonic.Unmarshal(data, &dataMap); err != nil {
		return nil, err
	}
	if err := t.Execute(buf, dataMap); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
