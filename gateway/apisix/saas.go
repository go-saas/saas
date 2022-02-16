package apisix

import (
	"context"
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
		log.Fatalf("failed to register plugin say: %s", err)
	}
}

//Saas resolve and validate tenant information
type Saas struct {
}

type SaasConf struct {
	TenantKey          string `json:"tenant_key"`
	NextHeader         string `json:"next_header"`
	PathRegex          string `json:"path_regex"`
	TenantNotFoundBody string `json:"tenant_not_found_body"`
}

//global variable to store tenants
var tenantStore common.TenantStore

func InitTenantStore(t common.TenantStore) {
	tenantStore = t
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
	tenantConfig, _, err := tenantConfigProvider.Get(context.Background(), true)
	if err != nil {
		//not found
		if errors.Is(err, common.ErrTenantNotFound) {
			w.WriteHeader(404)
			if len(cfg.TenantNotFoundBody) > 0 {
				w.Write([]byte(cfg.TenantNotFoundBody))
			}
		} else {
			log.Fatalf("%s", err)
			w.WriteHeader(500)
		}
	}
	w.Header().Set(nextHeader, tenantConfig.ID)
	return
}
