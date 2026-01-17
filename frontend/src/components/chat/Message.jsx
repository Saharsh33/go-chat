import React from 'react';

export default function Message({ message, currentUsername }) {
  const isOwnMessage = message.user === currentUsername;
  const isSystemMessage = message.type === 'system';

  if (isSystemMessage) {
    return (
      <div className="flex justify-center">
        <div className="bg-yellow-100 text-yellow-800 px-4 py-2 rounded-lg text-sm text-center max-w-md">
          {message.content}
        </div>
      </div>
    );
  }

  return (
    <div className={`flex ${isOwnMessage ? 'justify-end' : 'justify-start'}`}>
      <div className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
        isOwnMessage
          ? 'bg-blue-600 text-white'
          : 'bg-gray-200 text-gray-800'
      }`}>
        <div className="text-xs opacity-75 mb-1">{message.user}</div>
        <div className="break-words">{message.content}</div>
        {/* TODO: Add timestamp if available in message object */}
        {/* {message.timestamp && (
          <div className="text-xs opacity-60 mt-1">
            {new Date(message.timestamp).toLocaleTimeString()}
          </div>
        )} */}
      </div>
    </div>
  );
}
