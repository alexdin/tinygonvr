package main

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

const CodecH264 string = "h264"
const CodecH265 string = "h265"
const outFileName string = "out.mp4"

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
		fmt.Println("Need transcode:\n" + CodecH264 + " codec detected")
	}

	context.AVPacket.AVFrame = C.av_frame_alloc()
	if context.AVPacket.AVFrame == nil {
		log.Fatal("Error alloc frame data")
	}

	context.AVPacket.AVPacket = C.av_packet_alloc()
	if context.AVPacket.AVPacket == nil {
		log.Fatal("Error alloc packet data")
	}

	var outContext *C.AVFormatContext
	var outStream *C.AVStream

	C.avformat_alloc_output_context2(&outContext, nil, nil, C.CString(outFileName))
	if outContext == nil {
		log.Fatal("out fail initialize outContext")
	}

	outStream = C.avformat_new_stream(outContext, nil)
	if outStream == nil {
		log.Fatal("os stream")
	}

	// copy params
	ret := C.avcodec_parameters_copy(outStream.codecpar, context.AVStream.codecpar)
	if ret < 0 {
		log.Fatal("Fail copy codec params")
	}

	outStream.codecpar.codec_tag = 0

	// write format for outputfile
	C.av_dump_format(outContext, 0, C.CString(outFileName), 1)

	if C.avio_open(&outContext.pb, C.CString(outFileName), C.AVIO_FLAG_WRITE) < 0 {
		log.Fatal("Could not open file for write")
	}

	if C.avformat_write_header(outContext, nil) < 0 {
		log.Fatal("Could not write header data")
	}

	var response C.int = 0
	var seconds C.int = 5
	var i C.int = 0

	//	outStream.r_frame_rate = context.AVStream.r_frame_rate
	for i = 0; C.av_read_frame(context.AVFormatCtx, context.AVPacket.AVPacket) >= 0 && i < context.AVStream.r_frame_rate.num*seconds; {

		// check this is video stream (0) TODO refactor for true video check stream (not always 0 stream)
		if context.AVPacket.AVPacket.stream_index == 0 {

			fmt.Println(context.AVStream.time_base)
			fmt.Println(outStream.time_base)
			C.log_packet(context.AVFormatCtx, context.AVPacket.AVPacket, C.CString("in"))
			/* copy packet */
			//C.av_packet_rescale_ts(context.AVPacket.AVPacket, context.AVStream.time_base, outStream.time_base)
			//context.AVPacket.AVPacket.pts = C.av_rescale_q_rnd(context.AVPacket.AVPacket.pts, context.AVStream.time_base, outStream.time_base, C.AV_ROUND_NEAR_INF|45000)
			//context.AVPacket.AVPacket.dts = C.av_rescale_q_rnd(context.AVPacket.AVPacket.dts, context.AVStream.time_base, outStream.time_base, C.AV_ROUND_NEAR_INF|45000)
			//	context.AVPacket.AVPacket.duration = C.av_rescale_q(context.AVPacket.AVPacket.duration, context.AVStream.time_base, outStream.time_base)
			context.AVPacket.AVPacket.pos = -1
			//	context.AVPacket.AVPacket.stream_index = 0

			C.log_packet(context.AVFormatCtx, context.AVPacket.AVPacket, C.CString("out"))

			// here video sream
			response = C.av_interleaved_write_frame(outContext, context.AVPacket.AVPacket)
			if response < 0 {
				//log.Fatal("write file error")
				//	break
			}
			i++
		} else {
			C.av_packet_unref(context.AVPacket.AVPacket)
		}

	}
	C.av_write_trailer(outContext)
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
