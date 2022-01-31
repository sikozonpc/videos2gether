import { useState, useRef, useCallback, useEffect } from "react"
import { VideoData } from "../screens/room/types"
import ReactPlayer from "react-player"
import useWebsocket, { ActionType } from "hooks/useWebsocket"
import { useParams } from "react-router"
import axios from "axios"
import { WS_URL, API_URL } from "config"

const initialVideoData = {
  time: 0,
  url: "",
  playing: false,
}

export const useRoom = () => {
  const playerRef = useRef<ReactPlayer | null>(null);

  const [videoData, setVideoData] = useState<VideoData>(initialVideoData);
  const [isMediaReady, setMediaReady] = useState(false);
  const [playlist, setPlaylist] = useState<VideoData[]>([]);
  const [synced, setSynced] = useState(false);
  const [userName, setUsername] = useState('');
  const [usersConnected, setUsersConnected] = useState(0);
  const stateRef = useRef<VideoData>();

  const { roomID } = useParams<{ roomID: string }>();

  stateRef.current = videoData;

  const getPlaylist = useCallback(async () => {
    const { data } = (await axios.get<string[]>(`${API_URL}/room/${roomID}/playlist`)) as {
      data: VideoData[],
    };
    if (!data) return;

    setPlaylist(data);
  }, [roomID])


  useEffect(() => {
    (async () => setTimeout(getPlaylist, 1000))();
  }, [getPlaylist])

  const seekVideo = (durationTime: number) => {
    console.log('seeking to', durationTime);
    if (durationTime > 0 && playerRef?.current) {
      playerRef.current.seekTo(durationTime, 'seconds');
      return;
    }
    console.warn("Failed to seek", videoData)
  }

  const messageListener = (ev: MessageEvent) => {
    const res = JSON.parse(ev.data);
    console.log(JSON.parse(ev.data));

    const action: ActionType = res.action;

    if (!res || !action) {
      console.warn("No action to handle");
      return;
    }

    setUsersConnected(res.metadata?.usersConnected);

    switch (action) {
      case ActionType.REQUEST: {
        (async () => getPlaylist())()
        console.log(videoData)

        if (!res.data) return

        const isRequestingNewVideo = stateRef.current?.url !== res.data.url
        if (isRequestingNewVideo) return

        syncVideoWithServer(res.data)
        return
      }

      case ActionType.USER_DISCONNECTED:
      case ActionType.USER_JOINED: {
        return;
      }
      case ActionType.SEND_TIME_TO_SERVER: {
        if (!res.data || videoData.time === 0) return

        return sendMessage({
          action: ActionType.SEND_TIME_TO_SERVER,
          data: {
            time: playerRef.current?.getCurrentTime(),
            playing: videoData.playing,
            roomID,
          }
        });
      }

      case ActionType.END_VIDEO: {
        if (!res.data || !res.data.url) return

        (async () => getPlaylist())()
        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case ActionType.SYNC: {
        setUsername(res.metadata.actionFrom);
        console.log("Client username: ", res.metadata.actionFrom)
        syncVideoWithServer(res.data);

        if (!res.data.url) return
        seekVideo(res.data.time)
        return
      }

      case ActionType.PLAY_VIDEO: {
        if (!res.data) return

        syncVideoWithServer(res.data)
        seekVideo(res.data.time)
        return
      }

      case ActionType.PAUSE_VIDEO: {
        if (!res.data) return
        syncVideoWithServer(res.data)
        return
      }

      default: return console.log("Nothing", res.data)
    }
  }

  const { sendMessage } = useWebsocket(`${WS_URL}/ws/${roomID}`, messageListener)

  useEffect(() => {
    if (!roomID || synced) return;

    sendMessage({ action: ActionType.REQUEST_TIME, data: { roomID } });
    setSynced(true);
  }, [roomID, sendMessage, synced])

  const handleRequestVideo = (url: string) => {
    sendMessage({
      action: ActionType.REQUEST,
      data: { url }
    });
  }

  const syncVideoWithServer = useCallback((newVideoData: VideoData) => {
    setVideoData({
      url: newVideoData.url,
      time: newVideoData.time,
      playing: newVideoData.playing,
    })
  }, [])

  const handlePlay = () => {
    if (!playerRef?.current || videoData.playing) return;

    sendMessage({
      action: ActionType.PLAY_VIDEO,
      data: {
        time: playerRef.current.getCurrentTime(),
        url: videoData.url,
        playing: true,
      }
    });
  }

  const handlePause = () => {
    if (!playerRef?.current) return;

    console.log('SEND')
    sendMessage({
      action: ActionType.PAUSE_VIDEO,
      data: {
        url: videoData.url,
        playing: false,
      }
    });
  }

  const handleSeek = () => {
    if (!playerRef?.current) return;

    sendMessage({
      action: ActionType.PLAY_VIDEO,
      data: {
        time: playerRef?.current?.getCurrentTime() || 0,
        url: videoData.url,
        playing: true,
      }
    });
  }

  const handleMediaEnd = () => {
    sendMessage({
      action: ActionType.END_VIDEO,
      data: {}
    })
  }

  const handleMediaReady = (_player: ReactPlayer) => {
    setMediaReady(true)
  }

  return {
    usersConnected,
    userName,
    isMediaReady,
    videoData,
    playerRef,
    playlist,

    handleMediaEnd,
    handleMediaReady,
    handleSeek,
    handlePause,
    handlePlay,
    handleRequestVideo,
    syncVideoWithServer,
  }
}
