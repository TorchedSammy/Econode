package econetwork

type Item struct {
	CPS float64
	Price float64
	Name string
	Count int
}

var (
	ItemElectron = Item{0.01, 20, "electron", 0}
)

