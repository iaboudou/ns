"use client";

import { createContext, useContext, useEffect, useState, useCallback, useRef } from "react";

const WebSocketContext = createContext(null);

export const WebSocketProvider = ({ children }) => {
  const [port, setPort] = useState(null);
  const [portKey, setPortKey] = useState(null);
  const portKeyRef = useRef(null);
  const [onlineUsers, setOnlineUsers] = useState([]);
  const [messages, setMessages] = useState([]);
  const [hasMoreMap, setHasMoreMap] = useState({});
  const [unreadNotifCount, setUnreadNotifCount] = useState(0);
  const [notifications, setNotifications] = useState([]);

  useEffect(() => {
    const worker = new SharedWorker("/ws-worker.js");
    worker.port.start();

    worker.port.onmessage = (msg) => {
      const event = msg.data.event;

      switch (event) {
        case "connected": {
          setPortKey(msg.data.portKey);
          portKeyRef.current = msg.data.portKey;
          setPort(worker.port);
          worker.port.postMessage({ type: "connect" });
          worker.port.postMessage({ type: "send", payload: { type: "get_unread_notifications_count" } });
          break;
        }

        case "online_users": {
          setOnlineUsers(msg.data.users);
          break;
        }

        case "join": {
          setOnlineUsers((prev) => {
            if (!prev.includes(msg.data.newcomer)) {
              return [...prev, msg.data.newcomer];
            }
            return prev;
          });
          break;
        }

        case "leave": {
          setOnlineUsers((prev) => prev.filter(id => id !== msg.data.left));
          break;
        }

        case "own_message": {
          setMessages((oldMsg) => {
            const map = new Map(oldMsg.map((m) => [m.id, m]));
            map.set(msg.data.message.id, { ...msg.data.message, isNew: true });
            return Array.from(map.values()).sort((a, b) => a.id - b.id);
          });
          break;
        }

        case "other_message":
        case "new_group_message": {
          setMessages((oldMsg) => {
            const map = new Map(oldMsg.map((m) => [m.id, m]));
            map.set(msg.data.message.id, { ...msg.data.message, isNew: true });
            return Array.from(map.values()).sort((a, b) => a.id - b.id);
          });
          break;
        }

        case "history": {
          setMessages((oldMsg) => {
            const map = new Map(oldMsg.map((m) => [m.id, m]));
            msg.data.messages.forEach((m) => map.set(m.id, m));
            return Array.from(map.values()).sort((a, b) => a.id - b.id);
          });
          setHasMoreMap((prev) => ({
            ...prev,
            [msg.data.receiver_id]: msg.data.hasMore,
          }));
          break;
        }

        case "notification": {
          setUnreadNotifCount((prev) => prev + 1);
          setNotifications((prev) => {
            const filtered = prev.filter(
              (n) => !(n.type === msg.data.type && n.ref_id === msg.data.ref_id)
            );
            return [msg.data, ...filtered];
          });
          break;
        }

        case "notifications_list": {
          setNotifications(msg.data.data || []);
          break;
        }

        case "unread_notifications_count": {
          setUnreadNotifCount(msg.data.count || 0);
          break;
        }
      }
    };

    const handleBeforeUnload = () => {
      if (portKeyRef.current) {
        worker.port.postMessage({ type: "disconnect-tab", portKey: portKeyRef.current });
      }
    };

    window.addEventListener("beforeunload", handleBeforeUnload);
    return () => {
      window.removeEventListener("beforeunload", handleBeforeUnload);
    };
  }, []);

  const sendMessage = useCallback(
    (message) => {
      if (port) {
        port.postMessage(message);
      }
    },
    [port]
  );

  return (
    <WebSocketContext.Provider
      value={{
        port,
        portKey,
        onlineUsers,
        messages,
        setMessages,
        hasMoreMap,
        sendMessage: sendMessage,
        unreadNotifCount,
        setUnreadNotifCount,
        notifications,
        setNotifications,
      }}
    >
      {children}
    </WebSocketContext.Provider>
  );
};

export const useWebSocket = () => {
  const context = useContext(WebSocketContext);
  if (!context) {
    throw new Error("useWebSocket must be used within a WebSocketProvider");
  }
  return context;
};
