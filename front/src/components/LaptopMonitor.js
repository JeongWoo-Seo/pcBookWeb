import { useEffect, useState } from "react";
import useLaptopWS from "../hooks/useLaptopWS";

export default function LaptopMonitor({ laptopId }) {
  const [data, setData] = useState(null);

  // laptopId 변경 시 초기화
  useEffect(() => {
    setData(null);
  }, [laptopId]);

  useLaptopWS((msg) => {
    // 서버는 string 또는 JSON 보낸다고 했지
    const parsed = JSON.parse(msg);

    // 🔥 선택된 laptop만 처리
    if (parsed.id !== laptopId) return;

    setData(parsed);
  });

  if (!laptopId) {
    return <div>Select a laptop to monitor</div>;
  }

  if (!data) {
    return <div>Waiting for data from {laptopId}...</div>;
  }

  return (
    <div style={{ flex: 1 }}>
      <h3>Monitoring: {laptopId}</h3>

      <div>CPU: {data.cpu.toFixed(2)}%</div>
      <div>RAM: {data.ram.usage.toFixed(2)}%</div>
      <div>Storage: {data.storages.usage.toFixed(2)}%</div>
      <div>Battery: {data.battery}%</div>
      <div>
        Network: RX {data.network.rx} / TX {data.network.tx}
      </div>
    </div>
  );
}
