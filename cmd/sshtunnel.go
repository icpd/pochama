package cmd

import (
	"bytes"
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/whoisix/pochama/sshtunnel"
	"golang.org/x/crypto/ssh"
	"gopkg.in/yaml.v2"
)

var sshtunnelCmd = &cobra.Command{
	Use: "sshtunnel",
	Run: cmd,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("sshtunnel")
	},
}

func init() {
	rootCmd.AddCommand(sshtunnelCmd)
}

type Tunnel struct {
	Local      string `yaml:"local"`
	Remote     string `yaml:"remote"`
	Server     string `yaml:"server"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	PrivateKey string `yaml:"privateKey"`
}

type TunnelConfig struct {
	Server     string   `yaml:"server"`
	User       string   `yaml:"user"`
	Password   string   `yaml:"password"`
	PrivateKey string   `yaml:"privateKey"`
	Tunnels    []Tunnel `yaml:"tunnels"`
}

func cmd(cmd *cobra.Command, args []string) {
	if initCfg {
		err := createSSHTunnelCfgIFNotExists()
		cobra.CheckErr(err)
		return
	}

	var config TunnelConfig
	config.loadConfig()

	var globalAuth ssh.AuthMethod
	switch {
	case config.PrivateKey != "":
		globalAuth = ssh.PublicKeys(signer(config.PrivateKey))
	case config.Password != "":
		globalAuth = ssh.Password(config.Password)
	}

	var tunnels []*sshtunnel.SSHTunnel
	for _, t := range config.Tunnels {
		user := config.User
		if t.User != "" {
			user = t.User
		}

		auth := t.getAuth(globalAuth)
		if auth == (ssh.AuthMethod)(nil) {
			cobra.CheckErr("no auth method")
		}

		server := config.Server
		if t.Server != "" {
			server = t.Server
		}

		if server == "" {
			cobra.CheckErr("no server")
		}

		tunnel := sshtunnel.NewSSHTunnel(
			sshtunnel.WithLocal(t.Local),
			sshtunnel.WithServer(server),
			sshtunnel.WithRemote(t.Remote),
			sshtunnel.WithAuth(user, auth),
		)

		tunnels = append(tunnels, tunnel)
	}

	for _, t := range tunnels {
		go func(t *sshtunnel.SSHTunnel) {
			err := t.Start()
			cobra.CheckErr(err)
		}(t)
	}

	select {}
}

func (t Tunnel) getAuth(defaultAuth ssh.AuthMethod) ssh.AuthMethod {
	var auth ssh.AuthMethod
	switch {
	case t.PrivateKey != "":
		auth = ssh.PublicKeys(signer(t.PrivateKey))
	case t.Password != "":
		auth = ssh.Password(t.Password)
	default:
		auth = defaultAuth
	}

	return auth
}

func createSSHTunnelCfgIFNotExists() error {
	emptyCfg := &TunnelConfig{
		Tunnels: []Tunnel{
			{},
		},
	}
	yamlCfg, err := yaml.Marshal(emptyCfg)
	if err != nil {
		return err
	}

	err = viper.ReadConfig(bytes.NewBuffer(yamlCfg))
	if err != nil {
		return err
	}
	err = viper.SafeWriteConfig()
	if err != nil {
		return err
	}
	return nil
}

func (t *TunnelConfig) loadConfig() {
	err := viper.ReadInConfig()
	cobra.CheckErr(err)

	err = viper.Unmarshal(t)
	cobra.CheckErr(err)
}

func signer(privateKeyPath string) ssh.Signer {
	key, err := ioutil.ReadFile(privateKeyPath)
	cobra.CheckErr(err)

	gSigner, err := ssh.ParsePrivateKey(key)
	cobra.CheckErr(err)

	return gSigner
}
