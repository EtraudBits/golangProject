API# README — Student Management API (Golang) / API de Gerenciamento de Alunos (Golang)

---

> Português (PT-BR) — em seguida, English (EN).

---

# 🇧🇷 Versão em Português

## Visão geral

API simples em **Golang** para gerenciar alunos de um curso. Permite criar, listar, consultar, atualizar e deletar registros de estudantes. Ideal para demonstrar boas práticas (estrutura de projeto, validações, testes, docs), integrações com banco de dados e deploy em container.

### Objetivo

Ter um projeto de portfólio que:

- Mostre arquitetura e organização de código em Go;
- Tenha endpoints RESTful claros e testáveis;
- Conte com validações básicas (ex.: CPF, email);
- Esteja pronto para deploy (Docker) e com instruções para rodar localmente.

---

## Funcionalidades (routes / rotas)

- `GET /students` — Lista todos os estudantes.
- `POST /students` — Cadastrar (criar) um estudante.
- `GET /students/:id` — Obter informações específicas de um estudante.
- `PUT /students/:id` — Atualizar estudante.
- `DELETE /students/:id` — Deletar (apagar) um estudante.

---

## Estrutura do recurso `Student` (estudante)

```json
{
  "id": "uuid",
  "name": "string",
  "cpf": "string", // CPF sem formatação ou formatado (decida padrão)
  "email": "string",
  "age": 0,
  "active": true
}
```

Campos:

- `name` — nome completo.
- `cpf` — CPF do aluno (validar formato e dígito).
- `email` — validar formato de e-mail.
- `age` — número inteiro.
- `active` — booleano indicando se o aluno está ativo no curso.

---

## Requisitos / Tecnologias sugeridas

- Linguagem: **Go** (versão ≥ 1.20 recomendada)
- Router: `chi` ou `gorilla/mux` ou `gin` (sugestão: `chi` por ser leve e idiomático)
- Banco de dados: PostgreSQL (ou SQLite para versão local/simples)
- ORM/DB driver: `sqlx` ou `gorm` (sugestão: `sqlx` para mais controle SQL)
- Migrations: `golang-migrate/migrate`
- Validação: `go-playground/validator.v10` + validação específica de CPF
- Logger: `zap` ou `logrus`
- Testes: `testing` nativo + `httptest`
- Container: `Dockerfile` + `docker-compose` (opcional)
- Documentação: README + exemplos curl; opcionalmente Swagger / OpenAPI

---

## Estrutura de pastas sugerida

```
/student-api
├─ cmd/
│  └─ server/                # main da aplicação
├─ internal/
│  ├─ students/
│  │  ├─ handler.go          # handlers HTTP
│  │  ├─ service.go          # lógica de negócio
│  │  ├─ repository.go       # acesso ao DB
│  │  └─ model.go            # structs e validações
│  └─ pkg/                   # pacotes utilitários (logger, config)
├─ migrations/
├─ test/
├─ Dockerfile
├─ docker-compose.yml
└─ README.md
```

---

## Plano passo a passo (como começar / milestones)

1. **Inicialização**

   - Criar repositório Git, estrutura de pastas básica.
   - Definir `go.mod`.

2. **Modelo e DB**

   - Implementar `model.Student`.
   - Criar migrations iniciais (tabela students).
   - Configurar conexão com PostgreSQL (ou SQLite para dev).

3. **Repository**

   - Implementar `Create`, `FindAll`, `FindByID`, `Update`, `Delete` no repositório.

4. **Service / Business**

   - Adicionar validações (CPF, email, idade mínima se desejar).
   - Lógica de negócio (p.ex. impedir duplicidade de CPF).

5. **Handlers / Router**

   - Implementar endpoints REST conforme rotas.
   - Mapear códigos HTTP apropriados (200, 201, 400, 404, 422, 500).

