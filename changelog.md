# Changelog - Sistema de Propuestas ORGM

## Objetivo del Proyecto
Sistema automatizado para generar propuestas comerciales profesionales de ORGM (Consultoría Energética y Ambiental) con diferentes formatos y plantillas. El sistema permite crear propuestas en HTML optimizadas para impresión con WeasyPrint, manteniendo consistencia visual y profesionalismo en todos los documentos.

## Changelog

### [2024-12-19] - Separación de generación de texto y HTML
- **Objetivo**: Separar la generación de texto de la generación de HTML para mayor flexibilidad
- **Cambios realizados**:
  - **Modelos actualizados**:
    - Agregado campo `md_url` al modelo `Proposal` para almacenar URL del archivo markdown
    - Creados nuevos modelos `TextGenerationRequest` y `HTMLGenerationRequest`
    - Actualizado `ProposalUpdate` para incluir `md_url`
  - **Nuevos endpoints**:
    - `POST /generate-text`: Genera solo el contenido de texto y lo guarda en base de datos
    - `POST /generate-html`: Genera HTML desde texto existente o ID de propuesta
  - **Funcionalidades**:
    - Subida de archivos .md a carpeta `md/` en R2
    - Generación de texto usando `propuesta.yaml`
    - Generación de HTML usando `propuestahtml.yaml` con contenido existente
    - Soporte para generar HTML desde texto directo o desde propuesta existente
  - **Archivos modificados**:
    - `models.py`: Nuevos modelos y campo md_url
    - `main.py`: Nuevos endpoints y funciones de upload
    - `database.py`: Actualización para manejar md_url
  - **Archivos generados**:
    - `test_separated_endpoints.py`: Script de prueba para los nuevos endpoints
  - **Resultado**: Sistema más flexible que permite generar texto y HTML por separado

### [2024-12-19] - Corrección de errores en endpoints separados
- **Objetivo**: Corregir errores de validación y restricciones de base de datos
- **Cambios realizados**:
  - **Migraciones de base de datos**:
    - Creada migración para agregar campo `md_url` a tabla `proposals`
    - Creada migración para hacer `html_url` y `pdf_url` nullable
    - Aplicadas migraciones exitosamente
  - **Corrección de modelos**:
    - `ProposalResponse`: Campos `html_url` y `pdf_url` ahora son opcionales
    - `ProposalListItem`: Agregados campos `title` y `subtitle`, campos URL opcionales
  - **Corrección de endpoints**:
    - Corregida subida de archivos .md a carpeta `propuestas/md/`
    - Corregido manejo de metadatos en `ProposalUpdate`
    - Corregida función `list_proposals` para incluir todos los campos
  - **Archivos modificados**:
    - `alembic/versions/`: Nuevas migraciones de base de datos
    - `main.py`: Corrección de modelos y endpoints
    - `models.py`: Campos opcionales en modelos de respuesta
  - **Resultado**: Todos los endpoints funcionando correctamente, sistema completamente operativo

### [2024-12-19] - Aplicación de formato template a propuesta medidor
- **Objetivo**: Aplicar el formato del template.html a la propuesta_medidor.html para mantener consistencia visual
- **Cambios realizados**:
  - Análisis de diferencias entre propuesta_medidor.html y template.html
  - Adaptación del CSS de template.html para propuesta_medidor.html
  - Eliminación de elementos específicos del navegador (botón de exportación, animaciones complejas)
  - Optimización para impresión con WeasyPrint
  - Mantenimiento del contenido específico de la propuesta del medidor testigo
  - Creación de changelog.md para seguimiento del proyecto
  - **Archivos generados**:
    - `propuesta_medidor_template.html` - Versión optimizada con formato del template
    - `propuesta_medidor_template.pdf` - PDF generado exitosamente con WeasyPrint
  - **Resultado**: PDF de 136KB generado correctamente, manteniendo el diseño profesional y optimizado para impresión

### [2024-12-19] - Ajuste de ancho de cuadros en propuesta medidor
- **Objetivo**: Modificar el layout de los cuadros en las secciones "Alcance de Servicios" y "Entregables" para que ocupen todo el ancho disponible
- **Cambios realizados**:
  - Modificación del CSS `.service-grid` para usar una sola columna (`grid-template-columns: 1fr`)
  - Actualización de las reglas de impresión para mantener consistencia en el PDF
  - Los cuadros ahora ocupan todo el ancho de la página en lugar de estar en dos columnas
  - **Archivos modificados**:
    - `propuesta_medidor_template.html` - Ajuste de layout de cuadros

