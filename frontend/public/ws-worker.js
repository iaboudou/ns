let socket = null;
const ports = new Map();
const currenTab = new Map();
let pendingMessages = [];

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
              broadcast(data);
            } catch (err) {}
          };

          socket.onclose = () => {
            socket = null;
            pendingMessages = [];
            broadcast({ event: "ws-close" });
          };

          socket.onerror = (err) => {};
        }
        break;
      }

      case "send": {
        if (socket && socket.readyState === WebSocket.OPEN) {
        } else if (socket) {
          pendingMessages.push(msg.payload);
        }
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
        socket.close();
        break;
      }
    }
  };

  port.start();
});
