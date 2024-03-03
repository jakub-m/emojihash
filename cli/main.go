package main

import (
	"bufio"
	"emojihash/data"
	"emojihash/filter"
	"emojihash/loader"
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := mainerr(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainerr() error {
	opts := parseOptions()
	if opts.debug {
		log.SetOutput(os.Stderr)
	} else {
		log.SetOutput(io.Discard)
	}
	groupFilter := parseGroupFilter(opts.groupFilterString)
	emojis := loader.LoadEmojiFromFile(
		data.EmojiFile,
		filter.All(
			groupFilter,
			filter.IgnoreRunesContaining(
				filter.ZeroWidthJoiner,
				filter.SkinTones,
			),
		),
	)
	if opts.listGroups {
		printedGroups := make(map[string]bool)
		for _, e := range emojis {
			s := fmt.Sprintf("%s\t%s", e.Group, e.SubGroup)
			if printedGroups[s] {
				continue
			}
			fmt.Println(s)
			printedGroups[s] = true
		}
		return nil
	} else if opts.listEmojis {
		for _, e := range emojis {
			fmt.Println(e)
		}
	} else {
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			h := int(hash(scanner.Text() + opts.seed))
			selected := emojis[h%len(emojis)]
			log.Println(selected)
			fmt.Println(selected.Character)
		}
		if scanner.Err() != nil {
			return scanner.Err()
		}
	}
	return nil
}

type options struct {
	listGroups        bool
	listEmojis        bool
	debug             bool
	seed              string
	groupFilterString string
}

func parseOptions() options {
	var o options
	flag.BoolVar(&o.listGroups, "lg", false, "")
	flag.BoolVar(&o.listGroups, "list-groups", false, "List groups of emojis. Those groups can be later used to configure which groups are used")
	flag.BoolVar(&o.listEmojis, "l", false, "")
	flag.BoolVar(&o.listEmojis, "list", false, "list emojis")
	flag.BoolVar(&o.debug, "d", false, "Debug mode")
	flag.StringVar(&o.seed, "s", "", "Additional seed. This string is concatenated to every input string before hashing.")
	flag.StringVar(&o.groupFilterString, "g", "", `Filter groups. The syntax is: "alphanum,~flags". "~" means that the group or subgroup should not be included. The order does not matter.`)
	flag.Parse()
	return o
}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

const excludeMark = "~"

func parseGroupFilter(s string) filter.EmojiFilter {
	includeGroups := []string{}
	excludeGroups := []string{}
	for _, f := range strings.Split(s, ",") {
		f = strings.TrimSpace(f)
		if f == "" {
			continue
		}
		log.Printf("Group filter: %s", strconv.Quote(f))
		if strings.HasPrefix(f, excludeMark) {
			excludeGroups = append(excludeGroups, f[len(excludeMark):])
		} else {
			includeGroups = append(includeGroups, f)
		}
	}
	includeFilter := filter.UseAll
	if len(includeGroups) > 0 {
		log.Printf("Include groups: %s", includeGroups)
		includeFilter = filter.IncludeGroups(includeGroups)
	}
	excludeFilter := filter.UseAll
	if len(excludeGroups) > 0 {
		log.Printf("Exclude groups: %s", excludeGroups)
		excludeFilter = filter.Not(filter.IncludeGroups(excludeGroups))
	}
	return filter.All(includeFilter, excludeFilter)
}
