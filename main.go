package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/nabinthapaa/vantage/constants"
	"github.com/nabinthapaa/vantage/table"
)

func main() {
	user := os.Getenv("USER")

	if len(os.Args) < 2 {
		log.Fatalf("Please provide 'on' 'off' or 'table' as a flag.")
		os.Exit(1)
	}

	if user != "root" {
		log.Fatalf("Permission denied. Run sudo %s", strings.Join(os.Args, " "))
		os.Exit(1)
	}

	flag := os.Args[1]

	if flag == "table" {
		if err := table.Run(); err != nil {
			fmt.Printf("err %v", err)
			os.Exit(1)
		}
	} else {
		file, err := os.OpenFile(constants.CONSERVATION_MODE_FILE, os.O_RDWR, 0644)
		if err != nil {
			log.Fatalf("Failed to open conservation_mode file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		if !scanner.Scan() {
			log.Fatal("Failed to read from conservation_mode file")
		}
		fileContent := scanner.Text()

		mode, err := strconv.Atoi(fileContent)
		if err != nil {
			log.Fatalf("Error parsing mode: %v", err)
		}

		switch flag {
		case "on":
			if mode == 1 {
				fmt.Println("Conservation mode is already on.")
				return
			}
			fmt.Println("Turning on conservation mode...")
			if _, err := file.WriteString("1"); err != nil {
				log.Fatalf("Failed to write to conservation_mode file: %v", err)
			}
			fmt.Println("Conservation mode enabled.")

		case "off":
			if mode == 0 {
				fmt.Println("Conservation mode is already off.")
				return
			}
			fmt.Println("Turning off conservation mode...")
			if _, err := file.WriteString("0"); err != nil {
				log.Fatalf("Failed to write to conservation_mode file: %v", err)
			}
			fmt.Println("Conservation mode disabled.")

		default:
			fmt.Println("Invalid flag. Use 'on' or 'off'.")
			os.Exit(1)
		}

	}
}
