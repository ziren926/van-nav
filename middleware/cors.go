// middleware/cors.go
package middleware

import (
    "github.com/gin-gonic/gin"
    "github.com/ziren926/van-nav/logger" // 添加 logger 导入
)

// CORS 处理跨域请求的中间件
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 设置允许的来源
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

        // 允许携带凭证
        c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

        // 允许的请求方法
        c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")

        // 允许的请求头
        c.Writer.Header().Set("Access-Control-Allow-Headers", ""+
            "Authorization, "+
            "Content-Type, "+
            "Content-Length, "+
            "Accept-Encoding, "+
            "X-CSRF-Token, "+
            "X-Requested-With, "+
            "Accept, "+
            "Origin, "+
            "Cache-Control")

        // 处理预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        // 记录跨域请求
        origin := c.Request.Header.Get("Origin")
        if origin != "" {
            logger.LogInfo("CORS request from: %s, Path: %s", origin, c.Request.URL.Path)
        }

        c.Next()
    }
}