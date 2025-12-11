# 📋 Guía de Migraciones con Alembic

## 🎯 Objetivo
Esta guía te explica cómo usar Alembic para gestionar las migraciones de la base de datos en el proyecto de propuestas.

## 🚀 Comandos Principales

### 1. **Configuración Inicial**

```bash
# 1. Crear archivo .env (si no existe)
cp env_example.txt .env

# 2. Editar .env con tu DATABASE_URL real
# DATABASE_URL=postgresql://username:password@ep-xxx-xxx.us-east-1.aws.neon.tech/neondb?sslmode=require
```

### 2. **Comandos Básicos**

```bash
# Ver estado actual de migraciones
uv run alembic current

# Ver historial de migraciones
uv run alembic history

# Aplicar todas las migraciones pendientes
uv run alembic upgrade head

# Aplicar migración específica
uv run alembic upgrade <revision_id>

# Revertir última migración
uv run alembic downgrade -1

# Revertir a migración específica
uv run alembic downgrade <revision_id>
```

### 3. **Generar Nuevas Migraciones**

```bash
# Generar migración automática (detecta cambios en modelos)
uv run alembic revision --autogenerate -m "Descripción del cambio"

# Generar migración vacía (para cambios manuales)
uv run alembic revision -m "Descripción del cambio"
```

## 📁 Estructura de Archivos

```
alembic/
├── env.py              # Configuración de Alembic
├── script.py.mako      # Template para nuevas migraciones
└── versions/           # Archivos de migración
    ├── a9bc3a9ecd96_add_md_url_field_to_proposal_model.py
    └── ...
```

## 🔧 Migración Actual: Campo md_url

### **Archivo generado**: `a9bc3a9ecd96_add_md_url_field_to_proposal_model.py`

**Cambios incluidos**:
- ✅ Agregar columna `md_url` (String, nullable) a la tabla `proposals`
- ✅ Función de rollback para eliminar la columna

### **Para aplicar esta migración**:

```bash
# 1. Verificar estado actual
uv run alembic current

# 2. Aplicar migración
uv run alembic upgrade head

# 3. Verificar que se aplicó correctamente
uv run alembic current
```

## 🛠️ Script de Ayuda

Usa el script `migrate.py` para una interfaz interactiva:

```bash
python migrate.py
```

Este script te permite:
- ✅ Verificar configuración
- ✅ Ver estado de migraciones
- ✅ Aplicar migraciones
- ✅ Revertir migraciones
- ✅ Generar nuevas migraciones

## ⚠️ Consideraciones Importantes

### **Antes de aplicar migraciones**:
1. **Backup**: Siempre haz backup de tu base de datos
2. **Testing**: Prueba en un entorno de desarrollo primero
3. **Revisión**: Revisa el código SQL generado antes de aplicar

### **En producción**:
1. **Mantenimiento**: Programa ventanas de mantenimiento
2. **Rollback**: Ten plan de rollback preparado
3. **Monitoreo**: Monitorea la aplicación después de aplicar migraciones

## 🔍 Verificación Post-Migración

Después de aplicar la migración del campo `md_url`:

```sql
-- Verificar que la columna existe
\d proposals

-- Verificar estructura de la tabla
SELECT column_name, data_type, is_nullable 
FROM information_schema.columns 
WHERE table_name = 'proposals' 
ORDER BY ordinal_position;
```

## 🚨 Solución de Problemas

### **Error de conexión**:
```
psycopg2.OperationalError: connection to server at "localhost" failed
```
**Solución**: Verifica que `DATABASE_URL` esté configurada correctamente en `.env`

### **Migración ya aplicada**:
```
alembic.util.exc.CommandError: Can't locate revision identified by 'xxx'
```
**Solución**: Usa `alembic current` para ver el estado actual

### **Conflicto de migraciones**:
```
alembic.util.exc.CommandError: Multiple heads detected
```
**Solución**: Usa `alembic merge` para combinar migraciones

## 📚 Recursos Adicionales

- [Documentación oficial de Alembic](https://alembic.sqlalchemy.org/)
- [SQLModel con Alembic](https://sqlmodel.tiangolo.com/tutorial/create-db-and-table/)
- [PostgreSQL con Python](https://www.postgresql.org/docs/current/libpq-connect.html)
