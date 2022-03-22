package ffmpeg

// #cgo pkg-config: libavcodec libavutil libavformat libswscale
// #cgo CFLAGS: -std=c11 -g
// #include "ffmpeg.h"
import "C"
import (
	"fmt"
	"github.com/alexdin/tinygonvr/alarm"
	"log"
	"unsafe"
)

const CodecH264 string = "h264"
const CodecH265 string = "h265"

var supportedCodecs = [2]string{CodecH264, CodecH265}

type Stream struct {
	Url     string
	CamName string
	Context Context
}

type Context struct {
	AVFormatCtx    *C.AVFormatContext
	AVStream       *C.AVStream
	AVCodecContext *C.AVCodecContext
	AVPacket       Packet
	VideoIndex     C.uint
}

type Packet struct {
	AVPacket *C.AVPacket
	AVFrame  *C.AVFrame
}

func Boot(config []Stream) {
	for i, stream := range config {
		stream.Open()
		stream.Screen(i)
		stream.Close()
	}
}

func (s *Stream) Open() {
	s.Context.AVFormatCtx = C.openStream(C.CString(s.Url))
	s.Context.AVStream, s.Context.VideoIndex = getVideoStream(s.Context.AVFormatCtx)
	s.Context.AVCodecContext = getCodec(s.Context.AVStream)
}

func getVideoStream(ctx *C.AVFormatContext) (*C.AVStream, C.uint) {

	streams := (*[1 << 30]*C.AVStream)(unsafe.Pointer(ctx.streams))
	var stream *C.AVStream = nil
	var i C.uint = 0
	for i = 0; i < ctx.nb_streams; i++ {
		if streams[i].codecpar.codec_type == C.AVMEDIA_TYPE_VIDEO {
			stream = streams[i]
			break
		}
	}
	return stream, i
}

func getCodec(stream *C.AVStream) *C.AVCodecContext {
	codec := C.avcodec_find_decoder(stream.codec.codec_id)
	if codec == nil {
		log.Fatal("Cant find codec")
	}

	codecCtx := C.avcodec_alloc_context3(codec)
	if C.avcodec_parameters_to_context(codecCtx, stream.codecpar) < 0 {
		log.Fatal("fail init params to codec")
	}

	if C.avcodec_open2(codecCtx, codec, nil) < 0 {
		log.Fatal("Fatal avcodec_open2")
	}

	return codecCtx
}

func (s *Stream) Close() {
	C.avformat_close_input(&s.Context.AVFormatCtx)
	C.avcodec_free_context(&s.Context.AVCodecContext)
}

func (s *Stream) Screen(index int) {
	//s.Context.WriteVideoFile(s.CamName)
	s.Context.WaitForAlarm(index)
}

func (context *Context) WriteVideoFile(outFileName string) {

	outFileName = outFileName + ".mp4"

	if !context.canWatch() {
		log.Fatal("Video codec of stream not support")
	}

	if context.needTranscode() {
		fmt.Println("Need transcode:\n" + CodecH264 + " codec detected")
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

	// write format for output file
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

		if context.AVPacket.AVPacket.stream_index == C.int(context.VideoIndex) {

			//	fmt.Println(context.AVStream.time_base)
			//	fmt.Println(outStream.time_base)
			C.log_packet(context.AVFormatCtx, context.AVPacket.AVPacket, C.CString("in"))

			/* copy packet */
			context.AVPacket.AVPacket.pos = -1
			context.AVPacket.AVPacket.stream_index = C.int(context.VideoIndex)

			// correct first packet
			if context.AVPacket.AVPacket.dts > 0 && i == 0 {
				context.AVPacket.AVPacket.dts = 0
				context.AVPacket.AVPacket.pts = 0
			}

			//C.log_packet(context.AVFormatCtx, context.AVPacket.AVPacket, C.CString("out"))

			// here video stream
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
	C.av_packet_free(&context.AVPacket.AVPacket)
}

func (context *Context) WaitForAlarm(index int) {
	context.AVPacket.AVFrame = C.av_frame_alloc()
	if context.AVPacket.AVFrame == nil {
		log.Fatal("Error alloc frame data")
	}

	context.AVPacket.AVPacket = C.av_packet_alloc()
	if context.AVPacket.AVPacket == nil {
		log.Fatal("Error alloc packet data")
	}

	for C.av_read_frame(context.AVFormatCtx, context.AVPacket.AVPacket) >= 0 {
		if context.AVPacket.AVPacket.stream_index == C.int(context.VideoIndex) && C.avcodec_send_packet(context.AVCodecContext, context.AVPacket.AVPacket) >= 0 {
			for response := C.avcodec_receive_frame(context.AVCodecContext, context.AVPacket.AVFrame); !C.has_decode_error(response); {
				if isAlarm(index) {
					log.Fatal("alarm")
				} else {
					fmt.Println(".")
				}
			}
		} else {
			C.av_packet_unref(context.AVPacket.AVPacket)
		}
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
	return true
}

func (context *Context) needTranscode() bool {
	return CodecH264 == C.GoString(context.AVCodecContext.codec.name)
}

func isAlarm(index int) bool {
	conf := alarm.GetAlarm()
	outChan := conf.GetCamChannelByIndex(index)
	return len(outChan) != 0
}
