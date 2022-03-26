package apisix

import (
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
	TenantKey      string `json:"tenant_key"`
	NextHeader     string `json:"next_header"`
	NextInfoHeader string `json:"next_info_header"`
	PathRegex      string `json:"path_regex"`
}

type FormatError func(err error, w http.ResponseWriter)

//global variable to store tenants
var (
	tenantStore          common.TenantStore
	nextTenantHeader     string
	nextTenantInfoHeader string
)

var errFormat FormatError = func(err error, w http.ResponseWriter) {
	if errors.Is(err, common.ErrTenantNotFound) {
		w.WriteHeader(404)
	}
	w.WriteHeader(500)
}

func Init(t common.TenantStore, nextHeader, nextInfoHeader string, format FormatError) {
	tenantStore = t
	errFormat = format
	nextTenantHeader = nextHeader
	nextTenantInfoHeader = nextInfoHeader
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
		nextHeader = nextTenantHeader
	}
	if len(nextHeader) == 0 {
		nextHeader = key
	}
	ctx := r.Context()
	//get tenant config
	tenantConfigProvider := common.NewDefaultTenantConfigProvider(NewResolver(r, key, cfg.PathRegex), tenantStore)
	tenantConfig, ctx, err := tenantConfigProvider.Get(ctx)
	if err != nil {
		errFormat(err, w)
		return
	}
	resolveValue := common.FromTenantResolveRes(ctx)
	idOrName := ""
	if resolveValue != nil {
		idOrName = resolveValue.TenantIdOrName
	}
	log.Infof("resolve tenant: %s ,id: %s ,is host: %v", idOrName, tenantConfig.ID, len(tenantConfig.ID) == 0)
	r.Header().Set(nextHeader, tenantConfig.ID)
	nextInfoHeader := cfg.NextInfoHeader
	if len(nextInfoHeader) == 0 {
		nextInfoHeader = nextTenantInfoHeader
	}
	nextInfoHeader = InfoHeaderOrDefault(nextInfoHeader)
	b, _ := json.Marshal(tenantConfig)
	r.Header().Set(nextInfoHeader, base64.StdEncoding.EncodeToString(b))
	return
}

func InfoHeaderOrDefault(h string) string {
	if len(h) == 0 {
		return "X-TENANT-INFO"
	} else {
		return h
	}
}
