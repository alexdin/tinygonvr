test: clean
	GOROOT=~/sdk/go1.17.6 #gosetup
	GOPATH=~/sdk #gosetup
	~/sdk/go1.17.6/bin/go build && ./main
clean:
	rm -f frame-*.pgm