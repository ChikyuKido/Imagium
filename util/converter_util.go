package util

import "fmt"

func FormatBytesToString(bytes int64) string {
	prefixes := []string{"b", "kb", "mb", "gb", "tb"}
	var currentValue = float64(bytes)
	iterations := 0
	for currentValue >= 1024 && iterations < len(prefixes)-1 {
		currentValue = currentValue / 1024
		iterations++
	}

	return fmt.Sprintf("%.2f%s", currentValue, prefixes[iterations])
}
