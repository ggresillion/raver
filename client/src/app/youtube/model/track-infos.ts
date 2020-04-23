export interface TrackInfos {
  videoId?: string;
  author?: any;
  channelTitle?: string;
  description: string;
  liveBroadcastContent?: string;
  publishedAt?: Date;
  thumbnails?: {default: {url: string}, medium: {url: string}, high: {url: string}};
  title: string;
}
