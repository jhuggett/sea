import type { Inbound as GeneratedInbound } from "../generated/inbound";
import type { Outbound as GeneratedOutbound } from "../generated/outbound";

// Outbound as in what the backend will send, Inbound as in what the frontend will receive
type Inbound = {
  [K in keyof GeneratedOutbound]: {
    req: GeneratedOutbound[K]["req"];
    resp: GeneratedOutbound[K]["resp"];
  };
};

// Inbound as in what the backend will receive, Outbound as in what the frontend will send
type Outbound = {
  [K in keyof GeneratedInbound]: {
    req: GeneratedInbound[K]["req"];
    resp: GeneratedInbound[K]["resp"];
  };
};

type Request<T> = {
  method: string;
  params: T;
  id: string;
};

const isRequest = <T>(data: any): data is Request<T> => {
  return (
    data.method !== undefined &&
    data.params !== undefined &&
    data.id !== undefined
  );
};

type Response<T> = {
  result?: T;
  error?: {
    code: JSONRPCError;
    message: string;
    data?: any;
  };
  id: string;
};

// code	message	meaning
// -32700	Parse error	Invalid JSON was received by the server.
// An error occurred on the server while parsing the JSON text.
// -32600	Invalid Request	The JSON sent is not a valid Request object.
// -32601	Method not found	The method does not exist / is not available.
// -32602	Invalid params	Invalid method parameter(s).
// -32603	Internal error	Internal JSON-RPC error.
// -32000 to -32099	Server error

export enum JSONRPCError {
  ParseError = -32700,
  InvalidRequest = -32600,
  MethodNotFound = -32601,
  InvalidParams = -32602,
  InternalError = -32603,
  ServerError = -32000,
}

const isResponse = <T>(data: any): data is Response<T> => {
  return (
    data.id !== undefined &&
    (data.result !== undefined || data.error !== undefined)
  );
};

export class JSONRPC {
  constructor(private connection: WebSocket) {
    this.connection.onmessage = (message) => {
      const data = JSON.parse(message.data);
      if (isRequest(data)) {
        this.handleRequest(data);
      } else if (isResponse(data)) {
        this.handleResponse(data);
      }
    };
  }

  private sendResponse(data: Response<any>) {
    this.connection.send(JSON.stringify(data));
  }

  private sendRequest(data: Request<any>) {
    this.connection.send(JSON.stringify(data));
  }

  private async handleRequest(data: Request<unknown>) {
    const { method, params, id } = data;

    console.log("Handling request", data);

    const callback = this.requestRegistry.get(method);

    if (!callback) {
      return this.sendResponse({
        id,
        error: {
          code: JSONRPCError.MethodNotFound,
          message: "Method not found",
        },
      });
    }

    const resp = await callback(params);

    this.sendResponse({ id, ...resp });
  }

  private handleResponse(data: Response<unknown>) {
    const { id, result, error } = data;

    const cb = this.responseRegistry.get(id);
    if (!cb) {
      return;
    }

    if (error) {
      cb.reject(error);
    } else {
      cb.resolve(result);
    }
  }

  private requestRegistry: Map<
    string,
    (resp: any) => Promise<Omit<Response<any>, "id">>
  > = new Map();

  private responseRegistry: Map<string, { resolve: any; reject: any }> =
    new Map();

  /**
   * Calling register for the same method multiple times will overwrite the previous callback.
   */
  receive<T extends keyof Inbound>(
    method: T,
    callback: (
      params: Inbound[T]["req"]
    ) => Promise<Omit<Response<Inbound[T]["resp"]>, "id">>
  ) {
    this.requestRegistry.set(method.toString(), callback);

    return () => {
      this.requestRegistry.delete(method.toString());
    };
  }

  send<T extends keyof Outbound>(
    method: T,
    params: Outbound[T]["req"]
  ): Promise<Outbound[T]["resp"]> {
    const id = Math.random().toString(36).slice(2);

    const { promise, reject, resolve } =
      Promise.withResolvers<Outbound[T]["resp"]>();

    this.responseRegistry.set(id, { reject, resolve });

    this.sendRequest({ method: method.toString(), params, id });

    return promise;
  }
}
