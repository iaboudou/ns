let socket = null;
const ports = new Map(); // (key : port (or if you want tab) )

//send specific data to all the tabs
function broadcast(data) {
  ports.forEach((port) => port.postMessage(data));
}

//init the logic of the worker
onconnect = function (e) {
  const port = e.ports[0];
  const key = crypto.randomUUID(); // random key for the tab in the ports map

  ports.set(key, port);

  // send the the specified tab his key for future data exchange that require
  port.postMessage({
    event: "connected",
    portKey: key,
  });

  //depending on the "event type" do some logics
  port.onmessage = function (event) {
    const msg = event.data;

    switch (msg.type) {
      case "connect": {
        // first time we open the tab, check if there is no ws
        if (!socket) {
          socket = new WebSocket("ws://localhost:4001/ws");

          socket.onmessage = (e) => {
            try {
              const data = JSON.parse(e.data);
              console.log(data);
              broadcast(data);
            } catch (err) {
              console.error("WS parse error", err);
            }
          };

          socket.onclose = () => {
            socket = null;
            broadcast({ event: "ws-close" }); // send to all the tabs that the ws is closed
          };

          socket.onerror = (err) => {
            console.error("WS error", err);
          };
        }
        break;
      }

      case "send": {
        socket.send(JSON.stringify(msg.payload)); // send the data to the backend
        break;
      }

      case "disconnect-tab": {
        ports.delete(msg.portKey); // remove one tab
        break;
      }

      case "logout": {
        socket.close();
      }
    }
  };

  port.start();
};
