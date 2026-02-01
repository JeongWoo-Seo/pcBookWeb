import { useEffect, useState } from "react";
import { getLaptopList } from "../api/laptop";

export default function LaptopList({ selected, onSelect }) {
  const [laptops, setLaptops] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const load = async () => {
      try {
        const list = await getLaptopList();
        setLaptops(list || []);
      } catch (err) {
        console.error(err);
        setError("Failed to load laptop list");
      } finally {
        setLoading(false);
      }
    };

    load();
  }, []);

  if (loading) return <div>Loading laptops...</div>;
  if (error) return <div>{error}</div>;

  return (
    <div style={{ minWidth: 200 }}>
      <h3>Active Laptops</h3>
      <ul>
        {laptops.map((id) => (
          <li
            key={id}
            onClick={() => onSelect(id)}
            style={{
              cursor: "pointer",
              fontWeight: selected === id ? "bold" : "normal",
              color: selected === id ? "blue" : "black",
            }}
          >
            {id}
          </li>
        ))}
      </ul>
    </div>
  );
}