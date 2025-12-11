# ORGMPROP

> CLI para generar propuestas comerciales en HTML con IA (Anthropic Claude)

```
╔══════════════════════════════════════════════════════╗
║                    ORGMPROP CLI                      ║
║           Generador de Propuestas                    ║
╚══════════════════════════════════════════════════════╝
```

## Funcionalidades

- **📝 Generación de Propuestas** - Crea propuestas HTML profesionales con IA
- **📂 Gestión de Proyectos** - Estructura de carpetas para proyectos de ingeniería
- **🤖 Múltiples Modelos** - Soporte para varios modelos de Claude:
  - `claude-sonnet-4-5-20250929` (por defecto)
  - `claude-haiku-4-5-20251001`
  - `claude-3-5-sonnet-20241022`
  - `claude-3-opus-20240229`
- **🎨 Interfaz TUI** - Interfaz interactiva con colores y formularios
- **📊 Resumen de Propuestas** - Vista general de todas las propuestas generadas

## Instalación

### Instalación Rápida (Recomendada)

```bash
curl -fsSL https://custom.or-gm.com/orgmprop.sh | bash
```

### Compilar desde código fuente

**Requisitos:**
- Go 1.21 o superior

```bash
# Clonar repositorio
git clone https://github.com/osmargm1202/propuestas.git
cd propuestas

# Compilar e instalar
make install
```

## Uso

### Configuración Inicial

```bash
# Configurar tu API key de Anthropic
orgmprop config apikey

# Configurar carpeta base de proyectos
orgmprop config folder

# Seleccionar modelo (opcional)
orgmprop config model
```

### Comandos

| Comando | Descripción |
|---------|-------------|
| `orgmprop menu` | Menú principal interactivo |
| `orgmprop new` | Crear nueva propuesta |
| `orgmprop proyecto` | Crear estructura de carpetas de proyecto |
| `orgmprop list` | Listar proyectos existentes |
| `orgmprop resumen` | Ver resumen de todas las propuestas |
| `orgmprop config` | Menú de configuración |
| `orgmprop config apikey` | Configurar API key |
| `orgmprop config model` | Seleccionar modelo |
| `orgmprop config folder` | Configurar carpeta base |
| `orgmprop --debug [cmd]` | Ejecutar con logs de debug |

### Ejemplos

```bash
# Abrir menú principal
orgmprop menu

# Crear nueva propuesta directamente
orgmprop new

# Crear un nuevo proyecto (estructura de carpetas)
orgmprop proyecto

# Ver todas las propuestas generadas
orgmprop resumen

# Ejecutar con debug
orgmprop --debug menu
```

## Configuración

La configuración se guarda en `~/.config/orgmprop/`:

```yaml
# config.yaml
anthropic_api_key: "sk-ant-api..."
model: "claude-sonnet-4-5-20250929"
base_folder: "/home/user/proyectos"
```

### Archivos de configuración

- `config.yaml` - Configuración principal
- `template.css` - Estilos CSS para las propuestas
- `propuesta.yaml` - Prompt de generación de contenido
- `html_template.yaml` - Estructura HTML de la propuesta
- `logo.svg` / `logo.png` - Logo de la empresa

## Estructura de Proyectos

Al crear un proyecto con `orgmprop proyecto`, se genera la siguiente estructura:

```
[COT]-[NOMBRE_PROYECTO]/
├── Comunicacion/
├── Diseño/
├── Estudios/
├── Calculos/
├── Levantamientos/
├── Entregas/
├── Recibido/
└── Oferta/           <- Las propuestas se generan aquí
```

## Archivos Generados

Al crear una propuesta, se generan los siguientes archivos:

- `propuesta.json` - Datos de la propuesta (título, subtítulo, prompt)
- `propuesta.html` - HTML con CSS embebido, listo para imprimir
- `logo.svg` - Logo de la empresa

## Desarrollo

```bash
# Compilar
make build

# Instalar localmente
make install

# Ejecutar con debug
make debug

# Limpiar
make clean

# Compilar para todas las plataformas
make build-all
```

## Links

- **Repositorio:** [github.com/osmargm1202/propuestas](https://github.com/osmargm1202/propuestas)
- **Descarga Directa:** [custom.or-gm.com/orgmprop](https://custom.or-gm.com/orgmprop)
- **Anthropic API:** [console.anthropic.com](https://console.anthropic.com)

## Licencia

MIT License - Ver archivo LICENSE para más detalles.

---

**Desarrollado con Go y la API de Claude (Anthropic)**
