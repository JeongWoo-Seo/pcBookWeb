import { useEffect } from "react";

export default function useLaptopWS(onMessage) {
  useEffect(() => {
    const ws = new WebSocket("ws://localhost:8081/ws");

    ws.onmessage = (e) => {
      onMessage(e.data);
    };

    ws.onerror = (e) => {
      console.error("WebSocket error", e);
    };

    return () => {
      ws.close();
    };
  }, [onMessage]);
}