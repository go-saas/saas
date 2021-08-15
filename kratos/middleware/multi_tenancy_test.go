package middleware
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/go-kratos/kratos/v2/transport/grpc"
//	http2 "github.com/go-kratos/kratos/v2/transport/http"
//	"github.com/goxiaoy/go-saas/common"
//	"github.com/stretchr/testify/assert"
//	"io/ioutil"
//	"net/http"
//	"strings"
//	"testing"
//	"time"
//)
//
//func TestServer(t *testing.T){
//	ctx := context.Background()
//	m:= MultiTenancy(nil, nil, common.NewMemoryTenantStore(
//		[]common.TenantConfig{
//			{ID: "1", Name: "Test1"},
//			{ID: "2", Name: "Test3"},
//		}))
//	var hOpt = []http2.ServerOption{
//		http2.Middleware(
//			m,
//		),
//	}
//	var gOpt = []grpc.ServerOption{
//		grpc.Middleware(
//			m,
//		),
//	}
//	httpServer :=http2.NewServer(hOpt...)
//
//	httpServer.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
//		currentTenant := common.ContextCurrentTenant{}
//		rCtx := request.Context()
//		t := currentTenant.Id(rCtx)
//		trR := common.FromTenantResolveRes(rCtx)
//		json.NewEncoder(writer).Encode(map[string]interface{}{
//			"tenantId":  t,
//				"resolvers": trR.AppliedResolvers,
//		})
//	})
//
//	grpcServer := grpc.NewServer(gOpt...)
//
//	if e, err := httpServer.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
//		t.Fatal(e, err)
//	}
//
//	go func() {
//		if err := httpServer.Start(ctx); err != nil {
//			panic(err)
//		}
//	}()
//
//	if e, err := grpcServer.Endpoint(); err != nil || e == nil || strings.HasSuffix(e.Host, ":0") {
//		t.Fatal(e, err)
//	}
//
//	go func() {
//		if err := grpcServer.Start(ctx); err != nil {
//			panic(err)
//		}
//	}()
//
//	time.Sleep(time.Second)
//
//	httpClient := testHttpClient(t,httpServer)
//	testHostMultiTenancy(t,httpServer,httpClient)
//	//testNotFoundMultiTenancy(t,httpClient)
//	testHeaderMultiTenancy(t,httpServer,httpClient)
//}
//
//func testHttpClient(t *testing.T,httpServer *http2.Server)*http2.Client  {
//	e, err := httpServer.Endpoint()
//	if err != nil {
//		t.Fatal(err)
//	}
//	client, err := http2.NewClient(context.Background(),http2.WithEndpoint(e.Host))
//	if err != nil {
//		t.Fatal(err)
//	}
//	return client
//}
//
//func getW(t *testing.T,url string,httpServer *http2.Server,client *http2.Client, f func(r *http.Request)) *http.Response {
//	e, err := httpServer.Endpoint()
//	if err != nil {
//		t.Fatal(err)
//	}
//	url = fmt.Sprintf(e.String() + url)
//	req, err := http.NewRequest("GET", url, nil)
//	if err != nil {
//		t.Fatal(err)
//	}
//	if f !=nil{
//		f(req)
//	}
//	resp, err := client.Do(req)
//	if err != nil {
//		t.Fatal(err)
//	}
//	return  resp
//}
//
//func testHostMultiTenancy(t *testing.T,httpServer *http2.Server,httpClient *http2.Client) {
//	resp := getW(t,"/",httpServer,httpClient,nil)
//	assert.Equal(t, http.StatusOK, resp.StatusCode)
//	content, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		t.Fatalf("read resp error %v", err)
//	}
//	var response map[string]interface{}
//	err = json.Unmarshal(content, &response)
//	value, exists := response["tenantId"]
//	assert.True(t, exists)
//	assert.Equal(t, "", value)
//	assert.Nil(t, err)
//}
//
//func testNotFoundMultiTenancy(t *testing.T,httpServer *http2.Server,httpClient *http2.Client) {
//	resp := getW(t,"/",httpServer,httpClient,nil)
//	assert.Equal(t, http.StatusNotFound, resp.Status)
//}
//
//
//func testHeaderMultiTenancy(t *testing.T,httpServer *http2.Server,httpClient *http2.Client)  {
//	w :=  getW(t,"/",httpServer,httpClient, func(r *http.Request) {
//		r.Header.Set("__tenant", "1")
//	})
//	assert.Equal(t, http.StatusOK, w.StatusCode)
//	content, err := ioutil.ReadAll(w.Body)
//	var response map[string]interface{}
//	err = json.Unmarshal(content, &response)
//	value, exists := response["tenantId"]
//	assert.True(t, exists)
//	assert.Equal(t, "1", value)
//	r := response["resolvers"].([]interface{})
//	assert.Equal(t, "Header", r[len(r)-1])
//	assert.Nil(t, err)
//}
