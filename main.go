package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

type Stream struct {
	url     string
	camName string
	context Context
}

type Context struct {
	AVFormatCtx    *C.AVFormatContext
	AVStream       *C.AVStream
	AVCodecContext *C.AVCodecContext
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	stream := Stream{url: os.Getenv("TEST_CAM_URL"), camName: "Cam 1"}
	stream.Open()
	stream.Screen()
	stream.Close()
	fmt.Println("Done")
}

func (s *Stream) Open() {
	s.context.AVFormatCtx = C.openStream(C.CString(s.url))
	s.context.AVStream = C.getVideoStream(s.context.AVFormatCtx)
	s.context.AVCodecContext = C.getCodec(s.context.AVStream)
}

func (s *Stream) Close() {
	C.avformat_close_input(&s.context.AVFormatCtx)
	C.avcodec_free_context(&s.context.AVCodecContext)
}

func (s *Stream) Screen() {
	//C.read_context(s.context.AVCodecContext, s.context.AVFormatCtx)
	s.context.Read()
}

func (context *Context) Read() {
	fmt.Println("Start on go lang encode")

	avFrame := C.av_frame_alloc()
	if avFrame == nil {
		log.Fatal("Error alloc frame data")
		os.Exit(0)
	}
	avPacket := C.av_packet_alloc()
	if avPacket == nil {
		log.Fatal("Error alloc packet data")
		os.Exit(0)
	}
	var response C.int = 0
	for i := 5; C.av_read_frame(context.AVFormatCtx, avPacket) >= 0 && i > 0; {

		if avPacket.stream_index == 0 {
			response = C.decode_packet(avPacket, context.AVCodecContext, avFrame)
			if C.has_decode_error(response) {
				continue
			} else if response < 0 {
				break
			}
			i--
		}

		C.av_packet_unref(avPacket)
	}
}
