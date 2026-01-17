import React, { useState } from 'react';
import { ChatProvider, useChat } from './context/ChatContext';
import Login from './components/Auth/Login';
import Sidebar from './components/Sidebar/Sidebar';
import ChatArea from './components/Chat/ChatArea';
import CustomJsonModal from './components/Modals/CustomJsonModal';

function ChatApp() {
  const { connected, connect } = useChat();
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [showJsonModal, setShowJsonModal] = useState(false);

  const handleLogin = async (username) => {
    await connect(username);
    setIsLoggedIn(true);
  };

  if (!isLoggedIn) {
    return <Login onLogin={handleLogin} />;
  }

  return (
    <div className="h-screen bg-gray-100 flex">
      <Sidebar onOpenJsonModal={() => setShowJsonModal(true)} />
      <ChatArea />
      <CustomJsonModal 
        isOpen={showJsonModal} 
        onClose={() => setShowJsonModal(false)} 
      />
    </div>
  );
}

export default function App() {
  return (
    <ChatProvider>
      <ChatApp />
    </ChatProvider>
  );
}