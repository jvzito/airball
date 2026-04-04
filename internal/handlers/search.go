package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type PlayerSearchResult struct {
	PlayerID   int    `json:"player_id"`
	PlayerName string `json:"player_name"`
	Team       string `json:"team"`
	Position   string `json:"position"`
}

type SearchHandler struct{}

func NewSearchHandler() *SearchHandler { return &SearchHandler{} }

// GET /api/v1/players/search?q=curry
func (h *SearchHandler) Search(c *gin.Context) {
	query := strings.TrimSpace(c.Query("q"))
	if len(query) < 2 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "digite ao menos 2 caracteres"})
		return
	}

	players, err := fetchNBAPlayers(query)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": "falha ao buscar jogadores"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": players})
}

func fetchNBAPlayers(query string) ([]PlayerSearchResult, error) {
	url := "https://stats.nba.com/stats/commonallplayers?LeagueID=00&Season=2025-26&IsOnlyCurrentSeason=1"

	req, err := http.NewRequest(http.MethodGet, url, nil)
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

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("nba api status %d", resp.StatusCode)
	}

	var result struct {
		ResultSets []struct {
			Headers []string        `json:"headers"`
			RowSet  [][]interface{} `json:"rowSet"`
		} `json:"resultSets"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.ResultSets) == 0 {
		return []PlayerSearchResult{}, nil
	}

	rs := result.ResultSets[0]
	idx := make(map[string]int)
	for i, h := range rs.Headers {
		idx[h] = i
	}

	q := strings.ToLower(query)
	players := []PlayerSearchResult{}

	for _, row := range rs.RowSet {
		name := fmt.Sprintf("%v", row[idx["DISPLAY_FIRST_LAST"]])
		if !strings.Contains(strings.ToLower(name), q) {
			continue
		}
		players = append(players, PlayerSearchResult{
			PlayerID:   int(row[idx["PERSON_ID"]].(float64)),
			PlayerName: name,
			Team:       fmt.Sprintf("%v", row[idx["TEAM_ABBREVIATION"]]),
		})
		if len(players) >= 10 {
			break
		}
	}

	return players, nil
}
