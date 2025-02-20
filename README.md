# Projeto GO - Clean Architecture

## Inicialização do Projeto

### Passos para rodar o projeto:
1. Subir os containers do Docker:
   ```sh
   docker compose up -d
   ```
2. Gerar o banco de dados:
   ```sh
   make migrate
   ```
3. Criar um arquivo `.env` na raiz do projeto se for rodar o comando abaixo na raiz:
   ```sh
   go run cmd/server/main.go cmd/server/wire_gen.go
   ```
   
   Se for rodar dentro da pasta do `main.go`, coloque o `.env` dentro da pasta `server` e execute:
   ```sh
   go run main.go wire_gen.go
   ```

## Configuração do Ambiente
Crie um arquivo `.env` e adicione as seguintes configurações:
```env
DB_DRIVER=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=root
DB_NAME=orders
WEB_SERVER_PORT=:8000
GRPC_SERVER_PORT=50051
GRAPHQL_SERVER_PORT=8080
```

## Endpoints e Portas

### REST API
- Porta **8000**
- **Criar Order:**
  - Método: `POST`
  - URL: `http://localhost:8000/order`
- **Listar Orders:**
  - Método: `GET`
  - URL: `http://localhost:8000/order`
  - Parâmetros opcionais: `page` (int), `limit` (int), `sort` (string) -> asc ou desc
  - Saída de dados:
  ```json
  {
      "id": string,
      "price": float,
      "tax": float,
      "finalPrice": float,
      "issueDate": date,
      "typeRequisition": string,
      "deleteAt": date
  }
  ```

### GraphQL
- Porta **8080**
- **Criação de Order:** `createOrder`
- **Listagem de Orders:** `listOrders`
  - Parâmetros opcionais: `page` (int), `limit` (int), `sort` (string) -> asc ou desc
  - Saída de dados:
  ```json
  {
      "id": string,
      "price": float,
      "tax": float,
      "finalPrice": float,
      "issueDate": date,
      "typeRequisition": string,
      "deleteAt": date
  }
  ```

### gRPC
- Porta **50051**
- Recomendado o uso do Evans:
  ```sh
  evans -r repl
  ```
- **Listagem de Orders:**
  - Parâmetros opcionais: `page` (int), `limit` (int), `sort` (string) -> asc ou desc
  - Saída de dados:
  ```json
  {
      "id": string,
      "price": float,
      "tax": float,
      "finalPrice": float,
      "issueDate": date,
      "typeRequisition": string,
      "deleteAt": date
  }
  ```

## Estrutura do Payload
A criação de `orders` aceita apenas os seguintes dados:
```json
{
    "price": float,
    "tax": float,
    "issueDate": date
}
```

> **Nota:** `issueDate` deve estar formatado exatamente como `"2025-02-19T15:04:05Z"`.

