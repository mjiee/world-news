#define MINIAUDIO_IMPLEMENTATION
#include "miniaudio.h"
#include "audio.h"

/* Production-standard audio parameters */
#define TARGET_FORMAT      ma_format_s16
#define TARGET_CHANNELS    2
#define TARGET_SAMPLE_RATE 44100
#define FRAME_COUNT        1024
#ifndef M_PI
#define M_PI 3.14159265358979323846
#endif

/**
 * internal_stream_copy - Core loop to move PCM frames from decoder to encoder.
 * Handles MA_AT_END as a valid termination state during playback.
 */
static ma_result internal_stream_copy(ma_decoder* pDecoder, ma_encoder* pEncoder, ma_uint64* pFramesOut) {
    ma_int16 buffer[FRAME_COUNT * TARGET_CHANNELS];
    ma_uint64 framesRead;
    ma_result result;
    *pFramesOut = 0;

    while (1) {
        result = ma_decoder_read_pcm_frames(pDecoder, buffer, FRAME_COUNT, &framesRead);
        if (result != MA_SUCCESS && result != MA_AT_END) return result;
        if (framesRead == 0) break;
        
        result = ma_encoder_write_pcm_frames(pEncoder, buffer, framesRead, NULL);
        if (result != MA_SUCCESS) return result;
        
        *pFramesOut += framesRead;
    }
    return MA_SUCCESS;
}

int audio_transcode(const unsigned char* pData, size_t dataSize, const char* outputPath) {
    if (pData == NULL || dataSize == 0) return MA_INVALID_ARGS;

    ma_decoder decoder;
    ma_decoder_config dCfg = ma_decoder_config_init(TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
    
    ma_result res = ma_decoder_init_memory(pData, dataSize, &dCfg, &decoder);
    if (res != MA_SUCCESS) return res;

    ma_encoder encoder;
    ma_encoder_config eCfg = ma_encoder_config_init(ma_encoding_format_wav, TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
    res = ma_encoder_init_file(outputPath, &eCfg, &encoder);
    if (res != MA_SUCCESS) {
        ma_decoder_uninit(&decoder);
        return res;
    }

    ma_uint64 total = 0;
    res = internal_stream_copy(&decoder, &encoder, &total);
    
    ma_encoder_uninit(&encoder);
    ma_decoder_uninit(&decoder);
    return res;
}

int set_stereo(const char* inputPath, const char* outputPath, float pan) {
    if (inputPath == NULL || outputPath == NULL) return MA_INVALID_ARGS;
    if (pan < -1.0f || pan > 1.0f) return MA_INVALID_ARGS;

    ma_decoder decoder;
    ma_decoder_config dCfg = ma_decoder_config_init(TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
    ma_result res = ma_decoder_init_file(inputPath, &dCfg, &decoder);
    if (res != MA_SUCCESS) return res;

    ma_encoder encoder;
    ma_encoder_config eCfg = ma_encoder_config_init(ma_encoding_format_wav, TARGET_FORMAT, 2, TARGET_SAMPLE_RATE);
    res = ma_encoder_init_file(outputPath, &eCfg, &encoder);
    if (res != MA_SUCCESS) {
        ma_decoder_uninit(&decoder);
        return res;
    }

    float angle = (pan + 1.0f) * ((float)M_PI / 4.0f);
    float lG = cosf(angle);
    float rG = sinf(angle);

    ma_int16 buf[FRAME_COUNT * TARGET_CHANNELS];
    ma_uint64 fRead, total = 0;

    while (1) {
        res = ma_decoder_read_pcm_frames(&decoder, buf, FRAME_COUNT, &fRead);
        if (res != MA_SUCCESS && res != MA_AT_END) break;
        if (fRead == 0) { res = MA_SUCCESS; break; }

        total += fRead;

        for (ma_uint64 i = 0; i < fRead; i++) {
            float mono = ((float)buf[i*2] + (float)buf[i*2+1]) * 0.5f;
            float l = mono * lG;
            float r = mono * rG;

            if (l >  32767.0f) l =  32767.0f;
            if (l < -32768.0f) l = -32768.0f;
            if (r >  32767.0f) r =  32767.0f;
            if (r < -32768.0f) r = -32768.0f;

            buf[i*2]   = (ma_int16)(l >= 0.0f ? l + 0.5f : l - 0.5f);
            buf[i*2+1] = (ma_int16)(r >= 0.0f ? r + 0.5f : r - 0.5f);
        }

        res = ma_encoder_write_pcm_frames(&encoder, buf, fRead, NULL);
        if (res != MA_SUCCESS) break;
    }

    ma_encoder_uninit(&encoder);
    ma_decoder_uninit(&decoder);
    return res;
}

int audio_merge(char** inputPaths, int count, const char* outputPath) {
    if (count <= 0) return MA_INVALID_ARGS;

    ma_encoder encoder;
    ma_encoder_config eCfg = ma_encoder_config_init(ma_encoding_format_wav, TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
    ma_result res = ma_encoder_init_file(outputPath, &eCfg, &encoder);
    if (res != MA_SUCCESS) return res;

    for (int i = 0; i < count; i++) {
        ma_decoder decoder;
        ma_decoder_config dCfg = ma_decoder_config_init(TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
        res = ma_decoder_init_file(inputPaths[i], &dCfg, &decoder);
        if (res != MA_SUCCESS) {
            ma_encoder_uninit(&encoder);
            return res;
        }

        ma_uint64 written = 0;
        res = internal_stream_copy(&decoder, &encoder, &written);
        ma_decoder_uninit(&decoder);
        if (res != MA_SUCCESS) {
            ma_encoder_uninit(&encoder);
            return res;
        }
    }

    ma_encoder_uninit(&encoder);
    return MA_SUCCESS;
}

double audio_duration(const char* inputPath) {
    if (inputPath == NULL) return -1.0;

    ma_decoder decoder;
    ma_decoder_config dCfg = ma_decoder_config_init(TARGET_FORMAT, TARGET_CHANNELS, TARGET_SAMPLE_RATE);
    ma_result res = ma_decoder_init_file(inputPath, &dCfg, &decoder);
    if (res != MA_SUCCESS) return -1.0;

    ma_uint64 totalFrames = 0;
    res = ma_decoder_get_length_in_pcm_frames(&decoder, &totalFrames);
    
    ma_decoder_uninit(&decoder);

    if (res != MA_SUCCESS || totalFrames == 0) return -1.0;
    return (double)totalFrames / (double)TARGET_SAMPLE_RATE;
}

const char* audio_error_description(int res) {
    return ma_result_description((ma_result)res);
}