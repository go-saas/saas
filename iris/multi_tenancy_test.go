package iris

import (
	"github.com/go-saas/saas"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/httptest"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func SetUp() *iris.Application {
	r := iris.New()
	r.Use(MultiTenancy(saas.NewMemoryTenantStore(
		[]saas.TenantConfig{
			{ID: "1", Name: "Test1"},
			{ID: "2", Name: "Test3"},
		})))
	r.Get("/", func(c iris.Context) {
		rCtx := c.Request().Context()
		tenantInfo, _ := saas.FromCurrentTenant(rCtx)
		trR := saas.FromTenantResolveRes(rCtx)
		c.JSON(iris.Map{
			"tenantId":  tenantInfo.GetId(),
			"resolvers": trR.AppliedResolvers,
		})
	})
	return r
}

func TestHostMultiTenancy(t *testing.T) {
	e := httptest.New(t, SetUp())
	t1 := e.GET("/").Expect().Status(http.StatusOK)

	var response map[string]interface{}
	t1.JSON().Decode(&response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "", value)

}
func TestNotFoundMultiTenancy(t *testing.T) {

	e := httptest.New(t, SetUp())
	e.GET("/").WithHeader("__tenant", "1000").Expect().Status(http.StatusNotFound)
}

func TestCookieMultiTenancy(t *testing.T) {

	e := httptest.New(t, SetUp())
	t1 := e.GET("/").WithCookie("__tenant", "1").Expect().Status(http.StatusOK)

	var response map[string]interface{}
	t1.JSON().Decode(&response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "1", value)

}

func TestHeaderMultiTenancy(t *testing.T) {

	e := httptest.New(t, SetUp())
	t1 := e.GET("/").WithHeader("__tenant", "1").Expect().Status(http.StatusOK)

	var response map[string]interface{}
	t1.JSON().Decode(&response)
	value, exists := response["tenantId"]
	assert.True(t, exists)
	assert.Equal(t, "1", value)

}
