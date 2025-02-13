package main

import (
	"fmt"
	"time"
)

const layout = "2006-01-02T15:04:05"

func main() {
	timeString := "2024-12-12T13:30:00"

	t, err := time.Parse(layout, timeString)
	if err != nil {
		fmt.Println("Ошибка парсинга даты:", err)
		return
	}
	unixTime := t.Unix()
	fmt.Printf("Время %s в Unix формате: %d\n", timeString, unixTime)
}
