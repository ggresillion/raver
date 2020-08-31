export interface TrackInfos {
  live: boolean;
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
  duration: string | null;
  uploaded_at: string | null;
}
