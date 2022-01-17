import React from 'react'
import axios from 'axios'
import { useHistory } from 'react-router'
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
    <div className='h-screen w-screen bg-gray-900'>
      <button
        className='absolute top-2/4 left-2/4 transform -translate-x-2/4 -translate-y-2/4 py-4 px-6 border border-gray-600 rounded-md text-2xl text-white bg-gray-700'
        onClick={handleCreateRoom}
      >
        CREATE ROOM
      </button>
    </div>
  )
}

export default Hub
