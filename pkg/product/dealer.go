package product

type Dealer struct{}

// Provide provides store with necessary items
func (d Dealer) Provide(productName string, quantity, price int) (Product) {
	originalPrice := price / 2
	
	return Product{
		Name:          productName,
		Quantity:      quantity,
		Price:         price,
		OriginalPrice: originalPrice,
	}
}