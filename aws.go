// aws.go includes all the functions that make AWS API calls

package main

import (
	"./config"
	"encoding/csv"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/organizations"
	"os"
)

func get_organization_accounts(config config.Config) map[string]string {
	// Start AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)
	checkError("Could not start session", err)
	// Create organization service client
	c := organizations.New(sess)
	// Create variable for the list of accounts and initialize input
	list_organization_accounts := make(map[string]string)

	input := &organizations.ListAccountsInput{}
	// Start a do-while loop
	for {
		// Retrieve the accounts with a limit of 20 per call
		list_organization_paginated, err := c.ListAccounts(input)
		// Append the accounts from the current call to the total list
		for _, account := range list_organization_paginated.Accounts {
			list_organization_accounts[*account.Name] = *account.Id
		}
		checkError("Could not retrieve account list", err)
		// Check if more accounts need to be retrieved, otherwise break the loop
		if list_organization_paginated.NextToken == nil {
			break
		} else {
			input = &organizations.ListAccountsInput{NextToken: list_organization_paginated.NextToken}
		}
	}
	return list_organization_accounts
}

func get_account_ec2(config config.Config, account_name string, account_id string, result map[string][]string) map[string][]string {
	// Create EC2 service client
	var c Clients
	svc := c.EC2(config.Region, account_id, config.Organization_Role)
	// Get the EC2 list of the given account
	input := &ec2.DescribeInstancesInput{}
	list_instances, err := svc.DescribeInstances(input)
	checkError("Could not retrieve the EC2s", err)

	// Iterate over the EC2 instances and add elements to global list, if instances > 0
	if len(list_instances.Reservations) != 0 {
		for _, reservation := range list_instances.Reservations {
			// Loop through every individual EC2 instance
			for _, instance := range reservation.Instances {
				// Set the map key using the unique instance ID
				key := *instance.InstanceId
				// Retrieve account information
				result[key] = append(result[key], account_name)
				result[key] = append(result[key], account_id)
				// Check if the instance name is set using tags, otherwise use default
				for _, tag := range instance.Tags {
					if *tag.Key == "Name" {
						result[key] = append(result[key], *tag.Value)

					}
				}
				if len(result) == 2 {
					result[key] = append(result[key], "N/A")
				}
				// Retrieve instance information, use default is potentially null
				result[key] = append(result[key], *instance.InstanceType)
				result[key] = append(result[key], *instance.InstanceId)
				result[key] = append(result[key], *instance.ImageId)
				if instance.Platform != nil {
					result[key] = append(result[key], *instance.Platform)
				} else {
					result[key] = append(result[key], "linux")
				}
				if instance.PrivateIpAddress != nil {
					result[key] = append(result[key], *instance.PrivateIpAddress)
				} else {
					result[key] = append(result[key], "N/A")
				}
				result[key] = append(result[key], *instance.State.Name)
				result[key] = append(result[key], (*instance.LaunchTime).String())
			}
		}
	}
	fmt.Println("Account number " + account_id + " done")
	return result
}

func get_organization_ec2(config config.Config) {
	// Retrieve all accounts from the organization
	list_accounts := get_organization_accounts(config)
	// Create list variable to store every ec2 instances
	var list_ec2 = make(map[string][]string)
	// Loop over each account and get its instances via a function
	fmt.Println("Retrieving the instances...")
	for account_name, account_id := range list_accounts {
		list_ec2 = get_account_ec2(config, account_name, account_id, list_ec2)
	}
	fmt.Println("All the instances from the Organization were retrieved.")
	// Create the csv file using the os package
	fmt.Println("Creating a CSV file...")
	file, err := os.Create("result.csv")
	checkError("Cannot create file", err)
	defer file.Close()
	// Create the writer object
	writer := csv.NewWriter(file)
	defer writer.Flush()
	// Write headers
	var headers = []string{"Account Name", "Account ID", "Instance Name", "Instance Size", "Instance ID", "Image ID", "Platform", "Private IP", "State", "Timestamp"}
	writer.Write(headers)
	// Loop over the organization ec2 list and write them in rows in the csv file
	for _, value := range list_ec2 {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
	fmt.Println("CSV file created in " + "result.csv")
}
