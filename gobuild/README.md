# gobuild 🧰

**Status:** Em desenvolvimento (WIP) 🔧

Microserviço de exemplo para gerenciamento de produtos e estoque usando SQLite + Echo (Go).

---

## Sumário

1. Visão geral
2. Pré-requisitos
3. Como rodar
4. Seed / Popular banco
5. Endpoints principais (exemplos curl)
6. Banco de dados
7. Desenvolvimento e testes
8. Roadmap rápido
9. Contato

---

## 1) Visão geral

O `gobuild` é um serviço simples para demonstrar conceitos de API REST em Go: CRUD de produtos, controle de estoque (movimentações) e persistência em SQLite.

---

## 2) Pré-requisitos

- Go 1.24
- Make (opcional)
- (Opcional) `sqlite3` CLI para inspeção do arquivo `data.db`

---

## 3) Como rodar

- Via Make (recomendado local):

  ```bash
  cd gobuild
  make run
  ```

  ou diretamente:

  ```bash
  cd gobuild
  go run ./cmd/api
  ```

- O servidor inicia em `http://localhost:8080` por padrão.

> Nota: se precisar rodar `ApiStudents` e `gobuild` simultaneamente, altere a porta em `cmd/api` ou execute um dos serviços em outra porta.

---

## 4) Seed / Popular banco

- Script (Go) para inserir produtos de exemplo:

  ```bash
  cd gobuild
  go run ./cmd/seed
  ```

  Exemplo de saída:

  ```text
  Inseridos 3 produtos de exemplo no data.db
  ```

- Observação: `data.db` será criado no diretório onde o comando for executado.

---

## 5) Endpoints principais (exemplos curl) 🔬

Base: `http://localhost:8080`

### Produtos

- Criar produto (POST /api/products)

  ```bash
  curl -X POST http://localhost:8080/api/products \
    -H 'Content-Type: application/json' \
    -d '{"name":"Cimento CP-II 50kg","preco":25.5,"estoque":100,"unidade":"saco","categoria":"Materiais"}'
  ```

- Listar produtos (GET /api/products)

  ```bash
  curl http://localhost:8080/api/products
  ```

- Obter produto por id (GET /api/products/:id)

  ```bash
  curl http://localhost:8080/api/products/1
  ```

- Atualizar produto (PUT /api/products/:id)

  ```bash
  curl -X PUT http://localhost:8080/api/products/1 \
    -H 'Content-Type: application/json' \
    -d '{"name":"Cimento CP-II 50kg","preco":26.0,"estoque":120,"unidade":"saco","categoria":"Materiais"}'
  ```

- Deletar produto (DELETE /api/products/:id)

  ```bash
  curl -X DELETE http://localhost:8080/api/products/1
  ```

- Consultar estoque (GET /api/products/:id/stock)

  ```bash
  curl http://localhost:8080/api/products/1/stock
  ```

### Movimentações de estoque

- Entrada (POST /api/stock/entrada)

  ```bash
  curl -X POST http://localhost:8080/api/stock/entrada \
    -H 'Content-Type: application/json' \
    -d '{"product_id":1,"quantity":10}'
  ```

- Saída (POST /api/stock/saida)

  ```bash
  curl -X POST http://localhost:8080/api/stock/saida \
    -H 'Content-Type: application/json' \
    -d '{"product_id":1,"quantity":5}'
  ```

- Ajuste (POST /api/stock/ajuste)

  ```bash
  curl -X POST http://localhost:8080/api/stock/ajuste \
    -H 'Content-Type: application/json' \
    -d '{"product_id":1,"quantity":200}'
  ```

- Histórico (GET /api/stock/historico/:product_id)

  ```bash
  curl http://localhost:8080/api/stock/historico/1
  ```

---

## 6) Banco de dados

- SQLite com arquivo `data.db`.
- Tabelas criadas automaticamente na primeira execução:
  - `products` (id, name, price, stock, unit, category, created_at)
  - `stock_movements` (id, product_id, tipo, quantidade, created_at)

---

## 7) Desenvolvimento e testes

- Formatar: `go fmt ./...`
- Checar vet: `go vet ./...`
- Testes: `go test ./...`
- Sugestão: adicione testes unitários para `product` e `stock` e um teste de integração que rode contra um DB temporário.

---

## 8) Roadmap rápido

- Cobertura de testes (unit & integration)
- CI (GitHub Actions) para checks e testes
- Dockerfile / Devcontainer
- Melhorar logs / métricas

---

## 9) Contato

- Repo: `github.com/EtraudBits/golangProject`
- Issues / PRs bem-vindos — veja o README principal para guidelines de contribuição.

---

🔧 Se quiser, posso também adicionar um `Makefile` com meta `seed` e `test` para facilitar o fluxo. Quer que eu faça isso? (sim/não)
