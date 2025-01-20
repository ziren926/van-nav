package database

import (
    "time"
    "github.com/ziren926/van-nav/logger"
)

func migration_2024_12_13() {
	// 1. 首先更新现有的 NULL 值为 0
	sql_update_null_sort := `
        UPDATE nav_catelog
        SET sort = 0
        WHERE sort IS NULL;
    `

	_, err := DB.Exec(sql_update_null_sort)
	if err != nil {
		panic(err)
	}

	// 2. 创建新表
	sql_create_new_table := `
        CREATE TABLE nav_catelog_new (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            sort INTEGER NOT NULL DEFAULT 0,
						hide BOOLEAN
        );
    `

	_, err = DB.Exec(sql_create_new_table)
	if err != nil {
		panic(err)
	}

	// 3. 复制数据
	sql_copy_data := `
        INSERT INTO nav_catelog_new (id, name, sort, hide)
        SELECT id, name, sort, hide FROM nav_catelog;
    `

	_, err = DB.Exec(sql_copy_data)
	if err != nil {
		panic(err)
	}

	// 4. 删除旧表
	sql_drop_old_table := `DROP TABLE nav_catelog;`

	_, err = DB.Exec(sql_drop_old_table)
	if err != nil {
		panic(err)
	}

	// 5. 重命名新表
	sql_rename_table := `ALTER TABLE nav_catelog_new RENAME TO nav_catelog;`

	_, err = DB.Exec(sql_rename_table)
	if err != nil {
		panic(err)
	}
}

func migration_2025_01_19() {
    // 添加帖子相关字段到 nav_table 表
    if !columnExists("nav_table", "post_title") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN post_title TEXT;`)
    }
    if !columnExists("nav_table", "post_content") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN post_content TEXT;`)
    }
    if !columnExists("nav_table", "post_created_at") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN post_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`)
    }
    if !columnExists("nav_table", "post_updated_at") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN post_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP;`)
    }
}

func migration_2025_01_19_fix() {
    // 1. 添加 content 字段（如果不存在）
    if !columnExists("nav_table", "content") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN content TEXT;`)
    }

    // 2. 确保所有必需字段都存在
    if !columnExists("nav_table", "sort") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN sort INTEGER DEFAULT 0;`)
    }
    if !columnExists("nav_table", "hide") {
        DB.Exec(`ALTER TABLE nav_table ADD COLUMN hide BOOLEAN DEFAULT 0;`)
    }
}





func migration_2025_01_20() {
    logger.LogInfo("开始执行 2025-01-20 数据库迁移...")

    // 获取当前时间
    currentTime := time.Date(2025, 1, 20, 3, 15, 30, 0, time.UTC)

    // 1. 更新现有帖子的更新时间和创建时间
    updateSQL := `
        UPDATE nav_table
        SET post_updated_at = ?,
            post_created_at = ?
        WHERE post_title IS NOT NULL
          AND post_content IS NOT NULL
          AND (post_updated_at IS NULL OR post_created_at IS NULL)
    `

    if _, err := DB.Exec(updateSQL, currentTime, currentTime); err != nil {
        logger.LogError("更新帖子时间戳失败: %v", err)
    }

    // 2. 创建新表以确保字段的正确性和默认值
    createTempTable := `
        CREATE TABLE nav_table_new (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT,
            url TEXT,
            logo TEXT,
            catelog TEXT,
            desc TEXT,
            content TEXT,
            sort INTEGER DEFAULT 0,
            hide BOOLEAN DEFAULT 0,
            post_title TEXT,
            post_content TEXT,
            post_created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
            post_updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        )
    `

    // 3. 复制数据到新表
    migrateData := `
        INSERT INTO nav_table_new
        SELECT id, name, url, logo, catelog, desc, content,
               COALESCE(sort, 0) as sort,
               COALESCE(hide, 0) as hide,
               post_title, post_content,
               COALESCE(post_created_at, ?) as post_created_at,
               COALESCE(post_updated_at, ?) as post_updated_at
        FROM nav_table
    `

    // 执行迁移
    tx, err := DB.Begin()
    if err != nil {
        logger.LogError("开始事务失败: %v", err)
        return
    }

    // 执行创建新表
    if _, err := tx.Exec(createTempTable); err != nil {
        logger.LogError("创建新表失败: %v", err)
        tx.Rollback()
        return
    }

    // 复制数据
    if _, err := tx.Exec(migrateData, currentTime, currentTime); err != nil {
        logger.LogError("迁移数据失败: %v", err)
        tx.Rollback()
        return
    }

    // 替换旧表
    if _, err := tx.Exec("DROP TABLE nav_table"); err != nil {
        logger.LogError("删除旧表失败: %v", err)
        tx.Rollback()
        return
    }

    if _, err := tx.Exec("ALTER TABLE nav_table_new RENAME TO nav_table"); err != nil {
        logger.LogError("重命名表失败: %v", err)
        tx.Rollback()
        return
    }

    // 提交事务
    if err := tx.Commit(); err != nil {
        logger.LogError("提交事务失败: %v", err)
        return
    }

    logger.LogInfo("2025-01-20 03:15:30 数据库迁移完成")
}