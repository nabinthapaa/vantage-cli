package fn_lock

import (
	"bufio"
	"fmt"
	"os"
	"strconv"

	c "github.com/nabinthapaa/vantage-cli/constants"
)

const (
	On  int = 1
	Off int = 0
)

var ValueMap = map[string]int{
	"On":  On,
	"Off": Off,
}

func GetCurrentValue() (string, error) {
	file, err := os.OpenFile(c.FN_LOCK, os.O_RDONLY, 0644)
	if err != nil {
		return "", fmt.Errorf("Failed to open conservation_mode file: %v", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	if !scanner.Scan() {
		return "", fmt.Errorf("Failed to read conservation_mode file: %v", err)
	}
	fileContent := scanner.Text()

	mode, err := strconv.Atoi(fileContent)
	if err != nil {
		return "", fmt.Errorf("Error Parsing value %v", err)
	}

	if mode == On {
		return "On", nil
	} else if mode == Off {
		return "Off", nil
	} else {
		return "", fmt.Errorf("Invalid value ")
	}
}

func UpdateCurrentValue(value string) error {
	v, exists := ValueMap[value]
	if !exists {
		return fmt.Errorf("Invalid value passed")
	}
	file, err := os.OpenFile(c.FN_LOCK, os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("Failed to open conservation_mode file: %v", err)
	}
	defer file.Close()

	var mode string = "1"
	if v == On {
		mode = "0"
	}

	if _, err := file.WriteString(mode); err != nil {
		return fmt.Errorf("Failed to write %v", err)
	}
	return nil
}
