package proxy

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/kurogane/gateway/internal/config"
	"github.com/kurogane/gateway/internal/model"
)

type EngineProxy struct {
	client    *http.Client
	engineURL string
}

func New(cfg *config.Config) *EngineProxy {
	return &EngineProxy{
		client: &http.Client{
			Timeout: 120 * time.Second,
		},
		engineURL: cfg.EngineURL,
	}
}

func (p *EngineProxy) CreateRun(req *model.CreateRunRequest) (*model.RunResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	resp, err := p.client.Post(
		p.engineURL+"/runs",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return nil, fmt.Errorf("engine request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		var errResp model.ErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("engine returned status %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("engine error: %s", errResp.Error)
	}

	var result model.RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (p *EngineProxy) GetRun(id string) (*model.RunResponse, error) {
	resp, err := p.client.Get(p.engineURL + "/runs/" + id)
	if err != nil {
		return nil, fmt.Errorf("engine request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("run not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("engine returned status %d", resp.StatusCode)
	}

	var result model.RunResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return &result, nil
}

func (p *EngineProxy) ListRuns() ([]model.RunSummary, error) {
	resp, err := p.client.Get(p.engineURL + "/runs")
	if err != nil {
		return nil, fmt.Errorf("engine request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("engine returned status %d", resp.StatusCode)
	}

	var result []model.RunSummary
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}
	return result, nil
}

type SSEEvent struct {
	EventType string
	Data      string
}

func (p *EngineProxy) StreamRunEvents(id string) (<-chan SSEEvent, error) {
	req, err := http.NewRequest("GET", p.engineURL+"/runs/"+id+"/events", nil)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Accept", "text/event-stream")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("engine request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("engine returned status %d", resp.StatusCode)
	}

	ch := make(chan SSEEvent, 64)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		var eventType, data string

		for scanner.Scan() {
			line := scanner.Text()
			if len(line) == 0 {
				// Empty line = end of event
				if eventType != "" || data != "" {
					ch <- SSEEvent{EventType: eventType, Data: data}
					eventType = ""
					data = ""
				}
				continue
			}
			if len(line) > 6 && line[:6] == "event:" {
				eventType = line[7:]
			} else if len(line) > 5 && line[:5] == "data:" {
				if data != "" {
					data += "\n"
				}
				data += line[6:]
			}
		}
	}()

	return ch, nil
}