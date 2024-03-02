bin=bin/emojihash
gofiles=$(shell find . -name \*.go)
$(bin): $(gofiles)
	go build -o $(bin) $(gofiles)
clean:
	rm -rfv $(bin) out/
.phony: clean
