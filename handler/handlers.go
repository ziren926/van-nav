package handler

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
    "fmt"
	"github.com/gin-gonic/gin"
	"github.com/ziren926/van-nav/database"
	"github.com/ziren926/van-nav/logger"
	"github.com/ziren926/van-nav/service"
	"github.com/ziren926/van-nav/types"
	"github.com/ziren926/van-nav/utils"
)

func ExportToolsHandler(c *gin.Context) {
	tools := service.GetAllTool()
	c.JSON(200, gin.H{
		"success": true,
		"message": "导出工具成功",
		"data":    tools,
	})
}

func ImportToolsHandler(c *gin.Context) {
	var tools []types.Tool
	err := c.ShouldBindJSON(&tools)
	if err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	// 导入所有工具
	service.ImportTools(tools)
	c.JSON(200, gin.H{
		"success": true,
		"message": "导入工具成功",
	})
}

func DeleteApiTokenHandler(c *gin.Context) {
	// 删除 Token
	id := c.Param("id")
	sql_delete_api_token := `
		UPDATE nav_api_token
		SET disabled = 1
		WHERE id = ?;
		`
	stmt, err := database.DB.Prepare(sql_delete_api_token)
	utils.CheckErr(err)
	res, err := stmt.Exec(id)
	utils.CheckErr(err)
	_, err = res.RowsAffected()
	utils.CheckErr(err)
	c.JSON(200, gin.H{
		"success": true,
		"message": "删除 API Token 成功",
	})
}

func AddApiTokenHandler(c *gin.Context) {
	var token types.AddTokenDto
	err := c.ShouldBindJSON(&token)
	if err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	newId := utils.GenerateId()
	var signedJwt string
	signedJwt, err = utils.SignJWTForAPI(token.Name, newId)
	if err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	service.AddApiTokenInDB(types.Token{
		Name:     token.Name,
		Value:    signedJwt,
		Id:       newId,
		Disabled: 0,
	})
	// 签名 jwt
	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"id":    newId,
			"Value": signedJwt,
			"Name":  token.Name,
		},
		"message": "添加 Token 成功",
	})
}

func UpdateSettingHandler(c *gin.Context) {
	var data types.Setting
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	logger.LogInfo("更新配置: %+v", data)
	err := service.UpdateSetting(data)
	if err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"success": true,
		"message": "更新配置成功",
	})
}

func UpdateUserHandler(c *gin.Context) {
	var data types.UpdateUserDto
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	service.UpdateUser(data)
	c.JSON(200, gin.H{
		"success": true,
		"message": "更新用户成功",
	})
}

func GetAllHandler(c *gin.Context) {
    tools := service.GetAllTool()
    // 获取全部数据，包括帖子内容
    catelogs := service.GetAllCatelog()
    if !utils.IsLogin(c) {
        // 过滤掉隐藏工具
        tools = utils.FilterHideTools(tools, catelogs)
    }
    if !utils.IsLogin(c) {
        // 过滤掉隐藏分类
        catelogs = utils.FilterHideCates(catelogs)
    }
    setting := service.GetSetting()
    c.JSON(200, gin.H{
        "success": true,
        "data": gin.H{
            "tools":    tools,
            "catelogs": catelogs,
            "setting":  setting,
        },
    })
}

func GetLogoImgHandler(c *gin.Context) {
	url := c.Query("url")

	img := service.GetImgFromDB(url)
	imgBuffer, _ := base64.StdEncoding.DecodeString(img.Value)
	// 检测不同的格式发送不同的响应头
	l := strings.Split(url, ".")
	suffix := l[len(l)-1]
	var t string = "image/x-icon"
	if suffix == "svg" || strings.Contains(url, ".svg") {
		t = "image/svg+xml"
	}
	if suffix == "png" {
		t = "image/png"
	}
	c.Writer.Header().Set("content-type", t)
	c.Writer.WriteString(string(imgBuffer))
	// resStr := "data:image/x-icon;base64," + img.Value
	// c.Writer.WriteString(resStr)
}

