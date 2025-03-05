package event

type InkPublishedEvt struct {
	Id       int64  `json:"id"`
	Title    string `json:"title"`
	AuthorId int64  `json:"author_id"`
	Content  string `json:"content"`
	Status   int    `json:"status"`
}
