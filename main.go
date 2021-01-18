// Main cal

package main

import (
	"./config"
)

func main() {
	// Retrieve config file
	config := config.InitVariables()
	// Call main function
	get_organization_ec2(config)
}
