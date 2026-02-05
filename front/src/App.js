import { useState } from "react";
import LaptopList from "./components/LaptopList";
import LaptopMonitor from "./components/LaptopMonitor";

export default function App() {
  const [selectedLaptop, setSelectedLaptop] = useState(null);

  return (
    <div style={{ padding: 20 }}>
      <h1>PC Monitor Dashboard</h1>

      <LaptopList
        selected={selectedLaptop}
        onSelect={setSelectedLaptop}
      />

      <hr />

      <LaptopMonitor laptopID={selectedLaptop} />
    </div>
  );
}