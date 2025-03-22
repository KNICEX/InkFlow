package inkpub

import (
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

const ReviewSignal = "review_signal"

const bizInk = "ink"

type Workflow struct {
	activities *Activities
}

func (w *Workflow) InkPublish(ctx workflow.Context, inkId int64) {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var inkInfo ink.Ink
	l := workflow.GetLogger(ctx)
	err := workflow.ExecuteActivity(ctx, w.activities.SubmitReview, inkId).Get(ctx, &inkInfo)
	if err != nil {
		l.Error("ink publish workflow error", "error", err, "inkId", inkId)
		return
	}

	selector := workflow.NewSelector(ctx)
	var reviewResult review.Result
	selector.AddReceive(workflow.GetSignalChannel(ctx, ReviewSignal), func(c workflow.ReceiveChannel, more bool) {
		c.Receive(ctx, &reviewResult)
		l.Info("review result received", "result", reviewResult)
	})
	selector.Select(ctx)

	if reviewResult.Passed {
		// 通过
		inkInfo.AiTags = reviewResult.ReviewTags
		err = workflow.ExecuteActivity(ctx, w.activities.UpdateToPublished, inkInfo.Id, inkInfo.Author.Id).Get(ctx, nil)
		if err != nil {
			l.Error("update ink status to published error", "error", err, "inkId", inkInfo.Id)
			return
		}
		// 提前创建交互记录
		err = workflow.ExecuteActivity(ctx, w.activities.CreateIntr, bizInk, inkInfo.Id).Get(ctx, nil)
		if err != nil {
			l.Error("create interactive error", "error", err, "inkId", inkInfo.Id)
			return
		}

		searchInk := search.Ink{
			Id:       inkInfo.Id,
			Title:    inkInfo.Title,
			Tags:     inkInfo.Tags,
			AiTags:   inkInfo.AiTags,
			Cover:    inkInfo.Cover,
			AuthorId: inkInfo.Author.Id,
			Content:  inkInfo.ContentHtml, // TODO 去除html标签
		}
		// 同步到搜索引擎
		err = workflow.ExecuteActivity(ctx, w.activities.SyncToSearch, searchInk).Get(ctx, nil)
		if err != nil {
			l.Error("sync ink to search error", "error", err, "inkId", inkInfo.Id)
			return
		}

		recommendInk := recommend.Ink{
			Id:        inkInfo.Id,
			AuthorId:  inkInfo.Author.Id,
			Tags:      inkInfo.Tags, // TODO 这里合并 + 去重
			CreatedAt: inkInfo.CreatedAt,
		}
		// 同步到推荐引擎
		err = workflow.ExecuteActivity(ctx, w.activities.SyncToRecommend, recommendInk).Get(ctx, nil)
		if err != nil {
			l.Error("sync ink to recommend error", "error", err, "inkId", inkInfo.Id)
			return
		}
	} else {
		// 未通过审核

		// 更新文章状态为已拒绝
		err = workflow.ExecuteActivity(ctx, w.activities.UpdateInkToRejected, inkInfo.Id).Get(ctx, nil)
		if err != nil {
			l.Error("update ink status to rejected error", "error", err, "inkId", inkInfo.Id)
			return
		}

		// 通知作者拒绝原因
		err = workflow.ExecuteActivity(ctx, w.activities.NotifyRejected, inkInfo.Id, inkInfo.Author.Id, reviewResult.Reason).Get(ctx, nil)
		if err != nil {
			l.Error("notify ink rejected error", "error", err, "inkId", inkInfo.Id)
			return
		}
	}
}
