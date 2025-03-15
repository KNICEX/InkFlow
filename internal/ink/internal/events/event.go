package events

type InkPublishedEvt struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	AuthorId    int64  `json:"author_id"`
	ContentHtml string `json:"content"`
	ContentMeta string `json:"content_meta"`
	Status      int    `json:"status"`
}
