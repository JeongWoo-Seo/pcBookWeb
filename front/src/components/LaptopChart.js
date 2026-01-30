import {
  LineChart,
  Line,
  XAxis,
  YAxis,
  Tooltip,
  CartesianGrid,
  ResponsiveContainer,
} from "recharts";

export default function LaptopChart({ data }) {
  return (
    <ResponsiveContainer width="100%" height={300}>
      <LineChart data={data}>
        <CartesianGrid strokeDasharray="3 3" />
        <XAxis dataKey="time" />
        <YAxis />
        <Tooltip />
        <Line type="monotone" dataKey="cpu" stroke="#ff7300" dot={false} />
        <Line type="monotone" dataKey="ram" stroke="#387908" dot={false} />
        <Line type="monotone" dataKey="storage" stroke="#8884d8" dot={false} />
      </LineChart>
    </ResponsiveContainer>
  );
}