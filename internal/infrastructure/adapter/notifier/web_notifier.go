package notifier

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"weather-notification/internal/domain/entity"
)

type WebNotifier struct {
	webhookURL string
	client     *http.Client
	authToken  string
}

func NewWebNotifier(webhookURL string) *WebNotifier {
	authToken := os.Getenv("API_TOKEN")

	return &WebNotifier{
		webhookURL: webhookURL,
		client:     &http.Client{},
		authToken:  authToken,
	}
}

func (n *WebNotifier) Send(ctx context.Context, notification *entity.Notification) error {
	payload := map[string]interface{}{
		"id":        notification.ID,
		"user_id":   notification.UserID,
		"content":   notification.Content,
		"timestamp": notification.CreatedAt,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("erro ao serializar notificação: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", n.webhookURL, bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("erro ao criar request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	if n.authToken != "" {
		req.Header.Set("Authorization", "Bearer "+n.authToken)
	}

	resp, err := n.client.Do(req)
	if err != nil {
		return fmt.Errorf("erro ao enviar notificação: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("erro ao enviar notificação: status %d", resp.StatusCode)
	}

	return nil
}
