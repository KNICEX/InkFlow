package dao

type Ink struct {
	Id       int64    `json:"id"`
	Title    string   `json:"title"`
	AuthorId int64    `json:"author_id"`
	Content  string   `json:"content"`
	Tags     []string `json:"tags"`
}

// TODO ink写入es之前需要去除html标签
