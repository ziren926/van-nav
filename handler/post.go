package handler

import (
    "net/http"
    "time"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/ziren926/van-nav/types"
    "github.com/ziren926/van-nav/logger"
    "github.com/ziren926/van-nav/service"
    "github.com/ziren926/van-nav/utils"
)

// GetPostHandler 获取工具的帖子内容
func GetPostHandler(c *gin.Context) {
    id := c.Param("id")
    logger.LogInfo("获取工具帖子内容，ID: %s", id)

    numberId, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        logger.LogError("无效的ID格式: %s, 错误: %v", id, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的ID格式",
        })
        return
    }

    post, err := service.GetPost(numberId)
    if err != nil {
        logger.LogError("获取帖子失败, ID: %d, 错误: %v", numberId, err)
        if err.Error() == "工具不存在" {
            c.JSON(http.StatusNotFound, gin.H{
                "success": false,
                "errorMessage": "帖子不存在",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "errorMessage": "获取帖子失败",
        })
        return
    }

    // 检查用户权限
    tool, err := service.GetToolById(numberId)
    if err == nil && tool.Hide && !utils.IsLogin(c) {
        c.JSON(http.StatusForbidden, gin.H{
            "success": false,
            "errorMessage": "无权访问该帖子",
        })
        return
    }

    logger.LogInfo("成功获取工具帖子，ID: %d", numberId)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": post,
    })
}

// UpdatePostHandler 更新工具的帖子内容
// handler/post.go

func UpdatePostHandler(c *gin.Context) {
    // 从 URL 获取工具 ID
    id, err := strconv.ParseInt(c.Param("id"), 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的工具ID",
        })
        return
    }

    // 解析请求体
    var post types.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的请求数据",
        })
        return
    }

    // 更新帖子
    err = service.UpdatePost(id, &post)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "errorMessage": err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "帖子更新成功",
        "data": gin.H{
            "updated_at": time.Date(2025, 1, 20, 2, 54, 4, 0, time.UTC).Format("2006-01-02 15:04:05"),
            "updated_by": "ziren926",
        },
    })
}

// GetPostWithContentHandler 获取工具的所有内容（包括帖子和基本信息）
func GetPostWithContentHandler(c *gin.Context) {
    id := c.Param("id")
    logger.LogInfo("获取工具完整内容，ID: %s", id)

    numberId, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        logger.LogError("无效的ID格式: %s, 错误: %v", id, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的ID格式",
        })
        return
    }

    tool, err := service.GetToolById(numberId)
    if err != nil {
        logger.LogError("获取工具失败, ID: %d, 错误: %v", numberId, err)
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "errorMessage": "工具不存在",
        })
        return
    }

    // 检查用户权限
    if tool.Hide && !utils.IsLogin(c) {
        c.JSON(http.StatusForbidden, gin.H{
            "success": false,
            "errorMessage": "无权访问该内容",
        })
        return
    }

    logger.LogInfo("成功获取工具完整内容，ID: %d", numberId)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data": tool,
    })
}