package web

type NotificationVO struct {
	Id int64 `json:"id"`

	Read bool `json:"read"`
}
