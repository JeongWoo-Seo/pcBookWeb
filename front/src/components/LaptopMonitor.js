import { useEffect, useState } from "react";
import useLaptopWS from "../hooks/useLaptopWS";
import CpuChart from "./charts/CpuChart";
import RamChart from "./charts/RamChart";
import StorageChart from "./charts/StorageChart";
import NetworkChart from "./charts/NetworkChart";
import BatteryChart from "./charts/BatteryChart";

export default function LaptopMonitor({ laptopID }) {
  const [history, setHistory] = useState([]);

  useLaptopWS({
    mode: laptopID ? "single" : "broadcast",
    laptopID,
    onMessage: (msg) => {
      setHistory((prev) =>
        [...prev, {
          time: new Date().toLocaleTimeString(),
          cpu: msg.cpu,
          ram: msg.ram.usage,
          storage: msg.storages.usage,
          rx: msg.network.rx,
          tx: msg.network.tx,
          battery: msg.battery,
        }].slice(-30)
      );
    },
  });

  if (!laptopID) {
    return <div>노트북을 선택하세요</div>;
  }

  return (
    <>
      <h3>Laptop ID: {laptopID}</h3>
      <CpuChart data={history} />
      <RamChart data={history} />
      <StorageChart data={history} />
      <NetworkChart data={history} />
      <BatteryChart data={history} />
    </>
  );
}
