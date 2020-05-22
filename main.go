package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

var csvFile string
var limit int
var shuffle bool
var score int8

func init() {
	flag.StringVar(&csvFile, "problem", "./problems.csv", "Path to csv file")
	flag.IntVar(&limit, "limit", 30, "time limit in seconds")
	flag.BoolVar(&shuffle, "shuffle", false, "pass true if want to shuffle questions")
}

type Problem struct {
	question string
	answer   string
}

func main() {
	flag.Parse()

	records := getRecordsFromCsvFile(csvFile)

	fmt.Printf("Answer the %v questions:\n", len(records))

	reader := bufio.NewReader(os.Stdin)

	timer := time.NewTimer(time.Duration(limit) * time.Second)

	done := make(chan bool)

	go func() {
		for i, r := range records {
			fmt.Printf("%d. %s\nYour answer: ", i+1, r.question)

			answer := readInput(reader)

			if answer == r.answer {
				score++
			}
		}
		done <- true
	}()

	select {
	case <-done:
	case <-timer.C:
		fmt.Println("\nYour time is over!")
	}

	fmt.Printf("Final score %d out of %d", score, len(records))
}

func readInput(r *bufio.Reader) string {
	answer, err := r.ReadString('\n')

	if err != nil {
		log.Fatalln(err)
	}

	return strings.TrimSpace(answer)
}

func getRecordsFromCsvFile(path string) []Problem {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		log.Fatalln(err)
	}

	r := csv.NewReader(file)

	records, err := r.ReadAll()

	if err != nil {
		log.Fatalln(err)
	}

	p := make([]Problem, len(records))

	for i, r := range records {
		p[i] = Problem{
			question: r[0],
			answer:   strings.TrimSpace(r[1]),
		}
	}

	if shuffle {
		s := make([]Problem, len(p))
		perm := rand.Perm(len(p))
		for i, v := range perm {
			s[v] = p[i]
		}
		return s
	}

	return p
}
