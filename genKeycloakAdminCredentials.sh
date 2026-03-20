# Usage: ./genKeycloakAdminCredentials.sh <password>
#!/bin/bash
if [ -z "$1" ]; then
  echo "Usage: $0 <password>"
  exit 1
fi
# This script generates a Keycloak admin password hash using Argon2.
# It requires the 'argon2' command-line tool to be installed.
PASSWORD="$1"
SALT=$(openssl rand -base64 16)
echo $PASSWORD
HASH=$(echo -n "$PASSWORD" | argon2 "$SALT" -t 5 -k 7168 -p 1 -l 32 -id -v 13 -e)
echo "Generated hash: $HASH"
SALT_DECODED=$(echo "$HASH" | awk -F '$' '{print $5}')==
echo "SALT DECODED: $SALT_DECODED"
PASSWORD_DECODED=$(echo "$HASH" | awk -F '$' '{print $6}')=
echo "Password decoded: $PASSWORD_DECODED"