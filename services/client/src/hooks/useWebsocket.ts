import { useState, useEffect, useCallback } from "react"
import { VideoData } from "../screens/room/types"


export enum ActionType {
  PLAY_VIDEO = "PLAY_VIDEO",
  PAUSE_VIDEO = "PAUSE_VIDEO",
  END_VIDEO = "END_VIDEO",
  SYNC = "SYNC",
  REQUEST = "REQUEST",
  REQUEST_TIME = "REQUEST_TIME",
  SEND_TIME_TO_SERVER = "SEND_TIME_TO_SERVER",
}

export type Message = {
  action: ActionType,
  data: VideoData,
}

/** WebSocket wrapper */
const useWebsocket = (
  url: string,
  messageListener: (ev: MessageEvent) => void
) => {
  const [ws, setWebsocket] = useState<WebSocket | null>(null)

  let timeout = 250 // Initial timeout duration as a class variable

  /** Establishes the connect with the websocket and also ensures constant reconnection if connection closes */
  const connect = useCallback(() => {
    var ws = new WebSocket(url)
    let connectInterval: any

    // websocket onopen event listener
    ws.onopen = () => {
      console.log("connected websocket.")

      setWebsocket(ws)

      timeout = 250 // reset timer to 250 on open of websocket connection
      clearTimeout(connectInterval) // clear Interval on on open of websocket connection
    }

    // websocket onclose event listener
    ws.onclose = (e: any) => {
      console.log(
        `Socket is closed. Reconnect will be attempted in ${Math.min(
          10000 / 1000,
          (timeout + timeout) / 1000
        )} second.`,
        e.reason
      )

      timeout += timeout //increment retry interval
      connectInterval = setTimeout(check, Math.min(10000, timeout)) //call check function after timeout
    }

    // websocket onerror event listener
    ws.onerror = (wb: any) => {
      console.error(
        "Socket encountered error: ",
        wb,
        "Closing streaming"
      )

      ws.close()
    }

    ws.onmessage = (ev: MessageEvent) => {
      messageListener(ev)
    }
  }, [url])

  useEffect(() => {
    // single websocket instance for the own application and constantly trying to reconnect.
    connect()
  }, [connect])


  function sendMessage(message: Message) {
    if (!ws) {
      console.warn("No web streaming")
      return
    }

    const data = JSON.stringify(message)

    try {
      ws.send(data)
    } catch (error) {
      console.log(error)
    }
  }

  /**
   * utilited by the @function connect to check if the connection is close, if so attempts to reconnect
   */
  function check() {
    //check if websocket instance is closed, if so call `connect` function.
    if (!ws || ws.readyState === WebSocket.CLOSED) connect()
  }

  return {
    sendMessage,

  }
}

export default useWebsocket
