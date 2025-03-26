package inkpub

import (
	"fmt"
	"github.com/KNICEX/InkFlow/internal/ink"
	"github.com/KNICEX/InkFlow/internal/recommend"
	"github.com/KNICEX/InkFlow/internal/review"
	"github.com/KNICEX/InkFlow/internal/search"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
	"time"
)

const ReviewSignal = "review-signal"

const bizInk = "ink"

func InkPublish(ctx workflow.Context, inkId int64) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval: time.Second,
			MaximumAttempts: 3,
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var activities *Activities

	var inkInfo ink.Ink
	l := workflow.GetLogger(ctx)

	err := workflow.ExecuteActivity(ctx, activities.FindInkInfo, inkId).
		Get(ctx, &inkInfo)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx, activities.SubmitReview, review.Ink{
		Id:        inkId,
		AuthorId:  inkInfo.Author.Id,
		Cover:     inkInfo.Cover,
		Title:     inkInfo.Title,
		Content:   inkInfo.ContentHtml,
		CreatedAt: inkInfo.CreatedAt,
		UpdatedAt: inkInfo.UpdatedAt,
	}).Get(ctx, &inkInfo)
	if err != nil {
		l.Error("ink publish workflow error", "error", err, "inkId", inkId)
		return err
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
		err = workflow.ExecuteActivity(ctx, activities.UpdateToPublished, inkInfo.Id, inkInfo.Author.Id).Get(ctx, nil)
		if err != nil {
			l.Error("update ink status to published error", "error", err, "inkId", inkInfo.Id)
			return err
		}
		// 提前创建交互记录
		err = workflow.ExecuteActivity(ctx, activities.CreateIntr, bizInk, inkInfo.Id).Get(ctx, nil)
		if err != nil {
			l.Error("create interactive error", "error", err, "inkId", inkInfo.Id)
			return err
		}

		searchInk := search.Ink{
			Id:     inkInfo.Id,
			Title:  inkInfo.Title,
			Tags:   inkInfo.Tags,
			AiTags: inkInfo.AiTags,
			Cover:  inkInfo.Cover,
			Author: search.User{
				Id: inkInfo.Author.Id,
			},
			Content: inkInfo.ContentHtml,
		}
		// 同步到搜索引擎
		err = workflow.ExecuteActivity(ctx, activities.SyncToSearch, searchInk).Get(ctx, nil)
		if err != nil {
			l.Error("sync ink to search error", "error", err, "inkId", inkInfo.Id)
			return err
		}

		recommendInk := recommend.Ink{
			Id:        inkInfo.Id,
			AuthorId:  inkInfo.Author.Id,
			Tags:      inkInfo.Tags, // TODO 这里合并 + 去重
			CreatedAt: inkInfo.CreatedAt,
		}
		// 同步到推荐引擎
		err = workflow.ExecuteActivity(ctx, activities.SyncToRecommend, recommendInk).Get(ctx, nil)
		if err != nil {
			l.Error("sync ink to recommend error", "error", err, "inkId", inkInfo.Id)
			return err
		}
	} else {
		// 未通过审核

		// 更新文章状态为已拒绝
		err = workflow.ExecuteActivity(ctx, activities.UpdateInkToRejected, inkInfo.Id).Get(ctx, nil)
		if err != nil {
			l.Error("update ink status to rejected error", "error", err, "inkId", inkInfo.Id)
			return err
		}

		// 通知作者拒绝原因
		err = workflow.ExecuteActivity(ctx, activities.NotifyRejected, inkInfo.Id, inkInfo.Author.Id, reviewResult.Reason).Get(ctx, nil)
		if err != nil {
			l.Error("notify ink rejected error", "error", err, "inkId", inkInfo.Id)
			return err
		}
	}
	return nil
}

func WorkflowId(inkId int64, pubTime time.Time) string {
	return fmt.Sprintf("ink-pub-%d-%d", inkId, pubTime.UnixMilli())
}
