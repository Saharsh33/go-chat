import React, { useState } from 'react';
import { Users, MessageSquare, Settings } from 'lucide-react';
import { useChat } from '../../context/ChatContext';
import RoomList from './RoomList';
import DirectMessageList from './DirectMessageList';

export default function Sidebar({ onOpenJsonModal }) {
  const { username, connected } = useChat();
  const [activeTab, setActiveTab] = useState('rooms');

  return (
    <div className="w-80 bg-white border-r border-gray-200 flex flex-col">
      <div className="p-4 border-b border-gray-200">
        <div className="flex items-center justify-between mb-4">
          <h2 className="text-xl font-bold text-gray-800">{username}</h2>
          <div className={`w-3 h-3 rounded-full ${connected ? 'bg-green-500' : 'bg-red-500'}`} />
        </div>
        <div className="flex gap-2">
          <button
            onClick={() => setActiveTab('rooms')}
            className={`flex-1 py-2 px-3 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'rooms' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'
            }`}
          >
            <Users className="w-4 h-4 inline mr-1" />
            Rooms
          </button>
          <button
            onClick={() => setActiveTab('direct')}
            className={`flex-1 py-2 px-3 rounded-lg text-sm font-medium transition-colors ${
              activeTab === 'direct' ? 'bg-blue-100 text-blue-700' : 'bg-gray-100 text-gray-600'
            }`}
          >
            <MessageSquare className="w-4 h-4 inline mr-1" />
            Direct
          </button>
        </div>
      </div>

      <div className="flex-1 overflow-y-auto">
        {activeTab === 'rooms' ? <RoomList /> : <DirectMessageList />}
      </div>

      <div className="p-4 border-t border-gray-200">
        <button
          onClick={onOpenJsonModal}
          className="w-full py-2 px-4 bg-gray-800 text-white rounded-lg hover:bg-gray-900 transition-colors flex items-center justify-center gap-2"
        >
          <Settings className="w-4 h-4" />
          Send Custom JSON
        </button>
      </div>
    </div>
  );
}
