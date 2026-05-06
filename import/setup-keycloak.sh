#!/bin/bash

#####################################################################
# Keycloak Configuration Script
#
# This script uses the Keycloak Admin REST API to configure a local
# Keycloak instance for Disuko development.
#
# Steps:
#   1. Obtain an admin access token via the admin-cli client
#   2. Derive all unique FOSSDP group names from users.json
#      and create them as Keycloak groups. Groups are used instead of
#      realm roles so that the oidc-group-membership-mapper emits only
#      FOSSDP.* values in the entitlement_group claim, without any
#      Keycloak built-in roles (e.g. default-roles-master, offline_access).
#   3. Create test users defined in users.json and assign
#      the configured FOSSDP groups to each user. The mapping drives
#      both user creation and group assignment — to add or change a
#      user's entitlement groups, only users.json needs
#      to be updated.
#   4. Create custom OpenID Connect client scopes that simulate the
#      claims provided by the production IAM system:
#        - sub               (username as subject)
#        - first_name        (given name)
#        - last_name         (family name)
#        - group_type        (hardcoded: "0" = employee)
#        - object_class      (hardcoded LDAP object classes including
#                             dcxInternalEmployee)
#        - company_identifier (hardcoded: "0001")
#        - department        (hardcoded: "AG")
#        - department_description (empty string)
#        - authorization_group    (empty scope, reserved for future use)
#        - personal_data          (empty scope, reserved for future use)
#        - organizational_data    (empty scope, reserved for future use)
#        - entitlement_group (group membership mapper → emits only
#                             the user's FOSSDP.* groups as a JSON
#                             array)
#   5. Create the Disuko OIDC client with the required redirect URI,
#      client secret, and service account configuration.
#   6. Assign all custom scopes as default scopes and offline_access
#      as optional scope to the Disuko client.
#####################################################################

set -e  # Exit on error

# Configuration variables
KEYCLOAK_URL="${KEYCLOAK_URL:-http://keycloak:8080}"
KEYCLOAK_USER="${KEYCLOAK_USER:-admin}"
KEYCLOAK_PASSWORD="${KEYCLOAK_PASSWORD:-password}"
KEYCLOAK_REALM="${KEYCLOAK_REALM:-master}"
DISUKO_CLIENT_ID="${DISUKO_CLIENT_ID:-243e5c8-9b1a-4c3d-9f0e-7b2a1c8e5f6ac}"
DISUKO_HOST="${DISUKO_HOST:-https://localhost:3009}"
DISUKO_CLIENT_SECRET="${DISUKO_CLIENT_SECRET:-RST845JLOP8x9Z2n1QFDA25A1B2C3D4k}"

# Customer to entitlement group mapping — loaded from external config
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
USERS_JSON="${USERS_JSON:-${SCRIPT_DIR}/users.json}"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

if [ ! -f "$USERS_JSON" ]; then
    print_error "Users config file not found: ${USERS_JSON}"
    exit 1
fi

#####################################################################
# Step 1: Obtain Admin Access Token
#####################################################################
print_info "Obtaining admin access token..."

TOKEN_RESPONSE=$(curl -s -X POST "${KEYCLOAK_URL}/realms/master/protocol/openid-connect/token" \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=${KEYCLOAK_USER}" \
  -d "password=${KEYCLOAK_PASSWORD}" \
  -d "grant_type=password" \
  -d "client_id=admin-cli")

ACCESS_TOKEN=$(echo "$TOKEN_RESPONSE" | grep -o '"access_token":"[^"]*' | cut -d'"' -f4)

if [ -z "$ACCESS_TOKEN" ]; then
    print_error "Failed to obtain access token"
    echo "$TOKEN_RESPONSE"
    exit 1
fi

print_info "Access token obtained successfully"

#####################################################################
# Step 2: Create Groups for entitlement_group
#####################################################################
print_info "Creating groups..."

# Collect all unique group names from the JSON config
UNIQUE_GROUPS=$(jq -r '.users[].groups[]' "$USERS_JSON" | sort -u)

for GROUP in $UNIQUE_GROUPS; do
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/groups" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{
        \"name\": \"${GROUP}\"
      }" 2>/dev/null || print_warning "Group ${GROUP} may already exist"
done

print_info "Groups created"

#####################################################################
# Step 3: Create/Update Users
#####################################################################
print_info "Creating users..."

