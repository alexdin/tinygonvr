package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
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
	stream := Stream{url: "test.mp4", camName: "Cam 1"}
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
	C.read_context(s.context.AVCodecContext, s.context.AVFormatCtx)
}
