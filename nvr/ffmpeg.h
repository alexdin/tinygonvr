#pragma once

#include <math.h>

#include "libswscale/swscale.h"
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"
#include "stdbool.h"
#include <libswresample/swresample.h>
#include <libavutil/avassert.h>
#include <libavutil/channel_layout.h>
#include <libavutil/opt.h>
#include <libavutil/mathematics.h>
#include <libavutil/timestamp.h>

 AVFormatContext* openStream(const char *filename);

 AVStream* getVideoStream(AVFormatContext *av_context);

 AVCodecContext* getCodec(AVStream *stream);

 void read_context(AVCodecContext *codec_ctx, AVFormatContext *format_ctx);

 void save_gray_frame(unsigned char *buf, int wrap, int xsize, int ysize, char *filename);

 int decode_packet(AVPacket *pPacket, AVCodecContext *pCodecContext, AVFrame *pFrame);

 bool has_decode_error(int result);

 void log_packet(const AVFormatContext *fmt_ctx, const AVPacket *pkt, const char *tag);
