import { LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer } from "recharts";

export default function CpuChart({ data }) {
  return (
    <>
      <h4>CPU Usage (%)</h4>
      <ResponsiveContainer width="100%" height={200}>
        <LineChart data={data}>
          <XAxis dataKey="time" />
          <YAxis domain={[0, 100]} />
          <Tooltip />
          <Line dataKey="cpu" dot={false} />
        </LineChart>
      </ResponsiveContainer>
    </>
  );
}
