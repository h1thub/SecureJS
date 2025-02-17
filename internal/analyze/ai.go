package analyze

import (
	"SecureJS/internal/matcher"
	"context"
	"fmt"
	"time"

	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"
	"github.com/volcengine/volcengine-go-sdk/volcengine"
)

func FormatResultsToString(results []*matcher.MatchResult) string {
	var resultString string
	for _, mr := range results {
		if mr.Error != nil {
			resultString += fmt.Sprintf("\n[!] Parse error on %s: %v\n", mr.URL, mr.Error)
			continue
		}
		if len(mr.Items) == 0 {
			continue
		}
		resultString += fmt.Sprintf("\n[+] %s: found %d item(s)\n", mr.URL, len(mr.Items))
		for _, item := range mr.Items {
			resultString += fmt.Sprintf("    - Rule: %s, Matched: %s\n", item.RuleName, item.MatchedText)
		}
	}
	return resultString
}

func Analyze(resultString string, key string, id string) {
	var instructionString = `
您是资深网络安全专家，擅长识别代码中的敏感信息泄露。请保持专业严谨，区分测试数据与实际风险。
请按以下步骤处理输入内容（主要目标是从JS中寻找一些硬编码的”有具体数值的”敏感信息，不需要乱七八遭的代码如某token、key等值为+t.access_token+这样形式的）：
1. 提取条目：识别所有符合 字段名 + 赋值符 + 值  结构的条目
2. 过滤规则：
（1）排除项：无明确值的字段（例如只包含一个AccessKeyId而没有具体的值）；公开/测试密钥；通用配置；没有具体的敏感信息数值（有的只是代码？乱糟糟的？）；
3. 输出格式（严格遵循，不要多余的输出比如总结什么的，就只输出我下面三行内容）：
URL：{文件URL}
敏感信息：{字段名} = {值}
分析：{风险说明，包含用途、泄露后果、是否生产环境}
4. 如果此次分析整体并没有任何敏感信息，就直接表明 “无敏感信息”`

    client := arkruntime.NewClientWithApiKey(
        //通过 os.Getenv 从环境变量中获取 ARK_API_KEY
        key,
        //深度推理模型耗费时间会较长，请您设置较大的超时时间，避免超时导致任务失败。推荐30分钟以上
        arkruntime.WithTimeout(30*time.Minute),
    )
    // 创建一个上下文，通常用于传递请求的上下文信息，如超时、取消等
    ctx := context.Background()
    // 构建聊天完成请求，设置请求的模型和消息内容
    req := model.ChatCompletionRequest{
        // 需要替换 <YOUR_ENDPOINT_ID> 为您的推理接入点 ID
        Model: id,
        Messages: []*model.ChatCompletionMessage{
            {
                // 消息的角色为用户
                Role: model.ChatMessageRoleUser,
                Content: &model.ChatCompletionMessageContent{
                    StringValue: volcengine.String(instructionString + resultString),
                },
            },
        },
    }

    // 发送聊天完成请求，并将结果存储在 resp 中，将可能出现的错误存储在 err 中
    resp, err := client.CreateChatCompletion(ctx, req)
    if err != nil {
        // 若出现错误，打印错误信息并终止程序
        fmt.Printf("standard chat error: %v\n", err)
        return
    }
    // 检查是否触发深度推理，触发则打印思维链内容
    // if resp.Choices[0].Message.ReasoningContent != nil {
    //     fmt.Println(*resp.Choices[0].Message.ReasoningContent)
    // }
    // 打印聊天完成请求的响应结果
    fmt.Println(*resp.Choices[0].Message.Content.StringValue)
}