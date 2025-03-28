package domain

type Relation struct {
	Uid       int64
	TargetUid int64
	Relation  RelationType
}

type RelationType string

const (
	RelationTypeFollow   RelationType = "follow"
	RelationTypeUnFollow RelationType = "unfollow"
	RelationTypeBlock    RelationType = "block"
	RelationTypeUnBlock  RelationType = "unblock"
)