USER_COUNT=$(jq '.users | length' "$USERS_JSON")
i=0
while [ "$i" -lt "$USER_COUNT" ]; do
    USERNAME=$(jq -r ".users[$i].username" "$USERS_JSON")
    GROUPS_CSV=$(jq -r ".users[$i].groups | join(\",\")" "$USERS_JSON")
    INDEX=$(echo "$USERNAME" | sed 's/[^0-9]//g')
    PASSWORD=$(echo "$USERNAME" | tr '[:lower:]' '[:upper:]')

    # Create user
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/users" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${USERNAME}\",
        \"firstName\": \"Customer ${INDEX} Forename\",
        \"lastName\": \"Customer ${INDEX} Lastname\",
        \"email\": \"${USERNAME}@company.com\",
        \"emailVerified\": false,
        \"enabled\": true
      }" 2>/dev/null || print_warning "User ${USERNAME} may already exist"

    # Get user ID and set password
    USER_ID=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/users?username=${USERNAME}" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}" | grep -o '"id":"[^"]*' | head -1 | cut -d'"' -f4)

    if [ -n "$USER_ID" ]; then
        curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/users/${USER_ID}/reset-password" \
          -H "Authorization: Bearer ${ACCESS_TOKEN}" \
          -H "Content-Type: application/json" \
          -d "{
            \"type\": \"password\",
            \"value\": \"${PASSWORD}\",
            \"temporary\": false
          }"
        print_info "Password set for ${USERNAME}"

        # Get all groups once
        ALL_GROUPS=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/groups" \
          -H "Authorization: Bearer ${ACCESS_TOKEN}")

        GROUP_COUNT=0
        for GROUP_NAME in $(echo "$GROUPS_CSV" | tr ',' ' '); do
            GROUP_ID=$(echo "$ALL_GROUPS" | jq -r ".[] | select(.name==\"${GROUP_NAME}\") | .id")
            if [ -n "$GROUP_ID" ] && [ "$GROUP_ID" != "null" ]; then
                curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/users/${USER_ID}/groups/${GROUP_ID}" \
                  -H "Authorization: Bearer ${ACCESS_TOKEN}"
                GROUP_COUNT=$((GROUP_COUNT + 1))
            else
                print_warning "Group not found: ${GROUP_NAME}"
            fi
        done

        print_info "Assigned ${GROUP_COUNT} groups to ${USERNAME}"
    fi
    i=$((i + 1))
done

print_info "Users created"

#####################################################################
# Step 4: Create Custom Client Scopes
#####################################################################
print_info "Creating custom client scopes..."

# Create authorization_group scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "authorization_group",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    }
  }' 2>/dev/null || print_warning "Scope authorization_group may already exist"

# Create last_name scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "last_name",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "last_name",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-usermodel-attribute-mapper",
      "consentRequired": false,
      "config": {
        "aggregate.attrs": "false",
        "introspection.token.claim": "false",
        "multivalued": "false",
        "userinfo.token.claim": "false",
        "user.attribute": "lastName",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "last_name",
        "jsonType.label": "String"
      }
    }]
  }' 2>/dev/null || print_warning "Scope last_name may already exist"

# Create object_class scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "object_class",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "object_class",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "claim.value": "[\"top\",\"person\",\"organizationalPerson\",\"inetOrgPerson\",\"dcxPerson\",\"dcxEmployee\",\"dcxInternalEmployee\",\"dcxADPerson\"]",
        "userinfo.token.claim": "false",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "object_class",
        "jsonType.label": "JSON",
        "access.tokenResponse.claim": "false"
      }
    }]
  }' 2>/dev/null || print_warning "Scope object_class may already exist"

# Create group_type scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "group_type",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "group_type",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "claim.value": "0",
        "userinfo.token.claim": "false",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "group_type",
        "jsonType.label": "String",
        "access.tokenResponse.claim": "false"
      }
    }]
  }' 2>/dev/null || print_warning "Scope group_type may already exist"

# Create company_identifier scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "company_identifier",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "company_identifier",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "claim.value": "0001",
        "userinfo.token.claim": "false",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "company_identifier",
        "jsonType.label": "String",
        "access.tokenResponse.claim": "false"
      }
    }]
  }' 2>/dev/null || print_warning "Scope company_identifier may already exist"

# Create personal_data scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "personal_data",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    }
  }' 2>/dev/null || print_warning "Scope personal_data may already exist"

