// handler/post.go
package handler

import (
    "github.com/gin-gonic/gin"
    "github.com/ziren926/van-nav/database"
    "github.com/ziren926/van-nav/types"
    "net/http"
    "strconv"
    "time"
)

// AddPostHandler 添加新帖子
func AddPostHandler(c *gin.Context) {
    var post types.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 使用事务确保数据一致性
    tx, err := database.DB.Begin()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 准备SQL语句
    stmt, err := tx.Prepare(`
        INSERT INTO posts (title, content, create_time, update_time)
        VALUES (?, ?, datetime('now'), datetime('now'))
    `)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer stmt.Close()

    // 执行SQL语句
    result, err := stmt.Exec(post.Title, post.Content)
    if err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 提交事务
    if err := tx.Commit(); err != nil {
        tx.Rollback()
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // 获取新插入的ID
    id, _ := result.LastInsertId()
    post.ID = id
    post.CreateTime = time.Now()
    post.UpdateTime = time.Now()

    c.JSON(http.StatusOK, post)
}

// DeletePostHandler 删除帖子
func DeletePostHandler(c *gin.Context) {
    id := c.Param("id")
    postID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
        return
    }

    result, err := database.DB.Exec("DELETE FROM posts WHERE id = ?", postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdatePostHandler 更新帖子
func UpdatePostHandler(c *gin.Context) {
    id := c.Param("id")
    postID, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "无效的ID"})
        return
    }

    var post types.Post
    if err := c.ShouldBindJSON(&post); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    result, err := database.DB.Exec(`
        UPDATE posts
        SET title = ?, content = ?, update_time = datetime('now')
        WHERE id = ?
    `, post.Title, post.Content, postID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "帖子不存在"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// GetPostsHandler 获取所有帖子
func GetPostsHandler(c *gin.Context) {
    posts := []types.Post{}
    rows, err := database.DB.Query(`
        SELECT id, title, content, create_time, update_time
        FROM posts
        ORDER BY create_time DESC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer rows.Close()

    for rows.Next() {
        var post types.Post
        err := rows.Scan(&post.ID, &post.Title, &post.Content, &post.CreateTime, &post.UpdateTime)
        if err != nil {
            continue
        }
        posts = append(posts, post)
    }

    c.JSON(http.StatusOK, posts)
}