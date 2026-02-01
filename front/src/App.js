import { useState } from "react";
import LaptopList from "./components/LaptopList";
import LaptopMonitor from "./components/LaptopMonitor";

function App() {
  const [selectedLaptop, setSelectedLaptop] = useState(null);

  return (
    <div style={{ padding: 20 }}>
      <h1>PC Monitor Dashboard</h1>

      <div style={{ display: "flex", gap: 20 }}>
        <LaptopList
          selected={selectedLaptop}
          onSelect={setSelectedLaptop}
        />

        <LaptopMonitor laptopId={selectedLaptop} />
      </div>
    </div>
  );
}

export default App;