#pragma once
#include "libswscale/swscale.h"
#include "libavcodec/avcodec.h"
#include "libavformat/avformat.h"

 AVFormatContext* openStream(const char *filename);

 AVStream* getVideoStream(AVFormatContext *av_context);

 AVCodecContext* getCodec(AVStream *stream);

 void read_context(AVCodecContext *codec_ctx, AVFormatContext *format_ctx);

 void save_gray_frame(unsigned char *buf, int wrap, int xsize, int ysize, char *filename);

 int decode_packet(AVPacket *pPacket, AVCodecContext *pCodecContext, AVFrame *pFrame);
