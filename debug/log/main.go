package main

import (
	"log"
	"os"
	"strings"
	"time"
)

var version = "test"

var corpus = `
This planet has — or rather had — a problem, which was this:
most of the people living on it were unhappy for pretty much
all of the time. Many solutions were suggested for this problem,
but most of these were largely concerned with the movement of small
green pieces of paper, which was odd because on the whole it wasn't
the small green pieces of paper that were unhappy.

I am bored, that's all. From time to time I yawn so widely that tears
roll down my cheek. It is a profound boredom, profound, the profound
heart of existence, the very matter I am made of. I do not neglect
myself, quite the contrary: this morning I took a bath and shaved.
Only when I think back over those careful little actions, I cannot
understand how I was able to make them: they are so vain.
Habit, no doubt, made them for me. They aren't dead,
they keep on busying themselves, gently, insidiously weaving
their webs, they wash me, dry me, dress me, like nurses.
`

func main() {
	header := os.Getenv("LOGS_HEADER")

	log.Println(header)
	log.Println(version)
	splitCorpus := strings.Split(corpus, "\n")
	for _, line := range splitCorpus {
		log.Println(line)
		time.Sleep(time.Second * 2)
	}
}
