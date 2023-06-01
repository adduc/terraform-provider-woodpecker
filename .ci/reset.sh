#!/bin/bash

set -e -u -o pipefail

[ "${VERBOSE:-0}" -eq 0 ] || set -x

_log () { echo "[$(date -u +"%Y-%m-%dT%H:%M:%SZ")] $1"; }

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

cd "$SCRIPT_DIR"

_log "reseting .env file..."
cat << EOF > "$SCRIPT_DIR/.env"
WOODPECKER_GITEA_CLIENT=
WOODPECKER_GITEA_SECRET=
WOODPECKER_AGENT_SECRET=
EOF

_log "removing existing instances..."
docker compose kill
docker compose rm -f

_log "creating forgejo instance..."
docker compose up -d forgejo

_log "waiting until forgejo has initialized its database..."
while ! (docker compose logs forgejo | grep -q "ORM engine initialization successful!"); do sleep 0.2; done

_log "provisioning test user in forgejo..."
docker compose exec -u git forgejo gitea admin user create --username test --password test --email "test@localhost"

_log "waiting until forgejo has started its web server..."
while ! (docker compose logs forgejo | grep -q "Starting new Web server"); do sleep 0.2; done

_log "creating empty git repository in forgejo..."
rm -rf dummy-repo
git init dummy-repo
cd dummy-repo
git checkout -b main
touch .woodpecker.yml
git add .woodpecker.yml
git config user.email "test@example.com"
git config user.name "Test User"
git commit -m "initial commit"
git remote add origin http://test:test@127.0.0.1:3000/test/test.git
git push origin main
cd ..


_log "provisioning oauth2 app in forgejo for woodpecker..."
OAUTH_APP=$(docker compose exec forgejo curl -s -X POST \
    http://127.0.0.1:3000/api/v1/user/applications/oauth2 \
    --user test:test \
    --json '{
        "name": "woodpecker",
        "redirect_uris": ["http://127.0.0.1:8000/authorize"],
        "confidential_client": true
    }')

_log "extracting oauth2 app credentials, generating agent secret..."
CLIENT_ID=$(echo "$OAUTH_APP" | sed -r 's/.*client_id":"([^"]+)".*/\1/')
CLIENT_SECRET=$(echo "$OAUTH_APP" | sed -r 's/.*client_secret":"([^"]+)".*/\1/')
AGENT_SECRET=$(openssl rand -base64 24)

_log "provisioning .env file for woodpecker..."
cat << EOF > "$SCRIPT_DIR/.env"
WOODPECKER_GITEA_CLIENT=$CLIENT_ID
WOODPECKER_GITEA_SECRET=$CLIENT_SECRET
WOODPECKER_AGENT_SECRET=$AGENT_SECRET

# Used for testing
WOODPECKER_SERVER="http://127.0.0.1:8000"
EOF

_log "creating woodpecker instance..."
docker compose up -d woodpecker

COOKIE_JAR=cookie-jar.txt

# reset cookie jar
rm -f ${COOKIE_JAR}

_log "preparing csrf token for login..."
CSRF_TOKEN=$(curl -s \
    "http://127.0.0.1:3000/user/login" \
    --cookie-jar ${COOKIE_JAR} \
  | grep _csrf | sed -r "s/.*value=\"(.*)\".*/\1/" \
)

_log "logging in to forgejo..."
curl -s http://127.0.0.1:3000/user/login \
  -X POST --cookie ${COOKIE_JAR} --cookie-jar ${COOKIE_JAR} \
  --data "_csrf=${CSRF_TOKEN}&user_name=test&password=test"

_log "preparing csrf token for oauth2 authorize..."
RESPONSE=$(curl -s \
    "http://127.0.0.1:3000/login/oauth/authorize?client_id=${CLIENT_ID}&redirect_uri=http%3A%2F%2F127.0.0.1%3A8000%2Fauthorize&response_type=code&state=woodpecker")

echo "RESPONSE: $RESPONSE"

CSRF_TOKEN=$(echo "$RESPONSE" \
    --cookie ${COOKIE_JAR} --cookie-jar ${COOKIE_JAR} \
    | grep -m1 _csrf | sed -r "s/.*value=\"(.*)\".*/\1/" \
)

_log "authorizing forgejo access to woodpecker..."
curl -L -s http://127.0.0.1:3000/login/oauth/grant \
  --cookie ${COOKIE_JAR} --cookie-jar ${COOKIE_JAR} \
  --data "_csrf=${CSRF_TOKEN}" \
  --data "client_id=${CLIENT_ID}" \
  --data "redirect_uri=http%3A%2F%2F127.0.0.1%3A8000%2Fauthorize" \
  --data "response_type=code&state=woodpecker&scope=&nonce=" \
  > /dev/null

_log "preparing csrf token for woodpecker token request..."
CSRF_TOKEN=$(curl -s http://127.0.0.1:8000/web-config.js \
  --cookie ${COOKIE_JAR} --cookie-jar ${COOKIE_JAR} \
  | grep "WOODPECKER_CSRF" \
  | sed -r 's/.*WOODPECKER_CSRF = "(.*)".*/\1/' \
)

_log "requesting woodpecker token for api access..."
WOODPECKER_TOKEN=$(curl -s http://127.0.0.1:8000/api/user/token \
  -X POST --cookie ${COOKIE_JAR} --cookie-jar ${COOKIE_JAR} \
  -H "X-CSRF-TOKEN: ${CSRF_TOKEN}" \
)

_log "WOODPECKER_TOKEN: $WOODPECKER_TOKEN"
echo "WOODPECKER_TOKEN=$WOODPECKER_TOKEN" >> "$SCRIPT_DIR/.env"
