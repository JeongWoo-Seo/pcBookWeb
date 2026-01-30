import {
  ResponsiveContainer,
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
} from "recharts";

export default function NetworkChart({ data }) {
  return (
    <>
      <h4>Network (RX / TX)</h4>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={data}>
          <XAxis dataKey="time" />
          <YAxis />
          <Tooltip />
          <Line dataKey="rx" dot={false} />
          <Line dataKey="tx" dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </>
  );
}
