import { useState } from "react";
import useLaptopWS from "../hooks/useLaptopWS";
import CpuChart from "./charts/CpuChart";
import RamChart from "./charts/RamChart";
import StorageChart from "./charts/StorageChart";
import NetworkChart from "./charts/NetworkChart";
import BatteryChart from "./charts/BatteryChart";

const MAX_POINTS = 60;

export default function LaptopMonitor() {
  const [data, setData] = useState([]);

  useLaptopWS((msg) => {
    const time = new Date().toLocaleTimeString();

    const point = {
      time,
      cpu: msg.cpu,
      ram: msg.ram.usage,
      storage: msg.storages.usage,
      rx: msg.network.rx,
      tx: msg.network.tx,
      battery: msg.battery,
    };

    setData((prev) => [...prev, point].slice(-MAX_POINTS));
  });

  return (
    <div>
      <h2>Laptop Monitor</h2>

      <CpuChart data={data} />
      <RamChart data={data} />
      <StorageChart data={data} />
      <NetworkChart data={data} />
      <BatteryChart data={data} />
    </div>
  );
}