# Local Setup and Testing Guide

Follow these steps to get the `gocart-api` running on your local machine and start testing the endpoints.

## 1. Prerequisites
- **Go**: 1.20 or later
- **Docker & Docker Compose**: For database and LocalStack (S3)
- **Postman**: For API testing

## 2. Environment Configuration
Ensure your `.env` file is set up correctly. If not already present, copy `.env.example` to `.env`:
```bash
cp .env.example .env
```
The default values in `.env.example` are pre-configured for local development avec Docker.

## 3. Start Dependencies
Use Docker Compose to start PostgreSQL and LocalStack:
```bash
docker-compose -f docker/docker-compose.yml up -d
```

> [!IMPORTANT]
> **If you see "role 'gocart' does not exist" OR "connection refused":**
> 1. **Check for local Postgres**: You might have a local PostgreSQL service running on port 5432 (e.g., via Homebrew). Stop it with:
>    ```bash
>    brew services stop postgresql@14  # or your version
>    ```
> 2. **Reset Docker**: Run these to clear old data:
>    ```bash
>    docker-compose -f docker/docker-compose.yml down -v
>    docker-compose -f docker/docker-compose.yml up -d
>    ```

## 4. Run the Application
Start the Go server:
```bash
go run cmd/server/main.go
```
The server will start on `http://localhost:8080`.

## 5. Testing with Swagger / RapiDoc
You can access the interactive documentation directly in your browser:
- **Swagger UI**: [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)
- **RapiDoc UI**: [http://localhost:8080/api-docs](http://localhost:8080/api-docs)

## 6. Testing with Postman
You can use the pre-configured Postman collection included in this repository:

1.  **Download**: Locate the [gocart-api.postman_collection.json](file:///Users/kuldeepsingh/Desktop/Golang/gocart-api/gocart-api.postman_collection.json) file in the root directory.
2.  **Import**:
    - Open Postman.
    - Click **Import** (top left).
    - Drag and drop the `gocart-api.postman_collection.json` file.
3.  **Environment**: 
    - The collection uses a `base_url` variable set to `http://localhost:8080/api/v1`.
    - Under the **Variables** tab of the collection, you can update `access_token` after logging in.
## 7. Step-by-Step Testing Flow
Follow this order to test the full lifecycle of the API:

### Phase 1: Authentication
1.  **Register**: Run `POST /auth/register` (uses default JSON body).
2.  **Login**: Run `POST /auth/login` with the same credentials.
3.  **Set Token**: 
    - Copy the `access_token` from the login response.
    - Click the **GoCart API** collection root -> **Variables** tab.
    - Paste the token into the `access_token` CURRENT VALUE field and click **Save**.

### Phase 2: Catalog Exploration
1.  **List Categories**: Run `GET /categories` to see available groups.
2.  **List Products**: Run `GET /products` to see items.
3.  **Search**: Run `GET /search?q=...` to test search functionality.

### Phase 3: Cart Management (Protected)
1.  **Add to Cart**: Run `POST /cart/items` (uses product ID from catalog).
2.  **Get Cart**: Run `GET /cart/` to verify item was added.
3.  **Update Cart**: Run `PUT /cart/items/:id` to change quantity.

### Phase 4: Checkout
1.  **Create Order**: Run `POST /orders/`. This will convert your cart into an order.
2.  **List Orders**: Run `GET /orders/` to see your purchase history.
3.  **Order Detail**: Run `GET /orders/:id` using the ID from the previous step.

### Phase 5: Admin Tasks (Optional)
*Note: To test these, you must be an admin user. By default, the first user registered might not be an admin unless manually updated in the DB.*
1.  **Create Category**: `POST /categories/`.
2.  **Create Product**: `POST /products/`.
