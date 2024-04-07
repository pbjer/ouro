package prompt

import (
	"bytes"
	"text/template"
)

type Prompt struct {
	template string
	data     any
}

func New(template string, data any) *Prompt {
	return &Prompt{
		template: template,
		data:     data,
	}
}

func (p *Prompt) Render() (string, error) {
	var buf bytes.Buffer
	tmpl, err := template.New("prompt").Parse(p.template)
	if err != nil {
		return "", err
	}

	err = tmpl.Execute(&buf, p.data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
