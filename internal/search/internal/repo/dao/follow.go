package dao

type Follow struct {
	Id       int64 `json:"id"`
	Follower int64 `json:"follower"`
	Followee int64 `json:"followee"`
}
