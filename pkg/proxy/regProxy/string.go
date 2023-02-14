package regProxy

import "strings"

type Strings []string

// HeadT ...
func (ss Strings) HeadT() (string, Strings) {
	if len(ss) > 0 {
		return ss[0], Strings(ss[1:])
	}

	return "", Strings{}
}

// Head ...
func (ss Strings) Head() string {
	if len(ss) > 0 {
		return ss[0]
	}
	return ""
}

// Head2 ...
func (ss Strings) Head2() (h0, h1 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	return
}

// Head3 ...
func (ss Strings) Head3() (h0, h1, h2 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	if len(ss) > 2 {
		h2 = ss[2]
	}
	return
}

// Head4 ...
func (ss Strings) Head4() (h0, h1, h2, h3 string) {
	if len(ss) > 0 {
		h0 = ss[0]
	}
	if len(ss) > 1 {
		h1 = ss[1]
	}
	if len(ss) > 2 {
		h2 = ss[2]
	}
	if len(ss) > 3 {
		h3 = ss[3]
	}
	return
}

// Split ...
func Split(raw string, sep string) Strings {
	return Strings(strings.Split(raw, sep))
}
