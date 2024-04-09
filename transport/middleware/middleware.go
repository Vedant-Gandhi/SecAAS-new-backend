package middleware

import "github.com/gin-gonic/gin"

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Define the allowed CORS domains
		allowedOrigins := []string{"http://localhost:5173", "http://localhost:5174"}

		// Get the Origin header from the request
		origin := c.GetHeader("Origin")

		// Check if the request origin is in the list of allowed origins
		// If it is, set the Access-Control-Allow-Origin header to the request origin
		// Otherwise, set it to "*"
		var allowedOrigin string
		for _, o := range allowedOrigins {
			if o == origin {
				allowedOrigin = o
				break
			}
		}
		if allowedOrigin != "" {
			c.Writer.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		} else {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
