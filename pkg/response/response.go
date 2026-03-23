package response

import "github.com/gin-gonic/gin"

// JSON wraps successful payloads in a stable response envelope.
func JSON(c *gin.Context, status int, data any) {
	c.JSON(status, gin.H{
		"data": data,
	})
}

// Error wraps error messages in a stable response envelope.
func Error(c *gin.Context, status int, message string) {
	c.JSON(status, gin.H{
		"error": message,
	})
}
