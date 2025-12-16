package budget

// budget representa um orçamento (cabeçalho)
// Nota principal do orçamento
type Budget struct {
	ID int64 `json:"id"` //ID do Orçamento
	Customer string `json:"customer"` // Nome do cliente
	Total float64 `json:"total"` // Valor total do orçamento -> será calculado no service, não no handler
	CreatedAt string `json:"created_at"` // Data de criação
	Items []BudgetItem `json:"items"` // itens do orçamento
}

type BudgetItem struct {
	ID int64 `json:"id"` // ID do item
	BudgetID int64 `json:"budget_id"` // ID do orçamento (FK)
	ProductID int `json:"product_id"` // ID do produto
	Product string `json:"product"` // nome do produto (para exibição)
	Quantity float64 `json:"quantity"` //Quantidade
	UnitPrice float64 `json:"unit_price"` // Preço Unitário
	Subtotal float64 `json:"subtotal"` // Quantity * unitprice -> também será calculado no service
}