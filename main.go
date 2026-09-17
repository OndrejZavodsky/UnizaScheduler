package main

import (
	"fmt"
	"log"
)

func main() {
	groupInput := "5ZYS33"

	scheduleHTML, err := FetchScheduleBlocks(groupInput)
	if err != nil {
		log.Fatalf("Error fetching schedule: %v", err)
	}
	classes, err := ParseClasses(scheduleHTML)
	if err != nil {
		log.Fatalf("Error parsing classes: %v", err)
	}

	for _, c := range classes {
		fmt.Printf("Block %2d | Class: %-8s | Room: %-6s | Teacher: %-18s | ID: %s\n",
			c.Start, c.Name, c.Room, c.Teacher, c.Id)
	}
}
