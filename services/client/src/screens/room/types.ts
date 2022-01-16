export type RoomData = {
  ID: string,
  available: boolean,
}

export interface RoomProps {
  roomData: RoomData,
}

export interface VideoData {
  url?: string,
  time?: number,
  playing?: boolean,
  roomID?: string;
}
