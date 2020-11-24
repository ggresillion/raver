import * as ytdl from 'ytdl-core';
import * as prism from 'prism-media';
import { Readable } from 'stream';

function filter(format) {
  return format.codecs === 'opus' &&
    format.container === 'webm' &&
    format.audioSampleRate === 48000;
}

/**
 * Tries to find the highest bitrate audio-only format. Failing that, will use any available audio format.
 * @private
 * @param {Object[]} formats The formats to select from
 */
function nextBestFormat(formats) {
  formats = formats
    .filter(format => format.url)
    .filter(format => format.audioBitrate)
    .sort((a, b) => b.audioBitrate - a.audioBitrate);
  return formats.find(format => !format.bitrate) || formats[0];
}

async function stream(url: string, options = {})
  : Promise<{ stream: Readable, totalLengthSeconds: number }> {
  const info = await ytdl.getInfo(url);
  // Prefer opus
  const format = info.formats.find(filter);
  const canDemux = format && info['length_seconds'] !== '0';
  if (canDemux) {
    options = { ...options, filter };
  } else if (info['length_seconds'] !== '0') {
    options = { ...options, filter: 'audioonly' };
  }
  if (canDemux) {
    const demuxer = new prism.opus.WebmDemuxer();
    const s = ytdl.downloadFromInfo(info, options)
      .on('response', res => {
        const totalSize = res.headers['content-length'];
        return { stream: s, totalLengthSeconds: parseInt(info.videoDetails.lengthSeconds, 10) };
      })
      .pipe(demuxer)
      .on('end', () => demuxer.destroy());
  } else {
    const transcoder = new prism.FFmpeg({
      args: [
        '-reconnect', '1',
        '-reconnect_streamed', '1',
        '-reconnect_delay_max', '5',
        '-i', nextBestFormat(info.formats).url,
        '-analyzeduration', '0',
        '-loglevel', '0',
        '-f', 's16le',
        '-ar', '48000',
        '-ac', '2',
      ],
    });
    const opus = new prism.opus.Encoder({ rate: 48000, channels: 2, frameSize: 960 });
    const s = transcoder.pipe(opus);
    s.on('close', () => {
      transcoder.destroy();
      opus.destroy();
    });
    return { stream: s, totalLengthSeconds: parseInt(info.videoDetails.lengthSeconds, 10) };
  }
}

export { stream, ytdl };
