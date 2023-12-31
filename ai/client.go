package ai

import (
	"context"
	"github.com/sashabaranov/go-openai"
)

func Send(ctx context.Context, openaiClient *openai.Client, request string) (string, error) {
	resp, err := openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:            openai.GPT4,
			Temperature:      0,
			MaxTokens:        256,
			TopP:             1,
			FrequencyPenalty: 0,
			PresencePenalty:  0,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "Ты помощник который из текста находит какой продукт съеден и в каком количестве. В ответ ты выдаешь только название продукта и граммовку, калорийность, количество белков, жиров, углеводов продукта в расчете на данное количество. Формат ответа json структура, которая описывает массив структур в которых содержится информация по каждому продукту. Вче числовые значения представляй как float. Если не нашел продуктов то верни пустой массив в ответ.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "Я съел 200 грамм риса с 500 граммами кабачков и 300 грамм куриной грудки.",
				},
				{
					Role:    openai.ChatMessageRoleAssistant,
					Content: "[{\"product_name\": \"Рис\", \"weight\": 200.0, \"calories\": 340.0, \"proteins\": 4.2, \"fats\": 0.6, \"carbohydrates\": 43.5}, {\"product_name\": \"Кабачки\", \"weight\": 500.0, \"calories\": 75.0, \"proteins\": 1.6, \"fats\": 0.3, \"carbohydrates\": 6.1}, {\"product_name\": \"Куриная грудка\", \"weight\": 300.0, \"calories\": 495.0, \"proteins\": 29.4, \"fats\": 4.3, \"carbohydrates\": 0.0}]",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: request,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	content := resp.Choices[0].Message.Content

	return content, nil
}
