package config

import (
	// "fmt"

	"github.com/BurntSushi/toml"
)

type Cnf struct {
	SysCnf  `toml:"Sys"`
	WebCnf  `toml:"Web"`
	LogCnf  LogCnf `toml:"Log"`
	Web3Cnf Web3   `toml:"Web3"`
}

type SysCnf struct {
	Token string `toml:"token"`
}

type WebCnf struct {
	Host string `toml:"host"`
	Port string `toml:"port"`
}

type LogCnf struct {
	Output string `toml:"output"`
	Level  string `toml:"level"`
}

type Web3 struct {
	SolanaAPIDomain string `toml:"solanaAPIDomain"`
	SolanaWSDomain  string `toml:"solanaWSDomain"`
	SolanaRecipient string `toml:"solanaRecipient"`
	FromPrivateKey  string `toml:"fromPrivateKey"`
	TokenMint       string `toml:"tokenMint"`
}

var cnf Cnf

func LoadCnf(path string) error {

	_, err := toml.DecodeFile(path, &cnf)
	// fmt.Println(cnf)

	return err
}

func GetCnf() Cnf {
	return cnf
}
