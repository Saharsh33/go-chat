export const MESSAGE_TYPES = {
    // Client → Server
    JOIN_ROOM: 'join',
    LEAVE_ROOM: 'leave',
    CREATE_ROOM: 'create',
    ROOM_MESSAGE: 'messageRoom',
    DIRECT_MESSAGE: 'messageDirect',
    BROADCAST: 'broadcast',
    NEXT_ROOM_MESSAGES: 'nextRoomMsgs',
    NEXT_DIRECT_MESSAGES: 'nextDirectMsgs',
    // Server → Client
    SYSTEM: 'system'
  };
  
  export const WS_URL = 'ws://localhost:3000/ws';