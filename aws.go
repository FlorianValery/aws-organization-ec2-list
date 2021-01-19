// aws.go includes all the functions that make AWS API calls

package main

import (
	"./config"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/organizations"
)

func getOrganizationAccounts(config config.Config) map[string]string {
	// Start AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.Region)},
	)
	checkError("Could not start session", err)
	// Create organization service client
	c := organizations.New(sess)
	// Create variable for the list of accounts and initialize input
	listOrganizationAccounts := make(map[string]string)

	input := &organizations.ListAccountsInput{}
	// Start a do-while loop
	for {
		// Retrieve the accounts with a limit of 20 per call
		listOrganizationPaginated, err := c.ListAccounts(input)
		// Append the accounts from the current call to the total list
		for _, account := range listOrganizationPaginated.Accounts {
			listOrganizationAccounts[*account.Name] = *account.Id
		}
		checkError("Could not retrieve account list", err)
		// Check if more accounts need to be retrieved, otherwise break the loop
		if listOrganizationPaginated.NextToken == nil {
			break
		} else {
			input = &organizations.ListAccountsInput{NextToken: listOrganizationPaginated.NextToken}
		}
	}
	return listOrganizationAccounts
}

func getAccountEc2(config config.Config, accountName string, accountID string, result map[string][]string) map[string][]string {
	// Create EC2 service client
	var c Clients
	svc := c.EC2(config.Region, accountID, config.Organization_Role)
	// Get the EC2 list of the given account
	input := &ec2.DescribeInstancesInput{}
	listInstances, err := svc.DescribeInstances(input)
	checkError("Could not retrieve the EC2s", err)

	// Iterate over the EC2 instances and add elements to global list, if instances > 0
	if len(listInstances.Reservations) != 0 {
		for _, reservation := range listInstances.Reservations {
			// Loop through every individual EC2 instance
			for _, instance := range reservation.Instances {
				// Set the map key using the unique instance ID
				key := *instance.InstanceId
				// Retrieve account information
				result[key] = append(result[key], accountName)
				result[key] = append(result[key], accountID)
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
	fmt.Println("Account number " + accountID + " done")
	return result
}
