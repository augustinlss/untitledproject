### Add this script in root to run the service

```bash
set -euo pipefail

export __app_id=""

# Encode service account JSON to base64
export __firebase_config="$(cat serviceAccount.json | base64)"

export MS_APP_ID=""
export MS_APP_SECRET=""
export MS_REDIRECT_URI=""
export MS_SCOPES=""

export SERVER_PORT="8080"

# --- Run server ---
echo "Starting Go Mail Gateway on port $SERVER_PORT..."
go run ./cmd/api
```
