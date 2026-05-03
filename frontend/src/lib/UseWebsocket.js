// lib/WebSocketContext.js
"use client";
import { createContext, useContext, useEffect, useState } from "react";

const WebSocketContext = createContext(null);

export function WebSocketProvider({ children }) {
  const [port, setPort] = useState(null);
  const [portKey, setPortKey] = useState("");
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [messages, setMessages] = useState([]);

  useEffect(() => {
    const worker = new SharedWorker("/ws-worker.js");
    worker.port.start();

    worker.port.onmessage = (msg) => {
      const event = msg.data.event;
      console.log(event);

      switch (event) {
        case "connected": {
          setPortKey(msg.data.portKey);
          setPort(worker.port);
          worker.port.postMessage({ type: "connect" });
          break;
        }

        case "online_users": {
          setOnlineUsers(msg.data.users);
          break;
        }

        case "own_message": {
          setMessages((oldMsg) => [...new Set([...oldMsg, msg.data.message])]);

          break;
        }

        case "other_message": {
          setMessages((oldMsg) => [...new Set([...oldMsg, msg.data.message])]);

          break;
        }

        case "history": {
          setMessages((oldMsg) => [
            ...new Set([...msg.data.messages, ...oldMsg]),
          ]);
        }

        case "join": {
          setOnlineUsers((prev) => [...new Set([...prev, msg.data.newcomer])]);
          break;
        }

        case "leave": {
          setOnlineUsers((prev) => prev.filter((u) => u !== msg.data.left));
          break;
        }
      }
    };

    return () => {
      worker.port.postMessage({ type: "disconnect-tab", portKey });
      worker.port.close();
    };
  }, []);

  const sendMessage = (msg) => {
    port.postMessage(msg);
  };

  return (
    <WebSocketContext.Provider
      value={{ port, portKey, onlineUsers, messages, sendMessage: sendMessage }}
    >
      {children}
    </WebSocketContext.Provider>
  );
}

export function useWebSocket() {
  return useContext(WebSocketContext);
}
