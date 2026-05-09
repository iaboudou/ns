let socket = null;
const ports = new Map();
const currenTab = new Map();
let pendingMessages = [];
let onlineUsers = [];

function broadcast(data) {
  ports.forEach((port) => {
    port.postMessage(data);
  });
}

self.addEventListener("connect", (e) => {
  const port = e.ports[0];
  const key = crypto.randomUUID();

  ports.set(key, port);

  port.postMessage({
    event: "connected",
    portKey: key,
  });

  port.postMessage({
    event: "online_users",
    users: onlineUsers,
  });

  port.onmessage = (event) => {
    const msg = event.data;

    switch (msg.type) {
      case "connect": {
        if (!socket) {
          socket = new WebSocket("ws://localhost:4001/ws");

          socket.onopen = () => {
            pendingMessages.forEach((m) => socket.send(JSON.stringify(m)));
            pendingMessages = [];
          };

          socket.onmessage = (e) => {
            try {
              const data = JSON.parse(e.data);
              if (data.event === "online_users") onlineUsers = data.users || [];
              if (data.event === "join" && !onlineUsers.includes(data.newcomer)) onlineUsers.push(data.newcomer);
              if (data.event === "leave") onlineUsers = onlineUsers.filter((id) => id !== data.left);
              broadcast(data);
            } catch (err) {}
          };

          socket.onclose = () => {
            socket = null;
            pendingMessages = [];
            onlineUsers = [];
            broadcast({ event: "ws-close" });
          };

          socket.onerror = (err) => {};
        }
        break;
      }

      case "send": {
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify(msg.payload));
        } else if (socket) {
          pendingMessages.push(msg.payload);
        }
        break;
      }

      case "notifs_seen": {
        broadcast({ event: "notifs_seen" });
        break;
      }

      case "set_notif": {
        broadcast({ event: "set_notif", notif: msg.notif });
        break;
      }

      case "disconnect-tab": {
        ports.delete(msg.portKey);
        currenTab.delete(msg.portKey);
        if (ports.size === 0 && socket) {
          socket.close();
        }
        break;
      }

      case "logout": {
        onlineUsers = [];
        socket.close();
        break;
      }
    }
  };

  port.start();
});
