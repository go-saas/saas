package data

type ConnStrings map[string]string

const Default = "default"

type ConnStrOption struct {
	// Conn string map
	Conn ConnStrings
}

func NewConnStrOption(cs ConnStrings) *ConnStrOption {
	return &ConnStrOption{
		Conn: cs,
	}
}

func (c ConnStrings) Default() string {
	return c[Default]
}

func (c ConnStrings) GetOrDefault(k string) string {
	ret := c[k]
	if ret == "" {
		return c.Default()
	}
	return ret
}
func (c ConnStrings) SetDefault(value string) {
	c[Default] = value
}
