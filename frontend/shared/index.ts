import { JSONRPC, JSONRPCError } from "./backend";
export * from "./backend";
export * from "./generated/inbound";

export var logger = {
  log: console.info,
  error: console.error,
  warn: console.warn,
  info: console.info,
  debug: console.debug,
};

export function OverrideLogger(newLogger: Partial<typeof logger>) {
  logger = {
    ...logger,
    ...newLogger,
  };
}

export function SetupRPC(conn: WebSocket) {
  logger.log("Setting up RPC");

  return new JSONRPC(conn);
}
