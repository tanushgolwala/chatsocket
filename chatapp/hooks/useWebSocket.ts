import { useContext } from 'react';
import { WebSocketContext } from '../contexts/WebSocketContext';

const useWebSocket = () => {
  const socket = useContext(WebSocketContext);

  const sendMessage = (to: string, from: string, content: string) => {
    if (!socket) {
      console.error('WebSocket is not connected');
      return;
    }

    const message = {
      to,
      from,
      content,
    };

    socket.send(JSON.stringify(message));
  };

  return { sendMessage, socket };
};

export default useWebSocket;
