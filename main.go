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
	blocks := TransformClassesIntoBlocks(classes)
	for _, b := range blocks {
		fmt.Printf("Day %s | Block %2d | Class: %-8s | Room: %-6s | Teacher: %-18s | ID: %s\n",
			b.Classes[0].Day, b.Classes[0].Start, b.Classes[0].Name, b.Classes[0].Room, b.Classes[0].Teacher, b.Classes[0].ID)
	}
}
