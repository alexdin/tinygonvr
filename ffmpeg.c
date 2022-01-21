
#include "ffmpeg.h"

 AVFormatContext* openStream(const char *filename) {

    AVFormatContext *av_format_ctx = avformat_alloc_context();
    if (!av_format_ctx) {
        exit(1);
    }

    if (avformat_open_input(&av_format_ctx, filename, NULL, NULL) < 0) {
        exit(1);
    }

    if (avformat_find_stream_info(av_format_ctx, NULL) < 0) {
        exit(1);
    }

    return av_format_ctx;
}

 AVStream* getVideoStream(AVFormatContext* av_context) {
    for (int i = 0; i < av_context->nb_streams; i++) {
        AVCodecParameters *pLocalCodecParameters =  NULL;
        pLocalCodecParameters = av_context->streams[i]->codecpar;
        if (pLocalCodecParameters->codec_type == AVMEDIA_TYPE_VIDEO) {
             return av_context->streams[i];
        }
    }
    exit(1);
}

 AVCodecContext* getCodec(AVStream* stream) {

    AVCodec *av_codec = NULL;
    AVCodecContext *codec_ctx = NULL;
    av_codec = avcodec_find_decoder(AV_CODEC_ID_H264);
    if (!av_codec) {
        exit(1);
    }

    codec_ctx = avcodec_alloc_context3(av_codec);
    int ret = avcodec_parameters_to_context(codec_ctx, stream->codecpar);
    if (ret < 0) {
        exit(1);
    }
    if (avcodec_open2(codec_ctx, av_codec, NULL) < 0) {
        exit(1);
    }

    return codec_ctx;
}

 void save_gray_frame(unsigned char *buf, int wrap, int xsize, int ysize, char *filename) {
    FILE *f;
    int i;
    f = fopen(filename, "w");
    // writing the minimal required header for a pgm file format
    // portable graymap format -> https://en.wikipedia.org/wiki/Netpbm_format#PGM_example
    fprintf(f, "P5\n%d %d\n%d\n", xsize, ysize, 255);

    // writing line by line
    for (i = 0; i < ysize; i++)
        fwrite(buf + i * wrap, 1, xsize, f);
    fclose(f);
}

 int decode_packet(AVPacket *pPacket, AVCodecContext *pCodecContext, AVFrame *pFrame) {
    // Supply raw packet data as input to a decoder
    // https://ffmpeg.org/doxygen/trunk/group__lavc__decoding.html#ga58bc4bf1e0ac59e27362597e467efff3
    int response = avcodec_send_packet(pCodecContext, pPacket);

    if (response < 0) {
        return response;
    }

    while (response >= 0) {
        // Return decoded output data (into a frame) from a decoder
        // https://ffmpeg.org/doxygen/trunk/group__lavc__decoding.html#ga11e6542c4e66d3028668788a1a74217c
        response = avcodec_receive_frame(pCodecContext, pFrame);
        if (response == AVERROR(EAGAIN) || response == AVERROR_EOF) {
            break;
        } else if (response < 0) {
            return response;
        }

        if (response >= 0) {

            char frame_filename[1024];
            snprintf(frame_filename, sizeof(frame_filename), "%s-%d.pgm", "frame", pCodecContext->frame_number);
            // Check if the frame is a planar YUV 4:2:0, 12bpp
            // That is the format of the provided .mp4 file
            // RGB formats will definitely not give a gray image
            // Other YUV image may do so, but untested, so give a warning
            if (pFrame->format != AV_PIX_FMT_YUV420P) {
            }
            // save a grayscale frame into a .pgm file
            save_gray_frame(pFrame->data[0], pFrame->linesize[0], pFrame->width, pFrame->height, frame_filename);
        }
    }
    return 0;
}

 void read_context(AVCodecContext *codec_ctx, AVFormatContext *format_ctx) {
    printf("Start decode data\n");
    int response = 0;
    int fi = 5;
    AVFrame *av_frame = av_frame_alloc();
    if (!av_frame) {
        exit(1);
    }
    AVPacket *av_packet = av_packet_alloc();
    if (!av_packet) {
        exit(1);
    }
    //av_read_play(format_ctx);

    while (av_read_frame(format_ctx, av_packet) >= 0) {
        // if it's the video stream
        if (av_packet->stream_index == 0) {
            response = decode_packet(av_packet, codec_ctx, av_frame);
            if (response == AVERROR(EAGAIN) || response == AVERROR_EOF) {
                continue;
            } else if (response < 0)
                break;
            // stop it, otherwise we'll be saving hundreds of frames
            if (--fi <= 0) break;
        }
        // https://ffmpeg.org/doxygen/trunk/group__lavc__packet.html#ga63d5a489b419bd5d45cfd09091cbcbc2
        av_packet_unref(av_packet);
    }
     av_frame_free(&av_frame);
     av_packet_free(&av_packet);
}