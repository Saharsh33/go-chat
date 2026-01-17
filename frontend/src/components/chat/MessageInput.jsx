import React, { useState } from 'react';
import { Send } from 'lucide-react';
import { useChat } from '../../context/ChatContext';

export default function MessageInput() {
  const { selectedChat, sendRoomMessage, sendDirectMessage } = useChat();
  const [message, setMessage] = useState('');

  const handleSendMessage = () => {
    if (!message.trim() || !selectedChat) return;

    if (selectedChat.type === 'room') {
      sendRoomMessage(selectedChat.id, message);
    } else if (selectedChat.type === 'direct') {
      sendDirectMessage(selectedChat.id, message);
    }

    setMessage('');
  };

  const handleKeyPress = (e) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSendMessage();
    }
  };

  return (
    <div className="bg-white border-t border-gray-200 p-4">
      <div className="flex gap-2">
        <input
          type="text"
          value={message}
          onChange={(e) => setMessage(e.target.value)}
          onKeyPress={handleKeyPress}
          className="flex-1 px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-blue-500 focus:border-transparent"
          placeholder="Type a message..."
          disabled={!selectedChat}
        />
        <button
          onClick={handleSendMessage}
          disabled={!selectedChat || !message.trim()}
          className="bg-blue-600 text-white p-2 rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <Send className="w-5 h-5" />
        </button>
      </div>
    </div>
  );
}