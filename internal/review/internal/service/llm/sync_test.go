package llm

import (
	"context"
	"github.com/KNICEX/InkFlow/internal/ai"
	"github.com/KNICEX/InkFlow/internal/review/internal/domain"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
	"testing"
)

func TestService_ReviewInk(t *testing.T) {
	cli, err := genai.NewClient(context.Background(), option.WithAPIKey("AIzaSyDsK-uD5Y-mW17slUROmaA4kFpohM4V96Y"))
	if err != nil {
		t.Fatal(err)
	}
	llmSvc := ai.InitLLMService(cli)

	svc := NewLLMService(llmSvc)

	testCases := []struct {
		name     string
		ink      domain.Ink
		wantErr  bool
		wantPass bool
	}{
		{
			name: "a poetry",
			ink: domain.Ink{
				Content: "让我们出发！享受吧！电压正在上升\n跟着节拍拍手\n让我们出发！享受吧！要开始了 甚至连眨眼\n都不要错过那令人忘却的瞬间\n\n在波浪声的引导下\n在无尽的蓝天之下\n想要变得大胆的原因一定是\n多亏了海风和阳光吧\n\n以谁都无法比拟的热情\n你的目光 你的心也都\n被我全都带走了！\n我可等不及了\n必须要好好开始 这双手，来吧\n\n让我们出发！享受吧！电压达到了巅峰\n让我听听心跳的声音\n" +
					"让我们出发！享受吧！将最棒的舞台\n铭刻在你的记忆中\n被耀眼的细语包围\n新的世界正在觉醒\n脱去犹豫赤裸相对\n让我们和你一起尽情狂欢！\n\n眼中映出的海市蜃楼\n在永不沉睡的阳光下\n伴随着渗出的汗水涌现的\n无尽的热情都是因为你\n\n每一秒的心跳\n我会独占你情感中炽热的部分！\n绝不会让你感到无聊\n感受到节奏后 来吧，伸出你的手\n\nRise Up High!\n随心所欲地舞动\n在炙热的沙滩上\n我曾祈愿不要结束\n感受声音" +
					" 永远永远\n\n让我们出发！享受吧！电压正在上升\n跟着节拍拍手\n让我们出发！享受吧！要开始了\n一起分享那连眨眼都忘记的瞬间\n\n让我们出发！享受吧！电压达到了巅峰\n让我听听心跳的声音\n让我们出发！享受吧！将最棒的舞台\n铭刻在你的记忆中\n在最美好的记忆中 永远",
			},
			wantErr:  false,
			wantPass: true,
		},
		{
			name: "a news",
			ink: domain.Ink{
				Content: "今晚，AI 圈也地震了！谷歌深夜搞突袭，正式上线「最强推理大模型」Gemini 2.5 Pro！没错，就是我昨天发的文章谷歌大型推理模型曝光！击败 Claude-3.7-Thinking，泄漏的大模型，代号是「Nebula」，之前就被爆料这个新模型效果据说特别好，打败 o1、o3-mini、Claude 3.7 Thinking 等一众模型。没想到，" +
					"新模型兑现的这么迅速，24 号才被爆料，25 号谷歌就官宣上线！\n\nGemini 2.5 Pro 在大模型榜单 LMSYS Arena 上排名第一，而且是断层第一！分数比 Grok-3、GPT-4.5 整整高出了 40 分！要知道此前 LMSYS 上的顶流模型们的分数咬的特别紧，只差几分。Grok 前脚宣布突破 1400 分数大关，这次 Gemini 2.5 Pro 直接干到了 1443 分，创下最大 jump up 记录。\n\n首先 Gemini 2.5 Pro（模型版本是 gemini-2.5-pro-exp-03-25）是一个推理模型，谷歌称这是迄今为止最强大的模型。" +
					"不止是全面领先，而且是无短板。在所有评测类别（综合能力、编码、数学、创意写作等）中均排名第 1，尤其在带风格控制的复杂提示（Hard Prompts w/ Style Control）和多轮对话（Multi-Turn）表现突出。\n\nGemini 2.5 Pro 不止是谷歌目前最大的推理模型，而且还具备多模态能力，在 Vision Arena 视觉排行榜上也是第一。在网页开发榜单 WebDev Arena 上排名第二，仅次于 Claude-3.7，Claude 的编程地位依旧难以撼动。\n\n下面看下在各个 benchmark 上的具体得分 ——Gemini 2.5 Pro 综合表现拿下最佳" +
					"。尤其在科学（Science）、代码生成、视觉推理（MMMU）和长文本理解（MRCR）上均领先。在号称最难的测试「人类最后一次考试」中，Gemini 2.5 Pro 遥遥领先 OpenAI o3-mini。在号称最难的 AI 测试 “人类最后一次考试” 中，Gemini 2.5 Pro 遥遥领先其他模型。\n\nimage\n\nSWE-bench 代表编码能力，Aider Polyglot 则是代表代码编辑水平。等我看完所有的榜单之后，我只能说 “恐怖如斯”！现在，Gemini 2.5 Pro 已经可以在 Google AI Studio 和 Gemini APP 中使用了。传送门：Google AI Studio\n" +
					"\nimage\n\n接下来看下效果 ——\n\n第一个：曼德博集合演示效果#\n曼德博集合（Mandelbrot set）是一种在复平面上组成分形的点的集合，有人称它是人类有史以来做出的最奇异、最瑰丽的几何图形，曾被称为 “上帝的指纹”。看下 Gemini 2.5 Pro 生成的效果吧。\n\nimage\n\nimage\n\n第二个：网页小游戏#\n还记得这个再熟悉不过的恐龙跑酷游戏吗，记忆里的黑白版变成了有色版。生成地很带感。\n\nimage\n\nimage\n\nGemini 2.5 Pro 最大的优势是，依然具备原生多模态能力和超长上下文长度，目前支持到 1M 窗口，2M " +
					"的在路上了。但是目前尚未公布 API 价格。DeepSeek V3-0324 也刚刚发布，且是最宽松的 MIT 协议，究竟是闭源巨头巩固高地，还是开源阵营推动技术平权？",
			},
			wantPass: true,
		},
		{
			name: "some bad",
			ink: domain.Ink{
				Content: "@xss, 你等着, 你的手机号码是13235435342, 你是住在北京市朝阳区的金华大道45号吗,你看我来不来找你就完事了",
			},
			wantErr:  false,
			wantPass: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			resp, err := svc.ReviewInk(context.Background(), tc.ink)
			if (err != nil) != tc.wantErr {
				t.Errorf("ReviewInk() error = %v, wantErr %v", err, tc.wantErr)
				return
			}
			if resp.Passed != tc.wantPass {
				t.Errorf("ReviewInk() got = %v, want %v", resp.Passed, tc.wantPass)
			}

			t.Logf("%+v", resp)
		})
	}
}
