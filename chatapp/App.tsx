import React from 'react';
import ChatScreen from './components/Chat/ChatScreen';
import { WebSocketProvider } from './contexts/WebSocketContext';

const App = () => {
  return (
    <WebSocketProvider>
      <ChatScreen />
    </WebSocketProvider>
  );
};

export default App;
