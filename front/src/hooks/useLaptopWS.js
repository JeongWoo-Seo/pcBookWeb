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
        // batching 대응
        const messages = event.data.split("\n");

        messages.forEach((m) => {
          if (!m) return;

          const data = JSON.parse(m);
          onMessage?.(data);
        });

      } catch (e) {
        console.error("WS parse error", e, event.data);
      }
    };

    ws.onerror = (e) => console.error("WS error", e);
    ws.onclose = () => console.log("WS closed");

    return () => ws.close();
  }, [mode, laptopID]);
}