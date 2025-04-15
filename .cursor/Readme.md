# Cursor Project Prompt

This project consists of a **frontend** and a **backend**, stored in the following structure:

```
/frontend     # Frontend app (not detailed here)
/backend      # Backend app (REST API)
```

---

## 📦 Backend Structure

### 🗃️ `/database`

- `migrations/`: Contains database migration files
- `database.go`: Responsible for database initialization

### 🌐 `/http`

- `handlers/`: One file per resource; contains REST HTTP handlers (e.g., `user.go`, `team.go`)
- `middleware/`: One file per middleware (e.g., `auth.go`, `logging.go`)
- `server/`:
  - `routes.go`: Declares all HTTP routes
  - `server.go`: HTTP server setup and launch
- `types/`:
  - Request/response models shared across routes
  - Response types for each resource must be defined here in separate files.
  - ✅ **Use an enclosing object for all responses** (see format below).
  - These response types should be reused across multiple operations that return the same resource (e.g., `GetUser` and `CreateUser` both return a `user` object).
- `openapi.yml`:
  - The OpenAPI 3 spec for the API.
  - ⚠️ **Must be updated automatically** when routes, middleware, or payloads change.

### 🧠 `/services`

Each service is organized under its own `{resourceName}` directory:

```
/services/
  └── {resourceName}/
      ├── api/
      │   ├── api.go              # Service struct and initialization
      │   └── {method}.go         # One file per method on the service struct containing the business logic
      ├── database/
      │   └── {resourceName}.go   # DB access logic, all functions in one file
      └── types/
          └── {resourceName}.go   # Business logic structs & types
```

- Prefer **dependency injection** in all service initializations.
- Use **descriptive filenames** that match the API and method name.

### 🚀 `main.go`

- Entry point of the application
- Responsible for:
  - Configuring and initializing services
  - Injecting dependencies into handlers
  - Starting the HTTP server

---

## 🔐 Authentication Context

- The project uses an **auth middleware** to validate JWTs.
- The middleware extracts the `user_claims` from the token and injects them into the request context.
- `user_claims` contains fields such as `email`, which is used to identify the user.
- Endpoints that require user context should **extract `user_claims` from the context**, and then fetch the full user record from the database using the email.

Example logic:

```go
userClaims, exists := c.Get("user_claims")
if !exists {
  c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found in context"})
  return
}

user, err := h.userApi.FindUser(usertypes.UserSearchParams{Email: &userClaims.(*authTypes.UserClaims).Email})
if err != nil {
  c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
  return
}

if user == nil {
  c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
  return
}
```

---

## 📤 HTTP Response Format

- All HTTP responses must return data wrapped in an **enclosing object**.
- This ensures consistent and predictable API structure for frontend consumption and tooling.
- Example: `GetUser` should return:

```json
{
  "user": {
    "name": "mario"
  }
}
```

- When querying a **list of items**, the API must return an **empty array** (not null) if no results are found:

```json
{
  "users": []
}
```

- Response types should be defined in `http/types/{resourceName}/get_{resourceName}.go`.
- These types should be reused for all endpoints returning the same resource, such as `GetUser`, `CreateUser`, etc.

---

## 🧠 Cursor Instructions

- 🛠 When generating new handlers or services, use the structure described above.
- 🔄 If you generate or modify any handlers, middleware, or types, **also update `/http/openapi.yml`**.
- 💉 Always use **dependency injection** when creating service instances.
- 🔄 Keep logic separated: routing, business logic, and DB operations should live in their respective folders.
- 🔐 For protected routes, use `user_claims` from context to fetch user details.
- 📤 All HTTP responses must use an enclosing object format.
- 📁 Define response types in `http/types/{resourceName}` with filenames prefixed by `get_` for resource-returning responses.
- 📃 When returning lists, always return an empty array (`[]`) instead of `null` if no data is found.

---

Use this structure to maintain consistency and scalability in the backend architecture.
