package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Please provide 'on' or 'off' as a flag.")
		os.Exit(1)
	}
	flag := os.Args[1]

	file, err := os.OpenFile("/sys/bus/platform/devices/VPC2004:00/conservation_mode", os.O_RDWR, 0644)
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
