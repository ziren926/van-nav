package service

import (
    "fmt"
    "time"

    "github.com/ziren926/van-nav/database"
    "github.com/ziren926/van-nav/logger"
    "github.com/ziren926/van-nav/types"
    "github.com/ziren926/van-nav/utils"
)

func ImportTools(data []types.Tool) {
    var catelogs []string
    for _, v := range data {
        if !utils.In(v.Catelog, catelogs) {
            catelogs = append(catelogs, v.Catelog)
        }
        sql_add_tool := `
            INSERT INTO nav_table (id, name, catelog, url, logo, desc)
            VALUES (?, ?, ?, ?, ?, ?);
            `
        stmt, err := database.DB.Prepare(sql_add_tool)
        utils.CheckErr(err)
        res, err := stmt.Exec(v.Id, v.Name, v.Catelog, v.Url, v.Logo, v.Desc)
        utils.CheckErr(err)
        _, err = res.LastInsertId()
        utils.CheckErr(err)
    }
    for _, catelog := range catelogs {
        var addCatelogDto types.AddCatelogDto
        addCatelogDto.Name = catelog
        AddCatelog(addCatelogDto)
    }
    // 转存所有图片,异步
    go func(data []types.Tool) {
        for _, v := range data {
            UpdateImg(v.Logo)
        }
    }(data)
}

func UpdateTool(data types.UpdateToolDto) error {
    // 更新所有工具字段，包括帖子相关字段
    sql_update_tool := `
        UPDATE nav_table
        SET name = ?, url = ?, logo = ?, catelog = ?, desc = ?,
            sort = ?, hide = ?, content = ?,
            post_title = ?, post_content = ?,
            post_updated_at = CURRENT_TIMESTAMP
        WHERE id = ?;
    `
    stmt, err := database.DB.Prepare(sql_update_tool)
    if err != nil {
        return err
    }
    defer stmt.Close()

    result, err := stmt.Exec(
        data.Name, data.Url, data.Logo, data.Catelog, data.Desc,
        data.Sort, data.Hide, data.Content,
        data.PostTitle, data.PostContent,
        data.Id,
    )
    if err != nil {
        return err
    }

    // 检查是否有记录被更新
    affected, err := result.RowsAffected()
    if err != nil {
        return err
    }
    if affected == 0 {
        return fmt.Errorf("工具不存在")
    }

    // 更新图片
    if data.Logo != "" {
        UpdateImg(data.Logo)
    }

    return nil
}

func AddTool(data types.AddToolDto) {
    currentTime := time.Date(2025, 1, 20, 3, 27, 4, 0, time.UTC)
    currentUser := "ziren926"

    sql_add_tool := `
        INSERT INTO nav_table (
            name, url, logo, desc, catelog,
            sort, hide, content,
            post_title, post_content,
            post_created_at, post_updated_at
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
    `
    stmt, err := database.DB.Prepare(sql_add_tool)
    if err != nil {
        logger.LogError("准备添加工具语句失败: %v", err)
        return
    }
    defer stmt.Close()

    res, err := stmt.Exec(
        data.Name, data.Url, data.Logo, data.Desc, data.Catelog,
        data.Sort, data.Hide, data.Content,
        data.PostTitle, data.PostContent,
        currentTime, currentTime,
    )
    if err != nil {
        logger.LogError("执行添加工具失败: %v", err)
        return
    }

    id, err := res.LastInsertId()
    if err != nil {
        logger.LogError("获取插入ID失败: %v", err)
        return
    }

    logger.LogInfo("成功添加工具 ID: %d, 添加人: %s, 时间: %s",
        id, currentUser, currentTime.Format("2006-01-02 15:04:05"))
}

func GetAllTool() []types.Tool {
    sql_get_all := `
        SELECT id, name, url, logo, catelog, desc, sort, hide,
               COALESCE(content, '') as content,
               COALESCE(post_title, '') as post_title,
               COALESCE(post_content, '') as post_content,
               COALESCE(post_created_at, CURRENT_TIMESTAMP) as post_created_at,
               COALESCE(post_updated_at, CURRENT_TIMESTAMP) as post_updated_at
        FROM nav_table
        ORDER BY sort;
    `
    results := make([]types.Tool, 0)
    rows, err := database.DB.Query(sql_get_all)
    utils.CheckErr(err)
    defer rows.Close()

    for rows.Next() {
        var tool types.Tool
        var hide interface{}
        var sort interface{}
        err = rows.Scan(
            &tool.Id, &tool.Name, &tool.Url, &tool.Logo,
            &tool.Catelog, &tool.Desc, &sort, &hide,
            &tool.Content, &tool.PostTitle, &tool.PostContent,
            &tool.PostCreatedAt, &tool.PostUpdatedAt,
        )
        if hide == nil {
            tool.Hide = false
        } else {
            if hide.(int64) == 0 {
                tool.Hide = false
            } else {
                tool.Hide = true
            }
        }
        if sort == nil {
            tool.Sort = 0
        } else {
            i64 := sort.(int64)
            tool.Sort = int(i64)
        }
        utils.CheckErr(err)
        results = append(results, tool)
    }
    return results
}

func GetToolLogoUrlById(id int) string {
    sql_get_tool := `SELECT logo FROM nav_table WHERE id=?;`
    rows, err := database.DB.Query(sql_get_tool, id)
    utils.CheckErr(err)
    defer rows.Close()

    var tool types.Tool
    for rows.Next() {
        err = rows.Scan(&tool.Logo)
        utils.CheckErr(err)
    }
    return tool.Logo
}

func UpdateToolIcon(id int64, logo string) {
    sql_update_tool := `UPDATE nav_table SET logo=? WHERE id=?;`
    _, err := database.DB.Exec(sql_update_tool, logo, id)
    utils.CheckErr(err)
    UpdateImg(logo)
}

