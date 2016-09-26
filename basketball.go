package main

import (
	"fmt"
	"log"
	"strconv"
)

func totalPoints(ones string, twos string, threes string) string {

	x, err1 := strconv.ParseInt(ones, 0, 8)

	if err1 != nil {
		log.Println(err1)
	}

	y, err2 := strconv.ParseInt(twos, 0, 8)

	if err2 != nil {
		log.Println(err2)
	}

	z, err3 := strconv.ParseInt(threes, 0, 8)

	if err3 != nil {
		log.Println(err3)
	}

	total := x + y*2 + z*3

	return fmt.Sprintf("%d", total)

} // totalPoints

func fieldGoals(attempts string, made string) string {

	x, err1 := strconv.ParseInt(attempts, 0, 8)

	if err1 != nil {
		log.Println(err1)
	}

	y, err2 := strconv.ParseInt(made, 0, 8)

	if err2 != nil {
		log.Println(err2)
	}

	total := x + y

	return fmt.Sprintf("%d-%d", y, total)

} // fieldGoals

func rebounds(oreb string, dreb string) string {

	x, err1 := strconv.ParseInt(oreb, 0, 8)

	if err1 != nil {
		log.Println(err1)
	}

	y, err2 := strconv.ParseInt(dreb, 0, 8)

	if err2 != nil {
		log.Println(err2)
	}

	total := x + y

	return fmt.Sprintf("%d", total)

} // rebounds
