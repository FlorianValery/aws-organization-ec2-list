// Helper functions
// Includes handling error logic

package main

import (
	"log"
)

func checkError(message string, err error) {
	if err != nil {
		log.Fatal(message, err)
	}
}
