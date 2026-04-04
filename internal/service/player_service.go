package service

import (
	"context"
	"fmt"

	"github.com/jvzito/airball/internal/cache"
	"github.com/jvzito/airball/internal/httpclient"
	"github.com/jvzito/airball/internal/models"
	"github.com/jvzito/airball/pkg/logger"
	"go.uber.org/zap"
)

type PlayerService struct {
	nba   *httpclient.NBAClient
	cache *cache.Cache
}

func NewPlayerService(nba *httpclient.NBAClient, c *cache.Cache) *PlayerService {
	return &PlayerService{nba: nba, cache: c}
}

func (s *PlayerService) GetLeaders(ctx context.Context, season, category string) ([]models.LeaderEntry, error) {
	key := fmt.Sprintf("leaders:%s:%s", season, category)

	var cached []models.LeaderEntry
	if hit, _ := s.cache.Get(ctx, key, &cached); hit {
		return cached, nil
	}

	resp, err := s.nba.GetLeaders(ctx, season, category)
	if err != nil {
		return nil, err
	}

	leaders := []models.LeaderEntry{}
	if len(resp.ResultSets) > 0 {
		rs := resp.ResultSets[0]
		idx := httpclient.HeaderIndex(rs.Headers)
		for rank, row := range rs.RowSet {
			leaders = append(leaders, models.LeaderEntry{
				Rank:       rank + 1,
				PlayerID:   toInt(row[idx["PLAYER_ID"]]),
				PlayerName: toStr(row[idx["PLAYER_NAME"]]),
				Team:       toStr(row[idx["TEAM"]]),
				Value:      toFloat(row[idx[category]]),
			})
			if rank >= 9 {
				break
			}
		}
	}

	if err := s.cache.Set(ctx, key, leaders); err != nil {
		logger.Error("cache set", zap.Error(err))
	}
	return leaders, nil
}

func (s *PlayerService) GetShotChart(ctx context.Context, playerID, season string) ([]models.ShotChart, error) {
	key := fmt.Sprintf("shotchart:%s:%s", playerID, season)

	var cached []models.ShotChart
	if hit, _ := s.cache.Get(ctx, key, &cached); hit {
		return cached, nil
	}

	resp, err := s.nba.GetShotChart(ctx, playerID, season)
	if err != nil {
		return nil, err
	}

	shots := []models.ShotChart{}
	if len(resp.ResultSets) > 0 {
		rs := resp.ResultSets[0]
		idx := httpclient.HeaderIndex(rs.Headers)
		for _, row := range rs.RowSet {
			shots = append(shots, models.ShotChart{
				GameID:     toStr(row[idx["GAME_ID"]]),
				Period:     toInt(row[idx["PERIOD"]]),
				LocationX:  toFloat(row[idx["LOC_X"]]),
				LocationY:  toFloat(row[idx["LOC_Y"]]),
				ShotType:   toStr(row[idx["SHOT_TYPE"]]),
				ShotZone:   toStr(row[idx["SHOT_ZONE_BASIC"]]),
				Made:       toInt(row[idx["SHOT_MADE_FLAG"]]) == 1,
				Distance:   toInt(row[idx["SHOT_DISTANCE"]]),
				ActionType: toStr(row[idx["ACTION_TYPE"]]),
			})
		}
	}

	_ = s.cache.Set(ctx, key, shots)
	return shots, nil
}

func toStr(v interface{}) string {
	if v == nil { return "" }
	return fmt.Sprintf("%v", v)
}
func toFloat(v interface{}) float64 {
	if f, ok := v.(float64); ok { return f }
	return 0
}
func toInt(v interface{}) int { return int(toFloat(v)) }
