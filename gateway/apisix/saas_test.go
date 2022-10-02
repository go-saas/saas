package apisix

import (
	"io/ioutil"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	testIn = []byte(`{"tenant_key":"tenant","next_header":"next_tenant","path_regex":"http://(.*).test.com"}`)
)

func TestSaas(t *testing.T) {
	in := testIn
	saas := &Saas{}
	conf, err := saas.ParseConf(in)
	assert.Nil(t, err)

	w := httptest.NewRecorder()
	saas.RequestFilter(conf, w, nil)
	resp := w.Result()
	body, _ := ioutil.ReadAll(resp.Body)
	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, "", string(body))
}

func TestSaas_BadConf(t *testing.T) {
	in := []byte(``)
	saas := &Saas{}
	_, err := saas.ParseConf(in)
	assert.NotNil(t, err)
}

//
//func TestSaasHeader(t *testing.T) {
//	InitTenantStore(common.NewMemoryTenantStore([]common.TenantConfig{
//		{ID: "1", Name: "Test1"},
//		{ID: "2", Name: "Test2", Conn: map[string]string{
//			data.Default: ":memory:?cache=shared",
//		}},
//	}))
//	in := test_in
//	saas := &Saas{}
//	conf, err := saas.ParseConf(in)
//	assert.Nil(t, err)
//
//	w := httptest.NewRecorder()
//	//TODO
//	saas.Filter(conf, w, nil)
//	//_ := w.Result()
//}
//
//func TestSaasArgs(t *testing.T) {
//	InitTenantStore(common.NewMemoryTenantStore([]common.TenantConfig{
//		{ID: "1", Name: "Test1"},
//		{ID: "2", Name: "Test2", Conn: map[string]string{
//			data.Default: ":memory:?cache=shared",
//		}},
//	}))
//	in := test_in
//	saas := &Saas{}
//	conf, err := saas.ParseConf(in)
//	assert.Nil(t, err)
//
//	w := httptest.NewRecorder()
//	//TODO
//	saas.Filter(conf, w, nil)
//	//_ := w.Result()
//}
//
//func TestSaasPath(t *testing.T) {
//	InitTenantStore(common.NewMemoryTenantStore([]common.TenantConfig{
//		{ID: "1", Name: "Test1"},
//		{ID: "2", Name: "Test2", Conn: map[string]string{
//			data.Default: ":memory:?cache=shared",
//		}},
//	}))
//	in := test_in
//	saas := &Saas{}
//	conf, err := saas.ParseConf(in)
//	assert.Nil(t, err)
//
//	w := httptest.NewRecorder()
//	//TODO
//	saas.Filter(conf, w, nil)
//	//_ := w.Result()
//}
