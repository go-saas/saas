# go-saas

[English](./README.md) | [中文文档](./README_zh_Hans.md)

headless go framework for saas(multi-tenancy). `go-saas` targets to provide saas solution for go
this project suits for simple (web) project, which is also called monolithic. 
if you are finding microservice solution, please refer to [go-saas-kit](https://github.com/Goxiaoy/go-saas-kit)

# Overview

## Feature

* Different database architecture
  * [x] Single-tenancy:  Each database stores data from only one tenant.
  
  ![img.png](docs/mode1.png)

  * [x] Multi-tenancy:  Each database stores data from multiple separate tenants (with mechanisms to protect data privacy).
  
  ![img.png](docs/mode2.png)

  * [x] Hybrid tenancy models are also available.

  * [x] Implement your own resolver to achieve style like sharding


* Support multiple web framework
    * [x] [gin](https://github.com/gin-gonic/gin)
    * [x] net/http
    * [x] [kratos](https://github.com/go-kratos/kratos)
* Supported orm, which means all underlying database
    * [x] [gorm](https://github.com/go-gorm/gorm)
* Customizable tenant resolver
    * [x] Query String
    * [x] Form parameters
    * [x] Header
    * [x] Cookie
    * [x] Domain format
* Integration with gateway
  * [x] [apisix](https://github.com/apache/apisix)


## Install

```
go get github.com/goxiaoy/go-saas
```

## Design
```mermaid
graph TD
    A(InComming Request) -->|cookie,domain,form,header,query...|B(TenantResolver)
    B --> C(Tenant Context)  --> D(ConnectionString Resolver)
    D --> E(Tenant 1) --> J(Data Filter) -->  H(Shared Database)
    D --> F(Tenant 2) --> J
    D --> G(Tenant 3) --> I(Tenant 3 Database)
```
 
    
# Sample Project
* [example](https://github.com/Goxiaoy/go-saas/tree/main/examples) combination of `go-saas`,`gin`,`gorm(sqlite)`
* [go-saas-kit](https://github.com/Goxiaoy/go-saas-kit) Microservice architecture starter kit for golang sass project

# Documentation
 Refer to [wiki](https://github.com/Goxiaoy/go-saas/wiki)


# References

https://docs.microsoft.com/en-us/azure/azure-sql/database/saas-tenancy-app-design-patterns
