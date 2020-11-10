export interface TrackInfos {
  title: string;
  link: string;
  thumbnail: string;
  author: {
    name: string;
    ref: string;
    verified: boolean;
  };
  description: string;
  views: number;
  duration: number | null;
}
