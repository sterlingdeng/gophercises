package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

func main() {
	// set flags and their default values
	csvFileName := flag.String("csv", "problems.csv", "csv file of question and answers")
	timeLimit := flag.Int("time", 30, "time limit of the quiz in seconds")
	flag.Parse()

	file, err := os.Open(*csvFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	problems := parseRecords(records)

	correct := 0
	timer := time.NewTimer(time.Second * time.Duration((*timeLimit)))

	for _, p := range problems {
		answerCh := make(chan string)
		go func() {
			fmt.Printf("What is %s?\n", p.question)
			var ans string
			fmt.Scanf("%s\n", &ans)
			answerCh <- ans
		}()

		select {
		case <-timer.C:
			fmt.Println("Ran out of time")
			fmt.Printf("You got %d correct and %d incorrect!\n", correct, len(problems)-correct)
			return
		case answer := <-answerCh:
			if answer == p.answer {
				correct++
			}
		}
	}
	fmt.Printf("You got %d correct and %d incorrect!\n", correct, len(problems)-correct)
}

type problem struct {
	question string
	answer   string
}

func parseRecords(records [][]string) []problem {
	out := make([]problem, len(records))
	for i, record := range records {
		p := problem{
			question: record[0],
			answer:   record[1],
		}
		out[i] = p
	}
	return out
}
