# 📘 **API Students – Sistema de Cadastro de Estudantes em Go**

---

<div style="display: inline_block"><br>
  <img alt="Golang" width="48" src="https://cdn.jsdelivr.net/gh/devicons/devicon/icons/go/go-original.svg" />
</div>

## 📌 **Descrição do Projeto**

A **API Students** é uma API RESTful desenvolvida em **Golang**, utilizando o framework **Echo** e o ORM **GORM**, com banco de dados **SQLite**.
O objetivo é permitir o **cadastro, consulta, atualização e exclusão de estudantes** (CRUD completo).

Este projeto foi construído seguindo boas práticas de arquitetura, organização de pacotes, tratamento de erros e documentação com **Swagger**.

✔️ Ideal para estudos, portfólio e como base para aplicações maiores.
✔️ Suporte completo a validações, filtros, logs estruturados e responses padronizadas.

---

## 🗂️ **Estrutura do Projeto**

```
ApiStudents/
 ├── api/
 │   ├── api.go         # Inicialização do servidor e rotas
 │   ├── handler.go     # Funções que tratam as rotas (controllers)
 │   ├── request.go     # Validações de entrada
 ├── db/
 │   └── db.go          # Conexão e operações com o banco SQLite usando GORM
 ├── docs/              # Documentação Swagger gerada automaticamente
 ├── schemas/
 │   └── schemas.go     # Modelos e estruturas de resposta
 ├── main.go            # Início da aplicação
 ├── go.mod / go.sum    # Dependências
 └── Makefile           # Comandos automatizados
```

### 🧩 **Arquitetura**

A API segue uma arquitetura clara e organizada:

```
Request → Echo Router → Handler → Database Layer (GORM + SQLite)
                                → Responses (schemas)
```

---

## 🚀 **Tecnologias Utilizadas**

| Tecnologia           | Descrição                             |
| -------------------- | ------------------------------------- |
| **Go (Golang)**      | Linguagem principal do projeto        |
| **Echo Framework**   | Framework web rápido e minimalista    |
| **GORM**             | ORM para manipulação do banco         |
| **SQLite**           | Banco de dados leve e local           |
| **Swagger (swaggo)** | Documentação interativa da API        |
| **Zerolog**          | Logs estruturados de alta performance |
| **Postman**          | Utilizado para testar a API           |

---

# 🔧 **Como Rodar o Projeto**

### 📥 1. Clonar o repositório

```bash
git clone https://github.com/seu-usuario/ApiStudents.git
cd ApiStudents
```

### ▶️ 2. Rodar o projeto

Você pode iniciar manualmente:

```bash
go run main.go
```

Ou usando o **Makefile**:

```bash
make run
```

### 🗄️ 3. Banco de dados

O SQLite cria automaticamente o arquivo:

```
student.db
```

---

# 📚 **Documentação Swagger**

Após iniciar o servidor, acesse:

👉 **[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)**

![Swagger UI](https://i.imgur.com/5fx6n8W.png)

---

# 🔗 **Endpoints Disponíveis**

## 📌 GET /students

Retorna todos os estudantes.

### Filtro opcional:

```
/students?active=true
```

### Exemplo de Resposta:

```json
{
  "students:": [
    {
      "id": 1,
      "name": "Maria Souza",
      "email": "maria@email.com",
      "cpf": 123456789,
      "age": 22,
      "active": true
    }
  ]
}
```

---

## 📌 POST /students

Cria um novo estudante.

### Exemplo de Request:

```json
{
  "Name": "João Silva",
  "CPF": 11122233344,
  "Email": "joao@mail.com",
  "Age": 19,
  "Active": true
}
```

### Validações (request.go):

✔️ Campos obrigatórios
✔️ `Active` é ponteiro (`*bool`) → força preenchimento
✔️ Verificação de valores vazios

---

## 📌 GET /students/{id}

Retorna um único estudante.
Possíveis status HTTP:

| Status                        | Significado          |
| ----------------------------- | -------------------- |
| **200 OK**                    | Estudante encontrado |
| **404 Not Found**             | ID inexistente       |
| **500 Internal Server Error** | Erro interno         |

---

## 📌 PUT /students/{id}

Atualiza parcialmente ou totalmente os dados de um estudante.

A função `upDateStudentInfo()` trata de apenas atualizar campos enviados no body.

---

## 📌 DELETE /students/{id}

Remove o estudante do banco.

---

# 💾 **Camada de Banco de Dados**

A camada **db/** utiliza GORM para:

✔️ Criar banco automático com `AutoMigrate`
✔️ CRUD completo
✔️ Filtragem com `Where`
✔️ Tratamento de erros estruturados

Exemplo real usado no projeto:

```go
db, err := gorm.Open(sqlite.Open("./student.db"), &gorm.Config{})
db.AutoMigrate(&schemas.Student{})
```

---

# 🧱 **Schemas (Modelos)**

### Estrutura do Estudante (GORM)

```go
type Student struct {
    gorm.Model
    Name   string
    CPF    int
    Email  string
    Age    int
    Active bool
}
```

### Estrutura enviada ao usuário

```go
type StudentResponse struct {
    ID        int
    CreatedAt time.Time
    UpdatedAt time.Time
    Name      string
    Email     string
    CPF       int
    Age       int
    Active    bool
}
```

A função:

```go
NewResponse(students)
```

transforma os dados brutos em respostas limpas para a API.

---

# 🛠️ **Principais Recursos Implementados**

✔️ **CRUD completo**
✔️ **Tratamento de erros HTTP adequado**
✔️ **Validação robusta de entrada**
✔️ **Logs estruturados com zerolog**
✔️ **Documentação Swagger integrada**
✔️ **Filtro por campo (active)**
✔️ **Arquitetura clara e desacoplada**
✔️ **Uso correto de ponteiros, structs e contextos do Echo**

---

# 🎯 **Como Essa API Pode Ser Usada no Dia a Dia**

Esta API é ideal para:

- Sistemas internos escolares
- Gestão simples de alunos
- Painéis administrativos
- Treinamento para API REST
- Base para microsserviços em Go

Aplicações reais:

📌 Cadastrar alunos em cursos
📌 Consultar estudantes ativos
📌 Atualizar informações (email, idade, cpf)
📌 Excluir usuários inativos
📌 Integrar com front-ends ou aplicações mobile

---

# 🖼️ **Fluxo Completo de Funcionamento**

```
(Client) → Echo Router → Handler → Validação → DB Layer
                                       ↓
                                 Swagger Docs
```

---

# 📃 **Conclusão**

Este projeto demonstra domínio em:

- Programação Go
- Construção de APIs REST
- Arquitetura limpa
- Integração com banco usando GORM
- Validação, logs e documentação Swagger

[def]: https://i.imgur.com/UkB1cQX.png
