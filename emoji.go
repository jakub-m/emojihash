package main

import (
	"regexp"
	"strings"
)

const nl = "\n"

type EmojiManager map[string]map[string]Emoji

type Emoji struct {
}

func NewEmojiManager(s string) EmojiManager {
	em := EmojiManager(make(map[string]map[string]Emoji))
	currentGroup := ""
	currentSubGroup := ""
	for it := (emojiParser{data: s}); !it.end(); it.scan() {
		if g, ok := getGroup(it.line); ok {
			currentGroup = g
			em[currentGroup] = make(map[string]Emoji)
		}
		if sg, ok := getSubGroup(it.line); ok {
			if currentGroup == "" {
				panic("group missing")
			}
			currentSubGroup = sg
		}
		if e, ok := getEmojiWithMeta(it.line); ok {
			if currentGroup == "" || currentSubGroup == "" {
				panic("group or subgroup missing")
			}
			em[currentGroup][currentSubGroup] = e
		}
	}
	return em
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

func getEmojiWithMeta(line string) (Emoji, bool) {
	var zero Emoji
	return zero, false
}

func normalizeGroupName(s string) string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, "&", "and")
	return strings.ReplaceAll(s, " ", "-")
}
