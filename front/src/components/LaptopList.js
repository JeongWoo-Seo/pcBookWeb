import { useEffect, useState } from "react";
import { fetchLaptopList } from "../api/laptop";
import useLaptopWS from "../hooks/useLaptopWS";

export default function LaptopList() {
  const [laptops, setLaptops] = useState([]);

  // 초기 HTTP 요청
  useEffect(() => {
    fetchLaptopList().then(setLaptops);
  }, []);

  // WebSocket 수신
  useLaptopWS((msg) => {
    if (msg.type === "laptop") {
      setLaptops((prev) => {
        if (prev.includes(msg.id)) return prev;
        return [...prev, msg.id];
      });
    }
  });

  return (
    <div>
      <h2>Active Laptops</h2>
      <ul>
        {laptops.map((id) => (
          <li key={id}>{id}</li>
        ))}
      </ul>
    </div>
  );
}