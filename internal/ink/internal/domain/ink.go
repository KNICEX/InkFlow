package domain

import (
	"time"
)

type Ink struct {
	Id          int64
	Title       string
	Cover       string
	Summary     string
	Category    Category
	ContentType ContentType
	ContentHtml string
	// 可能引入块编辑器
	ContentMeta string
	// 手动添加的标签
	Tags []string
	// ai生成的标签
	AiTags    []string
	Status    Status
	Author    Author
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Abstract 返回摘要(大约前50个字符)
func (i Ink) Abstract() string {
	// TODO 这里要去掉各种标记符号
	if len(i.ContentHtml) > 50 {
		// TODO 这里需要一个转换函数只提取文字内容
		return i.ContentMeta[:50]
	}
	return i.ContentHtml
}

func (i Ink) CanPublish() bool {
	return i.Status != InkStatusPending
}

type Status int

func (s Status) ToInt() int {
	return int(s)
}
func StatusFromInt(i int) Status {
	return Status(i)
}

const (
	InkStatusUnKnown     Status = iota
	InkStatusUnPublished        // 未发布
	InkStatusPending            // 待审核
	InkStatusRejected           // 审核拒绝
	InkStatusPublished          // 已发布
	InkStatusPrivate            // 私密

	InKStatusDeleted // 已删除
)

type Author struct {
	Id      int64
	Name    string
	Account string
}

type Category struct {
	Id int64
}

type ContentType int

const (
	ContentTypeUnknown  ContentType = iota
	ContentTypeShare                // 图文分享
	ContentTypeLongForm             // 长文
	ContentTypeHelp                 // 求助
)

func ContentTypeFromInt(i int) ContentType {
	return ContentType(i)
}

func (c ContentType) ToInt() int {
	return int(c)
}

type TagStats struct {
	Name      string
	Reference int64
}
