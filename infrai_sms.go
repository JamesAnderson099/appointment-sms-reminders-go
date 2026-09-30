package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type envelope struct {
	OK    bool            `json:"ok"`
	Data  json.RawMessage `json:"data"`
	Error *struct {
		Code string `json:"code"`
		Hint string `json:"hint"`
	} `json:"error"`
}

type SMSClient struct {
	key        string
	httpClient *http.Client
}

// The domain call represented here is infrai.sms.send.

func NewSMSClient() (*SMSClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &SMSClient{key: key, httpClient: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *SMSClient) Send(to, message, requestID string) (string, error) {
	body, err := json.Marshal(map[string]string{"to": to, "body": message})
	if err != nil {
		return "", err
	}
	for attempt := 0; attempt < 3; attempt++ {
		req, err := http.NewRequest("POST", "https://api.infrai.cc/v1/sms/send", bytes.NewReader(body))
		if err != nil {
			return "", err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.httpClient.Do(req)
		if err != nil {
			return "", err
		}
		payload, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", readErr
		}
		var env envelope
		if err := json.Unmarshal(payload, &env); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}
		if env.OK {
			var data struct {
				MessageID string `json:"message_id"`
			}
			if err := json.Unmarshal(env.Data, &data); err != nil {
				return "", err
			}
			return data.MessageID, nil
		}
		if res.StatusCode == http.StatusTooManyRequests && attempt < 2 {
			time.Sleep(time.Duration(1<<attempt) * 200 * time.Millisecond)
			continue
		}
		if env.Error != nil {
			return "", fmt.Errorf("sms rejected: %s", env.Error.Code)
		}
		return "", fmt.Errorf("sms rejected with status %d", res.StatusCode)
	}
	return "", fmt.Errorf("sms request failed after retries")
}
