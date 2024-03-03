package filter

import (
	"emojihash/emoji"
)

type EmojiFilter func(emoji.Emoji) bool

func UseAll(e emoji.Emoji) bool {
	return true
}

func UseNone(e emoji.Emoji) bool {
	return false
}

func Not(f EmojiFilter) EmojiFilter {
	return func(e emoji.Emoji) bool {
		return !f(e)
	}
}

func All(filters ...EmojiFilter) EmojiFilter {
	return func(e emoji.Emoji) bool {
		for _, filter := range filters {
			if !filter(e) {
				return false
			}
		}
		return true
	}
}

// https://symbl.cc/en/unicode/blocks/miscellaneous-symbols-and-pictographs/#subblock-1F3FB
var SkinTones = []rune{
	'\U0001F3FB',
	'\U0001F3FC',
	'\U0001F3FD',
	'\U0001F3FE',
	'\U0001F3FF',
}

var ZeroWidthJoiner = []rune{'\U0000200D'}

func IgnoreRunesContaining(runesToIgnore ...[]rune) EmojiFilter {
	ignore := make(map[rune]bool)
	for _, runes := range runesToIgnore {
		for _, r := range runes {
			ignore[r] = true
		}
	}
	return func(e emoji.Emoji) bool {
		for _, r := range e.Runes {
			if ignore[r] {
				return false
			}
		}
		return true
	}
}

func IncludeGroups(groups []string) EmojiFilter {
	return func(e emoji.Emoji) bool {
		for _, g := range groups {
			if e.Group == g || e.SubGroup == g {
				return true
			}
		}
		return false
	}
}
