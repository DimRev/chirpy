# Chirpy

This is a tut project via boot.dev learning learning go by building a go REST API server.

# APIs

## Health check

### GET /api/healthz

Responds with a `200` status if the server is running

## Chirps

### POST /api/chirps

> Request header:
>
> ```json
> "Authorization":"Bearer ${Token}"
> ```
>
> Request body:
>
> ```json
> { "body": "This is a chirp" }
> ```

> Responds with:
>
> ```json
> {
>   "id": 1,
>   "body": "This is a chirp",
>   "author_id": 2
> }
> ```

### GET /api/chirps

> Request Search Params `200`:
>
> - `author_id` - filter by author id
> - `sort` - "asc" | "desc", ascending or descending order
> - `api/chirps?author_id=1&sort=desc`

> Responds with:
>
> ```json
> [
>   {
>     "id": 1,
>     "body": "This is a chirp",
>     "author_id": 2
>   },...
> ]
> ```

### GET /api/chirps/{id}

> Responds with:
>
> ```json
> {
>   "id": 1,
>   "body": "This is a chirp",
>   "author_id": 2
> }
> ```

### DELETE /api/chirps

> Request header:
>
> ```json
> "Authorization":"Bearer ${Token}"
> ```

> Resound with `200`

## Users

### POST /api/users

> Request body:
>
> ```json
> {
>   "email": "example@example.com",
>   "password": "123123"
> }
> ```

> Responds with:
>
> ```json
> {
>   "email": "example@example.com",
>   "id": 1,
>   "is_chirp_red": false
> }
> ```

### PUT /api/users

> Request header:
>
> ```json
> "Authorization":"Bearer ${Token}"
> ```
>
> Request body:
>
> ```json
> {
>   "email": "example@example.com",
>   "password": "123123"
> }
> ```

> Responds with:
>
> ```json
> {
>   "email": "example@example.com",
>   "id": 1,
>   "is_chirp_red": false
> }
> ```

## Auth

### POST /api/login

> Request body:
>
> ```json
> {
>   "email": "example@example.com",
>   "password": "123123"
> }
> ```

> Responds with:
>
> ```json
> {
>   "email": "example@example.com",
>   "id": 1,
>   "is_chirp_red": false,
>   "token": "{token}",
>   "refresh_token": "{refresh_token}"
> }
> ```

### POST /api/refresh

> Request header:
>
> ```json
> "Authorization":"Bearer ${Token}"
> ```

> Responds with:
>
> ```json
> {
>   "token": "{token}"
> }
> ```

### POST /api/revoke

> Request header:
>
> ```json
> "Authorization":"Bearer ${Token}"
> ```

> Responds with `200`

### POST /api/polka/webhook

> Request header:
>
> ```json
> "Authorization":"ApiKey ${Token}"
> ```

> Responds with:
>
> ```json
> {
>   "id": 1,
>   "email": "example@example.com",
>   "is_chirp_red": true
> }
> ```

## Deprecated

### POST /api/validate_chirp (Deprecated)

# .Env

```evn
JWT_SECRET=""
POLKA_API_KEY=""
```