### [2024-12-19] - Creación de changelog
- **Objetivo**: Establecer documentación de cambios del proyecto
- **Cambios realizados**:
  - Creación del archivo changelog.md
  - Documentación del objetivo general del proyecto
  - Establecimiento de estructura para seguimiento de cambios futuros

### [2024-12-19] - Creación de microservicio FastAPI
- **Objetivo**: Crear microservicio completo para generar propuestas con OpenAI y gestión en Cloudflare R2
- **Cambios realizados**:
  - Reemplazo completo de main.py con microservicio FastAPI
  - Integración de OpenAI para generación de contenido
  - Configuración de Cloudflare R2 para almacenamiento de archivos
  - Implementación de endpoints:
    - `POST /generate` - Generar nueva propuesta desde prompt
    - `GET /proposals` - Listar todas las propuestas
    - `GET /proposal/{id}/{type}` - Descargar propuesta (HTML o PDF)
    - `PUT /proposal/{id}` - Modificar propuesta existente
    - `GET /health` - Health check del servicio
  - Integración con stack definido: Rollbar para logs, manejo de errores
  - Función de conversión HTML a PDF adaptada de la versión original
  - Estructura de almacenamiento: `propuestas/html/` y `propuestas/pdf/`
  - Generación de IDs únicos para cada propuesta
  - Renderizado con template.html existente
  - **Dependencias agregadas**:
    - FastAPI, uvicorn, pydantic, openai, boto3, python-multipart, rollbar
  - **Variables de entorno requeridas**:
    - OPENAI_API_KEY, CLOUDFLARE_R2_ACCESS_KEY_ID, CLOUDFLARE_R2_SECRET_ACCESS_KEY, CLOUDFLARE_R2_ENDPOINT_URL
  - **Resultado**: Microservicio completo listo para deployment con uvicorn

### [2024-12-19] - Implementación de Base de Datos y Generación HTML Completa
- **Objetivo**: Mejorar el sistema con base de datos SQLite y generación HTML completa usando propuestahtml.yaml
- **Cambios realizados**:
  - **Base de datos SQLite**: Implementación completa para guardar metadatos de propuestas
    - Tabla `proposals` con campos: id, title, subtitle, prompt, html_url, pdf_url, created_at, updated_at, metadata
    - Funciones de CRUD: `save_proposal_to_db()`, `get_proposals_from_db()`, `get_proposal_from_db()`
  - **Generación HTML completa**: Nueva función `generate_html_with_openai()`
    - Usa `propuestahtml.yaml` como instrucciones para OpenAI
    - Genera documento HTML completo (no solo contenido para template)
    - Integración con `propuestahtml.yaml` para estilos y estructura avanzada
  - **Endpoints actualizados**:
    - `POST /generate` - Ahora usa generación HTML completa y guarda en BD
    - `GET /proposals` - Lista desde base de datos en lugar de R2
    - `PUT /proposal/{id}` - Modificación usando nueva función de generación
  - **Metadatos enriquecidos**: 
    - Tamaños de archivos, método de generación, historial de modificaciones
    - Tracking de prompts originales vs modificados
  - **Dependencias agregadas**:
    - `sqlalchemy>=2.0.0`, `alembic>=1.13.0` para gestión de BD
  - **Flujo mejorado**:
    1. OpenAI genera HTML completo usando propuestahtml.yaml
    2. HTML se convierte a PDF con WeasyPrint
    3. Archivos se suben a R2
    4. Metadatos se guardan en SQLite
    5. Listados se obtienen desde BD
  - **Resultado**: Sistema más robusto con persistencia de datos y generación HTML profesional

