package tvpn

import (
	"os"
	"encoding/json"
)

type Config struct {
	Name, Group string
	Friends map[string]Friend
	IPMan IPConfig
	Sig   SigConfig
	Stun  StunConfig
	VPN   VPNConfig
}


func ReadConfig(path string) (*Config,error) {
	file,err := os.Open(path)

	if err != nil {
		return nil,err
	}

	var config Config

	data := make([]byte,1024)

	n,err := file.Read(data)
	if err != nil {
		return nil,err
	}

	err = json.Unmarshal(data[:n],&config)

	if err == nil {
		return &config,nil
	}
	return nil,err
}

