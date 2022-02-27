test:clean
	clear
	go build && ./main
	totem out.mp4

clean:
	rm -f frame-*.pgm
	rm -f *.mp4
