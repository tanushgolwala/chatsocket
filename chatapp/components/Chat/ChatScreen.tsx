import React, { useEffect, useState } from 'react';
import { View, TextInput, Button, FlatList, Text, StyleSheet } from 'react-native';
import useWebSocket from '../../hooks/useWebSocket';

interface Message {
  from: string;
  content: string;
}

const ChatScreen = () => {
  const { sendMessage, socket } = useWebSocket();
  const [to, setTo] = useState('');
  const [message, setMessage] = useState('');
  const [chat, setChat] = useState<Message[]>([]);

  useEffect(() => {
    if (!socket) return;

    socket.onmessage = (event) => {
      const data = JSON.parse(event.data);
      setChat((prevChat) => [...prevChat, { from: data.from, content: data.content }]);
    };
  }, [socket]);

  const handleSendMessage = () => {
    sendMessage(to, 'user123', message); // Replace 'user123' with the current user's ID
    setMessage('');
  };

  return (
    <View style={styles.container}>
      <FlatList
        data={chat}
        keyExtractor={(_, index) => index.toString()}
        renderItem={({ item }) => (
          <Text style={styles.message}>
            {item.from}: {item.content}
          </Text>
        )}
      />
      <TextInput
        placeholder="Recipient ID"
        value={to}
        onChangeText={setTo}
        style={styles.input}
      />
      <TextInput
        placeholder="Type a message"
        value={message}
        onChangeText={setMessage}
        style={styles.input}
      />
      <Button title="Send" onPress={handleSendMessage} />
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 20,
    flex: 1,
  },
  input: {
    borderWidth: 1,
    marginBottom: 10,
    padding: 8,
    borderRadius: 5,
  },
  message: {
    padding: 10,
    borderBottomWidth: 1,
    borderColor: '#ddd',
  },
});

export default ChatScreen;
