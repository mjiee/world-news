#ifndef AUDIO_H
#define AUDIO_H


/**
 * audio_transcode - Converts input audio to standard WAV format.
 */
int audio_transcode(const unsigned char* pData, size_t dataSize, const char* outputPath);

/**
 * audio_render - Processes audio with stereo panning and exports to WAV.
 */
int audio_render(const char* inputPath, const char* outputPath, float pan);

/**
 * audio_merge - Concatenates multiple audio files into a single WAV file.
 */
int audio_merge(char** inputPaths, int count, const char* outputPath);

/**
 * audio_error_description - Maps miniaudio error codes to readable strings.
 */
const char* audio_error_description(int res);

#endif