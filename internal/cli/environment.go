package cli

import (
	"fmt"
	"os"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	deploymentEnvironmentsBaseDir = "deployments/environments"
)

type environmentConfiguration struct {
	Prefix       string `json:"prefix"`
	Region       string `json:"region"`
	AWSAccountID string `json:"aws_account_id"`
	DDBPrefix    string `json:"ddb_prefix"`
	StageName    string `json:"stage_name"`
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func retrieveConfig(cmd *cobra.Command) (*environmentConfiguration, error) {
	p, err := cmd.Flags().GetString("ENV")
	if err != nil {
		return nil, err
	}

	if p == "" {
		// Try to get prefix from env instead
		p = os.Getenv("ENV")
	}
	env = p

	if env == "" {
		return nil, errors.New("missing ENV variable")
	}

	// Check to see if this is running under the git repo root source
	tfConfigFileExists := fileExists(fmt.Sprintf(
		"%s/%s/config.auto.tfvars.json",
		deploymentEnvironmentsBaseDir,
		env,
	))

	v := viper.New()

	// If running under the current git repo, use the TF config.auto.tfvars file
	// Else use a configuration stored under the $HOME/.rudolph-cli or current working directory
	if tfConfigFileExists {
		v.SetConfigName("config.auto.tfvars")
		v.SetConfigType("json")
		v.AddConfigPath(fmt.Sprintf("%s/%s", deploymentEnvironmentsBaseDir, env))
	} else {
		v.SetConfigName(env)
		v.SetConfigType("json")
		v.AddConfigPath("$HOME/.rudolph-cli")
		v.AddConfigPath(".")
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			return nil, err
		} else {
			// Config file was found but another error was produced
			return nil, err
		}
	}

	config := &environmentConfiguration{}
	err = v.Unmarshal(config)

	if err != nil {
		return nil, err
	}

	cmd.Flags().Set("prefix", config.Prefix)
	cmd.Flags().Set("region", config.Region)
	cmd.Flags().Set("dynamodb_table", fmt.Sprintf("%s_rudolph_store", config.Prefix))

	return config, err

}
