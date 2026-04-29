// init the websocket inside the shared worker
let socket = null;
let ports = [];

// this function to broadcast to all the tabs 
function broadcast(msg) {
    ports.forEach((tab) => tab.postMessage(msg));
}

// this function the default for the sharedworker to work
onconnect = function (e) {
    // push the new tab
    const port = e.ports[0];
    ports.push(port);
    port.start();

    port.onmessage = (e) => {
        switch (e.data.type) {
            case "init":
                if (!socket) {
                    socket = new WebSocket("ws://localhost:4001/ws");

                    socket.onopen = () =>
                        broadcast({ type: "status", data: "connected" });

                    socket.onmessage = (msg) =>
                        broadcast({ type: "message", data: msg.data });

                    socket.onclose = () => (socket = null);
                }
                break;

            case "send":
                socket?.send(e.data.payload);
                break;
        }
    }
}