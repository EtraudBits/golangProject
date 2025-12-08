package product

// Product representa um produto no sistema
// cada campo tem tags `json` para mapear automaticamente entre JSON e struct
type Produto struct {
	ID int `json:"id"` // id auto-incremental (PK)
	Name string `json:"name"` // nome do produto (ex.: "cimento cp-II 50kg")
	Preco float64 `json:"preco"` // preço do produto (ex.: 25.50)
	Estoque float64 `json:"estoque"` // quantidade em estoque (ex.: 100.0)
	Unidade string `json:"unidade"` // unidade de medida (ex.: "kg", "m2", "un")
	Categoria string `json:"categoria"` // categoria do produto (ex.: "materiais de construção")
	DataCriacao string `json:"data_criacao"` // timestamp de criação do registro (ex.: "2024-06-01 12:00:00")
}