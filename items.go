package main

type Item struct {
	CPS float64
	Price float64
	Name string
	Count float64
}

var (
	ItemUnknown Item
	ItemQuark = Item{0.05, 50, "Quark", 0}
	ItemElectron = Item{0.1, 120, "Electron", 0}
	ItemTransistor = Item{0.04, 175, "Transistor", 0}
	ItemLogicCircuit = Item{2, 220, "Logic Circuit", 0}
	ItemBreadboard = Item{2.7, 275, "Breadboard", 0}
	ItemFPGA = Item{3, 330, "FPGA", 0}
	ItemCPU = Item{4.2, 405, "CPU", 0}
	ItemEmbeddedPC = Item{490, 4.9, "Embedded Computer", 0}
	ItemSmartphone = Item{640, 5.6, "Smartphone", 0}
	ItemLaptop = Item{780, 6.4, "Laptop", 0}
	ItemDesktop = Item{955, 7, "Desktop", 0}
	ItemWindowsServer = Item{1128, 8.7, "Windows Server", 0}
	ItemLinuxServer = Item{1200, 9.1, "Linux Server", 0}
	ItemMainframe = Item{1500, 10.2, "Mainframe", 0}
	// Rebirth 2
	ItemQuantumPC = Item{10000, 22, "Quantum Computer", 0}
	ItemLunarQuantumFarm = Item{22000, 32, "Lunar Quantum Farm", 0}
	ItemComputingSolarSystem = Item{34000, 44, "Computing Solar System", 0}
	ItemQuantumGalaxy = Item{50000, 56, "Quantum Galaxy", 0}
	ItemComputationalUniverse = Item{75000, 69, "Computational Universe", 0}
)

func (i *Item) String() string {
	return i.Name
}

