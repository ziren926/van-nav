// middleware/ratelimit.go
package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
    "golang.org/x/time/rate"
)

// RateLimit 创建速率限制中间件
func RateLimit() gin.HandlerFunc {
    // 创建一个限制器：每秒允许10个请求
    limiter := rate.NewLimiter(rate.Every(time.Second), 10)

    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{
                "success": false,
                "errorMessage": "请求过于频繁，请稍后再试",
            })
            c.Abort()
            return
        }
        c.Next()
    }
}