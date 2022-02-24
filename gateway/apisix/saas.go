package apisix

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/goxiaoy/go-saas/common"
	shttp "github.com/goxiaoy/go-saas/common/http"
	"net/http"

	pkgHTTP "github.com/apache/apisix-go-plugin-runner/pkg/http"
	"github.com/apache/apisix-go-plugin-runner/pkg/log"
	"github.com/apache/apisix-go-plugin-runner/pkg/plugin"
)

func init() {
	err := plugin.RegisterPlugin(&Saas{})
	if err != nil {
		log.Fatalf("failed to register plugin go-saas: %s", err)
	}
}

//Saas resolve and validate tenant information
type Saas struct {
}

type SaasConf struct {
	TenantKey          string `json:"tenant_key"`
	NextHeader         string `json:"next_header"`
	NextInfoHeader     string `json:"next_info_header"`
	PathRegex          string `json:"path_regex"`
	TenantNotFoundBody string `json:"tenant_not_found_body"`
}

//global variable to store tenants
var tenantStore common.TenantStore

type FormatError func(err error) (code int, body string)

var errFormat FormatError = func(err error) (code int, body string) {
	if errors.Is(err, common.ErrTenantNotFound) {
		return 404, ""
	}
	return 500, ""
}

func Init(t common.TenantStore, format FormatError) {
	tenantStore = t
	errFormat = format
}

func (p *Saas) Name() string {
	return "go-saas"
}

func (p *Saas) ParseConf(in []byte) (interface{}, error) {
	conf := SaasConf{}
	err := json.Unmarshal(in, &conf)
	return conf, err
}

func (p *Saas) Filter(conf interface{}, w http.ResponseWriter, r pkgHTTP.Request) {
	cfg := conf.(SaasConf)
	if tenantStore == nil {
		log.Warnf("fail to find tenant store. please call InitTenantStore first")
		return
	}
	key := shttp.KeyOrDefault(cfg.TenantKey)
	nextHeader := cfg.NextHeader
	if len(nextHeader) == 0 {
		nextHeader = key
	}
	//get tenant config
	tenantConfigProvider := common.NewDefaultTenantConfigProvider(NewResolver(r, key, cfg.PathRegex), tenantStore)
	tenantConfig, newCtx, err := tenantConfigProvider.Get(context.Background(), true)
	resolveValue := common.FromTenantResolveRes(newCtx)
	idOrName := ""
	if resolveValue != nil {
		idOrName = resolveValue.TenantIdOrName
	}
	if err != nil {
		code, body := errFormat(err)
		w.WriteHeader(code)
		if len(body) > 0 {
			w.Write([]byte(body))
		}
	}
	log.Infof("resolve tenant: %s ,is host: %v", idOrName, len(idOrName) == 0)
	w.Header().Set(nextHeader, tenantConfig.ID)
	nextInfoHeader := InfoHeaderOrDefault(cfg.NextInfoHeader)
	b, _ := json.Marshal(tenantConfig)
	w.Header().Set(nextInfoHeader, base64.StdEncoding.EncodeToString(b))
	return
}

func InfoHeaderOrDefault(h string) string {
	if len(h) == 0 {
		return "X-TENANT-INFO"
	} else {
		return h
	}
}
