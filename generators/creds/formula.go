package creds

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/briandowns/formatifier"

	"github.com/gen0cide/laforge/static"
)

// steps:
// - select random bad password 8
// - append a ficticous year 5
// - append a special character 2
// - append a number 6
// - prepend a number 2
// - capitalize the first letter 4
// - l33tspeak everything - 1

var (
	TotalWeight = 0

	Top500BadPasswords = []string{}

	specialChars = []rune("@!#$%^&")

	Steps = []Step{
		Step{
			Name:   "Append Random Top500 Password",
			Weight: 7,
			Func:   appendRandomBadPassword,
		},
		Step{
			Name:   "Append A Year",
			Weight: 5,
			Func:   appendYear,
		},
		Step{
			Name:   "Append a Special Char",
			Weight: 4,
			Func:   appendSpecial,
		},
		Step{
			Name:   "Append a Number",
			Weight: 6,
			Func:   appendNumber,
		},
		Step{
			Name:   "Prepend a Number",
			Weight: 3,
			Func:   prependNumber,
		},
		Step{
			Name:   "Capitalize First Letter",
			Weight: 5,
			Func:   capitalizeFirst,
		},
		// Step{
		// 	Name:   "Leetspeak Everything",
		// 	Weight: 2,
		// 	Func:   leetspeakEverything,
		// },
	}
)

type Step struct {
	Name   string
	Weight int
	Func   func(i string) string
}

func init() {
	tmpldata, err := static.ReadFile("badpasswords.txt")
	if err != nil {
		panic(err)
	}
	Top500BadPasswords = strings.Split(string(tmpldata), "\n")

	for _, x := range Steps {
		TotalWeight += x.Weight
	}
}

func RandomPassword(stepCount int) string {
	rand.Seed(time.Now().UnixNano())
	s := appendRandomBadPassword("")
	for i := 0; i < stepCount; i++ {
		beginningS := s
		for {
			s = RandomWeightedStep().Func(s)
			if s != beginningS {
				break
			}
		}
	}

	time.Sleep(10 * time.Millisecond)
	return s
}

func appendRandomBadPassword(s string) string {
	r := rand.Intn(len(Top500BadPasswords) - 1)
	w := Top500BadPasswords[r]
	return fmt.Sprintf("%s%s", s, w)
}

func appendYear(s string) string {
	return fmt.Sprintf("%s%d", s, random(1960, 2019))
}

func appendSpecial(s string) string {
	return fmt.Sprintf("%s%s", s, string(specialChars[rand.Intn(len(specialChars))]))
}

func appendNumber(s string) string {
	return fmt.Sprintf("%s%d", s, random(0, 9))
}

func prependNumber(s string) string {
	return fmt.Sprintf("%d%s", random(0, 9), s)
}

func capitalizeFirst(s string) string {
	return strings.Title(s)
}

func leetspeakEverything(s string) string {
	wat, err := formatifier.ToLeet(s)
	if err != nil {
		return s
	}
	return wat
}

func random(min, max int) int {
	return rand.Intn(max-min) + min
}

func RandomWeightedStep() Step {
	for {
		r := rand.Intn(TotalWeight)
		for _, g := range Steps {
			r -= g.Weight
			if r <= 0 {
				return g
			}
		}
	}
}
