package config

type CFGxAPI struct {
	Host        string `json:"Host"`
	Port        int    `json:"Port"`
	CorsFilter  string `json:"CorsFilter"`
	MainPath    string `json:"MainPath"`
	MetricsPath string `json:"MetricsPath"`
	RateLimit   string `json:"RateLimit"`
}

type CFG struct {
	APP CFGxAPP `json:"APP"`
	API CFGxAPI `json:"API"`
}

type CFGxAPP struct {
	Name          string `json:"Name"`
	Env           string `json:"Env"`
	ExitSyncTimer string `json:"ExitSyncTimer"`
}
