package data

import "context"

type ConnStrings map[string]string

const Default = "default"

func (c ConnStrings) Default() string {
	return c[Default]
}

func (c ConnStrings) Resolve(_ context.Context, key string) (string, error) {
	s := c.getOrDefault(key)
	return s, nil
}

func (c ConnStrings) getOrDefault(k string) string {
	if len(k) == 0 {
		return c.Default()
	}
	ret := c[k]
	if ret == "" {
		return c.Default()
	}
	return ret
}

func (c ConnStrings) SetDefault(value string) {
	c[Default] = value
}
