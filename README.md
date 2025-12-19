# golangProject ✅

[![Go 1.24](https://img.shields.io/badge/go-1.24-blue)](https://golang.org) [![CI](https://img.shields.io/github/actions/workflow/status/EtraudBits/golangProject/ci.yml?branch=main)](https://github.com/EtraudBits/golangProject/actions) [![License: MIT](https://img.shields.io/badge/license-MIT-green)](LICENSE)

**Status:** Em desenvolvimento (WIP) 🔧

Pequena coleção de projetos e experimentos em Go, com foco em aprendizado prático, organização e evolução contínua. Este README será atualizado conforme novas features e módulos forem adicionados.

---

## Sumário

1. Visão geral
2. Estrutura do repositório
3. Pré-requisitos
4. Como executar (por módulo)
5. Banco de dados
6. Documentação da API
7. Desenvolvimento e contribuições
8. Roadmap & tarefas futuras
9. Como manter o README atualizado

---

## 1) Visão geral

Este repositório contém projetos/tópicos:

- `ApiStudents` — API REST para gerenciamento de estudantes (com Swagger).
- `gobuild` — microserviço/experimentação com produto/estoque usando SQLite.
- Outros diretórios (ex.: `projetocimento`, `Comp_Tributaria`) para estudos/experimentos.

---

## 2) Estrutura do repositório (resumo)

- `ApiStudents/` — código fonte, docs (Swagger), `main.go`, `Makefile`.
  - Porta esperada: **8080**
  - Documentação: `ApiStudents/docs/swagger.yaml` e `swagger.json`
- `gobuild/` — `cmd/`, `internal/` (server, database, handlers), `Makefile`.
  - Executável: `go run ./cmd/api`
  - Banco SQLite: arquivo `data.db` (gerado automaticamente ao iniciar)
- `README.md` — este arquivo (iterativo)

---

## 3) Pré-requisitos

- Go 1.24
- Make (opcional, facilita executar `make run`)
- (Opcional) `swag` / `swaggo` se quiser regenerar docs: `go install github.com/swaggo/swag/cmd/swag@latest`

---

## 4) Como executar

### ApiStudents

- Rodar:
  - via Make:
    - `cd ApiStudents && make run`
  - ou:
    - `cd ApiStudents && go run main.go`
- Endpoints principais:
  - `GET /students`
  - `POST /students`
  - `GET /students/:id`
  - `PUT /students/:id`
  - `DELETE /students/:id`
  - Swagger: `/swagger/*` (ex.: `http://localhost:8080/swagger/index.html`)

### gobuild

- Rodar:
  - `cd gobuild && make run`
  - ou:
  - `cd gobuild && go run ./cmd/api`
- Observações:

  - O serviço cria/abre o arquivo `data.db` no diretório onde for executado.
  - Use `make fmt` / `make tidy` quando disponível para manter o código limpo.
  - Para popular o banco com dados de exemplo execute:

    - `cd gobuild && go run ./cmd/seed` (inserirá alguns produtos de exemplo em `data.db`). Exemplo de saída:

      `Inseridos 3 produtos de exemplo no data.db`

    - Nota: execute esse comando após iniciar o serviço se quiser que o servidor veja os dados ao iniciar, ou apenas use a cópia do DB (data.db) populada antes de iniciar.

- Nota sobre portas:
  - `ApiStudents` e `gobuild` usam, por padrão, a porta **8080**. Execute apenas um por vez ou altere a porta no código/variáveis de ambiente para executar simultaneamente.

---

### 4.1) Exemplos de requisições e scripts de seed 🔬

> Observação: os exemplos abaixo assumem que a API está rodando em `http://localhost:8080`.

#### ApiStudents (endpoints `/students`)

- Criar um estudante (POST):

  curl -X POST http://localhost:8080/students \
   -H 'Content-Type: application/json' \
   -d '{"Name":"João Silva","CPF":123456789,"Email":"joao@example.com","Age":21,"Active":true}'

- Listar estudantes (GET):

  curl http://localhost:8080/students

- Obter estudante por id (GET):

  curl http://localhost:8080/students/1

- Atualizar estudante (PUT):

  curl -X PUT http://localhost:8080/students/1 \
   -H 'Content-Type: application/json' \
   -d '{"Name":"João Silva","CPF":123456789,"Email":"joao.novo@example.com","Age":22,"Active":true}'

- Deletar estudante (DELETE):

  curl -X DELETE http://localhost:8080/students/1

- Script de exemplo: `ApiStudents/scripts/seed_students.sh` (executa 3 POSTs para criar estudantes). Para usar: inicie a API (`cd ApiStudents && make run`) e, em outro terminal, execute:

  `bash ApiStudents/scripts/seed_students.sh`

  Exemplo de saída:

  `Estudantes criados. Verifique com: curl http://localhost:8080/students`

#### gobuild (produtos e estoque)

- Criar produto (POST /api/products):

  curl -X POST http://localhost:8080/api/products \
   -H 'Content-Type: application/json' \
   -d '{"name":"Cimento CP-II 50kg","preco":25.5,"estoque":100,"unidade":"saco","categoria":"Materiais"}'

- Listar produtos (GET /api/products):

  curl http://localhost:8080/api/products

- Obter produto por id (GET /api/products/:id):

  curl http://localhost:8080/api/products/1

- Atualizar produto (PUT /api/products/:id):

  curl -X PUT http://localhost:8080/api/products/1 \
   -H 'Content-Type: application/json' \
   -d '{"name":"Cimento CP-II 50kg","preco":26.0,"estoque":120,"unidade":"saco","categoria":"Materiais"}'

- Deletar produto (DELETE /api/products/:id):

  curl -X DELETE http://localhost:8080/api/products/1

- Consultar estoque (GET /api/products/:id/stock):

  curl http://localhost:8080/api/products/1/stock

- Movimentações de estoque (exemplos):

  - Entrada: `POST /api/stock/entrada` → `{ "product_id": 1, "quantity": 10 }`
  - Saída: `POST /api/stock/saida` → `{ "product_id": 1, "quantity": 5 }`
  - Ajuste: `POST /api/stock/ajuste` → `{ "product_id": 1, "quantity": 200 }`

  Exemplo curl (entrada):

  curl -X POST http://localhost:8080/api/stock/entrada \
   -H 'Content-Type: application/json' \
   -d '{"product_id":1,"quantity":10}'

- Seed de produtos: `cd gobuild && go run ./cmd/seed` (inserirá alguns produtos de exemplo em `data.db`)

---

---

## 5) Banco de dados

- `gobuild` usa SQLite via driver `mattn/go-sqlite3`.
- Ao iniciar, o serviço cria as tabelas (se não existirem):
  - `products` (id, name, price, stock, unit, category, created_at)
  - `stock_movements` (id, product_id, tipo, quantidade, created_at)
- Arquivo padrão: `data.db` (criado no diretório onde o binário roda).

---

## 6) Documentação da API

- `ApiStudents` inclui Swagger (via `swaggo`). Arquivos em `ApiStudents/docs/`.
- Para atualizar/gerar docs: instale `swag` e execute na pasta do módulo (ex.: `swag init`).

---

## 7) Desenvolvimento e contribuição ✅

Agradeço contribuições! Para manter o projeto organizado e facilitar revisões, siga as guidelines abaixo antes de abrir PRs.

### Fluxo de trabalho (sugerido)

1. Crie uma branch a partir de `main` (ou `develop` se existir):
   - `git checkout -b feat/minha-feature`
2. Faça commits pequenos e com mensagem clara (veja padrão abaixo).
3. Garanta que o código está formatado e os testes passam localmente.
4. Abra um Pull Request com descrição do que foi feito e o motivo; vincule issues quando aplicável.

### Padrões de código & ferramentas 🔧

- Formate o código: `go fmt ./...`
- Organize/importe módulos: `go mod tidy`
- Execute checks básicos: `go vet ./...` e `go test ./...`
- (Opcional) Use `golangci-lint` para lint e verificação adicional: `golangci-lint run` (se optar por adicionar o config ao repositório, inclua `.golangci.yml`).
- Para regenerar docs do `ApiStudents`: instale `swag` e execute `swag init` dentro de `ApiStudents/`.

### Commit messages (padrão simples)

- prefira mensagens no formato: `tipo: descrição curta`
  - Ex.: `feat: adicionar endpoint de movimentação de estoque`
  - Tipos comuns: `feat`, `fix`, `chore`, `docs`, `test`, `refactor`

### Pull Request (PR)

- Título claro e objetivo
- Descreva o que foi feito, por quê e como testar localmente
- Relacione issues (`Fixes #123`) quando pertinente

#### Exemplo de template de PR

```
Title: feat: adicionar endpoint POST /api/stock/entrada

Descrição:
- Adiciona endpoint para criar movimentação de tipo "Entrada" e atualizar estoque

Como testar:
1. `cd gobuild && make run`
2. `curl -X POST http://localhost:8080/api/stock/entrada -H 'Content-Type: application/json' -d '{"product_id":1,"quantity":10}'`

Checklist:
- [ ] Testes adicionados
- [ ] Documentação atualizada (README/docs)
```

### Testes

- Execute `go test ./...` para rodar todos os testes.
- Para TDD/integração simples, crie testes na pasta do pacote e use subtests para cenários.

### Submissão de issues

- Abra um issue descrevendo o problema ou feature proposta, passo-a-passo para reproduzir (quando aplicável) e ambiente.

---

---

## 8) Roadmap & tarefas futuras 🚀

Abaixo está um backlog inicial com prioridades, estimativas de esforço e responsáveis sugeridos. Use-o como ponto de partida — atualize prioridades e atribuições conforme avançarmos.

| ID  | Tarefa                                                 | Prioridade | Esforço |    Assignee |   Status    | Notas                                                                                                  |
| --- | ------------------------------------------------------ | ---------: | ------: | ----------: | :---------: | ------------------------------------------------------------------------------------------------------ |
| R1  | Cobertura de testes (unit/integration)                 |       Alta |       M | @maintainer |    Todo     | Cobrir `product` e `stock` (unit) + um teste de integração básico que roda contra `data.db` temporário |
| R2  | Documentar endpoints e exemplos de requisição/resposta |       Alta |       S | @maintainer |    Todo     | Expandir README com exemplos e respostas reais; adicionar exemplos no Swagger se aplicável             |
| R3  | Adicionar CI (build, test, lint)                       |       Alta |       M | @maintainer |    Todo     | GitHub Actions para `go fmt`, `go vet`, `go test`, e linter (`golangci-lint`)                          |
| R4  | Melhorar logs e métricas                               |      Média |       M | @maintainer |    Todo     | Padronizar logs (zerolog) e expor métricas (Prometheus) em endpoints /metrics                          |
| R5  | Separar configuração (env / config file)               |      Média |       M | @maintainer |    Todo     | Usar env vars e possibilitar configuração da porta, DB path, etc.                                      |
| R6  | Dockerfile / Dev container                             |      Média |       M | @maintainer |    Todo     | Facilitar execução local e em CI com container dev e images minimalistas                               |
| R7  | Scripts de seed e exemplos automatizados               |      Baixa |       S | @maintainer |    Todo     | Transformar seeds em comandos de `make` e documentar saída esperada                                    |
| R8  | Testes de integração end-to-end                        |      Baixa |       L | @maintainer | Not started | Usar DB temporário ou docker-compose para testes E2E                                                   |

> Como usar esse backlog: escolha uma tarefa de alta prioridade (R1-R3), crie uma issue com referência ao ID (ex.: `R1`), atribua e siga o fluxo de contribuição descrito em **Desenvolvimento e contribuição**.

---
