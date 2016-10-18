package main

import (
	"fmt"
	"log"
	"strconv"

)

func calcStatAvg(stat string, gp string) string {

  res1, err1 := strconv.ParseInt(stat, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return "0.0"
	}

  res2, err2 := strconv.ParseInt(gp, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return "0.0"
	}

	if res2 != 0 {
    return fmt.Sprintf("%.1f", float64(res1)/float64(res2))
  } else {
    return "0.0"
  }

} // calcStatAvg

func calcPtsAvg(ftm string, fg2m string, fg3m string, gp string) string {

  res1, err1 := strconv.ParseInt(ftm, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return "0.0"
	}

  res2, err2 := strconv.ParseInt(fg2m, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return "0.0"
	}

  res3, err3 := strconv.ParseInt(fg3m, 10, 64)

	if err3 != nil {
		log.Println(err3)
    return "0.0"
	}

  res4, err4 := strconv.ParseInt(gp, 10, 64)

	if err4 != nil {
		log.Println(err4)
    return "0.0"
	}

  if res4 != 0 {
    return fmt.Sprintf("%.1f", float64(res1+res2*2+res3*3)/float64(res4))
  } else {
    return "0.0"
  }

} // calcPtsAvg

func calcRebAvg(oreb string, dreb string, gp string) string {

  res1, err1 := strconv.ParseInt(oreb, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return "0.0"
	}

  res2, err2 := strconv.ParseInt(dreb, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return "0.0"
	}

  res3, err3 := strconv.ParseInt(gp, 10, 64)

	if err3 != nil {
		log.Println(err3)
    return "0.0"
	}

  if res3 != 0 {
    return fmt.Sprintf("%.1f", float64(res1+res2)/float64(res3))
  } else {
    return "0.0"
  }

} // calcRebAvg

func calcPctAvg(missed string, made string) string {

  res1, err1 := strconv.ParseInt(missed, 10, 64)

	if err1 != nil {
		log.Println(err1)
    return "0.0"
	}

  res2, err2 := strconv.ParseInt(made, 10, 64)

	if err2 != nil {
		log.Println(err2)
    return "0.0"
	}

  if res1 != 0 && res2 != 0 {
     return fmt.Sprintf("%.1f", (float64(res2)/float64(res1+res2))*100.0)
  } else {
    return "0.0"
  }

} // calcPctAvg

