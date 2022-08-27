package econode

type Auction struct {
	Type AuctionType
	BasePrice float64
	SellingPrice float64
}

type AuctionType int
const (
	AuctionGem AuctionType = 1
)