### [2024-12-19] - Migración a Neon PostgreSQL con SQLModel y Alembic
- **Objetivo**: Migrar de SQLite a Neon PostgreSQL siguiendo las reglas del stack ORGM
- **Cambios realizados**:
  - **Neon PostgreSQL**: Configuración completa con SQLModel
    - Reemplazo de SQLite por PostgreSQL serverless
    - Configuración de conexión con `DATABASE_URL`
    - Soporte para JSONB para metadatos complejos
  - **SQLModel**: Implementación de modelos modernos
    - `Proposal`, `ProposalCreate`, `ProposalRead`, `ProposalUpdate`
    - Integración nativa con FastAPI
    - Validación automática con Pydantic
  - **Alembic**: Sistema de migraciones profesional
    - Configuración automática con `alembic.ini`
    - Migraciones automáticas desde modelos SQLModel
    - Soporte para evolución de esquema
  - **Funciones de base de datos actualizadas**:
    - `save_proposal_to_db()` - Usa SQLModel
    - `get_proposals_from_db()` - Retorna objetos SQLModel
    - `get_proposal_from_db()` - Búsqueda por ID
    - `update_proposal_in_db()` - Actualización con SQLModel
  - **Dependencias actualizadas**:
    - `sqlmodel>=0.0.14` - ORM moderno
    - `psycopg2-binary>=2.9.0` - Driver PostgreSQL
    - `alembic>=1.13.0` - Migraciones
  - **Variables de entorno**:
    - `DATABASE_URL` - Conexión a Neon PostgreSQL
  - **Archivos creados**:
    - `models.py` - Modelos SQLModel
    - `database.py` - Configuración y utilidades de BD
    - `alembic/` - Directorio de migraciones
    - `alembic.ini` - Configuración de Alembic
    - `test_neon_integration.py` - Script de prueba
  - **Resultado**: Sistema profesional con base de datos PostgreSQL serverless y migraciones automáticas

### [2024-12-19] - Implementación de Flujo de 2 Prompts Secuenciales
- **Objetivo**: Implementar flujo de 2 prompts para generación de propuestas más estructurada
- **Cambios realizados**:
  - **Prompt 1 - Contenido**: Usa `propuesta.yaml` para generar contenido estructurado
    - Genera alcance completo, condiciones económicas, tiempos de entrega
    - Estructura profesional con secciones numeradas
    - Incluye información de empresa consultora y logo
  - **Prompt 2 - HTML**: Usa `propuestahtml.yaml` + contenido del Prompt 1
    - Genera HTML completo con estilos embebidos
    - Aplica paleta de colores y tipografía profesional
    - Incluye efectos visuales modernos y botón de exportación PDF
  - **Flujo secuencial**:
    - Contenido → HTML → PDF → R2 → Base de datos
    - Aplicado tanto en generación como en modificación
    - Logging detallado de cada paso
  - **Archivos actualizados**:
    - `main.py` - Endpoints `/generate` y `/proposal/{id}` con flujo de 2 prompts
    - `test_two_prompts.py` - Script de prueba del nuevo flujo
  - **Beneficios**:
    - Contenido más estructurado y profesional
    - Separación clara entre lógica de negocio y presentación
    - Mejor control sobre la calidad del output
    - Reutilización de instrucciones YAML
  - **Resultado**: Sistema de generación de propuestas más robusto y estructurado

### [2024-12-19] - Implementación de Celery + Redis para Procesamiento Asíncrono
- **Objetivo**: Implementar procesamiento asíncrono con Celery y Redis para workflows de generación de propuestas
- **Cambios realizados**:
  - **Celery + Redis**: Configuración completa para procesamiento asíncrono
    - Redis como broker y backend para Celery
    - Configuración mediante variable de entorno `REDIS_URL`
    - Integración con el cliente Redis existente en `cache.py`
  - **Tareas de Celery**:
    - `generate_text_task`: Genera contenido de texto usando `propuesta.yaml`
    - `generate_html_task`: Genera HTML usando `propuestahtml.yaml` + contenido
    - `generate_pdf_task`: Convierte HTML a PDF con WeasyPrint
    - Workflow en cadena: texto → HTML → PDF
  - **Nuevos endpoints**:
    - `POST /process`: Inicia workflow completo y retorna `task_id`
    - `GET /status/{task_id}`: Consulta estado de tarea (PENDING, STARTED, SUCCESS, FAILURE)
    - `GET /result/{task_id}`: Obtiene resultado final cuando la tarea completa
  - **Flujo asíncrono**:
    1. Cliente envía request a `/process`
    2. API retorna `task_id` inmediatamente
    3. Cliente consulta estado con `/status/{task_id}`
    4. Cuando completa, cliente obtiene resultado con `/result/{task_id}`
  - **Archivos creados/modificados**:
    - `tasks.py` - Tareas de Celery con workflow completo
    - `main.py` - Nuevos endpoints de Celery
    - `test.py` - Pruebas para workflow asíncrono
    - `pyproject.toml` - Dependencias de Celery
  - **Dependencias agregadas**:
    - `celery>=5.3.0` - Framework de tareas asíncronas
    - `celery[redis]>=5.3.0` - Soporte Redis para Celery
  - **Variables de entorno**:
    - `REDIS_URL` - URL de conexión a Redis (requerida)
  - **Beneficios**:
    - Procesamiento no bloqueante para requests largos
    - Escalabilidad horizontal con múltiples workers
    - Monitoreo de progreso en tiempo real
    - Mejor experiencia de usuario para operaciones pesadas
  - **Resultado**: Sistema híbrido con endpoints síncronos (existentes) y asíncronos (nuevos) para máxima flexibilidad

