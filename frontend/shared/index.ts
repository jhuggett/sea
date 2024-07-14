import { JSONRPC, JSONRPCError } from "./backend";
export * from "./backend";

export const logger = {
  log: console.info,
  error: console.error,
  warn: console.warn,
  info: console.info,
  debug: console.debug,
};

export function SetupRPC(conn: WebSocket) {
  logger.log("Setting up RPC");

  return new JSONRPC(conn);
}