6. **Testes**

   - Unit tests para service e repositório (usar DB em memória ou dockerized test DB).
   - Integration tests para endpoints com `httptest`.

7. **Documentação & exemplos**

   - Exemplos `curl` no README, e talvez OpenAPI/Swagger.

8. **Docker & CI**

   - Criar `Dockerfile` e `docker-compose` para desenvolvimento.
   - Configurar CI (GitHub Actions): build, run tests, lint.

9. **Melhorias**

   - Autenticação (JWT), paginação, filtros, ordenação, rate-limit.

---

## Exemplos de uso (cURL)

### Criar estudante

```bash
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "cpf": "12345678909",
    "email": "joao@example.com",
    "age": 25,
    "active": true
  }'
```

Resposta esperada (201):

```json
{
  "id": "uuid-gerado",
  "name": "João Silva",
  "cpf": "12345678909",
  "email": "joao@example.com",
  "age": 25,
  "active": true
}
```

### Listar estudantes

```bash
curl http://localhost:8080/students
```

### Obter por id

```bash
curl http://localhost:8080/students/<id>
```

### Atualizar

```bash
curl -X PUT http://localhost:8080/students/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"João Atualizado","active":false}'
```

### Deletar

```bash
curl -X DELETE http://localhost:8080/students/<id>
```

---

## Contratos e códigos HTTP (resumo)

- `GET /students` — `200 OK` com lista (pode paginar).
- `POST /students` — `201 Created` + resource; `400/422` para validação.
- `GET /students/:id` — `200 OK` ou `404 Not Found`.
- `PUT /students/:id` — `200 OK` com resource atualizado ou `404`.
- `DELETE /students/:id` — `204 No Content` ou `404`.

---

## Validações importantes

- CPF: validar formato e dígitos verificadores (não aceitar CPFs inválidos).
- Email: válido conforme regex aceitável.
- Duplicidade: não permitir dois alunos com o mesmo CPF.
- Age: número inteiro não-negativo (pode definir mínimo).
- Campos obrigatórios: name, cpf, email.

---

## Observações e boas práticas

- Use `context.Context` nas chamadas de DB para timeout e trace.
- Retorne erros padronizados (JSON) com `message` e `code`.
- Separação clara entre handler (HTTP) e service (lógica) para testabilidade.
- Logging estruturado (zap) e configuração via variáveis de ambiente.
- Expor métricas (Prometheus) se quiser demonstrar observability.

---

## Deploy (rápido)

- Dockerfile básico + docker-compose com DB.
- Em produção: build multi-stage para reduzir imagem final, use variáveis de ambiente seguras e orquestração (Heroku, AWS ECS, GCP Cloud Run, DigitalOcean App Platform ou Kubernetes).

---

## Checklist mínimo para o portfólio

- [ ] Código bem estruturado e comentado.
- [ ] README claro (este documento).
- [ ] Dockerfile & docker-compose.
- [ ] Migrations.
- [ ] Testes unitários e de integração.
- [ ] Um exemplo de chamada (curl) e resposta.
- [ ] (Opcional) GitHub Actions para CI.

---

## Possíveis extensões / features "pioneiras"

- Paginação + filtros (por nome, ativo).
- Importação CSV para criar vários alunos.
- Endpoint para buscar por CPF.
- Integração com serviço externo para validação de CPF.
- Painel admin simples (frontend) usando React + Tailwind (apenas para portfólio).
- CQRS simples: separar leitura e escrita para demonstrar padrões arquiteturais.

---

## License / Créditos

Sugestão: escolher uma licença permissiva (MIT) para facilitar exibição no portfólio.

---

# 🇬🇧 English Version — Student Management API (Golang)

## Overview

Simple **Golang** API to manage course students. Create, list, retrieve, update and delete student records. Great for a portfolio to demonstrate project structure, validations, testing, and deployment.

### Goal

Provide a portfolio-ready project that:

