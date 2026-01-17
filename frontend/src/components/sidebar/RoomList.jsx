import React from 'react';
import { Plus, X } from 'lucide-react';
import { useChat } from '../../context/ChatContext';

export default function RoomList() {
  const { rooms, selectedChat, joinRoom, leaveRoom, createRoom } = useChat();

  const handleCreateRoom = () => {
    const roomId = prompt('Enter room ID:');
    if (roomId) {
      const id = parseInt(roomId);
      createRoom(id);
      // TODO: After successful creation, join the room
      // joinRoom(id);
    }
  };

  return (
    <div className="p-2">
      <button
        onClick={handleCreateRoom}
        className="w-full mb-2 p-3 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors flex items-center justify-center gap-2"
      >
        <Plus className="w-4 h-4" />
        Create/Join Room
      </button>
      
      {/* TODO: Map through rooms array */}
      {rooms.map(roomId => (
        <div
          key={roomId}
          onClick={() => joinRoom(roomId)}
          className={`p-3 mb-2 rounded-lg cursor-pointer transition-colors ${
            selectedChat?.type === 'room' && selectedChat?.id === roomId
              ? 'bg-blue-100 border-l-4 border-blue-600'
              : 'bg-gray-50 hover:bg-gray-100'
          }`}
        >
          <div className="flex items-center justify-between">
            <span className="font-medium">Room {roomId}</span>
            <button
              onClick={(e) => { 
                e.stopPropagation(); 
                leaveRoom(roomId); 
              }}
              className="text-red-500 hover:text-red-700"
            >
              <X className="w-4 h-4" />
            </button>
          </div>
        </div>
      ))}
    </div>
  );
}
