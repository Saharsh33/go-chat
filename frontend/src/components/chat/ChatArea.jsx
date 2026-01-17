import React from 'react';
import { MessageSquare } from 'lucide-react';
import { useChat } from '../../context/ChatContext';
import MessageList from './MessageList';
import MessageInput from './MessageInput';

export default function ChatArea() {
  const { selectedChat } = useChat();

  if (!selectedChat) {
    return (
      <div className="flex-1 flex items-center justify-center text-gray-400">
        <div className="text-center">
          <MessageSquare className="w-16 h-16 mx-auto mb-4 opacity-50" />
          <p className="text-lg">Select a chat to start messaging</p>
        </div>
      </div>
    );
  }

  return (
    <div className="flex-1 flex flex-col">
      <div className="bg-white border-b border-gray-200 p-4">
        <h3 className="text-lg font-semibold text-gray-800">
          {selectedChat.type === 'room' 
            ? `Room ${selectedChat.id}` 
            : selectedChat.id}
        </h3>
      </div>
      <MessageList />
      <MessageInput />
    </div>
  );
}