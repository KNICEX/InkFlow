package web

type UriIdReq struct {
	Id int64 `json:"id" form:"id" uri:"id" binding:"required"` // uri id
}

type OffsetPagedReq struct {
	Offset int `json:"offset" form:"offset"`                  // offset
	Limit  int `json:"limit" form:"limit" binding:"required"` // limit
}

type MaxIdPagedReq struct {
	MaxId int64 `json:"maxId" form:"maxId"`
	Limit int   `json:"limit" form:"limit"`
}
