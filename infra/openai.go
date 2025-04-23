package infra

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gong023/umi/domain"
)

const (
	openAIAPIBaseURL = "https://api.openai.com/v1"
	chatCompletionsEndpoint = "/chat/completions"
)

// OpenAIClient implements the domain.OpenAIClient interface
type OpenAIClient struct {
	apiKey     string
	httpClient *http.Client
	logger     domain.Logger
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(apiKey string, logger domain.Logger) *OpenAIClient {
	return &OpenAIClient{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger,
	}
}

// CreateChatCompletion sends a request to the OpenAI chat completions API
func (c *OpenAIClient) CreateChatCompletion(req *domain.ChatCompletionRequest) (*domain.ChatCompletionResponse, error) {
	c.logger.Info("Sending request to OpenAI chat completions API")
	
	// Convert the request to JSON
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal request: %v", err)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP request
	httpReq, err := http.NewRequest(
		"POST",
		openAIAPIBaseURL+chatCompletionsEndpoint,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		c.logger.Error("Failed to create HTTP request: %v", err)
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

	// Send the request
	c.logger.Info("Sending HTTP request to OpenAI API")
	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		c.logger.Error("Failed to send HTTP request: %v", err)
		return nil, fmt.Errorf("failed to send HTTP request: %w", err)
	}
	defer httpResp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		c.logger.Error("Failed to read response body: %v", err)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for non-200 status code
	if httpResp.StatusCode != http.StatusOK {
		c.logger.Error("OpenAI API returned non-200 status code: %d, body: %s", httpResp.StatusCode, string(body))
		return nil, fmt.Errorf("OpenAI API returned status code %d: %s", httpResp.StatusCode, string(body))
	}

	// Parse the response
	var response domain.ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		c.logger.Error("Failed to unmarshal response: %v", err)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	c.logger.Info("Successfully received response from OpenAI API")
	return &response, nil
}
