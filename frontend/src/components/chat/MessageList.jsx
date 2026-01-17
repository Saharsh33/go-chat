import React, { useRef, useEffect } from 'react';
import { useChat } from '../../context/ChatContext';
import Message from './Message';

export default function MessageList() {
  const { 
    username, 
    selectedChat, 
    roomMessages, 
    directMessages,
    loadMoreRoomMessages,
    loadMoreDirectMessages
  } = useChat();
  
  const messagesEndRef = useRef(null);
  const messagesContainerRef = useRef(null);

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  useEffect(() => {
    scrollToBottom();
  }, [roomMessages, directMessages, selectedChat]);

  const handleScroll = () => {
    if (messagesContainerRef.current) {
      const { scrollTop } = messagesContainerRef.current;
      
      // TODO: Load more messages when scrolled to top
      if (scrollTop === 0) {
        if (selectedChat?.type === 'room') {
          // loadMoreRoomMessages(selectedChat.id);
        } else if (selectedChat?.type === 'direct') {
          // loadMoreDirectMessages(selectedChat.id);
        }
      }
    }
  };

  const getCurrentMessages = () => {
    if (!selectedChat) return [];
    
    if (selectedChat.type === 'room') {
      return roomMessages[selectedChat.id] || [];
    } else {
      return directMessages[selectedChat.id] || [];
    }
  };

  const messages = getCurrentMessages();

  return (
    <div 
      ref={messagesContainerRef}
      onScroll={handleScroll}
      className="flex-1 overflow-y-auto p-4 space-y-3"
    >
      {/* TODO: Add loading indicator when fetching more messages */}
      {/* {isLoadingMore && (
        <div className="text-center py-2">
          <span className="text-gray-500">Loading more messages...</span>
        </div>
      )} */}
      
      {messages.map((msg, idx) => (
        <Message 
          key={msg.id || idx} 
          message={msg} 
          currentUsername={username} 
        />
      ))}
      
      <div ref={messagesEndRef} />
    </div>
  );
}