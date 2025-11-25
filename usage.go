package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var usageMutex sync.Mutex

func logUsage(userID int64, tokens int) {
	usageMutex.Lock()
	defer usageMutex.Unlock()

	f, err := os.OpenFile("usage.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening usage log:", err)
		return
	}
	defer f.Close()

	entry := fmt.Sprintf("%s,%d,%d\n", time.Now().Format("2006-01-02"), userID, tokens)
	if _, err := f.WriteString(entry); err != nil {
		log.Println("Error writing to usage log:", err)
	}
}

func flushUsageRoutine() {
	ticker := time.NewTicker(24 * time.Hour)
	for range ticker.C {
		flushUsage()
	}
}

func flushUsage() {
	usageMutex.Lock()
	defer usageMutex.Unlock()

	file, err := os.Open("usage.log")
	if err != nil {
		if !os.IsNotExist(err) {
			log.Println("Error opening usage log for flushing:", err)
		}
		return
	}

	usageMap := make(map[string]map[int64]int) // date -> user -> tokens

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			continue
		}
		date := parts[0]
		userID, _ := strconv.ParseInt(parts[1], 10, 64)
		tokens, _ := strconv.Atoi(parts[2])

		if _, ok := usageMap[date]; !ok {
			usageMap[date] = make(map[int64]int)
		}
		usageMap[date][userID] += tokens
	}
	file.Close()

	for date, userUsage := range usageMap {
		for userID, tokens := range userUsage {
			if err := UpdateDailyUsage(date, userID, tokens); err != nil {
				log.Println("Error updating daily usage:", err)
			}
		}
	}

	// Clear the log file
	if err := os.Truncate("usage.log", 0); err != nil {
		log.Println("Error truncating usage log:", err)
	}
}
