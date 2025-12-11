# API de Propuestas - Documentación

## Versión 2.0.0

La API ha sido simplificada para generar propuestas directamente en HTML, eliminando completamente la generación de Markdown y PDF. Ahora todo se maneja en una sola llamada que genera el HTML completo.

## Cambios Principales

### Eliminado
- ❌ Generación de Markdown (MD)
- ❌ Generación de PDF
- ❌ Endpoints `/generate-text`, `/generate-html`, `/generate-pdf`
- ❌ Endpoints `/proposal/{id}/pdf` y `/proposal/{id}/md`

### Nuevo
- ✅ Endpoint único `/generate-proposal` que genera HTML directamente
- ✅ Generación en una sola llamada combinando ambos YAML (prompt + html)
- ✅ CSS incluido directamente en el HTML generado
- ✅ Estructura homogénea con clases simplificadas

## Endpoints

### POST /generate-proposal

Genera una propuesta completa en HTML en una sola llamada.

**Request:**
```json
{
  "title": "Título de la Propuesta",
  "subtitle": "Subtítulo de la Propuesta",
  "prompt": "Descripción detallada de los servicios, alcances, costos, etc.",
  "model": "gpt-5-chat-latest"  // Opcional, por defecto gpt-5-chat-latest
}
```

**Response:**
```json
{
  "id": "abc12345",
  "html_url": "https://r2.endpoint.com/bucket/propuestas/html/abc12345.html",
  "created_at": "2024-01-15T10:30:00"
}
```

**Ejemplo:**
```bash
curl -X POST "http://localhost:8000/generate-proposal" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Propuesta de Servicios",
    "subtitle": "Instalación de Medidor Testigo",
    "prompt": "Servicio de instalación de medidor testigo para certificación..."
  }'
```

### GET /proposals

Lista todas las propuestas.

**Response:**
```json
[
  {
    "id": "abc12345",
    "title": "Título",
    "subtitle": "Subtítulo",
    "prompt": "Prompt original",
    "html_url": "/proposal/abc12345/html",
    "created_at": "2024-01-15T10:30:00",
    "size_html": 45678
  }
]
```

### GET /proposals/search?q=query&limit=10

Busca propuestas por título o subtítulo.

**Parámetros:**
- `q`: Query de búsqueda (tokens separados por espacios)
- `limit`: Límite de resultados (default: 10)

### GET /proposal/{proposal_id}/html

Descarga el archivo HTML de una propuesta.

**Response:** Archivo HTML descargable

### PUT /proposal/{proposal_id}

Modifica una propuesta existente regenerando el HTML.

**Request:**
```json
{
  "title": "Nuevo Título",
  "subtitle": "Nuevo Subtítulo",
  "prompt": "Instrucciones de modificación...",
  "model": "gpt-5-chat-latest"  // Opcional
}
```

**Response:**
```json
{
  "id": "abc12345",
  "html_url": "https://r2.endpoint.com/bucket/propuestas/html/abc12345.html",
  "created_at": "2024-01-15T10:30:00"
}
```

### POST /proposal/{proposal_id}/regenerate

Regenera una propuesta desde cero con nuevos datos.

**Request:**
```json
{
  "title": "Nuevo Título",
  "subtitle": "Nuevo Subtítulo",
  "prompt": "Nuevo prompt completo...",
  "model": "gpt-5-chat-latest"  // Opcional
}
```

**Response:**
```json
{
  "id": "abc12345",
  "html_url": "https://r2.endpoint.com/bucket/propuestas/html/abc12345.html",
  "created_at": "2024-01-15T10:30:00"
}
```

### PATCH /proposal/{proposal_id}/title-subtitle

Actualiza solo el título y subtítulo sin regenerar el contenido.

**Request:**
```json
{
  "title": "Nuevo Título",
  "subtitle": "Nuevo Subtítulo"
}
```

**Response:**
```json
{
  "id": "abc12345",
  "html_url": "https://r2.endpoint.com/bucket/propuestas/html/abc12345.html",
  "created_at": "2024-01-15T10:30:00"
}
```

### GET /

Health check endpoint.

**Response:**
```json
{
  "status": "healthy",
  "service": "propuestas-api",
  "version": "2.0.0"
}
```

## Flujo de Trabajo

### Crear una Nueva Propuesta

1. **Generar propuesta:**
   ```bash
   POST /generate-proposal
   ```
   - Recibe: title, subtitle, prompt
   - Genera: HTML completo con CSS incluido
   - Retorna: proposal_id y html_url

