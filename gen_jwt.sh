#!/bin/bash

# Generate key ID (KID)
JWKS_KID=$(uuidgen)

# Generate private key in-memory, suppress verbose output
PRIVATE_KEY=$(openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 2>/dev/null)

# Extract public key in-memory, suppress verbose output
PUBLIC_KEY=$(echo "$PRIVATE_KEY" | openssl rsa -pubout 2>/dev/null)

# Generate modulus from the public key in-memory, suppress verbose output
MODULUS_HEX=$(echo "$PUBLIC_KEY" | openssl rsa -pubin -modulus -noout 2>/dev/null | cut -d '=' -f2)
JWKS_MODULUS=$(echo -n "$MODULUS_HEX" | xxd -r -p | openssl base64 -A | tr '+/' '-_' | tr -d '=')

# Set JWKS exponent
JWKS_EXPONENT="AQAB"

# Create JWT_PRIVATE_KEY by replacing newlines with \n
JWT_PRIVATE_KEY=$(echo "$PRIVATE_KEY" | sed ':a;N;$!ba;s/\n/\\n/g')

# Create JWT_PUBLIC_KEY by replacing newlines with \n
JWT_PUBLIC_KEY=$(echo "$PUBLIC_KEY" | sed ':a;N;$!ba;s/\n/\\n/g')

# Print all generated values
echo "JWT_PRIVATE_KEY=\"$JWT_PRIVATE_KEY\""
echo "JWT_PUBLIC_KEY=\"$JWT_PUBLIC_KEY\""
echo "JWKS_KID=\"$JWKS_KID\""
echo "JWKS_MODULUS=\"$JWKS_MODULUS\""
echo "JWKS_EXPONENT=\"$JWKS_EXPONENT\""
