# emojihash

Hash a string to emoji. Ideal to enrich you terminal prompt with emoji based on directory or day. A recipe to emit
PWD-based emoji, changed daily:

```bash
% echo $PWD | ./bin/emojihash -s $(date +%y%m%d) -g animals-and-nature,~animal-bug
🌷
```

File [emoji-test.txt](emoji-test.txt) downloaded from [unicode.org][ref_unicode].

[ref_unicode]:https://unicode.org/Public/emoji/15.1/emoji-test.txt
