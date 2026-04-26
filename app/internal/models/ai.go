package models

type YandexAIType string

const (
	TypeVoice YandexAIType = "voice"
	TypeText  YandexAIType = "txt"
)

type AITextRequest struct {
	Msg    string `json:"folderId" binding:"omitempty,min=2,max=100"`
	Prompt string `json:"oAuthToken" binding:"omitempty,min=2,max=255"`
}

type AITextResponse struct {
	Text string `json:"text"`
}

type YandexAIResponse struct {
	Text string `json:"text"`
}

func (y *AITextResponse) ToResponse() YandexAIResponse {
	return YandexAIResponse{
		Text: y.Text,
	}
}

type YandexGPTRequest struct {
	ModelURI string `json:"modelUri"`
	Messages []struct {
		Role string `json:"role"`
		Text string `json:"text"`
	} `json:"messages"`
}

type YandexGPTResponse struct {
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
