import { useState } from "react";
import useLaptopWS from "../hooks/useLaptopWS";
import CpuChart from "./charts/CpuChart";
import RamChart from "./charts/RamChart";
import StorageChart from "./charts/StorageChart";
import NetworkChart from "./charts/NetworkChart";
import BatteryChart from "./charts/BatteryChart";

const MAX_POINTS = 60;

export default function LaptopMonitor({ laptopId }) {
  if (!laptopId) {
    return <div>Select a laptop to monitor</div>;
  }


  return (
    <div style={{ flex: 1 }}>
      <h3>Monitoring: {laptopId}</h3>

      {/* WebSocket / Chart 컴포넌트 */}
      {/* CPUChart, RAMChart, BatteryChart ... */}
    </div>
  );
}