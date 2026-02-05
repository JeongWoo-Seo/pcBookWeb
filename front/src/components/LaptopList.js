import { useEffect, useState } from "react";
import { getLaptopList } from "../api/laptop";

export default function LaptopList({ selected, onSelect }) {
  const [laptops, setLaptops] = useState([]);

  useEffect(() => {
    getLaptopList().then(setLaptops);
  }, []);

  return (
    <div>
      <h3>Active Laptops</h3>
      <ul>
        {laptops.map((id) => (
          <li
            key={id}
            style={{
              cursor: "pointer",
              fontWeight: selected === id ? "bold" : "normal",
            }}
            onClick={() => onSelect(id)}
          >
            {id}
          </li>
        ))}
      </ul>
    </div>
  );
}