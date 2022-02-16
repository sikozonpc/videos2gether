import React, { useState } from 'react'
import ReactPlayer from 'react-player'
import { RoomProps as Props } from './types'
import { useRoom } from 'hooks/useRoom'
import Playlist from 'components/Playlist'


const Room: React.FC<Props> = () => {
  const [videoUrl, setVideo] = useState('');
  const {
    isMediaReady, playerRef, videoData, playlist, usersConnected,
    handleMediaReady, handlePause, handlePlay, handleSeek, handleRequestVideo, handleMediaEnd,
    handleSkipVideo,
  } = useRoom()

  const onVideoRequestClick = () => {
    if (!videoUrl.includes('https://')) return;

    handleRequestVideo(videoUrl);
    setVideo('');
  }

  return (
    <div className='flex bg-gray-900'>
      <div id='left-side'>
        <div className='w-full p-4 text-white'>
          {/* <span className='text-white text-sm'>{userName}</span> */}
          <div className='flex mb-2'>
            <span className='text-red-500 font-bold flex bg-black rounded-full py-1 px-2 text-sm'>
              <figure>
                <svg fill='currentcolor' width="20px" height="20px" version="1.1" viewBox="0 0 20 20" x="0px" y="0px" className="text-red-600"><g><path fillRule="evenodd" d="M5 7a5 5 0 116.192 4.857A2 2 0 0013 13h1a3 3 0 013 3v2h-2v-2a1 1 0 00-1-1h-1a3.99 3.99 0 01-3-1.354A3.99 3.99 0 017 15H6a1 1 0 00-1 1v2H3v-2a3 3 0 013-3h1a2 2 0 001.808-1.143A5.002 5.002 0 015 7zm5 3a3 3 0 110-6 3 3 0 010 6z" clipRule="evenodd"></path></g></svg>
              </figure>
              {usersConnected}
            </span>
          </div>

          <div className='w-full '>
            <input
              className='mb-2.5 flex flex-col text-gray-900'
              value={videoUrl}
              onChange={({ target: { value } }) => setVideo(value)}
            />
            <button
              className='w-full px-1.5 py-2.5 border-none rounded-md text-base text-white bg-gray-700'
              onClick={onVideoRequestClick}
            >
              Add to playlist
            </button>
          </div>
        </div>
      </div>

      <div id='right-side' className='w-full flex-col h-screen '>
        <div className='w-full' style={{ height: '78%' }}>
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
          <Playlist
            className='mt-6 list-none'
            videosUrls={playlist}
            onSkipClick={handleSkipVideo}
          />
        </div>
      </div>
    </div>
  )
}

export default Room
