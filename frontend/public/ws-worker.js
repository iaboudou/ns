let socket = null;
const ports = new Map(); // (key : port (or if you want tab) )
const currenTab = new Map();
let latestOnlineUsers = [];
let hasOnlineUsersSnapshot = false;

console.log("whe are in worker ....")
//send specific data to all the tabs
function broadcast(data) {
  ports.forEach((port) => {
    port.postMessage(data);
  });
}

//init the logic of the worker
self.addEventListener("connect", (e) => {
  const port = e.ports[0];
  const key = crypto.randomUUID(); // random key for the tab in the ports map

  ports.set(key, port);

  // send the the specified tab his key for future data exchange that require
  port.postMessage({
    event: "connected",
    portKey: key,
  });

  if (hasOnlineUsersSnapshot) {
    port.postMessage({
      event: "online_users",
      users: latestOnlineUsers,
    });
  }

  //depending on the "event type" do some logics
  port.onmessage = (event) => {
    const msg = event.data;

    switch (msg.type) {
      case "connect": {
        // first time we open the tab, check if there is no ws
        if (!socket) {
          socket = new WebSocket("ws://localhost:4001/ws");

          socket.onmessage = (e) => {
            try {
              const data = JSON.parse(e.data);
              switch (data.event) {
                case "online_users":
                  latestOnlineUsers = data.users || [];
                  hasOnlineUsersSnapshot = true;
                  break;
                case "join":
                  hasOnlineUsersSnapshot = true;
                  if (!latestOnlineUsers.includes(data.newcomer)) {
                    latestOnlineUsers = [...latestOnlineUsers, data.newcomer];
                  }
                  break;
                case "leave":
                  hasOnlineUsersSnapshot = true;
                  latestOnlineUsers = latestOnlineUsers.filter(
                    (id) => id !== data.left,
                  );
                  break;

                case "notification": {
                  // if the user already open notification will not update the notif counter 
                  data.showNotif = true;
                  for (let tabType of currenTab.values()) {
                    if (tabType === "notification") {
                      data.showNotif = false;
                    }
                  }
                  break;
                }
              }
              broadcast(data);
            } catch (err) { }
          };

          socket.onclose = () => {
            socket = null;
            latestOnlineUsers = [];
            hasOnlineUsersSnapshot = false;
            broadcast({ event: "ws-close" }); // send to all the tabs that the ws is closed
          };

          socket.onerror = (err) => { };
        }
        break;
      }

      case "send": {
        if (socket && socket.readyState === WebSocket.OPEN) {
          socket.send(JSON.stringify(msg.payload));
          
          if (msg.payload && msg.payload.type === "mark_read") {
            broadcast({
              event: "messages_read",
              receiver_Id: msg.payload.receiver_Id
            });
          }
        }
        break;
      }

      case "disconnect-tab": {
        ports.delete(msg.portKey); // remove one tab
        currenTab.delete(msg.portKey);
        if (ports.size === 0 && socket) {
          socket.close();
          socket = null;
        }
        break;
      }

      case "logout": {
        socket.close();
        break;
      }

      case "focus": {
        currenTab.set(msg.payload.portKey, msg.payload.tab);
        //
        switch (msg.payload.tab) {
          case "notification":
            broadcast({
              event: "unread_notifications_count",
              count: 0
            });
            break;
        }
      }
    }
  };

  port.start();
});
