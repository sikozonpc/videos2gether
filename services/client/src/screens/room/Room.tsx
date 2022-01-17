import React, { useState } from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from '../../hooks/useRoom'
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
    <div className='flex bg-gray-900'>
      <div className='w-1/5 p-4 text-white'>
        <div>
          <input
            className='mb-2.5 flex flex-col text-gray-900'
            value={videoUrl}
            onChange={({ target: { value } }) => setVideo(value)}
          />
          <button
            className='px-1.5 py-2.5 border-none rounded-md text-base text-white bg-gray-700'
            onClick={onVideoRequestClick}
          >
            Add to playlist
          </button>
        </div>

        <Playlist
          className='mt-6 list-none'
          videosUrls={playlist}
        />
      </div>

      <div className='w-4/5 h-screen'>
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
