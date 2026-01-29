import { useEffect, useState } from "react";
import { getLaptopList } from "../api/laptop";
import useLaptopWS from "../hooks/useLaptopWS";

export default function LaptopList() {
  const [laptops, setLaptops] = useState([]);
  const [messages, setMessages] = useState([]);
  // 초기 HTTP 요청
  useEffect(() => {
    getLaptopList().then(setLaptops);
  }, []);

  // WebSocket 수신
  useLaptopWS((msg) => {
    // 1️⃣ 원본 메시지 출력용으로 저장
    setMessages((prev) => [
      ...prev,
      {
        time: new Date().toLocaleTimeString(),
        data: msg,
      },
    ]);
  });

  return (
    <div style={{ display: "flex", gap: 40 }}>
      {/* Laptop 리스트 */}
      <div>
        <h2>Active Laptops</h2>
        <ul>
          {laptops.map((id) => (
            <li key={id}>{id}</li>
          ))}
        </ul>
      </div>

      {/* 서버 메시지 출력 */}
      <div>
        <h2>Server Messages</h2>
        <ul style={{ maxHeight: 300, overflowY: "auto" }}>
          {messages.map((m, idx) => (
            <li key={idx}>
              <strong>{m.time}</strong>
              <pre style={{ margin: 0 }}>
                {JSON.stringify(m.data, null, 2)}
              </pre>
            </li>
          ))}
        </ul>
      </div>
    </div>
  );
}