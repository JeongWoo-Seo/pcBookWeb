import { useEffect, useRef } from "react";

export default function useLaptopWS({ mode, laptopID, onMessage }) {
  const wsRef = useRef(null);

  useEffect(() => {
    let url = "ws://localhost:8081/ws";

    if (mode === "single" && laptopID) {
      url += `?mode=single&id=${laptopID}`;
    } else {
      url += "?mode=broadcast";
    }

    const ws = new WebSocket(url);
    wsRef.current = ws;

    ws.onmessage = (event) => {
      try {
        const data = JSON.parse(event.data);
        onMessage?.(data);
      } catch (e) {
        console.error("WS parse error", e);
      }
    };

    ws.onerror = (e) => console.error("WS error", e);
    ws.onclose = () => console.log("WS closed");

    return () => ws.close();
  }, [mode, laptopID]);
}