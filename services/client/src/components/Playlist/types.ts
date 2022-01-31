import { VideoData } from 'screens/room/types';

export interface PlaylistProps extends React.HTMLAttributes<HTMLDivElement> {
  videosUrls: VideoData[],
}
