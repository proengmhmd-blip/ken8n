// Temporary stub for development
export function createOpencodeClient(config?: any) {
  console.log("createOpencodeClient called with config:", config);
  return {
    // Add any methods that might be called
    fetch: async (...args: any[]) => {
      console.log("OpencodeClient fetch called with args:", args);
      return new Response("{}", { status: 200 });
    }
  };
}

export type Config = {
  baseUrl?: string;
  fetch?: any;
};

export class OpencodeClient {
  constructor(options: { client: any }) {
    console.log("OpencodeClient created with options:", options);
  }
}

export type ServerConfig = {
  host?: string;
  port?: number;
};

export async function createOpencodeServer(config?: ServerConfig) {
  config = Object.assign(
    {
      host: "127.0.0.1",
      port: 4096,
    },
    config ?? {},
  );

  console.log("createOpencodeServer called with config:", config);
  
  return {
    url: `http://${config.host}:${config.port}`,
    close() {
      console.log("OpencodeServer closed");
    },
  };
}

// Export any other types that might be needed
export * from "./gen/types.gen.js";