// middleware/audit.go
package middleware

import (
    "time"
    "github.com/gin-gonic/gin"
    "github.com/ziren926/van-nav/logger"
)

var (
    CurrentUser = "ziren926"
    CurrentTime = time.Date(2025, 1, 20, 2, 45, 19, 0, time.UTC)
)

type AuditInfo struct {
    UserLogin  string
    Timestamp  time.Time
    RequestID  string
}

func AuditMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        audit := &AuditInfo{
            UserLogin: CurrentUser,
            Timestamp: CurrentTime,
            RequestID: c.GetString("X-Request-ID"),
        }

        c.Set("audit", audit)

        // 记录请求开始
        logger.LogInfo("Request started - User: %s, Time: %s, Path: %s",
            audit.UserLogin,
            audit.Timestamp.Format("2006-01-02 15:04:05"),
            c.Request.URL.Path)

        c.Next()

        // 记录请求结束
        logger.LogInfo("Request completed - User: %s, Path: %s",
            audit.UserLogin,
            c.Request.URL.Path)
    }
}