package types

type CloudResponse struct {
	Clouds []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	} `json:"clouds"`
}

type FolderResponse struct {
	Folders []struct {
		ID string `json:"id"`
	} `json:"folders"`
}

type CompletionOptions struct {
	Stream      bool    `json:"stream"`
	Temperature float64 `json:"temperature"`
	MaxTokens   int     `json:"maxTokens"`
}

type Message struct {
	Role string `json:"role"`
	Text string `json:"text"`
}

type ModelRequest struct {
	ModelUri          string            `json:"modelUri"`
	CompletionOptions CompletionOptions `json:"completionOptions"`
	Messages          []Message         `json:"messages"`
}

type Alternative struct {
	Message Message `json:"message"`
	Status  string  `json:"status"`
}

type ResponseGpt struct {
	Result struct {
		Alternatives []Alternative `json:"alternatives"`
		Usage        struct {
			InputTextTokens  string `json:"inputTextTokens"`
			CompletionTokens string `json:"completionTokens"`
			TotalTokens      string `json:"totalTokens"`
		} `json:"usage"`
		ModelVersion string `json:"modelVersion"`
	} `json:"result"`
}

type Auth struct {
	YandexPassportOauthToken string `json:"yandexPassportOauthToken"`
}

type IamTokenAuth struct {
	IamToken  string `json:"iamToken"`
	ExpiresAt string `json:"expiresAt"`
}
