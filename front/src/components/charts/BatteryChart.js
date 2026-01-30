import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

export default function BatteryChart({ data }) {
  return (
    <>
      <h4>Battery (%)</h4>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={data}>
          <XAxis dataKey="time" />
          <YAxis domain={[0, 100]} />
          <Tooltip />
          <Line dataKey="battery" dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </>
  );
}