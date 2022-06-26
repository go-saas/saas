package http

import (
	"encoding/json"
	"errors"
	"github.com/go-saas/saas"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func SetUp() *mux.Router {
	r := mux.NewRouter()

	r.Use(Middleware(saas.NewMemoryTenantStore(
		[]saas.TenantConfig{
			{ID: "1", Name: "Test1"},
			{ID: "2", Name: "Test3"},
		})))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		tenantInfo, _ := saas.FromCurrentTenant(r.Context())
		trR := saas.FromTenantResolveRes(r.Context())
		json.NewEncoder(w).Encode(map[string]interface{}{
			"tenantId":  tenantInfo.GetId(),
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
	assert.Nil(t, err)
}

func TestTerminate(t *testing.T) {
	r := mux.NewRouter()

	r.Use(Middleware(saas.NewMemoryTenantStore(
		[]saas.TenantConfig{
			{ID: "1", Name: "Test1"},
			{ID: "2", Name: "Test3"},
		}),
		WithErrorFormatter(func(w http.ResponseWriter, err error) {
			if err == ErrForbidden {
				http.Error(w, "Forbidden", 403)
			}
		}),
		WithResolveOption(saas.AppendContribs(&TerminateContrib{}))))

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	})

	req, _ := http.NewRequest("GET", "/", nil)

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusForbidden, w.Code)
}

var (
	ErrForbidden = errors.New("forbidden")
)

type TerminateContrib struct {
}

func (t *TerminateContrib) Name() string {
	return "Terminate"
}

func (t TerminateContrib) Resolve(_ *saas.Context) error {
	return ErrForbidden
}
