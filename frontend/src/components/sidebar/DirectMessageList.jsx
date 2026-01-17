import React from 'react';
import { Plus } from 'lucide-react';
import { useChat } from '../../context/ChatContext';

export default function DirectMessageList() {
  const { directChats, selectedChat, setSelectedChat } = useChat();

  const handleNewDirectMessage = () => {
    const receiver = prompt('Enter username to message:');
    if (receiver) {
      setSelectedChat({ type: 'direct', id: receiver });
      // TODO: Optionally add to directChats list
    }
  };

  return (
    <div className="p-2">
      <button
        onClick={handleNewDirectMessage}
        className="w-full mb-2 p-3 bg-blue-50 text-blue-700 rounded-lg hover:bg-blue-100 transition-colors flex items-center justify-center gap-2"
      >
        <Plus className="w-4 h-4" />
        New Message
      </button>
      
      {/* TODO: Map through directChats array */}
      {directChats.map(username => (
        <div
          key={username}
          onClick={() => setSelectedChat({ type: 'direct', id: username })}
          className={`p-3 mb-2 rounded-lg cursor-pointer transition-colors ${
            selectedChat?.type === 'direct' && selectedChat?.id === username
              ? 'bg-blue-100 border-l-4 border-blue-600'
              : 'bg-gray-50 hover:bg-gray-100'
          }`}
        >
          <span className="font-medium">{username}</span>
        </div>
      ))}
    </div>
  );
}
