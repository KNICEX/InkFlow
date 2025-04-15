package llm

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/KNICEX/InkFlow/internal/review/internal/service"
	"strings"
	"text/template"
)

const reviewPrompt = `
你是一个内容审核助手，专注于社交平台的内容合规性判断。
请根据以下标准判断文本是否符合社区规范（例如：不得包含暴力、色情、歧视、诈骗、明显广告推销[可以接受商品推荐]等内容）。
你的任务是判断是否通过审核，如果不通过,给出简洁明了的理由(reason)。
如果通过,你还需要给内容打一个0-100的分数(reviewScore)，表示内容的充实度, 
并且为文章内容打上标签(reviewTags),要求从大分类到小分类尽量全面, 10个左右, 例如：科技->AI->ChatGPT。
请按照json格式输入,如下：
{  
	"passed": true | false,  
	"reason": "如不通过，请说明原因；如通过，为空",
	"reviewScore": 0-100,
	"reviewTags": ["tag1", "tag2"...]
输出必须是合法 JSON，不要添加额外解释，不要有多余文本。

内容：{{.content}}
`

type Service struct {
	llm      ai.LLMService
	template *template.Template
}

func NewLLMService(llm ai.LLMService) service.Service {
	return &Service{
		llm:      llm,
		template: template.Must(template.New("review").Parse(reviewPrompt)),
	}
}

func (s *Service) trimJsonResp(content string) string {
	content = strings.Trim("\n", content)
	lines := strings.Split(content, "\n")
	if len(lines) > 3 {
		lines = lines[1 : len(lines)-1]
	}
	return strings.Join(lines, "\n")
}

func (s *Service) ReviewInk(ctx context.Context, ink domain.Ink) (domain.ReviewResult, error) {
	var bs bytes.Buffer
	if err := s.template.Execute(&bs, map[string]any{
		"content": ink.Content,
	}); err != nil {
		return domain.ReviewResult{}, err
	}

	resp, err := s.llm.AskOnce(ctx, bs.String())
	if err != nil {
		return domain.ReviewResult{}, err
	}

	var result domain.ReviewResult
	err = json.Unmarshal([]byte(s.trimMarkdown(resp.Content)), &result)
	if err != nil {
		return domain.ReviewResult{}, err
	}

	result.Reason = strings.Trim(result.Reason, "\"")

	return result, nil
}

func (s *Service) trimMarkdown(content string) string {
	lines := strings.Split(content, "\n")
	if len(lines) > 3 {
		lines = lines[1 : len(lines)-1]
	}
	return strings.Join(lines, "\n")
}
