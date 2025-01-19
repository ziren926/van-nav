package service

import (
	"sync"
	"fmt"


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
    // 除了更新工具本身之外，也要更新 img 表
    sql_update_tool := `
        UPDATE nav_table
        SET name = ?, url = ?, logo = ?, catelog = ?, desc = ?, sort = ?, hide = ?
        WHERE id = ?;
        `
    stmt, err := database.DB.Prepare(sql_update_tool)
    if err != nil {
        return err
    }
    defer stmt.Close()

    _, err = stmt.Exec(data.Name, data.Url, data.Logo, data.Catelog, data.Desc, data.Sort, data.Hide, data.Id)
    if err != nil {
        return err
    }

    // 更新 img
    if data.Logo != "" {
        UpdateImg(data.Logo)
    }

    return nil
}

func AddTool(data types.AddToolDto) (int64, error) {
	// 创建一个互斥锁来保护数据库操作
	var mu sync.Mutex
	mu.Lock()
	defer mu.Unlock()

	tx, err := database.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	sql_add_tool := `
		INSERT INTO nav_table (name, url, logo, catelog, desc, sort, hide)
		VALUES (?, ?, ?, ?, ?, ?, ?);
		`
	stmt, err := tx.Prepare(sql_add_tool)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(data.Name, data.Url, data.Logo, data.Catelog, data.Desc, data.Sort, data.Hide)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	logger.LogInfo("新增工具: %s", data.Name)

	// 在事务完成后再异步更新图片
	if data.Logo != "" {
		UpdateImg(data.Logo)
	}

	return id, nil
}

func GetAllTool() []types.Tool {
	sql_get_all := `
		SELECT id,name,url,logo,catelog,desc,sort,hide FROM nav_table order by sort;
		`
	results := make([]types.Tool, 0)
	rows, err := database.DB.Query(sql_get_all)
	utils.CheckErr(err)
	for rows.Next() {
		var tool types.Tool
		var hide interface{}
		var sort interface{}
		err = rows.Scan(&tool.Id, &tool.Name, &tool.Url, &tool.Logo, &tool.Catelog, &tool.Desc, &sort, &hide)
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
	defer rows.Close()
	return results
}

func GetToolLogoUrlById(id int) string {
	sql_get_tool := `
		SELECT logo FROM nav_table WHERE id=?;
		`
	rows, err := database.DB.Query(sql_get_tool, id)
	utils.CheckErr(err)
	var tool types.Tool
	for rows.Next() {
		err = rows.Scan(&tool.Logo)
		utils.CheckErr(err)

	}
	defer rows.Close()
	return tool.Logo
}

func UpdateToolIcon(id int64, logo string) {
	sql_update_tool := `
		UPDATE nav_table SET logo=? WHERE id=?;
		`
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
    // 简化 SQL 查询，暂时不包括时间相关字段
    sql := `SELECT id, name, url, logo, catelog, desc, content, sort, hide
            FROM nav_table WHERE id = ?`

    row := database.DB.QueryRow(sql, id)

    var (
        hide, sort interface{}
    )

    err := row.Scan(
        &tool.Id,
        &tool.Name,
        &tool.Url,
        &tool.Logo,
        &tool.Catelog,
        &tool.Desc,
        &tool.Content,
        &sort,
        &hide,
    )

    if err != nil {
        if err.Error() == "sql: no rows in result set" {
            logger.LogError("工具不存在, ID: %d", id)
            return tool, fmt.Errorf("工具不存在")
        }
        logger.LogError("查询工具出错: %v", err)
        return tool, err
    }

    // 处理可空字段
    if sort != nil {
        tool.Sort = int(sort.(int64))
    }
    if hide != nil {
        tool.Hide = hide.(int64) != 0
    }

    logger.LogInfo("成功获取工具信息: %+v", tool)
    return tool, nil
}