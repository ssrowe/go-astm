package astm1384

type ASTMMessage struct {
	Header       *Header
	Manufacturer *Manufacturer
	Records      []*Record
}

type Record struct {
	Patient          *Patient
	OrdersAndResults []*OrderResults
	Comments         []*Comment
}

type OrderResults struct {
	Order    *Order
	Results  []*CommentedResult
	Comments []*Comment
}
