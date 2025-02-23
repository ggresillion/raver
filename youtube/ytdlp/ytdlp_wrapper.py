import sys
import json
from yt_dlp import YoutubeDL

def get_video_info(url):
    ydl_opts = {
        'quiet': True,
        'extract_flat': False,  # Ensure full metadata is extracted
    }

    with YoutubeDL(ydl_opts) as ydl:
        try:
            info = ydl.extract_info(url, download=False)
            formats = info.get("formats", [])

            # Filter for audio-only, Opus-encoded, WebM formats
            opus_webm_formats = [
                f for f in formats
                if f.get("vcodec") == "none"  # Audio-only
                and f.get("acodec") == "opus"  # Opus-encoded
                and f.get("ext") == "webm"     # WebM format
            ]

            # Select the best format (e.g., highest bitrate)
            best_format = None
            if opus_webm_formats:
                best_format = max(opus_webm_formats, key=lambda f: f.get("abr", 0))  # Sort by audio bitrate

            return {
                "id": info.get("id"),
                "url": best_format.get("url"),
                "title": info.get("title"),
                "channel": info.get("uploader"),
                "duration": info.get("duration", 0),
                "is_live": info.get("is_live", False),
            }
        except Exception as e:
            return {"error": str(e)}

if __name__ == "__main__":
    while True:
        url = sys.stdin.readline().strip()
        if not url:
            break

        video_info = get_video_info(url)
        print(json.dumps(video_info), flush=True)
