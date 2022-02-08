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

const CodecH264 string = "h264"
const CodecH265 string = "h265"

var supportedCodecs = [2]string{CodecH264, CodecH265}

type Stream struct {
	url     string
	camName string
	context Context
}

type Context struct {
	AVFormatCtx    *C.AVFormatContext
	AVStream       *C.AVStream
	AVCodecContext *C.AVCodecContext
	AVPacket       Packet
}

type Packet struct {
	AVPacket *C.AVPacket
	AVFrame  *C.AVFrame
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

	if !context.canWatch() {
		log.Fatal("Video codec of stream not support")
	}

	if context.needTranscode() {
		fmt.Println("Need transcode:\n" + CodecH265 + " codec detected")
	}

	context.AVPacket.AVFrame = C.av_frame_alloc()
	if context.AVPacket.AVFrame == nil {
		log.Fatal("Error alloc frame data")
	}
	context.AVPacket.AVPacket = C.av_packet_alloc()
	if context.AVPacket.AVPacket == nil {
		log.Fatal("Error alloc packet data")
	}

	var response C.int = 0
	for i := 5; C.av_read_frame(context.AVFormatCtx, context.AVPacket.AVPacket) >= 0 && i > 0; {

		// check this is video stream (0) TODO refactor for true video check stream (not always 0 stream)
		if context.AVPacket.AVPacket.stream_index == 0 {
			response = C.decode_packet(context.AVPacket.AVPacket, context.AVCodecContext, context.AVPacket.AVFrame)
			if C.has_decode_error(response) {
				continue
			} else if response < 0 {
				break
			}
			i--
		}
		C.av_packet_unref(context.AVPacket.AVPacket)
	}
	C.av_frame_free(&context.AVPacket.AVFrame)
	C.av_packet_free(&context.AVPacket.AVPacket)
}

func (context *Context) decodePacket() {

}

func (context *Context) canWatch() bool {
	for _, codec := range supportedCodecs {
		if codec == C.GoString(context.AVCodecContext.codec.name) {
			return true
		}
	}
	return false
}

func (context *Context) needTranscode() bool {
	return CodecH264 == C.GoString(context.AVCodecContext.codec.name)
}
