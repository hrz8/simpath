# simpath

it's a simple auth identity provider

## How to exchange access token

### Using authorization code
```sh
curl -X POST "localhost:5001/v1/oauth/tokens" \
     -u "600ef080-d02c-426d-bf79-64247ba0fc90:test_secret" \
     -H "Content-Type: application/json" \
     -d '{
           "grant_type": "authorization_code",
           "code": "66a97c2b-c3e7-4ab8-bd0b-2dbffb5e70b9",
           "redirect_uri": "https://www.example.com"
         }'
```

### Using refresh token
```sh
curl -X POST "localhost:5001/v1/oauth/tokens" \
     -u "600ef080-d02c-426d-bf79-64247ba0fc90:test_secret" \
     -H "Content-Type: application/json" \
     -d '{
           "grant_type": "refresh_token",
           "refresh_token": "0070bcd5-278b-4c58-9d85-1b9b0afdc3c9"
         }'
```

### Using user's credentials
```sh
curl -X POST "localhost:5001/v1/oauth/tokens" \
     -u "600ef080-d02c-426d-bf79-64247ba0fc90:test_secret" \
     -H "Content-Type: application/json" \
     -d '{
           "grant_type": "refresh_token",
           "email": "test@root",
           "password": "test_password",
           "scope": "read_write"
         }'
```

## Introspect Token

## Using Access Token
```sh
curl -X POST "localhost:5001/v1/oauth/introspect" \
     -u "600ef080-d02c-426d-bf79-64247ba0fc90:test_secret" \
     -H "Content-Type: application/json" \
     -d '{
           "token": "b59f9c78-57be-44f2-8cf2-5e2506d2d3bb",
           "token_type_hint": "access_token"
         }'
```

## Using Refresh Token
```sh
curl -X POST "localhost:5001/v1/oauth/introspect" \
     -u "600ef080-d02c-426d-bf79-64247ba0fc90:test_secret" \
     -H "Content-Type: application/json" \
     -d '{
           "token": "ced70797-3566-464f-ba36-344ace1811f6",
           "token_type_hint": "refresh_token"
         }'
```
