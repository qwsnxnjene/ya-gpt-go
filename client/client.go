package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/qwsnxnjene/gpt-ya-try/types"
)

const (
	cloudsURL   = "https://resource-manager.api.cloud.yandex.net/resource-manager/v1/clouds"
	foldersURL  = "https://resource-manager.api.cloud.yandex.net/resource-manager/v1/folders"
	yaGptProURL = "gpt://%s/yandexgpt/rc"
	requestURL  = "https://llm.api.cloud.yandex.net/foundationModels/v1/completion"
)

func cloudID(token string) (string, error) {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodGet, cloudsURL, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res types.CloudResponse
	if err = json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.Clouds[0].ID, nil
}

func folderID(token string, cloudID string) (string, error) {
	c := http.Client{}

	req, err := http.NewRequest(http.MethodGet, foldersURL, nil)
	if err != nil {
		return "", err
	}

	q := url.Values{}
	q.Add("cloudId", cloudID)

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.URL.RawQuery = q.Encode()

	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var folders types.FolderResponse
	if err = json.Unmarshal(body, &folders); err != nil {
		return "", err
	}

	return folders.Folders[0].ID, nil
}

func initial(token string) (string, error) {
	cloudID, err := cloudID(token)
	if err != nil {
		return "", err
	}

	folderID, err := folderID(token, cloudID)
	if err != nil {
		return "", err
	}

	modelURL := fmt.Sprintf(yaGptProURL, folderID)
	return modelURL, nil
}

func Request(token string, request string) (string, error) {
	modelURL, err := initial(token)
	if err != nil {
		return "", err
	}

	toPostQ := types.ModelRequest{
		ModelUri: modelURL,
		CompletionOptions: types.CompletionOptions{
			Temperature: 0.3,
			MaxTokens:   1000,
		},
		Messages: []types.Message{
			{
				Role: "system",
				Text: "Помоги решить проблему",
			},
			{
				Role: "user",
				Text: request,
			},
		},
	}

	pq, err := json.Marshal(toPostQ)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", requestURL, bytes.NewBuffer(pq))
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	req.Header.Add("Accept", "application/json")

	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var answer types.ResponseGpt

	if err = json.Unmarshal(body, &answer); err != nil {
		return "", err
	}

	return answer.Result.Alternatives[0].Message.Text, nil
}
