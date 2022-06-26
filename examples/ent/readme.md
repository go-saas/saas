# Example project

combination of `go-saas`,`gin`,`ent(sqlite)`

```shell
go run github.com/go-saas/saas/examples/ent
```
---
Host side ( use shared database):

Open `http://localhost:8090/posts`

---
Multi-tenancy ( use shared database):

Open http://localhost:8090/posts?__tenant=1

Open http://localhost:8090/posts?__tenant=2

---
Single-tenancy ( use separate database):

Open http://localhost:8090/posts?__tenant=3