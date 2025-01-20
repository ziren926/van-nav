// middleware/error.go
func ErrorHandler() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()

        // 检查是否有错误
        if len(c.Errors) > 0 {
            c.JSON(http.StatusInternalServerError, gin.H{
                "success": false,
                "errorMessage": c.Errors.Last().Error(),
            })
            return
        }
    }
}

// main.go 中添加
router.Use(middleware.ErrorHandler())