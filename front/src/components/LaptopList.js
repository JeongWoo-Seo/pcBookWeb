import { useEffect, useState } from "react";
import { getLaptopList } from "../api/laptop";

export default function LaptopList({ selected, onSelect }) {
  const [laptops, setLaptops] = useState([]);

  useEffect(() => {
    getLaptopList().then(setLaptops);
  }, []);

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