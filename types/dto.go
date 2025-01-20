package types

type ResUserDto struct {
	Name string `json:"name"`
}

type UpdateUserDto struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type LoginDto struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
type AddTokenDto struct {
	Name string `json:"name"`
}

type UpdateCatelogDto struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Sort int    `json:"sort"`
	Hide bool   `json:"hide"`
}

type AddCatelogDto struct {
	Name string `json:"name"`
	Sort int    `json:"sort"`
	Hide bool   `json:"hide"`
}
type UpdateToolDto struct {
    Id          int    `json:"id"`
    Name        string `json:"name"`
    Url         string `json:"url"`
    Logo        string `json:"logo"`
    Catelog     string `json:"catelog"`
    Desc        string `json:"desc"`
    Sort        int    `json:"sort"`
    Hide        bool   `json:"hide"`
    Content     string `json:"content"`
    PostTitle   string `json:"post_title"`
    PostContent string `json:"post_content"`
}


type UpdateToolsSortDto struct {
	Id   int `json:"id"`
	Sort int `json:"sort"`
}


type AddToolDto struct {
    ID          int64  `json:"id"`
    Name        string `json:"name"`
    Url         string `json:"url"`
    Logo        string `json:"logo"`
    Desc        string `json:"desc"`
    Catelog     string `json:"catelog"`
    Sort        int    `json:"sort"`
    Hide        bool   `json:"hide"`
    // 新增字段
    Content     string `json:"content"`
    PostTitle   string `json:"post_title"`
    PostContent string `json:"post_content"`
}