package types

import "time"

// 默认是 0
type Setting struct {
    Id              int    `json:"id"`
    Favicon         string `json:"favicon"`
    Title           string `json:"title"`
    GovRecord       string `json:"govRecord"`
    Logo192         string `json:"logo192"`
    Logo512         string `json:"logo512"`
    HideAdmin       bool   `json:"hideAdmin"`
    HideGithub      bool   `json:"hideGithub"`
    JumpTargetBlank bool   `json:"jumpTargetBlank"`
}

type Tool struct {
    Id            int64     `json:"id"`
    Name          string    `json:"name"`
    Url           string    `json:"url"`
    Logo          string    `json:"logo"`
    Desc          string    `json:"desc"`
    Catelog       string    `json:"catelog"`
    Content       string    `json:"content,omitempty"`
    Sort          int       `json:"sort"`
    Hide          bool      `json:"hide"`
    PostTitle     string    `json:"post_title,omitempty"`
    PostContent   string    `json:"post_content,omitempty"`
    PostCreatedAt time.Time `json:"post_created_at"`  // 帖子创建时间
    PostUpdatedAt time.Time `json:"post_updated_at"`  // 帖子更新时间
}

type Token struct {
    Id       int    `json:"id"`
    Name     string `json:"name"`
    Value    string `json:"value"`
    Disabled int    `json:"disabled"`
}

type User struct {
    Id       int    `json:"id"`
    Name     string `json:"name"`
    Password string `json:"password"`
}

type Img struct {
    Id    int    `json:"id"`
    Url   string `json:"url"`
    Value string `json:"value"`
}

type Catelog struct {
    Id   int    `json:"id"`
    Name string `json:"name"`
    Sort int    `json:"sort"`
    Hide bool   `json:"hide"`
}

type Post struct {
    ID          int64     `json:"id"`
    ToolID      int64     `json:"tool_id"`
    Title       string    `json:"title"`
    Content     string    `json:"content"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
    CreatedBy   string    `json:"created_by"`
    UpdatedBy   string    `json:"updated_by"`
}