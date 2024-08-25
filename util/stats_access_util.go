package util

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var statsDir string = "./data/stats"

type AccessStats struct {
	AllTime    int32
	ThisYear   int32
	L1Year     int32
	L6Month    int32
	L1Month    int32
	L7Days     int32
	L3Days     int32
	L1Days     int32
	L12Hours   int32
	L6Hours    int32
	L1Hours    int32
	L30Minutes int32
	L15Minutes int32
}

var CurrentAccessStats *AccessStats

var AggregationTime int64 = 60 * 15

var (
	accessLogFile   *os.File
	logChannel      chan string
	wg              sync.WaitGroup
	nextAggregation time.Time
)

func init() {
	CurrentAccessStats = &AccessStats{}
	createAccessStats()
	logChannel = make(chan string, 2500)
	wg.Add(1)
	go logWriter()
}

func logWriter() {
	defer wg.Done()
	for entry := range logChannel {
		if accessLogFile == nil {
			err := os.MkdirAll(statsDir, 0755)
			if err != nil {
				fmt.Println("Error creating directory:", err)
				return
			}
			openAccessLogFile()
			nextAggregation = time.UnixMilli(0)
		}
		if time.Now().After(nextAggregation) {
			aggregateLogs()
			nextAggregation = time.Now().Add(time.Duration(AggregationTime) * time.Second)
		}
		_, err := accessLogFile.WriteString(entry)
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
}
func aggregateLogs() {
	if accessLogFile != nil {
		accessLogFile.Close()
	}
	readFile, err := os.Open(statsDir + "/access.log")
	if err != nil {
		fmt.Println("Error opening file for reading:", err)
		return
	}
	file, err := os.OpenFile(statsDir+"/aggregated.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	// the first log entry
	var firstTime int64 = 0
	// the last log entry
	var lastTime int64 = 0
	scanner := bufio.NewScanner(readFile)
	var count = 0
	for scanner.Scan() {
		line := scanner.Text()
		currentTime, _ := strconv.ParseInt(line, 10, 64)
		if firstTime == 0 || currentTime > firstTime {
			if firstTime != 0 {
				_, err := file.WriteString(fmt.Sprintf("%d,%d,%d\n", firstTime, lastTime, count))
				if err != nil {
					fmt.Println("Error writing to file:", err)
					return
				}
			}
			firstTime = currentTime + AggregationTime
			count = 0
		}
		lastTime = currentTime
		count++
	}
	if firstTime != 0 {
		_, err = file.WriteString(fmt.Sprintf("%d,%d,%d\n", firstTime, lastTime, count))
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}
	}
	file.Close()
	readFile.Close()
	os.Remove(statsDir + "/access.log")

	createAccessStats()
	openAccessLogFile()
}
func createAccessStats() {
	CurrentAccessStats = &AccessStats{}
	now := time.Now().Unix()
	yearStart := time.Date(time.Now().Year(), time.January, 1, 0, 0, 0, 0, time.UTC).Unix()
	oneYearAgo := now - 365*24*3600
	sixMonthsAgo := now - 180*24*3600
	oneMonthAgo := now - 30*24*3600
	sevenDaysAgo := now - 7*24*3600
	threeDaysAgo := now - 3*24*3600
	oneDayAgo := now - 24*3600
	twelveHoursAgo := now - 12*3600
	sixHoursAgo := now - 6*3600
	oneHourAgo := now - 3600
	thirtyMinutesAgo := now - 30*60
	fifteenMinutesAgo := now - 15*60
	file, err := os.Open(statsDir + "/aggregated.log")
	if err != nil {
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, ",")
		if len(parts) != 3 {
			return
		}
		end, err := strconv.ParseInt(parts[1], 10, 64)
		if err != nil {
			return
		}
		count, err := strconv.ParseInt(parts[2], 10, 64)
		if err != nil {
			return
		}
		if end >= fifteenMinutesAgo {
			CurrentAccessStats.L15Minutes += int32(count)
		}
		if end >= thirtyMinutesAgo {
			CurrentAccessStats.L30Minutes += int32(count)
		}
		if end >= oneHourAgo {
			CurrentAccessStats.L1Hours += int32(count)
		}
		if end >= sixHoursAgo {
			CurrentAccessStats.L6Hours += int32(count)
		}
		if end >= twelveHoursAgo {
			CurrentAccessStats.L12Hours += int32(count)
		}
		if end >= oneDayAgo {
			CurrentAccessStats.L1Days += int32(count)
		}
		if end >= threeDaysAgo {
			CurrentAccessStats.L3Days += int32(count)
		}
		if end >= sevenDaysAgo {
			CurrentAccessStats.L7Days += int32(count)
		}
		if end >= oneMonthAgo {
			CurrentAccessStats.L1Month += int32(count)
		}
		if end >= sixMonthsAgo {
			CurrentAccessStats.L6Month += int32(count)
		}
		if end >= oneYearAgo {
			CurrentAccessStats.L1Year += int32(count)
		}
		if end >= yearStart {
			CurrentAccessStats.ThisYear += int32(count)
		}
		CurrentAccessStats.AllTime += int32(count)
	}

}
func openAccessLogFile() {
	file, err := os.OpenFile(statsDir+"/access.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening or creating file:", err)
		return
	}
	accessLogFile = file
}
func LogAccess(uuid string) {
	//entry := fmt.Sprintf("%s,%d\n", uuid, time.Now().Unix())
	entry := fmt.Sprintf("%d\n", time.Now().Unix())
	logChannel <- entry
}

func CloseLog() {
	close(logChannel)
	wg.Wait()

	if accessLogFile != nil {
		accessLogFile.Close()
	}
}
