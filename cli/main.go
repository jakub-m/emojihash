package main

import (
	"bufio"
	"emojihash/data"
	"emojihash/emoji"
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
	} else if opts.listEmojisCompact {
		type group [2]string
		groupList := []group{}
		groupMap := make(map[group][]emoji.Emoji)
		for _, e := range emojis {
			eg := group([]string{e.Group, e.SubGroup})
			if _, ok := groupMap[eg]; !ok {
				groupMap[eg] = []emoji.Emoji{}
				groupList = append(groupList, eg)
			}
			groupMap[eg] = append(groupMap[eg], e)
		}
		for _, g := range groupList {
			fmt.Printf("%s\t%s\t", g[0], g[1])
			for _, e := range groupMap[g] {
				fmt.Printf("\t%s", e.Character)
			}
			fmt.Printf("\n")
		}
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
			if opts.verbose {
				fmt.Println(selected)
			} else {
				fmt.Println(selected.Character)
			}
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
	listEmojisCompact bool
	debug             bool
	seed              string
	groupFilterString string
	verbose           bool
}

func parseOptions() options {
	var o options
	flag.BoolVar(&o.listGroups, "lg", false, "alias list-groups")
	flag.BoolVar(&o.listGroups, "list-groups", false, "List groups of emojis. Those groups can be later used to configure which groups are used")
	flag.BoolVar(&o.listEmojis, "l", false, "alias list")
	flag.BoolVar(&o.listEmojis, "list", false, "List emojis")
	flag.BoolVar(&o.listEmojisCompact, "lc", false, "alias list-compact")
	flag.BoolVar(&o.listEmojisCompact, "list-compact", false, "List emojis in a compact way")
	flag.BoolVar(&o.debug, "d", false, "Debug mode")
	flag.StringVar(&o.seed, "s", "", "Additional seed. This string is concatenated to every input string before hashing.")
	flag.StringVar(&o.groupFilterString, "g", "", `Filter groups. The syntax is: "alphanum,~flags". "~" means that the group or subgroup should not be included. The order does not matter.`)
	flag.BoolVar(&o.verbose, "v", false, "Verbose output when hashing.")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Hash input to an emoji.\n\n")
		flag.PrintDefaults()
	}
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
