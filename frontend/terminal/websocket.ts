import { sleep } from "bun";

class FailedToConnectError extends Error {
  constructor() {
    super("Failed to connect to server after 10 attempts, exiting.");
  }
}

export const connectToBackend = async () => {
  var Connection: WebSocket | undefined;
  var attemptedConnections = 0;
  while (attemptedConnections < 8) {
    const socket = new WebSocket("ws://localhost:8080/ws");

    const { promise, reject, resolve } = Promise.withResolvers();

    socket.onopen = () => {
      resolve(socket);
    };

    socket.onerror = (error) => {
      reject(error);
    };

    try {
      await promise;
      Connection = socket;
      break;
    } catch (error) {
      attemptedConnections++;
      await sleep(250);
    }
  }

  if (!Connection) {
    throw new FailedToConnectError();
  }

  return Connection;
};
