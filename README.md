# go-saas
go framework for saas(multi-tenancy). `go-saas` targets to provide saas solution for go

# Overview

* Different database architecture
  * [x] Single-tenancy:  Each database stores data from only one tenant.
  * [x] Multi-tenancy:  Each database stores data from multiple separate tenants (with mechanisms to protect data privacy).
  * [x] Hybrid tenancy models are also available.
* Domain driven design (DDD)
* Support multiple web framework
    * [x] [gin](https://github.com/gin-gonic/gin)
    * [x] [mux](https://github.com/gorilla/mux)
* Support multiple orms
    * [x] [gorm](https://github.com/go-gorm/gorm)
* Customizable tenant resolver
    * [x] Query String
    * [x] Form parameters
    * [x] Header
    * [x] Cookie
    * [ ] Route
    * [x] Domain format
    
    
# Sample Project

 * [gin-vue-admin](https://github.com/Goxiaoy/gin-vue-admin) This is a fork repo from original one, which uses `go-saas` to build `saas` pattern
 
# Documentation
 Refer to [wiki](https://github.com/Goxiaoy/go-saas/wiki)


#Referene

https://docs.microsoft.com/en-us/azure/azure-sql/database/saas-tenancy-app-design-patterns
