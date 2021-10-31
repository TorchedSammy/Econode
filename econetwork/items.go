package econetwork

type Item struct {
	CPS float64
	Price float64
	Name string
	Count float64
}

var (
	ItemElectron = Item{0.01, 20, "electron", 0}
)

var itemMap = map[string]*Item{
	"electron": &ItemElectron,
}

func (i *Item) String() string {
	return i.Name
}

