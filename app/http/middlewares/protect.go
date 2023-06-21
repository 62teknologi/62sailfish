package middlewares

//func ProtectMiddleware(cfg utils.Config) gin.HandlerFunc {
//	return func(c *gin.Context) {
//		authorizationHeader := c.GetHeader("Authorization")
//		authorizationFields := strings.Fields(authorizationHeader)
//		if len(authorizationFields) == 0 {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "Unauthorized", nil))
//			return
//		}
//
//		auth, err := utils.Decode(authorizationFields[1])
//		if err != nil {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "Error when decode token", nil))
//			return
//		}
//		if !strings.Contains(auth, cfg.ApiKey) {
//			c.AbortWithStatusJSON(http.StatusUnauthorized, utils.ResponseData("error", "Invalid API Key", nil))
//			return
//		}
//
//		c.Next()
//	}
//}
