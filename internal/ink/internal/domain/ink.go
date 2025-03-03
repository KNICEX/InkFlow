package domain

type Ink struct {
	Id          int64
	Title       string
	Summary     string
	ContentHtml string
	// 可能引入块编辑器
	ContentMeta string
	Author      Author
}

type Author struct {
	Id      int64
	Name    string
	Account string
}
