package domain

type Comment struct {
	Id      int64
	Biz     string
	BizId   int64
	Content string

	Commentator   Commentator
	ReplyUser     Commentator
	RootComment   *Comment
	ParentComment *Comment
	Children      []Comment
	Liked         bool

	CreatedAt int64
	UpdatedAt int64
}

type Commentator struct {
	Id   int64
	Name string
}
