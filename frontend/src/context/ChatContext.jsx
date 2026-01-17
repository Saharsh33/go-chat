import React, { createContext, useContext, useState, useCallback } from 'react';
import websocketService from '../services/websocketService';
import { MESSAGE_TYPES } from '../utils/constants';

const ChatContext = createContext();

export const ChatProvider = ({ children }) => {
  const [connected, setConnected] = useState(false);
  const [username, setUsername] = useState('');
  
  // Room data
  const [rooms, setRooms] = useState([]);
  const [roomMessages, setRoomMessages] = useState({});
  const [roomLastMessageIds, setRoomLastMessageIds] = useState({});
  
  // Direct message data
  const [directChats, setDirectChats] = useState([]);
  const [directMessages, setDirectMessages] = useState({});
  const [directLastMessageIds, setDirectLastMessageIds] = useState({});
  
  const [selectedChat, setSelectedChat] = useState(null); // { type: 'room'|'direct', id: roomId|username }

  const handleIncomingMessage = useCallback((msg) => {
    console.log('Received message:', msg);

    // TODO: HANDLE DIFFERENT MESSAGE TYPES
    switch (msg.type) {
      case MESSAGE_TYPES.ROOM_MESSAGE:
        // TODO: Add message to room
        // setRoomMessages(prev => ({
        //   ...prev,
        //   [msg.room]: [...(prev[msg.room] || []), msg]
        // }));
        // TODO: Update last message ID for pagination
        // if (msg.id) {
        //   setRoomLastMessageIds(prev => ({
        //     ...prev,
        //     [msg.room]: msg.id
        //   }));
        // }
        break;

      case MESSAGE_TYPES.DIRECT_MESSAGE:
        // TODO: Add message to direct chat
        // const chatWith = msg.user === username ? msg.receiver : msg.user;
        // setDirectMessages(prev => ({
        //   ...prev,
        //   [chatWith]: [...(prev[chatWith] || []), msg]
        // }));
        // TODO: Update last message ID for pagination
        // if (msg.id) {
        //   setDirectLastMessageIds(prev => ({
        //     ...prev,
        //     [chatWith]: msg.id
        //   }));
        // }
        break;

      case MESSAGE_TYPES.SYSTEM:
        // TODO: Handle system messages
        // if (msg.room !== undefined) {
        //   setRoomMessages(prev => ({
        //     ...prev,
        //     [msg.room]: [...(prev[msg.room] || []), msg]
        //   }));
        // }
        break;

      default:
        console.log('Unknown message type:', msg.type);
    }
  }, [username]);

  const connect = async (name) => {
    try {
      await websocketService.connect(name);
      websocketService.onMessage(handleIncomingMessage);
      setUsername(name);
      setConnected(true);
    } catch (error) {
      console.error('Failed to connect:', error);
    }
  };

  const sendMessage = (message) => {
    websocketService.send(message);
  };

  // Room actions
  const createRoom = (roomId) => {
    sendMessage({
      type: MESSAGE_TYPES.CREATE_ROOM,
      user: username,
      room: roomId
    });
    // TODO: Add room to list after successful creation
    // setRooms(prev => [...new Set([...prev, roomId])]);
  };

  const joinRoom = (roomId) => {
    sendMessage({
      type: MESSAGE_TYPES.JOIN_ROOM,
      user: username,
      room: roomId
    });
    setSelectedChat({ type: 'room', id: roomId });
  };

  const leaveRoom = (roomId) => {
    sendMessage({
      type: MESSAGE_TYPES.LEAVE_ROOM,
      user: username,
      room: roomId
    });
    // TODO: Remove room from list
    // setRooms(prev => prev.filter(r => r !== roomId));
    if (selectedChat?.type === 'room' && selectedChat.id === roomId) {
      setSelectedChat(null);
    }
  };

  const sendRoomMessage = (roomId, content) => {
    sendMessage({
      type: MESSAGE_TYPES.ROOM_MESSAGE,
      user: username,
      room: roomId,
      content
    });
  };

  // Direct message actions
  const sendDirectMessage = (receiver, content) => {
    sendMessage({
      type: MESSAGE_TYPES.DIRECT_MESSAGE,
      user: username,
      receiver,
      content
    });
  };

  // Pagination
  const loadMoreRoomMessages = (roomId) => {
    const lastMessageId = roomLastMessageIds[roomId];
    // TODO: Send request for more messages
    // sendMessage({
    //   type: MESSAGE_TYPES.NEXT_ROOM_MESSAGES,
    //   user: username,
    //   room: roomId,
    //   lastMessageId: lastMessageId // Adjust based on your backend API
    // });
  };

  const loadMoreDirectMessages = (receiver) => {
    const lastMessageId = directLastMessageIds[receiver];
    // TODO: Send request for more messages
    // sendMessage({
    //   type: MESSAGE_TYPES.NEXT_DIRECT_MESSAGES,
    //   user: username,
    //   receiver: receiver,
    //   lastMessageId: lastMessageId // Adjust based on your backend API
    // });
  };

  const value = {
    connected,
    username,
    rooms,
    roomMessages,
    directChats,
    directMessages,
    selectedChat,
    connect,
    sendMessage,
    createRoom,
    joinRoom,
    leaveRoom,
    sendRoomMessage,
    sendDirectMessage,
    setSelectedChat,
    loadMoreRoomMessages,
    loadMoreDirectMessages
  };

  return <ChatContext.Provider value={value}>{children}</ChatContext.Provider>;
};

export const useChat = () => {
  const context = useContext(ChatContext);
  if (!context) {
    throw new Error('useChat must be used within ChatProvider');
  }
  return context;
};