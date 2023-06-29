package main

import (
	"encoding/csv" // reading csv file
	"flag"         // cli flags in commands
	"fmt"          // formattin
	"log"          // logging
	"math/rand"    // getrandom number
	"os"           // access file
	"strings"      // string formatting
	"time"
)

func main() {
	// cli flags
	csvFileName := flag.String("csv", "problems.csv", "a csv file in the format of question, answer")

	timeLimit := flag.Int("limit", 10, "time limit for the quiz in seconds")

	shuffle := flag.Bool("shuffle", false, "shuffles if the question list is set to true")

	flag.Parse()

	// reading file
	file, err := os.Open(*csvFileName)
	if err != nil {
		fmt.Printf("Failed to open csv file: %s\n", *csvFileName)
		os.Exit(1)
	}

	defer file.Close()
	_ = file

	// reading csv
	r := csv.NewReader(file)

	// gives the 2d slice
	lines, err := r.ReadAll()

	if err != nil {
		log.Fatal(err)
	}

	problems := parseLines(lines)

	if *shuffle {
		rand.Shuffle(len(problems), func(i, j int) {
			problems[i], problems[j] = problems[j], problems[i]
		})
	}

	fmt.Println("Quiz started")

	timer := time.NewTimer(time.Duration(*timeLimit) * time.Second)

	correct := 0
	for i, p := range problems {
		fmt.Printf("Problem #%d: %s = ", i+1, p.q)
		answerChan := make(chan string)

		// since scanf is blocking we will clown it and run it seperately as we have to run the timer
		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerChan <- formatString(answer)
		}()
		select {
		case <-timer.C:
			fmt.Printf("Time up!!!.You scored %d out of %d.\n", correct, len(problems))
			return
		case answer := <-answerChan:
			if answer == p.a {
				correct++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", correct, len(problems))
}

func parseLines(lines [][]string) []Problem {
	ret := make([]Problem, len(lines))
	for i, line := range lines {
		ret[i] = Problem{
			q: line[0],
			a: formatString(line[1]),
		}
	}
	return ret
}

func formatString(s string) string {
	return strings.ToLower(strings.TrimSpace(s))
}

type Problem struct {
	q string
	a string
}
