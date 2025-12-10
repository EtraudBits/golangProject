package stock

// Model representa uma movimentação de estoque
type Movement struct {
	ID int `json:"id"` // ID da movimentação
	ProductID int `json:"product_id"` // ID do produto relacionado
	Type string `json:"type"` // Tipo de movimentação: "Entrada", "Saida", "Ajuste"
	Quantity float64 `json:"quantity"` // Quantidade movimentada
	CreatedAt string `json:"created_at"` // Timestamp da movimentação pelo SQLite
}