2. **Ver propuesta:**
   ```bash
   GET /proposal/{proposal_id}/html
   ```
   - Descarga el HTML generado

### Modificar una Propuesta Existente

1. **Modificar contenido:**
   ```bash
   PUT /proposal/{proposal_id}
   ```
   - Regenera HTML con modificaciones
   - Mantiene el mismo proposal_id

2. **Regenerar completamente:**
   ```bash
   POST /proposal/{proposal_id}/regenerate
   ```
   - Regenera desde cero con nuevos datos

3. **Solo cambiar título/subtítulo:**
   ```bash
   PATCH /proposal/{proposal_id}/title-subtitle
   ```
   - Actualiza solo metadatos sin regenerar HTML

## Estructura del HTML Generado

El HTML generado incluye:

- **CSS completo** dentro de `<style>` tag (no links externos)
- **Bootstrap Icons** para iconos
- **Estructura homogénea** con clases:
  - `.intro-section` - Introducción
  - `.service-grid` - Grid de servicios
  - `.service-card` - Tarjetas de servicio
  - `.pricing-simple` - Precios
  - `.payment-simple` - Forma de pago
  - `.timeline-simple` - Cronograma
  - `.info-box` - Información adicional
  - `.note-box` - Notas importantes

## Base de Datos

### Columnas Mantenidas (No Usadas)
- `md_url` - Se mantiene en DB pero no se usa
- `pdf_url` - Se mantiene en DB pero no se usa

### Columnas Activas
- `id` - ID único de la propuesta
- `title` - Título
- `subtitle` - Subtítulo
- `prompt` - Prompt original
- `html_url` - URL del HTML en R2
- `proposal_metadata` - Metadatos JSON
- `created_at` - Fecha de creación
- `updated_at` - Fecha de actualización

## Almacenamiento

- **R2 (Cloudflare)**: Los archivos HTML se almacenan en R2
- **Ruta**: `propuestas/html/{proposal_id}.html`
- **Base de datos**: SQLite/PostgreSQL con metadatos

## Configuración

### Variables de Entorno Requeridas

```bash
OPENAI_API_KEY=sk-...
ROLLBAR_ACCESS_TOKEN=...
ROLLBAR_ENVIRONMENT=production
CLOUDFLARE_R2_ACCESS_KEY_ID=...
CLOUDFLARE_R2_SECRET_ACCESS_KEY=...
CLOUDFLARE_R2_ENDPOINT_URL=https://...
CLOUDFLARE_R2_BUCKET=propuestas-bucket
DATABASE_URL=postgresql://...
```

## Templates

La generación combina dos archivos YAML:

1. **`templates/prompt/propuesta.yaml`**: Define la estructura del contenido
2. **`templates/html/propuesta.yaml`**: Define la estructura HTML
3. **`templates/css/template.css`**: CSS incluido en el HTML

La función `generate_proposal_html()` combina ambos YAML en una sola instrucción para OpenAI, generando el HTML completo en una sola llamada.

## Ejemplos de Uso

### Python

```python
import requests

# Generar propuesta
response = requests.post(
    "http://localhost:8000/generate-proposal",
    json={
        "title": "Propuesta de Servicios",
        "subtitle": "Instalación de Medidor",
        "prompt": "Servicio completo de instalación..."
    }
)
proposal = response.json()
print(f"Propuesta creada: {proposal['id']}")
print(f"HTML: {proposal['html_url']}")

# Descargar HTML
html_response = requests.get(
    f"http://localhost:8000/proposal/{proposal['id']}/html"
)
with open("propuesta.html", "wb") as f:
    f.write(html_response.content)
```

### JavaScript/Node.js

```javascript
// Generar propuesta
const response = await fetch('http://localhost:8000/generate-proposal', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    title: 'Propuesta de Servicios',
    subtitle: 'Instalación de Medidor',
    prompt: 'Servicio completo de instalación...'
  })
});

const proposal = await response.json();
console.log(`Propuesta creada: ${proposal.id}`);
console.log(`HTML: ${proposal.html_url}`);
```

## Notas Importantes

1. **Una sola llamada**: La generación de HTML se hace en una sola llamada a OpenAI combinando ambos YAML
2. **CSS incluido**: El CSS se incluye directamente en el HTML, no hay dependencias externas
3. **Sin MD/PDF**: Ya no se generan archivos Markdown ni PDF
4. **Base de datos**: Las columnas md_url y pdf_url se mantienen pero no se usan
5. **Compatibilidad**: Los endpoints antiguos han sido eliminados



