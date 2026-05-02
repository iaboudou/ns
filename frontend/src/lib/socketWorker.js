let worker;

// this function is to get et the socket worker when ever needed or initlize it
export function getSocketWorker() {
    if (typeof window === "undefined") return null;

    if (!worker) {
        worker = new SharedWorker("/ws-worker.js");
        worker.port.start();
        worker.port.postMessage({ type: "init" });
    }

    return worker;
}