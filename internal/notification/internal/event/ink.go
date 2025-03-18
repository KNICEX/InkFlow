package event

type InkInteractiveEvent struct {
	InkId    int64
	AuthorId int64
	Uid      int64

	Type string
}
