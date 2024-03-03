# emojihash

Hash a string to emoji. Ideal to enrich you terminal prompt with emoji based on directory or day. A recipe to emit
PWD-based emoji, changed daily:

```bash
% echo $PWD | ./bin/emojihash -s $(date +%y%m%d) -g animals-and-nature,~animal-bug
🌷
```

To plant this to zsh use:

```zsh
function __prompt_emoji {
	echo "$PWD" | ~/bin/emojihash -s $(date +%Y%m%d)
}

function __prompt_dir {
	basename "$PWD"
}

setopt PROMPT_SUBST
PS1='$(__prompt_emoji) $(__prompt_dir) %'
```

File [emoji-test.txt](emoji-test.txt) downloaded from [unicode.org][ref_unicode].

[ref_unicode]:https://unicode.org/Public/emoji/15.1/emoji-test.txt
