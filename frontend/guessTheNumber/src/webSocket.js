import { ref, onUnmounted } from 'vue';

const ws = new WebSocket('ws://localhost:8080/ws');

export const receivedMessage = ref(null);

// Function to send action to the server
export const sendAction = (action) => {
  if (ws.readyState === WebSocket.OPEN) {
    ws.send(JSON.stringify(action));
  } else {
    console.error('WebSocket connection not open.');
  }
};

// Hook to manage WebSocket connection
export function useWebSocket() {
  const connected = ref(false);

  ws.onopen = () => {
    console.log('WebSocket connected');
    connected.value = true;
  };

  ws.onmessage = (event) => {
    console.log('Message received:', event.data);
    // Handle incoming messages from the server
    receivedMessage.value = JSON.parse(event.data);

  };

  ws.onclose = () => {
    console.log('WebSocket disconnected');
    connected.value = false;
    // Reconnect logic can be implemented here if needed
  };

  ws.onerror = (error) => {
    console.error('WebSocket error:', error);
  };

  onUnmounted(() => {
    // Close WebSocket connection when component is unmounted
    if (ws.readyState === WebSocket.OPEN) {
      ws.close();
    }
  });

  return {
    connected
  };
}
