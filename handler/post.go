package handler

import (
    "net/http"
    "time"

    "github.com/gin-gonic/gin"
    "github.com/mereith/nav/database"
)

// GetPostHandler 获取工具的帖子内容
func GetPostHandler(c *gin.Context) {
    id := c.Param("id")
    var post struct {
        PostTitle     string    `json:"post_title"`
        PostContent   string    `json:"post_content"`
        PostCreatedAt time.Time `json:"post_created_at"`
        PostUpdatedAt time.Time `json:"post_updated_at"`
    }

    err := database.DB.QueryRow(`
        SELECT post_title, post_content, post_created_at, post_updated_at
        FROM nav_table WHERE id = ?`, id).Scan(
        &post.PostTitle, &post.PostContent, &post.PostCreatedAt, &post.PostUpdatedAt)

    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
        return
    }

    c.JSON(http.StatusOK, post)
}

// UpdatePostHandler 更新工具的帖子内容
func UpdatePostHandler(c *gin.Context) {
    id := c.Param("id")

    var input struct {
        PostTitle   string `json:"post_title"`
        PostContent string `json:"post_content"`
    }

    if err := c.BindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    _, err := database.DB.Exec(`
        UPDATE nav_table
        SET post_title = ?, post_content = ?, post_updated_at = CURRENT_TIMESTAMP
        WHERE id = ?`, input.PostTitle, input.PostContent, id)

    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "帖子更新成功"})
}