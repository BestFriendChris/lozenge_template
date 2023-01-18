package go_format

import "mvdan.cc/gofumpt/format"

func Format(s string) (string, error) {
	var opts = format.Options{}
	formatted, err := format.Source([]byte(s), opts)
	if err != nil {
		return "", err
	}

	return string(formatted), nil
}
