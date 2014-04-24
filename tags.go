package jsarray

import (
	"strings"
)

type tagOptions struct {
	quote  bool
	ignore bool
}

func parseTag(tag string) (o tagOptions) {
	if tag == "-" {
		o.ignore = true
		return
	}

	if tag == "" {
		return
	}

	spl := strings.Split(tag, ",")

	for _, v := range spl[1:] {
		switch v {
		case "string":
			o.quote = true
		}
	}
	return
}
