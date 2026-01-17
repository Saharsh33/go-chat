class WebSocketService {
    constructor() {
      this.ws = null;
      this.messageHandlers = [];
    }
  
    connect(username) {
      return new Promise((resolve, reject) => {
        this.ws = new WebSocket(WS_URL);
  
        this.ws.onopen = () => {
          console.log('WebSocket connected');
          // SEND USERNAME ON FIRST CONNECTION
          this.send({ user: username });
          resolve();
        };
  
        this.ws.onmessage = (event) => {
          try {
            const message = JSON.parse(event.data);
            this.messageHandlers.forEach(handler => handler(message));
          } catch (error) {
            console.error('Failed to parse message:', error);
          }
        };
  
        this.ws.onerror = (error) => {
          console.error('WebSocket error:', error);
          reject(error);
        };
  
        this.ws.onclose = () => {
          console.log('WebSocket disconnected');
          // TODO: Add reconnection logic here
        };
      });
    }
  
    send(message) {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        this.ws.send(JSON.stringify(message));
      }
    }
  
    onMessage(handler) {
      this.messageHandlers.push(handler);
    }
  
    disconnect() {
      if (this.ws) {
        this.ws.close();
      }
    }
  }
  
  export default new WebSocketService();