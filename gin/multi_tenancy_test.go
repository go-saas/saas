package gin

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-saas/saas"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func SetUp() *gin.Engine {
	r := gin.Default()
	r.Use(MultiTenancy(saas.NewMemoryTenantStore(
		[]saas.TenantConfig{
			{ID: "1", Name: "Test1"},
			{ID: "2", Name: "Test3"},
		})))
	r.GET("/", func(c *gin.Context) {
		rCtx := c.Request.Context()
		tenantInfo, _ := saas.FromCurrentTenant(rCtx)
		trR := saas.FromTenantResolveRes(rCtx)
		c.JSON(200, gin.H{
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
