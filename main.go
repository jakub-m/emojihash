package main

import (
	"flag"
	"fmt"
	"log"
	"strings"
)

func main() {
	if err := mainerr(); err != nil {
		log.Fatal(err)
	}
}

func mainerr() error {
	opts := parseOptions()
	em := NewEmojiManager(emojiFile)
	if opts.listGroups {
		fmt.Println(strings.Join(em.GetEmojiGroups(), "\n"))
		return nil
	}
	panic("todo default action")
}

type options struct {
	listGroups bool
}

func parseOptions() options {
	var o options
	flag.BoolVar(&o.listGroups, "list-groups", false, "List groups of emojis. Those groups can be later used to configure which groups are used")
	flag.Parse()
	return o
}
