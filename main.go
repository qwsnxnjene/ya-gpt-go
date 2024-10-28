package main

import (
	"bytes"
	"encoding/json"

	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/qwsnxnjene/gpt-ya-try/client"
	"github.com/qwsnxnjene/gpt-ya-try/types"
)

func main() {
	iamtoken, err := mustToken()
	if err != nil {
		log.Fatal(err)
	}
	toks, err := iamToken(iamtoken)
	if err != nil {
		log.Fatal(err)
	}

	req := "какие важные события произошли в 2004 году?"

	ans, err := client.Request(toks, req)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(ans)
}

func mustToken() (string, error) {
	token := flag.String(
		"oauth-token",
		"",
		"token for access to yandex gpt",
	)

	flag.Parse()

	if *token == "" {
		return "", errors.New("[mustToken]: flag value is empty")
	}

	return *token, nil
}

func iamToken(ouath string) (string, error) {
	q := types.Auth{YandexPassportOauthToken: ouath}
	marshalledQ, err := json.Marshal(q)
	if err != nil {
		return "", fmt.Errorf("[iamToken]: can't get IamToken: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, "https://iam.api.cloud.yandex.net/iam/v1/tokens", bytes.NewBuffer(marshalledQ))
	if err != nil {
		return "", fmt.Errorf("[iamToken]: can't get IamToken: %w", err)
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", fmt.Errorf("[iamToken]: can't get IamToken: %w", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("[iamToken]: can't get IamToken: %w", err)
	}
	defer resp.Body.Close()

	var forToken types.IamTokenAuth

	if err = json.Unmarshal(body, &forToken); err != nil {
		return "", fmt.Errorf("[iamToken]: can't get IamToken: %w", err)
	}

	return forToken.IamToken, nil
}
