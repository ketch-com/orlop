package config

import (
	"fmt"
	"strings"
)

type configTag struct {
	Name         *string
	Encoding     *string
	DefaultValue *string
	Required     bool
}

func (t configTag) String() string {
	var s []string
	if t.Name != nil {
		s = append(s, fmt.Sprintf("name=%s", *t.Name))
	}
	if t.Encoding != nil {
		s = append(s, fmt.Sprintf("encoding=%s", *t.Encoding))
	}
	if t.DefaultValue != nil {
		s = append(s, fmt.Sprintf("defaultValue=%s", *t.DefaultValue))
	}
	if t.Required {
		s = append(s, "required")
	} else {
		s = append(s, "optional")
	}

	return strings.Join(s, ", ")
}

func parseConfigTag(tag string) *configTag {
	t := &configTag{}

	if len(tag) == 0 {
		return t
	}

	parts := strings.Split(tag, ",")

	t.Name = &parts[0]

	for _, part := range parts[1:] {
		elemParts := strings.Split(part, "=")
		switch elemParts[0] {
		case "default":
			t.DefaultValue = &elemParts[1]

		case "required":
			t.Required = true

		case "encoding":
			switch elemParts[1] {
			case "base64", "hex":
				t.Encoding = &elemParts[1]

			default:
				panic("config: unsupported encoding in " + tag)
			}
		}
	}

	return t
}
