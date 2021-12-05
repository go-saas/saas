package saas

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/goxiaoy/go-saas/common"
	http2 "github.com/goxiaoy/go-saas/common/http"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func SetUp() *mux.Router {
	r := mux.NewRouter()
	wOpt := http2.NewDefaultWebMultiTenancyOption()
	m := MultiTenancy{
		wOpt,
		nil,
		common.NewMemoryTenantStore(
			[]common.TenantConfig{
				{ID: "1", Name: "Test1"},
				{ID: "2", Name: "Test3"},
			}),
	}

	r.Use(m.Middleware)

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		currentTenant := common.ContextCurrentTenant{}
		trR := common.FromTenantResolveRes(r.Context())
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tenantId":  currentTenant.Id(r.Context()),
			"resolvers": trR.AppliedResolvers,
		})
	})
	return r
}

func getW(url string, f func(r *http.Request)) *httptest.ResponseRecorder {
	r := SetUp()
	req, _ := http.NewRequest("GET", url, nil)
	f(req)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestHostMultiTenancy(t *testing.T) {
	w := getW("/", func(r *http.Request) {
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "", value)
	assert.Nil(t, err)
}
func TestNotFoundMultiTenancy(t *testing.T) {
	w := getW("/", func(r *http.Request) {
		r.Header.Set("__tenant", "1000")
	})
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestCookieMultiTenancy(t *testing.T) {
	w := getW("/", func(r *http.Request) {
		r.AddCookie(&http.Cookie{
			Name:       "__tenant",
			Value:      "1",
			Path:       "",
			Domain:     "",
			Expires:    time.Time{},
			RawExpires: "",
			MaxAge:     0,
			Secure:     false,
			HttpOnly:   false,
			SameSite:   0,
			Raw:        "",
			Unparsed:   nil,
		})
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "1", value)
	r := response["resolvers"].([]interface{})
	assert.Equal(t, "Cookie", r[len(r)-1])
	assert.Nil(t, err)
}

func TestHeaderMultiTenancy(t *testing.T) {
	w := getW("/", func(r *http.Request) {
		r.Header.Set("__tenant", "1")
	})
	assert.Equal(t, http.StatusOK, w.Code)
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "1", value)
	r := response["resolvers"].([]interface{})
	assert.Equal(t, "Header", r[len(r)-1])
	assert.Nil(t, err)
}