# Create department_description scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "department_description",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "department_description",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "userinfo.token.claim": "false",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "department_description",
        "jsonType.label": "String",
        "access.tokenResponse.claim": "false"
      }
    }]
  }' 2>/dev/null || print_warning "Scope department_description may already exist"

# Create organizational_data scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "organizational_data",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    }
  }' 2>/dev/null || print_warning "Scope organizational_data may already exist"

# Create entitlement_group scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "entitlement_group",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "entitlement_group",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-group-membership-mapper",
      "consentRequired": false,
      "config": {
        "full.path": "false",
        "introspection.token.claim": "true",
        "userinfo.token.claim": "true",
        "id.token.claim": "true",
        "lightweight.claim": "true",
        "access.token.claim": "true",
        "claim.name": "entitlement_group",
        "access.tokenResponse.claim": "true"
      }
    }]
  }' 2>/dev/null || print_warning "Scope entitlement_group may already exist"

# Create department scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "department",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "department",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "claim.value": "AG",
        "userinfo.token.claim": "false",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "department",
        "jsonType.label": "String",
        "access.tokenResponse.claim": "false"
      }
    }]
  }' 2>/dev/null || print_warning "Scope department may already exist"

# Create first_name scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "first_name",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "first_name",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-usermodel-attribute-mapper",
      "consentRequired": false,
      "config": {
        "aggregate.attrs": "false",
        "introspection.token.claim": "false",
        "multivalued": "false",
        "userinfo.token.claim": "false",
        "user.attribute": "firstName",
        "id.token.claim": "true",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "first_name",
        "jsonType.label": "String"
      }
    }]
  }' 2>/dev/null || print_warning "Scope first_name may already exist"

# Create sub scope
curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "sub",
    "description": "",
    "protocol": "openid-connect",
    "attributes": {
      "include.in.token.scope": "false",
      "display.on.consent.screen": "false"
    },
    "protocolMappers": [{
      "name": "sub",
      "protocol": "openid-connect",
      "protocolMapper": "oidc-usermodel-attribute-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "false",
        "userinfo.token.claim": "true",
        "user.attribute": "username",
        "id.token.claim": "false",
        "lightweight.claim": "false",
        "access.token.claim": "false",
        "claim.name": "sub",
        "jsonType.label": "String"
      }
    }]
  }' 2>/dev/null || print_warning "Scope sub may already exist"

print_info "Custom client scopes created"

#####################################################################
# Step 5: Create Disuko Client
#####################################################################
print_info "Creating Disuko client..."

curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/clients" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" \
  -H "Content-Type: application/json" \
  -d "{
    \"clientId\": \"${DISUKO_CLIENT_ID}\",
    \"name\": \"Disuko\",
    \"description\": \"\",
    \"rootUrl\": \"${DISUKO_HOST}\",
    \"adminUrl\": \"${DISUKO_HOST}\",
    \"baseUrl\": \"${DISUKO_HOST}\",
    \"enabled\": true,
    \"clientAuthenticatorType\": \"client-secret\",
    \"secret\": \"${DISUKO_CLIENT_SECRET}\",
    \"redirectUris\": [\"${DISUKO_HOST}/api/v1/login\"],
    \"webOrigins\": [\"${DISUKO_HOST}\"],
    \"bearerOnly\": false,
    \"consentRequired\": false,
    \"standardFlowEnabled\": true,
    \"implicitFlowEnabled\": false,
    \"directAccessGrantsEnabled\": false,
    \"serviceAccountsEnabled\": true,
    \"authorizationServicesEnabled\": true,
    \"publicClient\": false,
    \"frontchannelLogout\": true,
    \"protocol\": \"openid-connect\",
    \"attributes\": {
      \"oidc.ciba.grant.enabled\": \"false\",
      \"backchannel.logout.session.required\": \"true\",
      \"standard.token.exchange.enabled\": \"true\",
      \"frontchannel.logout.session.required\": \"true\",
      \"display.on.consent.screen\": \"false\",
      \"oauth2.device.authorization.grant.enabled\": \"false\",
      \"backchannel.logout.revoke.offline.tokens\": \"false\"
    },
    \"fullScopeAllowed\": true,
    \"protocolMappers\": [
      {
        \"name\": \"client roles\",
        \"protocol\": \"openid-connect\",
        \"protocolMapper\": \"oidc-usermodel-client-role-mapper\",
        \"consentRequired\": false,
        \"config\": {
          \"introspection.token.claim\": \"true\",
          \"multivalued\": \"true\",
          \"userinfo.token.claim\": \"false\",
          \"user.attribute\": \"foo\",
          \"id.token.claim\": \"true\",
          \"lightweight.claim\": \"false\",
          \"access.token.claim\": \"true\",
          \"claim.name\": \"resource_access.\${client_id}.roles\",
          \"jsonType.label\": \"String\"
        }
      },
      {
        \"name\": \"Client ID\",
        \"protocol\": \"openid-connect\",
        \"protocolMapper\": \"oidc-usersessionmodel-note-mapper\",
        \"consentRequired\": false,
        \"config\": {
          \"user.session.note\": \"client_id\",
          \"id.token.claim\": \"true\",
          \"introspection.token.claim\": \"true\",
          \"access.token.claim\": \"true\",
          \"claim.name\": \"client_id\",
          \"jsonType.label\": \"String\"
        }
      },
      {
        \"name\": \"Client Host\",
        \"protocol\": \"openid-connect\",
        \"protocolMapper\": \"oidc-usersessionmodel-note-mapper\",
        \"consentRequired\": false,
        \"config\": {
          \"user.session.note\": \"clientHost\",
          \"id.token.claim\": \"true\",
          \"introspection.token.claim\": \"true\",
          \"access.token.claim\": \"true\",
          \"claim.name\": \"clientHost\",
          \"jsonType.label\": \"String\"
        }
      },
      {
        \"name\": \"Client IP Address\",
        \"protocol\": \"openid-connect\",
        \"protocolMapper\": \"oidc-usersessionmodel-note-mapper\",
        \"consentRequired\": false,
        \"config\": {
          \"user.session.note\": \"clientAddress\",
          \"id.token.claim\": \"true\",
          \"introspection.token.claim\": \"true\",
          \"access.token.claim\": \"true\",
          \"claim.name\": \"clientAddress\",
          \"jsonType.label\": \"String\"
        }
      }
    ]
  }" 2>/dev/null || print_warning "Client may already exist"

print_info "Disuko client created"

#####################################################################
# Step 6: Assign Client Scopes to Disuko Client
#####################################################################
print_info "Assigning client scopes..."

# Get Disuko client UUID
DISUKO_CLIENT_UUID=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/clients?clientId=${DISUKO_CLIENT_ID}" \
  -H "Authorization: Bearer ${ACCESS_TOKEN}" | jq -r '.[0].id')

if [ -z "$DISUKO_CLIENT_UUID" ] || [ "$DISUKO_CLIENT_UUID" = "null" ]; then
    print_warning "Could not find Disuko client UUID"
else
    print_info "Disuko client UUID: $DISUKO_CLIENT_UUID"

    # Get all client scope IDs
    ALL_SCOPES=$(curl -s -X GET "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/client-scopes" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}")

    # Default client scopes to assign
    for scope_name in "sub" "profile" "roles" "authorization_group" "last_name" "group_type" \
                      "company_identifier" "web-origins" "acr" "personal_data" \
                      "department_description" "organizational_data" "entitlement_group" \
                      "basic" "department" "object_class" "first_name" "email"; do

        SCOPE_ID=$(echo "$ALL_SCOPES" | jq -r ".[] | select(.name==\"$scope_name\") | .id")

        if [ -n "$SCOPE_ID" ] && [ "$SCOPE_ID" != "null" ]; then
            curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/clients/${DISUKO_CLIENT_UUID}/default-client-scopes/${SCOPE_ID}" \
              -H "Authorization: Bearer ${ACCESS_TOKEN}" 2>/dev/null
            print_info "  Added default scope: $scope_name"
        else
            print_warning "  Scope not found: $scope_name"
        fi
    done

    # Optional client scopes to assign
    for scope_name in "address" "phone" "offline_access" "microprofile-jwt"; do
        SCOPE_ID=$(echo "$ALL_SCOPES" | jq -r ".[] | select(.name==\"$scope_name\") | .id")

        if [ -n "$SCOPE_ID" ] && [ "$SCOPE_ID" != "null" ]; then
            curl -s -X PUT "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/clients/${DISUKO_CLIENT_UUID}/optional-client-scopes/${SCOPE_ID}" \
              -H "Authorization: Bearer ${ACCESS_TOKEN}" 2>/dev/null
            print_info "  Added optional scope: $scope_name"
        fi
    done
fi