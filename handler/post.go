package handler

import (
    "net/http"
    "time"
    "strconv"

    "github.com/gin-gonic/gin"
    "github.com/ziren926/van-nav/database"
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
func UpdatePostHandler(c *gin.Context) {
    id := c.Param("id")
    logger.LogInfo("更新工具帖子内容，ID: %s", id)

    numberId, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        logger.LogError("无效的ID格式: %s, 错误: %v", id, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的ID格式",
        })
        return
    }

    var input struct {
        PostTitle   string `json:"post_title" binding:"required"`
        PostContent string `json:"post_content" binding:"required"`
    }

    if err := c.ShouldBindJSON(&input); err != nil {
        logger.LogError("绑定请求数据失败: %v", err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success": false,
            "errorMessage": "无效的请求数据",
        })
        return
    }

    // 检查工具是否存在
    _, err = service.GetToolById(numberId)
    if err != nil {
        logger.LogError("工具不存在, ID: %d", numberId)
        c.JSON(http.StatusNotFound, gin.H{
            "success": false,
            "errorMessage": "工具不存在",
        })
        return
    }

    // 更新帖子内容
    err = service.UpdatePost(numberId, input.PostTitle, input.PostContent)
    if err != nil {
        logger.LogError("更新帖子失败, ID: %d, 错误: %v", numberId, err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "success": false,
            "errorMessage": "更新帖子失败",
        })
        return
    }

    logger.LogInfo("成功更新工具帖子，ID: %d, 标题: %s", numberId, input.PostTitle)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": "帖子更新成功",
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