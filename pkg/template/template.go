package template

import (
	"bytes"
	"text/template"

	"github.com/bytedance/sonic"
	"github.com/bytedance/sonic/encoder"
)

// orderedMarshalString marshals a given value into JSON with ordered keys
func orderedMarshalString(v any) (string, error) {
	b, err := encoder.Encode(v, encoder.SortMapKeys)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// ApplyTemplate gets template of an eventID and applies it
func ApplyTemplate(templateData, data []byte) ([]byte, error) {
	t, err := template.New("template").Funcs(
		template.FuncMap{
			"fjson": orderedMarshalString,
		},
	).Parse(string(templateData))
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
