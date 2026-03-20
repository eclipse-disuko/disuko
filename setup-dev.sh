#!/usr/bin/env bash
set -e

CERT=certs/localhost.pem
KEY=certs/localhost-key.pem
CERT_COPY_TO=backend/server.crt
KEY_COPY_TO=backend/server.key

mkdir -p certs && chmod 757 certs
mkdir -p uploads && chmod 757 uploads

if [ -f "$CERT" ] && [ -f "$KEY" ]; then
  echo "TLS certificate already exists, skipping generation."
  exit 0
fi

echo ""
echo "Choose TLS setup method:"
echo ""
echo "1) mkcert (recommended, trusted by browser)"
echo "2) self-signed certificate"
echo ""

read -p "Selection [1/2]: " choice

if [ "$choice" = "1" ]; then

  if ! command -v mkcert >/dev/null; then
    echo ""
    echo "mkcert is not installed."
    echo ""
    echo "Install mkcert:"
    echo ""
    echo "macOS:"
    echo "  brew install mkcert"
    echo "  brew install nss"
    echo ""
    echo "Linux:"
    echo "  https://github.com/FiloSottile/mkcert"
    echo ""
    echo "Windows:"
    echo "  choco install mkcert"
    echo ""
    exit 1
  fi

  echo "Installing local CA..."
  mkcert -install

  echo "Generating trusted certificate..."

  mkcert \
    -cert-file "$CERT" \
    -key-file "$KEY" \
    localhost 127.0.0.1 ::1

  echo "Trusted TLS certificate created."
  cp $CERT $CERT_COPY_TO
  cp $KEY $KEY_COPY_TO
elif [ "$choice" = "2" ]; then

  echo "Generating self-signed certificate..."

  openssl req \
    -x509 \
    -nodes \
    -days 3650 \
    -newkey rsa:2048 \
    -keyout "$KEY" \
    -out "$CERT" \
    -subj "/CN=localhost"

  echo "Self-signed TLS certificate created."
  cp $CERT $CERT_COPY_TO
  cp $KEY $KEY_COPY_TO
else
  echo "Invalid selection."
  exit 1
fi
