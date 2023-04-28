package types

type ProbeNode struct {
	Name string                 `json:"name"`
	Info map[string]interface{} `json:"info"`
}

type ProbeEdge struct {
	Src    string `json:"src"`
	Dst    string `json:"dst"`
	Status string `json:"status"`
}

type ProbeResponse struct {
	Traceid string       `json:"traceid"`
	Nodes   []*ProbeNode `json:"nodes"`
	Edges   []*ProbeEdge `json:"edges"`
}
