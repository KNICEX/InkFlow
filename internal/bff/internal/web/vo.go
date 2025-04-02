package web

type UriIdReq struct {
	Id int64 `json:"id" form:"id" uri:"id" binding:"required"` // uri id
}

type OffsetPagedReq struct {
	Offset int `json:"offset" form:"offset"`                  // offset
	Limit  int `json:"limit" form:"limit" binding:"required"` // limit
}
