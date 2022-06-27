package cmd

import (
	"encoding/json"
	"io/ioutil"

	"github.com/go-resty/resty/v2"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const importDataPath = "/api/open/import_data"

var yapiCmd = &cobra.Command{
	Use:   "yapi",
	Short: "yapi cli tool",
	Run:   yapi,
	PreRun: func(cmd *cobra.Command, args []string) {
		viper.SetConfigName("yapi")
	},
}

var yapiProjectName string
var yapiUpload bool
var yapiPrintConfig bool
var yapiMerge string

func init() {
	yapiCmd.Flags().StringVarP(&yapiProjectName, "projectname", "n", "", "project name")
	yapiCmd.Flags().BoolVar(&yapiUpload, "upload", false, "upload swagger file to yapi")
	yapiCmd.Flags().BoolVar(&yapiPrintConfig, "printconfig", false, "print config")
	yapiCmd.Flags().StringVarP(&yapiMerge, "merge", "m", "good", "sync model: normal (普通模式) , good (智能合并), merge (完全覆盖)")
	yapiCmd.MarkFlagsRequiredTogether("upload", "projectname")

	rootCmd.AddCommand(yapiCmd)
}

type YapiProjects map[string]string

type YapiConfig struct {
	Host     string       `yaml:"host" json:"host"`
	Projects YapiProjects `yaml:"projects" json:"projects"`
}

func (c *YapiConfig) loadConfig() {
	cobra.CheckErr(viper.ReadInConfig())
	cobra.CheckErr(viper.Unmarshal(c))
}

func yapi(cmd *cobra.Command, args []string) {
	var config YapiConfig

	if initCfg {
		config.Host = "http://localhost.yapi.com"
		config.Projects = YapiProjects{
			"project1": "example0269f91f0dab3b6624f33873f9164a2d81efaf24dcfdf1cc50d4f8f5f",
		}
		err := createCfgIFNotExists(config)
		cobra.CheckErr(err)
		return
	}

	config.loadConfig()

	if yapiPrintConfig {
		bytes, err := json.MarshalIndent(config, "", "\t")
		cobra.CheckErr(err)
		cmd.Print(string(bytes))
		return
	}

	if yapiUpload {
		file, err := ioutil.ReadFile("./docs/swagger.json")
		cobra.CheckErr(err)

		cmd.Println("upload merge: ", yapiMerge)

		resp, err := resty.New().SetBaseURL(config.Host).R().
			SetBody(map[string]any{
				"type":  "swagger",
				"json":  string(file),
				"merge": yapiMerge,
				"token": config.Projects[yapiProjectName],
			}).Post(importDataPath)
		cobra.CheckErr(err)

		cmd.Println(resp.String())
	}
}
