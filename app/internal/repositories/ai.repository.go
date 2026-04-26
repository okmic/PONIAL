package repositories

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type YandexGPTResponseResult struct {
	Text       string
	UsedTokens int
}

func YandexGPTLite(msg, prompt string) (YandexGPTResponseResult, error) {
	requestBody := map[string]interface{}{
		"modelUri": fmt.Sprintf("gpt://%s/yandexgpt-lite/latest", os.Getenv("YANDEX_FOLDER_ID")),
		"messages": []map[string]string{
			{"role": "system", "text": msg},
			{"role": "user", "text": prompt},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка маршалинга: %w", err)
	}

	url := "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Api-Key %s", os.Getenv("YANDEX_API_KEY")))
	req.Header.Set("x-folder-id", os.Getenv("YANDEX_FOLDER_ID"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка запроса: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка API: %s", string(body))
	}

	var gptResponse struct {
		Result struct {
			Alternatives []struct {
				Message struct {
					Text string `json:"text"`
				} `json:"message"`
			} `json:"alternatives"`
			Usage struct {
				TotalTokens int `json:"totalTokens"`
			} `json:"usage"`
		} `json:"result"`
	}

	if err := json.Unmarshal(body, &gptResponse); err != nil {
		return YandexGPTResponseResult{}, fmt.Errorf("ошибка парсинга: %w", err)
	}

	if len(gptResponse.Result.Alternatives) == 0 {
		return YandexGPTResponseResult{}, fmt.Errorf("пустой ответ от GPT")
	}

	return YandexGPTResponseResult{
		Text:       gptResponse.Result.Alternatives[0].Message.Text,
		UsedTokens: gptResponse.Result.Usage.TotalTokens,
	}, nil
}
