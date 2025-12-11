### Template mínimo para backend Python en contenedores
Estructura orientada a: FastAPI + Uvicorn + Rollbar + SQLModel + Alembic + psycopg2-binary + boto3 + uv. Incluye `.dockerignore`. Pensado para clonarse como base y que una IA pueda completarlo para cualquier objetivo de API.

- Contiene:
  - `pyproject.toml` con deps clave
  - `Dockerfile` basado en `uv`
  - `.dockerignore`
  - `.env.example`
  - `app/` con FastAPI, DB y Rollbar
  - `storage/` con cliente `boto3`
  - `alembic/` listo para migraciones
  - `scripts/` con `deploy.sh` y `test.sh`
  - `README.md` con comandos esenciales

```bash
proyecto-template/
  app/
    __init__.py
    main.py
    config.py
    db.py
    models.py
    deps.py
    routers/
      __init__.py
      health.py
  storage/
    __init__.py
    s3.py
  alembic/
    env.py
    script.py.mako
    versions/  # vacío inicialmente
  .dockerignore
  .env.example
  alembic.ini
  Dockerfile
  pyproject.toml
  README.md
  scripts/
    deploy.sh
    test.sh
```

```toml
# pyproject.toml
[project]
name = "proyecto-template"
version = "0.1.0"
description = "Template base para APIs Python en contenedores (FastAPI + SQLModel + Alembic + Rollbar + S3)."
readme = "README.md"
requires-python = ">=3.11"
dependencies = [
  "fastapi>=0.110.0",
  "uvicorn[standard]>=0.27.0",
  "pydantic>=2.6.0",
  "python-dotenv>=1.0.0",
  "rollbar>=0.16.0",
  "sqlmodel>=0.0.14",
  "alembic>=1.13.0",
  "psycopg2-binary>=2.9.0",
  "boto3>=1.34.0",
  "jinja2>=3.1.0",
  "weasyprint>=60.0",
  "openai>=1.0.0"
]
```

### variables de entorno
.env

ROLLBAR_ACCESS_TOKEN=""
ROLLBAR_ENVIRONMENT=""
APP_REVISION=""
#base de datos""
DATABASE_URL=""
#openai api key
OPENAI_API_KEY=""
#cloudflare r2
CLOUDFLARE_R2_ACCESS_KEY_ID=""
CLOUDFLARE_R2_SECRET_ACCESS_KEY=""
CLOUDFLARE_R2_BUCKET=""
CLOUDFLARE_R2_ENDPOINT_URL=""
DATABASE_URL_OLD=""