export interface TrackInfos {
  videoId?: string;
  channelId?: string;
  channelTitle?: string;
  description: string;
  publishedAt?: string;
  thumbnails?: {default: {url: string}, medium: {url: string}, high: {url: string}};
  title: string;
}
