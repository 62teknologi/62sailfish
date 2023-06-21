package middlewares

import (
	"fmt"
	"net/http"

	"github.com/62teknologi/62sailfish/62golib/utils"
	"github.com/gin-gonic/gin"
)

func IPWhitelist(allowedIPs []string) gin.HandlerFunc {
	allowedIPSet := make(map[string]struct{})
	for _, ip := range allowedIPs {
		allowedIPSet[ip] = struct{}{}
	}

	return func(c *gin.Context) {
		clientIP := c.ClientIP()
		fmt.Println("clientIP", clientIP)
		if _, ok := allowedIPSet[clientIP]; !ok {
			c.JSON(http.StatusForbidden, utils.ResponseData("error", "access denied", nil))
			c.Abort()
			return
		}
		c.Next()
	}
}
