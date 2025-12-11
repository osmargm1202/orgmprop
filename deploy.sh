#!/bin/bash

# Crear archivo temporal en formato YAML
echo "Convirtiendo .env a YAML..."
cat .env | grep -v '^#' | grep -v '^$' | while IFS='=' read -r key value; do
  # Reunir el resto de la línea como valor (por si hay '=' en el valor)
  if [[ "$key" != "" && "$value" != "" ]]; then
    rest="${value}"
    # Leer el resto de la línea si hay más '='
    if [[ "$rest" == *"="* ]]; then
      rest="${rest#*=}"
      value="${value%%=*}=${rest}"
    fi
    echo "$key: $value"
  fi
done > .env.yaml

echo "Desplegando propuesta-api..."

gcloud run deploy propuesta-api \
  --image gcr.io/orgm-797f1/propuesta-api:latest \
  --region us-east4 \
  --no-allow-unauthenticated \
  --platform managed \
  --env-vars-file .env.yaml \
  --memory 256Mi \
  --cpu 1 \
  --timeout 300 \
  --min-instances 0 \
  --max-instances 1 \
  --cpu-boost \
  --port 8000

# Limpiar archivo temporal
# rm .env.yaml

echo "✅ Despliegue completado"


