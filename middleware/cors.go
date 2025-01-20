// middleware/cors.go

package middleware

import (
    "time"

    "github.com/gin-gonic/gin"
)

// CORS 处理跨域请求的中间件
// 支持：
// - 允许所有来源
// - 允许常用 HTTP 方法
// - 允许常用请求头
// - 允许携带凭证
// - 预检请求缓存时间设置
// - OPTIONS 请求处理
func CORS() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 设置允许的来源
        c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

        // 允许携带凭证（cookies等）
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

        // 允许暴露的响应头
        c.Writer.Header().Set("Access-Control-Expose-Headers", ""+
            "Content-Length, "+
            "Access-Control-Allow-Origin, "+
            "Access-Control-Allow-Headers, "+
            "Cache-Control, "+
            "Content-Language, "+
            "Content-Type")

        // 预检请求缓存时间（秒）
        c.Writer.Header().Set("Access-Control-Max-Age", "3600")

        // 设置常用安全响应头
        c.Writer.Header().Set("X-Content-Type-Options", "nosniff")
        c.Writer.Header().Set("X-Frame-Options", "DENY")
        c.Writer.Header().Set("X-XSS-Protection", "1; mode=block")

        // 处理预检请求（OPTIONS）
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204) // 返回 204 No Content
            return
        }

        // 记录跨域请求日志
        if origin := c.Request.Header.Get("Origin"); origin != "" {
            logger.LogInfo("收到跨域请求，来源: %s, 路径: %s, 方法: %s",
                origin,
                c.Request.URL.Path,
                c.Request.Method)
        }

        // 继续处理请求
        c.Next()
    }
}

// CustomCORS 自定义CORS配置的中间件
// 参数：
// - allowOrigins: 允许的来源列表
// - allowMethods: 允许的HTTP方法列表
// - maxAge: 预检请求缓存时间（秒）
// - allowCredentials: 是否允许携带凭证
func CustomCORS(allowOrigins []string, allowMethods []string, maxAge int, allowCredentials bool) gin.HandlerFunc {
    // 将来源列表转换为 map，便于快速查找
    originsMap := make(map[string]bool)
    for _, origin := range allowOrigins {
        originsMap[origin] = true
    }

    // 将方法列表转换为字符串
    methodsStr := ""
    if len(allowMethods) > 0 {
        methodsStr = allowMethods[0]
        for _, method := range allowMethods[1:] {
            methodsStr += ", " + method
        }
    }

    return func(c *gin.Context) {
        origin := c.Request.Header.Get("Origin")

        // 检查来源是否允许
        if origin != "" {
            if len(originsMap) == 0 || originsMap[origin] {
                c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
            } else {
                // 来源不在允许列表中
                c.AbortWithStatus(403)
                return
            }
        }

        // 设置允许的方法
        if methodsStr != "" {
            c.Writer.Header().Set("Access-Control-Allow-Methods", methodsStr)
        }

        // 设置是否允许携带凭证
        if allowCredentials {
            c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
        }

        // 设置预检请求缓存时间
        if maxAge > 0 {
            c.Writer.Header().Set("Access-Control-Max-Age", string(maxAge))
        }

        // 设置允许的请求头
        c.Writer.Header().Set("Access-Control-Allow-Headers", ""+
            "Authorization, "+
            "Content-Type, "+
            "Content-Length, "+
            "Accept-Encoding, "+
            "X-CSRF-Token, "+
            "X-Requested-With")

        // 处理预检请求
        if c.Request.Method == "OPTIONS" {
            c.AbortWithStatus(204)
            return
        }

        c.Next()
    }
}

// 使用示例：
/*
func main() {
    router := gin.Default()

    // 使用默认CORS配置
    router.Use(CORS())

    // 或者使用自定义CORS配置
    router.Use(CustomCORS(
        []string{"http://localhost:3000", "https://your-domain.com"},
        []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        3600,
        true,
    ))

    // ... 其他路由配置
}
*/