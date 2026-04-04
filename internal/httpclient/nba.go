package httpclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/jvzito/airball/internal/config"
)

type NBAResponse struct {
	ResultSets []struct {
		Name    string          `json:"name"`
		Headers []string        `json:"headers"`
		RowSet  [][]interface{} `json:"rowSet"`
	} `json:"resultSets"`
}

type NBAClient struct {
	client  *http.Client
	baseURL string
}

func NewNBAClient(cfg *config.Config) *NBAClient {
	return &NBAClient{
		baseURL: cfg.NBABaseURL,
		client:  &http.Client{Timeout: time.Duration(cfg.NBATimeout) * time.Second},
	}
}

func (c *NBAClient) get(ctx context.Context, endpoint string, params map[string]string) (*NBAResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+endpoint, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64)")
	req.Header.Set("Referer", "https://www.nba.com/")
	req.Header.Set("Origin", "https://www.nba.com")
	req.Header.Set("Accept", "application/json, text/plain, */*")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("x-nba-stats-origin", "stats")
	req.Header.Set("x-nba-stats-token", "true")

	q := req.URL.Query()
	for k, v := range params {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("nba request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nba api %d for %s", resp.StatusCode, endpoint)
	}

	var result NBAResponse
	return &result, json.NewDecoder(resp.Body).Decode(&result)
}

func (c *NBAClient) GetLeaders(ctx context.Context, season, category string) (*NBAResponse, error) {
	return c.get(ctx, "/leagueleaders", map[string]string{
		"LeagueID": "00", "PerMode": "PerGame", "Scope": "S",
		"Season": season, "SeasonType": "Regular Season", "StatCategory": category,
	})
}

func (c *NBAClient) GetShotChart(ctx context.Context, playerID, season string) (*NBAResponse, error) {
	return c.get(ctx, "/shotchartdetail", map[string]string{
		"LeagueID": "00", "Season": season, "SeasonType": "Regular Season",
		"TeamID": "0", "PlayerID": playerID, "GameID": "",
		"ContextMeasure": "FGA", "Period": "0", "LastNGames": "0",
		"Month": "0", "OpponentTeamID": "0", "RangeType": "0",
		"StartPeriod": "1", "EndPeriod": "10",
		"StartRange": "0", "EndRange": "28800",
	})
}

func (c *NBAClient) GetPlayerStats(ctx context.Context, playerID, season string) (*NBAResponse, error) {
	return c.get(ctx, "/playerdashboardbyyearoveryear", map[string]string{
		"LeagueID": "00", "MeasureType": "Base", "PerMode": "PerGame",
		"Season": season, "SeasonType": "Regular Season",
		"PlayerID": playerID, "Month": "0", "OpponentTeamID": "0",
		"Period": "0", "LastNGames": "0",
	})
}

func HeaderIndex(headers []string) map[string]int {
	m := make(map[string]int, len(headers))
	for i, h := range headers {
		m[h] = i
	}
	return m
}
