// Docker Desktop Extension API type declarations

interface DDClient {
  host: {
    openExternal: (url: string) => void;
  };
  // Add other Docker Desktop extension APIs as needed
}

declare global {
  interface Window {
    ddClient?: DDClient;
  }
}

export {};