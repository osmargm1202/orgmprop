## Propuestas API

Servicio en Python (FastAPI) para generar y servir propuestas/documentos HTML/PDF a partir de plantillas y datos. Está preparado para despliegue en Google Cloud Run y usa `uv` para manejo de dependencias y ejecución local.

### Características
- Generación y render de propuestas desde plantillas en `propuesta/`.
- Preparado para contenedor Docker y despliegue en Cloud Run.
- Scripts incluidos para prueba de endpoint autenticado con Identity Token.

---

### Requisitos
- `uv` instalado
- `docker` y acceso a `gcr.io`
- `gcloud` CLI autenticado con el proyecto correspondiente

---

### Ejecución local
```bash
uv run uvicorn main:app --host 0.0.0.0 --port 8000
```

---

### Pruebas (Cloud Run protegido con IAP/Identity Token)
Usa el script `test.sh` que:
- Activa la service account desde `orgmdev_google.json`.
- Genera un Identity Token para la `audience` del servicio.
- Realiza una petición `curl` con `Authorization: Bearer <TOKEN>`.

```bash
bash test.sh
```

Si deseas ejecutarlo manualmente:
```bash
SERVICE_ACCOUNT_KEY="orgmdev_google.json"
BASE_URL="https://<tu-servicio>.run.app"
gcloud auth activate-service-account --key-file=$SERVICE_ACCOUNT_KEY
TOKEN=$(gcloud auth print-identity-token --audiences=$BASE_URL)
curl -H "Authorization: Bearer $TOKEN" $BASE_URL
```

---

### Construcción de imagen y push a Google Container Registry (GCR)
```bash
docker build -t gcr.io/orgm-797f1/propuesta-api:latest ./
docker push gcr.io/orgm-797f1/propuesta-api:latest
```

Ejecución local del contenedor (opcional):
```bash
docker rm -f propuestas_api 2>/dev/null && \
docker run --restart always --rm -d -p 8001:8000 --name propuestas_api gcr.io/orgm-797f1/propuesta-api:latest
```

---

### Despliegue en Google Cloud Run
Puedes usar el script `deploy.sh` que:
1) Convierte `.env` a `.env.yaml`.
2) Ejecuta `gcloud run deploy` con la imagen `gcr.io/orgm-797f1/propuesta-api:latest` en `us-east4`.

```bash
bash deploy.sh
```

Parámetros relevantes en `deploy.sh`:
- `--no-allow-unauthenticated` (requiere Identity Token)
- `--env-vars-file .env.yaml` (carga variables de entorno)
- `--port 8000` (puerto de la app)

---

### Estado pendiente
- Agregar soporte multi-tenant en la base de datos para separar clientes (campo `tenant_id` o esquemas por cliente según necesidad). Se gestionará vía Alembic/Django ORM según definición final.
