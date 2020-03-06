export interface TrackInfos {
  id?: string;
  channelId?: string;
  channelTitle?: string;
  description: string;
  liveBroadcastContent?: string;
  publishedAt?: Date;
  thumbnails?: {default: {url: string}, medium: {url: string}, high: {url: string}};
  title: string;
}
