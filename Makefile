test: clean
	GOROOT=~/go/go1.17.6 #gosetup
	GOPATH=~/go #gosetup
	~/go/go1.17.6/bin/go build && ./main
clean:
	rm -f frame-*.pgm