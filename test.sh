#!/bin/bash

SERVICE_ACCOUNT_KEY="orgmdev_google.json"
BASE_URL="https://propuesta-api-645601619292.us-east4.run.app"

# Activar service account
gcloud auth activate-service-account --key-file=$SERVICE_ACCOUNT_KEY

# Generar token con audiences
TOKEN=$(gcloud auth print-identity-token --audiences=$BASE_URL)

echo "🔍 Probando endpoint..."
echo ""

curl -H "Authorization: Bearer $TOKEN" $BASE_URL

echo -e "\n✅ Request completado"