### [2024-12-19] - Refactorización: Centralización de Funciones Comunes
- **Objetivo**: Eliminar duplicación de código entre FastAPI y Celery, centralizando funciones comunes
- **Cambios realizados**:
  - **Archivo `proposal_services.py`**: Servicios centralizados compartidos
    - `generate_text_content()`: Generación de contenido de texto
    - `generate_html_content()`: Generación de contenido HTML
    - `generate_pdf_content()`: Generación de contenido PDF
    - Funciones auxiliares: `OpenAiApiResponse()`, `download_from_r2()`, `upload_to_r2()`, `html_to_pdf()`
  - **Refactorización de `tasks.py`**: Simplificación usando servicios centralizados
    - Tareas de Celery ahora son wrappers ligeros de las funciones centralizadas
    - Eliminación de código duplicado
    - Configuración SSL mejorada para Redis
  - **Refactorización de `main.py`**: Endpoints simplificados
    - Endpoints FastAPI ahora usan las funciones centralizadas
    - Eliminación de funciones duplicadas
    - Limpieza de importaciones no utilizadas
  - **Configuración SSL Redis**: Solución al error de SSL
    - Detección automática de URLs `rediss://`
    - Configuración SSL con `CERT_NONE` para entornos de desarrollo
    - Soporte para Redis con SSL en producción
  - **Archivos modificados**:
    - `proposal_services.py` - Nuevo archivo con servicios centralizados
    - `tasks.py` - Simplificado usando servicios centralizados
    - `main.py` - Refactorizado para usar servicios centralizados
    - `_celery_config.py` - Configuración SSL mejorada
  - **Beneficios**:
    - Eliminación de duplicación de código
    - Mantenimiento más fácil y consistente
    - Mejor separación de responsabilidades
    - Configuración SSL robusta para Redis
    - Código más limpio y mantenible
  - **Resultado**: Arquitectura más limpia con servicios centralizados y configuración SSL funcional

### [2024-12-19] - Corrección de Problemas de Celery y Logging
- **Objetivo**: Resolver problemas de worker de Celery y optimizar logging
- **Cambios realizados**:
  - **Problema identificado**: Worker de Celery no procesaba tareas (PENDING permanente)
    - Causa: Configuración incorrecta de registro de tareas
    - Solución: Uso directo de instancia `celery` en lugar de `celery_app`
  - **Optimización de logging**:
    - Configuración `--loglevel=error` para mostrar solo errores
    - Reducción de concurrencia a 1 worker
    - Deshabilitación de gossip, mingle y heartbeat
    - Formato de log simplificado
  - **Configuración mejorada de Celery**:
    - Task routing correcto para cola 'default'
    - Configuración de expiración de resultados (1 hora)
    - Resultados persistentes habilitados
  - **Corrección de endpoint modify proposal**:
    - Agregados campos requeridos `title` y `subtitle`
    - Actualizado modelo por defecto a `gpt-5-chat-latest`
  - **Scripts de prueba**:
    - `test_celery_simple.py`: Prueba tarea individual
    - `test_celery_workflow.py`: Prueba workflow completo
  - **Archivos modificados**:
    - `run_celery_worker.py` - Configuración optimizada de worker
    - `tasks.py` - Configuración mejorada de Celery
    - `test.py` - Corrección de endpoint modify proposal
  - **Resultado**: 
    - ✅ Worker de Celery funcionando correctamente
    - ✅ Workflow completo completándose en ~3 minutos
    - ✅ Logging reducido a solo errores
    - ✅ Endpoint modify proposal corregido
    - ✅ Generación de propuestas asíncrona operativa

