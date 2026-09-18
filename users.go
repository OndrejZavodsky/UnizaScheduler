package main

import (
	"flag"
	"strings"
)

type Faculty int

// pridaj dalsie fakulty a uprav mapovanie podla api vzdelavanie
const (
	FRI Faculty = iota
	PEDAS
	FBI
	UnknownFaculty = -1
)

var facultyMap = map[string]Faculty{
	"fri":   FRI,
	"pedas": PEDAS,
	"fbi":   FBI,
}

type User struct {
	Group   string
	Faculty Faculty
}

func ParseUserFlags() User {
	var groupFlag string
	var facultyFlag string

	flag.StringVar(&groupFlag, "skupina", "", "Názov skupiny (napr. 5ZY011)")
	flag.StringVar(&facultyFlag, "fakulta", "", "Názov fakulty pis skratkov podla vzdelavania")

	flag.Parse()

	facultyEnum, exists := facultyMap[strings.ToLower(facultyFlag)]
	if !exists {
		facultyEnum = UnknownFaculty
	}

	return User{
		Group:   groupFlag,
		Faculty: facultyEnum,
	}
}
