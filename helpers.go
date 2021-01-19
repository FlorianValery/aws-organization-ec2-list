// Helper functions
// Includes handling error logic

package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
)

// Function that log error if not null
func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}

// Function that writes a map of slices to a CSV File
func writeToCSV(listEc2 map[string][]string) {
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
	for _, value := range listEc2 {
		err := writer.Write(value)
		checkError("Cannot write to file", err)
	}
	fmt.Println("CSV file created in " + "result.csv")
}
