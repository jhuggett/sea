import { overrideConsole } from "./log";
import { BunShell, UnknownKeyCodeError } from "@jhuggett/terminal";
import { startBackend } from "../backend/process";
import { connectToBackend } from "../backend/websocket";
import { JSONRPC } from "@sea/shared";

overrideConsole();

console.log("Starting frontend...");
console.warn("This is a warning.");
console.error("This is an error.");
console.info("This is info.");
console.debug("This is debug.");

const { stopBackend } = startBackend();

const connection = await connectToBackend();

const rpc = new JSONRPC<
  {
    ping: { req: { message: string }; resp: { message: string } };
  },
  {
    hello: { req: { message: string }; resp: { message: string } };
  }
>(connection);

rpc.receive("ping", async ({ method, params }) => {
  console.log("Received ping:", params);
  return {
    result: {
      message: "Pong!",
    },
  };
});

rpc.send("hello", { message: "Hello, world!" }).then((response) => {
  console.log("Received response:", response);
});

const shell = new BunShell();

shell.rootElement.renderer = ({ cursor }) => {
  cursor.write("Hello, world!");
};

shell.clear();

shell.rootElement.render();

shell.render();

shell.rootElement.focus();

let shouldExit = false;

shell.rootElement.on("q", () => {
  shouldExit = true;
});

shell.rootElement.on("Any key", (key) => {
  console.log("Pressed key:", key);
});

while (!shouldExit) {
  try {
    await shell.userInteraction();
  } catch (error) {
    if (error instanceof UnknownKeyCodeError) {
      console.error("Unknown key code:", error);
    } else {
      throw error;
    }
  }
}

stopBackend();
