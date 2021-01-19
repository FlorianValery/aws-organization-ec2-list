// Main cal

package main

import (
	"./config"
	"fmt"
)

func main() {
	// Retrieve config file
	config := config.InitVariables()

	// Get all accounts names & ID from the organization
	listAccounts := getOrganizationAccounts(config)

	// Create list variable to store every ec2 instances
	var listEc2 = make(map[string][]string)

	// Loop over each account and get its instances via a function
	fmt.Println("Retrieving the instances...")
	for accountName, accountID := range listAccounts {
		listEc2 = getAccountEc2(config, accountName, accountID, listEc2)
	}
	fmt.Println("All the instances from the Organization were retrieved.")

	// Write results to a CSV file
	writeToCSV(listEc2)
}
