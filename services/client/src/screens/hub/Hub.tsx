import React from 'react'
import axios from 'axios'
import { useHistory } from 'react-router'
import classes from './Hub.module.scss'
import { API_URL } from 'config'

const Hub: React.FC = () => {
  const history = useHistory()

  const handleCreateRoom = () => {
    axios.get(`${API_URL}/room`)
      .then(d => {
        history.push(`/room/${d.data.ID}`)
      })
  }

  return (
    <div className={classes.Root}>
      <button onClick={handleCreateRoom}>CREATE ROOM</button>
    </div>
  )
}

export default Hub
