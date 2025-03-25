package domain

import "time"

type Action struct {
	Id         int64
	UserId     int64
	TargetType TargetType
	TargetId   int64
	ActionType ActionType
	CreatedAt  time.Time
}

type ActionType string

type TargetType string

type Statistics struct {
	TargetType TargetType
	ActionCnt  int64
	StartTime  time.Time
	EndTime    time.Time
}

type User struct {
	Id           int64
	LastActionAt time.Time
}
