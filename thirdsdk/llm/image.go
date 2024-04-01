package llm

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

func ImageDescribe(openaiApiKey string, openaiUrl string, imageUrl string) (string, error) {
	ctx := context.Background()
	config := openai.DefaultConfig(openaiApiKey)
	config.BaseURL = openaiUrl
	client := openai.NewClientWithConfig(config)

	// 用户问题，携带上下文
	dialogue := []openai.ChatCompletionMessage{
		{Role: openai.ChatMessageRoleSystem, Content: "你是一个图片识别助手"},
		{Role: openai.ChatMessageRoleUser, MultiContent: []openai.ChatMessagePart{
			{
				Type: openai.ChatMessagePartTypeText,
				Text: "图片的内容是什么？",
			},
			{
				Type: openai.ChatMessagePartTypeImageURL,
				ImageURL: &openai.ChatMessageImageURL{
					URL:    imageUrl,
					Detail: openai.ImageURLDetailAuto,
				},
			},
		},
		},
	}

	resp, err := client.CreateChatCompletion(ctx,
		openai.ChatCompletionRequest{
			// 新模型可以返回两个函数 gpt6 1106模型，gpt3.5 1106 模型
			Model:    openai.GPT4VisionPreview,
			Messages: dialogue,
		},
	)
	if err != nil || len(resp.Choices) == 0 {
		return "", err
	}

	return resp.Choices[0].Message.Content, nil
}