- Shows Go project architecture and structure;
- Exposes RESTful endpoints;
- Implements basic validations (e.g. CPF, email);
- Is ready to run in Docker and has clear run instructions.

---

## Routes

- `GET /students` — List all students.
- `POST /students` — Create a student.
- `GET /students/:id` — Get specific student info.
- `PUT /students/:id` — Update student.
- `DELETE /students/:id` — Delete a student.

---

## Student structure

```json
{
  "id": "uuid",
  "name": "string",
  "cpf": "string",
  "email": "string",
  "age": 0,
  "active": true
}
```

Fields: `name`, `cpf` (Brazilian ID), `email`, `age`, `active`.

---

## Suggested tech stack

- Go ≥ 1.20
- Router: `chi` (suggested) / `gin` / `gorilla/mux`
- DB: PostgreSQL (or SQLite for local)
- DB layer: `sqlx` or `gorm` (`sqlx` suggested)
- Migrations: `migrate`
- Validation: `go-playground/validator.v10`
- Logger: `zap` or `logrus`
- Tests: built-in `testing` + `httptest`
- Docker & docker-compose
- Optional: OpenAPI / Swagger

---

## Project layout (suggested)

```
/student-api
├─ cmd/
├─ internal/
│  ├─ students/
│  └─ pkg/
├─ migrations/
├─ Dockerfile
├─ docker-compose.yml
└─ README.md
```

---

## Step-by-step roadmap

1. Initialize repo and `go.mod`.
2. Model & DB migrations.
3. Implement repository CRUD.
4. Implement service with validations (CPF/email).
5. Handlers + router.
6. Write unit and integration tests.
7. Add Docker + CI (GitHub Actions).
8. Add docs (README + curl examples or Swagger).

---

## cURL examples

### Create student

```bash
curl -X POST http://localhost:8080/students \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "cpf": "12345678909",
    "email": "jane@example.com",
    "age": 30,
    "active": true
  }'
```

### List students

```bash
curl http://localhost:8080/students
```

### Get by id

```bash
curl http://localhost:8080/students/<id>
```

### Update

```bash
curl -X PUT http://localhost:8080/students/<id> \
  -H "Content-Type: application/json" \
  -d '{"name":"Jane Updated","active":false}'
```

### Delete

```bash
curl -X DELETE http://localhost:8080/students/<id>
```

---

## HTTP codes summary

- `GET /students` — `200 OK`
- `POST /students` — `201 Created` (validation errors `400/422`)
- `GET /students/:id` — `200 OK` or `404 Not Found`
- `PUT /students/:id` — `200 OK` or `404`
- `DELETE /students/:id` — `204 No Content` or `404`

---

## Validations to implement

- CPF validation (check digits)
- Email format validation
- Prevent duplicate CPF
- Age must be non-negative integer
- Required fields: `name`, `cpf`, `email`

---

## Best practices

- Use `context.Context` for DB calls.
- Keep handler thin; service contains business rules.
- Standardize error responses (JSON).
- Use structured logging.
- Make config via env variables.

---

## Deployment hints

- Multi-stage Docker build for small images.
- Use `docker-compose` for dev (app + db).
- Use GitHub Actions to run tests on PR.

---

## Portfolio checklist

- [ ] Clean, commented code
- [ ] README (this file)
- [ ] Docker + docker-compose
- [ ] Migrations
- [ ] Tests
- [ ] Examples (curl)
- [ ] CI pipeline

---

Se quiser, eu já **posso** gerar a estrutura inicial do projeto (main.go, handlers, model, repository, Dockerfile, docker-compose e exemplos de testes) com comentários linha a linha — me diga se prefere `chi` ou `gin` como router e se quer PostgreSQL ou SQLite para começo. Vou preparar tudo pronto para você colocar no GitHub. 🚀

(Escolhi não perguntar mais nada para não atrasar seu fluxo — se preferir, responda qual router/DB prefere e eu gero os arquivos já comentados.)
