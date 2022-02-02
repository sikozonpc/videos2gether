import React from 'react'
import { PlaylistProps as Props } from './types'
import ReactPlayer from 'react-player'
import classnames from 'util/classnames'

const YT_VARS = {
  youtube: {
    playerVars: {
      autoplay: 0,
      modestbranding: 0,
      rel: 0,
    },
  },
}

const Playlist: React.FC<Props> = ({
  videosUrls,
  className,
  onSkipClick,
  ...rest
}) => {
  const hasVideos = Array.isArray(videosUrls) && videosUrls.length > 0

  const currentStyles: React.CSSProperties = {
    border: '2px solid greenyellow',
  };

  return (
    <div className={classnames('w-full flex flex-row overflow-x-scroll md:max-w-screen-md lg:max-w-screen-2xl h-40 bg-black', className)} {...rest}>
      {!hasVideos && <p className='w-full text-sm text-gray-200 text-center pt-20'>
        Add a video to the playlist to start watching
      </p>}

      {videosUrls.map((video, index) => {
        const isCurrentPlaying = index === 0;

        return (
          <div
            className='relative w-full h-40 mr-5'
            style={{
              ...(isCurrentPlaying ? currentStyles : {}),
              minWidth: 200,
            }}
            key={index}
          >
            <ReactPlayer
              url={video.url}
              light
              controls={false}
              playIcon={<div />}
              height='100%'
              width='100%'
              playing={false}
              config={YT_VARS}
            />
            {isCurrentPlaying && (
              <div
                className='text-white text-xl absolute top-1/2 left-1/2 bg-slate-700 p-2 rounded cursor-pointer'
                onClick={onSkipClick}
              >
                {videosUrls.length === 1 ? "End video" : "Skip"}
              </div>
            )}
          </div>
        )
      })}
    </div>
  )
}

export default Playlist
