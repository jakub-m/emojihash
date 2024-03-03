package loader

import (
	"emojihash/emoji"
	"emojihash/filter"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"
)

const nl = "\n"

func LoadEmojiFromFile(s string, filter filter.EmojiFilter) []emoji.Emoji {
	emojis := []emoji.Emoji{}
	currentGroup := ""
	currentSubGroup := ""
	for it := (emojiParser{data: s}); !it.end(); it.scan() {
		if it.line == "" {
			continue
		} else if g, ok := getGroup(it.line); ok {
			currentGroup = g
		} else if sg, ok := getSubGroup(it.line); ok {
			if currentGroup == "" {
				panic("group missing")
			}
			currentSubGroup = sg
		} else if e, ok := getEmojiWithMeta(it.line); ok {
			if currentGroup == "" || currentSubGroup == "" {
				panic("group or subgroup missing")
			}
			e.Group = currentGroup
			e.SubGroup = currentSubGroup
			// TODO filterig here
			if filter(e) {
				emojis = append(emojis, e)
			}
		} else {
			_ = "ignore this branch, no need to print failed parsing, input won't change"
			// log.Printf("Failed to handle: %s", strconv.Quote(it.line))
		}
	}
	return emojis
}

type emojiParser struct {
	data string
	line string
	pos  int
}

func (p *emojiParser) scan() {
	i := strings.Index(p.remaining(), nl)
	if i == -1 {
		// last line
		p.line = p.remaining()
		p.pos = len(p.data)
		return
	}
	p.line = p.remaining()[:i]
	p.pos += i + len(nl)
}

func (p *emojiParser) remaining() string {
	if p.end() {
		return ""
	}
	return p.data[p.pos:]
}

func (p *emojiParser) end() bool {
	return p.pos >= len(p.data)
}

var regexGroup = regexp.MustCompile("# group: (.+)")

func getGroup(line string) (string, bool) {
	matches := regexGroup.FindStringSubmatch(line)
	if len(matches) == 0 {
		return "", false
	}
	groupName := matches[1]
	return normalizeGroupName(groupName), true
}

var regexSubGroup = regexp.MustCompile("# subgroup: (.+)")

func getSubGroup(line string) (string, bool) {
	matches := regexSubGroup.FindStringSubmatch(line)
	if len(matches) == 0 {
		return "", false
	}
	groupName := matches[1]
	return normalizeGroupName(groupName), true
}

var regexEmoji = regexp.MustCompile(`(\w+(?:\s\w+)*)\s+;.*?#.*?E\d+(?:\.\d+)\s+(.*)`)

func getEmojiWithMeta(line string) (emoji.Emoji, bool) {
	var zero emoji.Emoji
	if strings.HasPrefix(line, "#") {
		return zero, false
	}
	matches := regexEmoji.FindStringSubmatch(line)
	if len(matches) == 0 {
		return zero, false
	}
	c, rr, err := decodeCharacterFromEncodedRunes(matches[1])
	if err != nil {
		log.Printf("Failed to decode rune for line: %s", strconv.Quote(line))
		return zero, false
	}
	description := matches[2]
	return emoji.Emoji{
		Character:   c,
		Runes:       rr,
		Description: description,
	}, true
}

func decodeCharacterFromEncodedRunes(s string) (string, []rune, error) {
	c := ""
	rr := []rune{}
	for _, p := range strings.Split(s, " ") {
		i, err := strconv.ParseInt(p, 16, 64)
		if err != nil {
			return c, rr, fmt.Errorf("cannot parse %s: %s", strconv.Quote(p), err)
		}
		r := rune(i)
		c += string(r)
		rr = append(rr, r)
	}
	return c, rr, nil
}

func normalizeGroupName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", "and")
	return strings.ReplaceAll(s, " ", "-")
}
