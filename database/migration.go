package database

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