func GetAdminAllDataHandler(c *gin.Context) {
    // 1. 获取所有工具数据
    tools := []types.Tool{}
    toolRows, err := database.DB.Query(`
        SELECT id, name, desc, logo, url, catelog, sort, hide
        FROM tools
        ORDER BY sort ASC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer toolRows.Close()

    for toolRows.Next() {
        var tool types.Tool
        err := toolRows.Scan(&tool.ID, &tool.Name, &tool.Desc, &tool.Logo, &tool.Url, &tool.Catelog, &tool.Sort, &tool.Hide)
        if err != nil {
            continue
        }
        tools = append(tools, tool)
    }

    // 2. 获取所有分类数据
    catelogs := []types.Catelog{}
    catelogRows, err := database.DB.Query(`
        SELECT id, name, sort, hide
        FROM catelogs
        ORDER BY sort ASC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer catelogRows.Close()

    for catelogRows.Next() {
        var catelog types.Catelog
        err := catelogRows.Scan(&catelog.ID, &catelog.Name, &catelog.Sort, &catelog.Hide)
        if err != nil {
            continue
        }
        catelogs = append(catelogs, catelog)
    }

    // 3. 获取所有 API Token 数据
    tokens := []types.ApiToken{}
    tokenRows, err := database.DB.Query(`
        SELECT id, name, value
        FROM api_tokens
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer tokenRows.Close()

    for tokenRows.Next() {
        var token types.ApiToken
        err := tokenRows.Scan(&token.ID, &token.Name, &token.Value)
        if err != nil {
            continue
        }
        tokens = append(tokens, token)
    }

    // 4. 获取系统设置数据
    var setting types.Setting
    err = database.DB.QueryRow(`
        SELECT
            favicon,
            title,
            gov_record,
            jump_target_blank,
            logo192,
            logo512,
            hide_admin,
            hide_github
        FROM settings
        LIMIT 1
    `).Scan(
        &setting.Favicon,
        &setting.Title,
        &setting.GovRecord,
        &setting.JumpTargetBlank,
        &setting.Logo192,
        &setting.Logo512,
        &setting.HideAdmin,
        &setting.HideGithub,
    )
    if err != nil {
        // 如果没有设置数据，使用默认值
        setting = types.Setting{
            Favicon:         "/favicon.ico",
            Title:          "VanNav",
            GovRecord:      "",
            JumpTargetBlank: true,
            Logo192:        "/logo192.png",
            Logo512:        "/logo512.png",
            HideAdmin:      false,
            HideGithub:     false,
        }
    }

    // 5. 获取所有帖子数据
    posts := []types.Post{}
    postRows, err := database.DB.Query(`
        SELECT id, title, content, create_time, update_time
        FROM posts
        ORDER BY create_time DESC
    `)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    defer postRows.Close()

    for postRows.Next() {
        var post types.Post
        err := postRows.Scan(&post.ID, &post.Title, &post.Content, &post.CreateTime, &post.UpdateTime)
        if err != nil {
            continue
        }
        posts = append(posts, post)
    }

    // 返回所有数据
    c.JSON(http.StatusOK, gin.H{
        "tools":    tools,
        "catelogs": catelogs,
        "tokens":   tokens,
        "setting":  setting,
        "posts":    posts,
    })
}

func LoginHandler(c *gin.Context) {
	var data types.LoginDto
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	user := service.GetUser(data.Name)
	if user.Name == "" {
		c.JSON(200, gin.H{
			"success":      false,
			"errorMessage": "用户名不存在",
		})
		return
	}
	if user.Password != data.Password {
		c.JSON(200, gin.H{
			"success":      false,
			"errorMessage": "密码错误",
		})
		return
	}
	// 生成 token
	token, err := utils.SignJWT(user)
	utils.CheckErr(err)

	c.JSON(200, gin.H{
		"success": true,
		"message": "登录成功",
		"data": gin.H{
			"user":  user,
			"token": token,
		},
	})

}

// 退出登录
func LogoutHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"success": true,
		"message": "登出成功",
	})
}

func AddToolHandler(c *gin.Context) {
    // 添加工具，支持帖子相关字段
    var data types.AddToolDto
    if err := c.ShouldBindJSON(&data); err != nil {
        utils.CheckErr(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success":      false,
            "errorMessage": err.Error(),
        })
        return
    }

    logger.LogInfo("新增工具: %s, 帖子標題: %s", data.Name, data.PostTitle)
    // 修改這裡，直接調用 AddTool
    service.AddTool(data)

    if data.Logo == "" {
        // 由於我們不再使用 id，需要修改這裡的邏輯
        go service.LazyFetchLogo(data.Url, 0) // 或者從 data 中獲取 ID
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "添加成功",
    })
}

func DeleteToolHandler(c *gin.Context) {
	// 删除工具
	id := c.Param("id")
	sql_delete_tool := `
		DELETE FROM nav_table WHERE id = ?;
		`
	stmt, err := database.DB.Prepare(sql_delete_tool)
	utils.CheckErr(err)
	res, err := stmt.Exec(id)
	utils.CheckErr(err)
	_, err = res.RowsAffected()
	utils.CheckErr(err)
	// 删除工具的 logo，如果有
	numberId, err := strconv.Atoi(id)
	utils.CheckErr(err)
	url1 := service.GetToolLogoUrlById(numberId)
	urlEncoded := url.QueryEscape(url1)
	sql_delete_tool_img := `
		DELETE FROM nav_img WHERE url = ?;
		`
	stmt, err = database.DB.Prepare(sql_delete_tool_img)
	utils.CheckErr(err)
	res, err = stmt.Exec(urlEncoded)
	utils.CheckErr(err)
	_, err = res.RowsAffected()
	utils.CheckErr(err)
	c.JSON(200, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

func AddCatelogHandler(c *gin.Context) {
	// 添加分类
	var data types.AddCatelogDto
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	service.AddCatelog(data)

	c.JSON(200, gin.H{
		"success": true,
		"message": "增加分类成功",
	})
}

func DeleteCatelogHandler(c *gin.Context) {
	// 删除分类
	id := c.Param("id")
	sql_delete_catelog := `
		DELETE FROM nav_catelog WHERE id = ?;
		`
	stmt, err := database.DB.Prepare(sql_delete_catelog)
	utils.CheckErr(err)
	res, err := stmt.Exec(id)
	utils.CheckErr(err)
	_, err = res.RowsAffected()
	utils.CheckErr(err)
	c.JSON(200, gin.H{
		"success": true,
		"message": "删除分类成功",
	})
}

func UpdateCatelogHandler(c *gin.Context) {
	// 更新分类
	var data types.UpdateCatelogDto
	if err := c.ShouldBindJSON(&data); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}
	service.UpdateCatelog(data)

	c.JSON(200, gin.H{
		"success": true,
		"message": "更新分类成功",
	})
}

func ManifastHanlder(c *gin.Context) {

	setting := service.GetSetting()
	title := setting.Title

	var icons = []gin.H{}

	logo192 := setting.Logo192
	if logo192 == "" {
		logo192 = "logo192.png"
	}

	logo512 := setting.Logo512
	if logo512 == "" {
		logo512 = "logo512.png"
	}

	icons = append(icons, gin.H{
		"src":   logo192,
		"type":  "image/png",
		"sizes": "192x192",
	})
	icons = append(icons, gin.H{
		"src":   logo512,
		"type":  "image/png",
		"sizes": "512x512",
	})

	if title == "" {
		title = "Van nav"
	}
	c.JSON(200, gin.H{
		"short_name":       title,
		"name":             title,
		"icons":            icons,
		"start_url":        "/",
		"display":          "standalone",
		"scope":            "/",
		"theme_color":      "#000000",
		"background_color": "#ffffff",
	})
}

func UpdateToolsSortHandler(c *gin.Context) {
	var updates []types.UpdateToolsSortDto
	if err := c.ShouldBindJSON(&updates); err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	err := service.UpdateToolsSort(updates)
	if err != nil {
		utils.CheckErr(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success":      false,
			"errorMessage": err.Error(),
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "更新排序成功",
	})
}

func UpdateToolHandler(c *gin.Context) {
    var data types.UpdateToolDto
    if err := c.ShouldBindJSON(&data); err != nil {
        utils.CheckErr(err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success":      false,
            "errorMessage": err.Error(),
        })
        return
    }

    logger.LogInfo("更新工具: %s, 帖子标题: %s", data.Name, data.PostTitle)
    err := service.UpdateTool(data)
    if err != nil {
        utils.CheckErr(err)
        c.JSON(http.StatusInternalServerError, gin.H{
            "success":      false,
            "errorMessage": err.Error(),
        })
        return
    }

    if data.Logo == "" {
        logger.LogInfo("%s 获取 logo: %s", data.Name, data.Logo)
        go service.LazyFetchLogo(data.Url, int64(data.Id))
    }

    c.JSON(200, gin.H{
        "success": true,
        "message": "更新成功",
    })
}

// 添加获取工具详情的处理函数
func GetToolDetailHandler(c *gin.Context) {
    id := c.Param("id")
    logger.LogInfo("收到获取工具详情请求，ID: %s", id)

    numberId, err := strconv.ParseInt(id, 10, 64)
    if err != nil {
        logger.LogError("无效的ID格式: %s, 错误: %v", id, err)
        c.JSON(http.StatusBadRequest, gin.H{
            "success":      false,
            "errorMessage": fmt.Sprintf("无效的ID格式: %s", id),
        })
        return
    }

    tool, err := service.GetToolById(numberId)
    if err != nil {
        logger.LogError("获取工具失败, ID: %d, 错误: %v", numberId, err)
        if err.Error() == "工具不存在" {
            c.JSON(http.StatusNotFound, gin.H{
                "success":      false,
                "errorMessage": "工具不存在",
            })
            return
        }
        c.JSON(http.StatusInternalServerError, gin.H{
            "success":      false,
            "errorMessage": fmt.Sprintf("获取工具信息失败: %v", err),
        })
        return
    }

    // 根据是否登录返回不同级别的信息
    if !utils.IsLogin(c) && tool.Hide {
        c.JSON(http.StatusForbidden, gin.H{
            "success":      false,
            "errorMessage": "无权访问该工具",
        })
        return
    }

    logger.LogInfo("成功获取工具详情，ID: %d", numberId)
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "data":    tool,
    })
}