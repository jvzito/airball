package models

type Player struct {
	ID       int    `db:"id"       json:"id"`
	Name     string `db:"name"     json:"name"`
	Team     string `db:"team"     json:"team"`
	Position string `db:"position" json:"position"`
}

type ShotChart struct {
	GameID     string  `json:"game_id"`
	Period     int     `json:"period"`
	LocationX  float64 `json:"location_x"`
	LocationY  float64 `json:"location_y"`
	ShotType   string  `json:"shot_type"`
	ShotZone   string  `json:"shot_zone"`
	Made       bool    `json:"made"`
	Distance   int     `json:"distance"`
	ActionType string  `json:"action_type"`
}

type PlayerStats struct {
	PlayerID int     `json:"player_id"`
	Season   string  `json:"season"`
	GP       int     `json:"gp"`
	PPG      float64 `json:"ppg"`
	APG      float64 `json:"apg"`
	RPG      float64 `json:"rpg"`
	FGPct    float64 `json:"fg_pct"`
	ThreePct float64 `json:"three_pct"`
	Minutes  float64 `json:"minutes"`
}

type LeaderEntry struct {
	Rank       int     `json:"rank"`
	PlayerID   int     `json:"player_id"`
	PlayerName string  `json:"player_name"`
	Team       string  `json:"team"`
	Value      float64 `json:"value"`
}
