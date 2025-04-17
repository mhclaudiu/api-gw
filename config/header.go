package config

type CFGxAPI struct {
	Host       string `json:"Host"`
	Port       int    `json:"Port"`
	CorsFilter string `json:"CorsFilter"`
	Path       string `json:"Path"`
	RateLimit  string `json:"RateLimit"`
	Auth       bool   `json:"Authorization"`
}

type CFG struct {
	APP CFGxAPP `json:"APP"`
	API CFGxAPI `json:"API"`
	LOG CFGxLOG `json:"LOG"`
}

type CFGxAPP struct {
	Name          string `json:"Name"`
	Env           string `json:"Environment"`
	ExitSyncTimer string `json:"ExitSyncTimer"`
}

type CFGxLOG struct {
	Enabled bool   `json:"Enabled"`
	Dir     string `json:"Directory"`
	MaxSize int    `json:"MaxSize"`
	MaxDays int    `json:"MaxDays"`
}

type CFGxAPI_Ratelimit struct {
	BurstRate int
	Seconds   int
}
