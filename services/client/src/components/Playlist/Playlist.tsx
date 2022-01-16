import React from 'react'
import { PlaylistProps as Props } from './types'
import ReactPlayer from 'react-player'
import classnames from 'util/classnames'
import c from './Playlist.module.scss'

const YT_VARS = {
  youtube: {
    playerVars: {
      autoplay: 0,
      modestbranding: 0,
      rel: 0,
    },
  },
}

const Playlist: React.FC<Props> = ({ videosUrls, className, ...rest }) => {
  const hasVideos = Array.isArray(videosUrls) && videosUrls.length > 0

  const currentStyles: React.CSSProperties = {
    border: '2px solid greenyellow',
  };

  return (
    <div className={classnames(c.List, className)} {...rest}>
      {hasVideos && <h3>Playlist</h3>}

      {videosUrls.map((video, index) => {
        const isCurrentPlaying = index === 0;

        return (
          <div
            className={c.Video}
            style={isCurrentPlaying ? currentStyles : {}}
            key={index}
          >
            <ReactPlayer
              url={video}
              light
              controls={false}
              playIcon={<div />}
              height='100%'
              width='100%'
              playing={false}
              config={YT_VARS}
            />
          </div>
        )
      })}
    </div>
  )
}

export default Playlist