func UpdateToolsSort(updates []types.UpdateToolsSortDto) error {
    tx, err := database.DB.Begin()
    if err != nil {
        return err
    }

    sql := `UPDATE nav_table SET sort = ? WHERE id = ?`
    stmt, err := tx.Prepare(sql)
    if err != nil {
        tx.Rollback()
        return err
    }
    defer stmt.Close()

    for _, update := range updates {
        _, err = stmt.Exec(update.Sort, update.Id)
        if err != nil {
            tx.Rollback()
            return err
        }
    }

    return tx.Commit()
}

func GetToolById(id int64) (types.Tool, error) {
    logger.LogInfo("正在查询工具ID: %d", id)

    var tool types.Tool
    sql := `
        SELECT id, name, url, logo, catelog, desc,
               COALESCE(content, '') as content,
               COALESCE(sort, 0) as sort,
               COALESCE(hide, 0) as hide,
               COALESCE(post_title, '') as post_title,
               COALESCE(post_content, '') as post_content,
               COALESCE(post_created_at, CURRENT_TIMESTAMP) as post_created_at,
               COALESCE(post_updated_at, CURRENT_TIMESTAMP) as post_updated_at
        FROM nav_table
        WHERE id = ?
    `

    row := database.DB.QueryRow(sql, id)

    var (
        hide, sort int64
    )

    err := row.Scan(
        &tool.Id, &tool.Name, &tool.Url, &tool.Logo,
        &tool.Catelog, &tool.Desc, &tool.Content,
        &sort, &hide,
        &tool.PostTitle, &tool.PostContent,
        &tool.PostCreatedAt, &tool.PostUpdatedAt,
    )

    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            logger.LogError("工具不存在, ID: %d", id)
            return tool, fmt.Errorf("工具不存在")
        }
        logger.LogError("查询工具出错: %v", err)
        return tool, fmt.Errorf("数据库查询错误: %v", err)
    }

    tool.Sort = int(sort)
    tool.Hide = hide != 0

    logger.LogInfo("成功获取工具信息: %+v", tool)
    return tool, nil
}

// UpdatePost 更新工具的帖子内容
func UpdatePost(id int64, post *types.Post) error {
    // SQL 查询，添加更新人和更新时间
    sql := `
        UPDATE nav_table
        SET post_title = ?,
            post_content = ?,
            post_updated_at = ?,
            updated_by = ?
        WHERE id = ?
    `

    // 使用提供的时间戳和用户信息
    updateTime := time.Date(2025, 1, 20, 2, 54, 4, 0, time.UTC)
    updatedBy := "ziren926"

    // 执行更新
    result, err := database.DB.Exec(sql,
        post.Title,           // 帖子标题
        post.Content,         // 帖子内容
        updateTime,           // 更新时间
        updatedBy,           // 更新人
        id,                  // 工具 ID
    )
    if err != nil {
        return fmt.Errorf("更新帖子失败: %v", err)
    }

    // 检查是否有行被更新
    affected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("获取影响行数失败: %v", err)
    }
    if affected == 0 {
        return fmt.Errorf("工具不存在")
    }

    // 记录审计日志
    logger.LogInfo("帖子已更新 - ID: %d, 更新人: %s, 更新时间: %s",
        id,
        updatedBy,
        updateTime.Format("2006-01-02 15:04:05"))

    return nil
}

func GetPost(id int64) (*types.Post, error) {
    currentTime := time.Date(2025, 1, 20, 3, 27, 4, 0, time.UTC)

    sql := `
        SELECT post_title, post_content, post_created_at, post_updated_at
        FROM nav_table
        WHERE id = ?
    `

    var post types.Post
    err := database.DB.QueryRow(sql, id).Scan(
        &post.Title,
        &post.Content,
        &post.CreatedAt,
        &post.UpdatedAt,
    )

    if err != nil {
        if err != nil && err.Error() == "sql: no rows in result set" {
            logger.LogInfo("未找到ID为 %d 的帖子", id)
            return nil, fmt.Errorf("帖子不存在")
        }
        logger.LogError("查询帖子失败: %v", err)
        return nil, fmt.Errorf("查询帖子失败: %v", err)
    }

    logger.LogInfo("成功获取帖子 - ID: %d, 时间: %s",
        id, currentTime.Format("2006-01-02 15:04:05"))

    return &post, nil
}

// AddPost 添加新帖子
func AddPost(toolId int64, post *types.Post) error {
    sql := `
        UPDATE nav_table
        SET post_title = ?,
            post_content = ?,
            post_created_at = ?,
            post_updated_at = ?,
            created_by = ?,
            updated_by = ?
        WHERE id = ?
    `

    createdTime := time.Date(2025, 1, 20, 2, 54, 4, 0, time.UTC)
    createdBy := "ziren926"

    result, err := database.DB.Exec(sql,
        post.Title,
        post.Content,
        createdTime,    // 创建时间
        createdTime,    // 更新时间（初始与创建时间相同）
        createdBy,      // 创建人
        createdBy,      // 更新人（初始与创建人相同）
        toolId,
    )
    if err != nil {
        return fmt.Errorf("添加帖子失败: %v", err)
    }

    affected, err := result.RowsAffected()
    if err != nil {
        return fmt.Errorf("获取影响行数失败: %v", err)
    }
    if affected == 0 {
        return fmt.Errorf("工具不存在")
    }

    logger.LogInfo("新帖子已添加 - 工具ID: %d, 创建人: %s, 创建时间: %s",
        toolId,
        createdBy,
        createdTime.Format("2006-01-02 15:04:05"))

    return nil
}