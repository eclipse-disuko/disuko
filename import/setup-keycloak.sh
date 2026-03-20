#!/bin/bash

#####################################################################
# Keycloak Configuration Script
#
# This script uses Keycloak Admin REST API to configure
# - the users
# - the Disuko OIDC client
# - custom client scopes
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

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print colored messages
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

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
# Step 2: Create/Update Users
#####################################################################
print_info "Creating users..."

# Create customer users in a loop
for i in 1 2; do
    USERNAME="customer${i}"
    PASSWORD="CUSTOMER${i}"

    # Create user
    curl -s -X POST "${KEYCLOAK_URL}/admin/realms/${KEYCLOAK_REALM}/users" \
      -H "Authorization: Bearer ${ACCESS_TOKEN}" \
      -H "Content-Type: application/json" \
      -d "{
        \"username\": \"${USERNAME}\",
        \"firstName\": \"Customer ${i} Forename\",
        \"lastName\": \"Customer ${i} Lastname\",
        \"email\": \"${USERNAME}@company.com\",
        \"emailVerified\": false,
        \"enabled\": true,
        \"realmRoles\": [\"default-roles-master\"]
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
    fi
done

print_info "Users created"

#####################################################################
# Step 3: Create Custom Client Scopes
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
      "protocolMapper": "oidc-hardcoded-claim-mapper",
      "consentRequired": false,
      "config": {
        "introspection.token.claim": "true",
        "claim.value": "[ \"FOSSDP.policy_admin\", \"FOSSDP.domain_admin\", \"FOSSDP.application_admin\", \"FOSSDP.license_admin\", \"FOSSDP.project_analyst\" ]",
        "userinfo.token.claim": "true",
        "id.token.claim": "true",
        "lightweight.claim": "true",
        "access.token.claim": "true",
        "claim.name": "entitlement_group",
        "jsonType.label": "JSON",
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
# Step 4: Create Disuko Client
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
# Step 5: Assign Client Scopes to Disuko Client
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