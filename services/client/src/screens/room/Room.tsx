import React, { useState } from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from '../../hooks/useRoom'
import classes from "./Room.module.scss"
import Playlist from 'components/Playlist'


const Room: React.FC<Props> = () => {
  const [videoUrl, setVideo] = useState('');
  const {
    isMediaReady, playerRef, videoData, playlist,
    handleMediaReady, handlePause, handlePlay, handleSeek, handleRequestVideo, handleMediaEnd,
  } = useRoom()

  const onVideoRequestClick = () => {
    handleRequestVideo(videoUrl);
    setVideo('');
  }

  return (
    <div className={classes.Root}>
      <div className={classes.PlaylistContainer}>
        <div className={classes.VideoInput}>
          <input value={videoUrl} onChange={({ target: { value } }) => setVideo(value)} />
          <button onClick={onVideoRequestClick}>Add to playlist</button>
        </div>

        <Playlist
          className={classes.List}
          videosUrls={playlist}
        />
      </div>

      <div className={classes.VideoContainer}>
        <ReactPlayer
          ref={playerRef}
          playing={videoData.playing && isMediaReady}
          url={videoData.url || ''}
          onSeek={handleSeek}
          onPlay={handlePlay}
          onPause={handlePause}
          onEnded={handleMediaEnd}
          onReady={handleMediaReady}
          controls
          height='100%'
          width='100%'
        />
      </div>
    </div>
  )
}

export default Room
