package emoji

import (
	"fmt"
	"strings"
)

type Emoji struct {
	Character   string
	Description string
	Group       string
	Runes       []rune
	SubGroup    string
}

func (e Emoji) String() string {
	return fmt.Sprintf("%s\t%s\t%s\t%s\t%s", e.Character, e.Group, e.SubGroup, e.Description, formatCharacter(e.Runes))
}

func formatCharacter(rr []rune) string {
	h := []string{}
	for _, r := range rr {
		h = append(h, fmt.Sprintf("%X", r))
	}
	return strings.Join(h, " ")
}
