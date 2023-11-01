package saas

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTenantConfig(t *testing.T) {
	in := []byte(`{"id":"1","name":"1","region":"1","planKey":"","conn":{"a":"a","b":"b"}}`)
	conf := &TenantConfig{}
	err := json.Unmarshal(in, &conf)
	assert.NoError(t, err)

	s, err := json.Marshal(conf)
	assert.NoError(t, err)
	assert.Equal(t, string(in), string(s))
}
