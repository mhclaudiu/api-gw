package stats

type STATSxOBJ struct {
	CPU  STATSxCPU
	MEM  STATSxMEM
	Time string
}

type INFOxOBJ struct {
	CPU  INFOxCPU
	MEM  INFOxMEM
	HOST INFOxHOST
	APP  INFOxAPP
}

type STATSxCPU struct {
	Usage string
	Load  string
}

type INFOxCPU struct {
	Cores string
	Freq  string
	Model string
	Usage string
	Load  string
}

type STATSxMEM struct {
	Usage     string
	Available string
	Cached    string
	Free      string
	Active    string
	Buffers   string
	Inactive  string
	Used      string
}

type INFOxMEM struct {
	Total     string
	Usage     string
	Available string
	Cached    string
	Free      string
	Active    string
	Buffers   string
	Inactive  string
	Used      string
}

type INFOxHOST struct {
	Name            string
	OS              string
	Platform        string
	PlatformVersion string `json:"Platform Version"`
	KernelVersion   string `json:"Kernel Version"`
	KernelArch      string `json:"Kernel Arch"`
	Uptime          string
	Proccesses      string
}

type INFOxAPP struct {
	Uptime string
}

type InfoObj INFOxOBJ

var appUptime int

type StatsObj STATSxOBJ
