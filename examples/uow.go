package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/goxiaoy/uow"
)

func Uow(um uow.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" || c.Request.Method == "PATCH" {
			um.WithNew(c, func(ctx context.Context) error {
				c.Request = c.Request.WithContext(ctx)
				c.Next()
				if len(c.Errors) > 0 {
					//return err to fallback
					return fmt.Errorf("UOW: %v", c.Errors.Last())
				}
				return nil
			})
		}
		c.Next()
	}
}
