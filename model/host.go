package model

// Host 主机信息
type Host struct {
	IP      string `json:"IP"`
	MAC     string `json:"MAC"`
	MACInfo string `json:"MACInfo"`
}
