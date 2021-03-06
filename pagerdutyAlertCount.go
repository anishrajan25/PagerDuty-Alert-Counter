package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

func getFileName() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Please provide the csv file name (incidents.csv): ")

	fileName, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("Failed while taking input: \n%v\n\n\n", err)
		os.Exit(1)
	}

	fileName = strings.Trim(fileName, " \n")
	if len(fileName) == 0 || !strings.HasSuffix(fileName, ".csv") {
		fileName = "incidents.csv"
	}

	return fileName
}

func getAlertData(fileName string) (map[string]int, int) {
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Printf("Failed reading data from file: \n%v\n\n\n", err)
		os.Exit(1)
	}
	defer file.Close()

	rows, err := csv.NewReader(file).ReadAll()

	return countDuplicates(rows), len(rows) - 1
}

func getIndexOf(list []string, key string) int {
	index := -1
	for k, v := range list {
		if v == key {
			index = k
			break
		}
	}

	return index
}

func countDuplicates(rows [][]string) map[string]int {
	index := getIndexOf(rows[0], "description")

	if index == -1 {
		fmt.Print("\n\nInvalid CSV\n\n\n")
		os.Exit(1)
	}

	alerts := map[string]int{}

	for i := 1; i < len(rows); i += 1 {
		name := rows[i][index]
		name = name[len("[FIRING:1]") : 1+strings.LastIndex(name, "]")]
		if c, ok := alerts[name]; ok {
			alerts[name] = c + 1
		} else {
			alerts[name] = 0
		}
	}

	return alerts
}

func printAlertData(alerts map[string]int, totalAlerts int) {
	fmt.Println("\n\n\nTotal PagerDuty Alerts: ", totalAlerts, "\nTotal Unique PagerDuty Alerts: ", len(alerts), "\n\n\n", "Count of duplicates of each alert: \n")
	for name, count := range alerts {
		if count == 0 {
			continue
		}
		fmt.Printf("%v - %v\n", name, count)
	}
}

func deleteFileIfApplicable(fileName string) {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("\n\nDo you want to delete the CSV file (Y/N) ?: ")

	choice, err := reader.ReadString('\n')

	if err != nil {
		fmt.Printf("\n\nFailed while taking input: \n%v\n\n\n", err)
		os.Exit(1)
	}

	choice = strings.ToUpper(strings.Trim(choice, " \n"))
	if choice == "Y" || choice == "YES" {
		deleteFile(fileName)
	}
}

func deleteFile(fileName string) {
	err := os.Remove(fileName)

    if err != nil {
        fmt.Printf("\n\nFailed while deleting the file: \n%v\n\n\n", err)
		os.Exit(1)
    }

	fmt.Println("File deleted successfully")
}

func main() {
	fileName := getFileName()
	alerts, totalAlerts := getAlertData(fileName)
	printAlertData(alerts, totalAlerts)
	deleteFileIfApplicable(fileName)
}