### [2024-12-19] - Migración a FastAPI Background Tasks
- **Objetivo**: Simplificar el sistema eliminando Redis y Celery, usando Background Tasks nativos de FastAPI
- **Cambios realizados**:
  - **Eliminación de dependencias**:
    - Removidas dependencias de `celery`, `celery[redis]` y `redis`
    - Eliminada variable de entorno `REDIS_URL`
    - Simplificación del stack tecnológico
  - **Implementación de Background Tasks**:
    - Uso de `BackgroundTasks` nativo de FastAPI
    - Almacenamiento en memoria con `task_storage` dict
    - Progreso detallado: pending → started → in_progress → completed
    - Tracking de pasos: text_generation → html_generation → pdf_generation
  - **Nuevos endpoints optimizados**:
    - `POST /process`: Inicia workflow con Background Tasks
    - `GET /status/{task_id}`: Estado detallado con progreso y pasos
    - `GET /result/{task_id}`: Resultado final con URLs de archivos
  - **Función de background**:
    - `background_generate_proposal()`: Workflow completo asíncrono
    - Logging detallado de cada paso
    - Manejo de errores robusto
    - Actualización de progreso en tiempo real
  - **Archivos eliminados**:
    - `tasks.py` - Tareas de Celery
    - `run_celery_worker.py` - Script de worker
    - `_celery_config.py` - Configuración de Celery
    - `CELERY_USAGE.md` - Documentación de Celery
  - **Archivos modificados**:
    - `main.py` - Background Tasks implementados
    - `test.py` - Pruebas actualizadas para Background Tasks
    - `pyproject.toml` - Dependencias simplificadas
  - **Beneficios**:
    - ✅ **Simplicidad**: Sin dependencias externas (Redis)
    - ✅ **Mantenimiento**: Código más simple y directo
    - ✅ **Performance**: Sin overhead de broker/worker
    - ✅ **Escalabilidad**: Background Tasks escalan con FastAPI
    - ✅ **Monitoreo**: Progreso detallado en tiempo real
    - ✅ **Confiabilidad**: Manejo de errores integrado
  - **Resultado**: Sistema más simple, eficiente y fácil de mantener con Background Tasks nativos

## [2024-12-19] - Corrección de Conexión SSL y Pooling de Base de Datos

### Cambios
- **Corrección de conexión SSL**:
  - Configuración SSL mejorada con `sslmode=require`
  - Pool de conexiones con `pool_pre_ping=True` para verificar conexiones
  - Configuración de timeout y aplicación name
  - Pool size: 5 conexiones base, 10 adicionales cuando sea necesario
- **Reconexión automática**:
  - Función `execute_with_retry()` para reintentos automáticos
  - Detección de errores de conexión SSL
  - Backoff exponencial (1s, 2s, 4s) para reintentos
  - Máximo 3 intentos por operación
- **Corrección de modelo OpenAI**:
  - `gpt-5-chat-latest` no es un modelo válido de OpenAI
- **Funciones de base de datos actualizadas**:
  - Todas las funciones usan `execute_with_retry()`
  - Manejo robusto de errores de conexión
  - Logging mejorado para debugging
- **Archivos modificados**:
  - `_database_config.py` - Pooling y reconexión automática
  - `main.py` - Modelo corregido
  - `proposal_services.py` - Modelo corregido
  - `test.py` - Modelo corregido
  - `_models_config.py` - Modelo corregido

## [2024-12-19] - Endpoint de Búsqueda y Correcciones de Logger

### Cambios
- **Nuevo endpoint de búsqueda**:
  - `GET /proposals/search?q=...` - Búsqueda por tokens en título/subtítulo
  - Búsqueda case-insensitive, parcial y orden-agnóstica
  - Cada token debe existir en título O subtítulo
  - Sin duplicados, ordenado por fecha de creación desc
  - Función `search_proposals_by_title_subtitle()` en `_database_config.py`
- **Corrección de Logger**:
  - Arreglado `Logger.report_exc_info()` para aceptar `extra_data` directamente
  - Corregidas todas las llamadas en `main.py`, `_database_config.py`, `proposal_services.py`
  - Eliminado parámetro `extra_data=` de todas las llamadas
- **DeployLogger mejorado**:
  - Retorna dict con `status_code` y `body`
  - Integrado en arranque de `main.py` con try/except
  - Usa variable `APP_REVISION` para versión
- **Corrección de modelo OpenAI**:
  - Actualizado en todos los archivos: `main.py`, `proposal_services.py`, `test.py`, `_models_config.py`
- **Test de búsqueda**:
  - Función `test_search_proposals()` en `test.py`
  - Prueba con query "inter 69" para verificar funcionalidad
- **Archivos modificados**:
  - `main.py` - Endpoint de búsqueda y DeployLogger
  - `_database_config.py` - Función de búsqueda y correcciones de logger
  - `_logs_config.py` - Corrección de Logger.report_exc_info()
  - `test.py` - Test de búsqueda
  - `_models_config.py` - Modelo por defecto actualizado
  - `proposal_services.py` - Modelo por defecto actualizado
  - `changelog.md` - Documentación de cambios
