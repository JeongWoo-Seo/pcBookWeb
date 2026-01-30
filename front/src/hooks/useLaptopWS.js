import { useEffect, useRef } from "react";

export default function useLaptopWS(onMessage) {
  const wsRef = useRef(null);

  useEffect(() => {
    wsRef.current = new WebSocket("ws://localhost:8081/ws");

    wsRef.current.onmessage = (event) => {
      const data = JSON.parse(event.data);
      onMessage(data);
    };

    return () => {
      wsRef.current.close();
    };
  }, []);
}