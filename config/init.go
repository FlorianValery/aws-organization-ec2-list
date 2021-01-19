// Initialize the variables from the json confile file
// Forked from https://github.com/tkanos/gonfig

package config

import (
	"github.com/tkanos/gonfig"
	"fmt"
)

// Config structure
type Config struct {
	Region string
	OrganizationRole string
	MasterAccountID string
}

// InitVariables based on json file
func InitVariables(input ...string) Config {
	configuration := Config{}
	fileName := fmt.Sprintf("./config/default.json")
	gonfig.GetConf(fileName, &configuration)

	return configuration
}