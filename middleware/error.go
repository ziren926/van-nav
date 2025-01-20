// middleware/error.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "net/http"
)

func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 檢查是否有錯誤
        if len(c.Errors) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "errorMessage": c.Errors.Last().Error(),
            })
            return
        }
    }
}