# ORGMAI

> CLI minimalista para interactuar con Claude (Anthropic) desde la terminal

```
╔════════════════════════════════════════╗
║           ORGMAI CLI                   ║
╚════════════════════════════════════════╝
```

## ✨ Funcionalidades

- **💬 Chat Interactivo** - Conversaciones fluidas con Claude directamente desde la terminal
- **⚡ Streaming** - Respuestas en tiempo real con streaming de texto
- **🔄 Historial de Conversaciones** - Guarda y recupera conversaciones anteriores
- **🤖 Múltiples Modelos** - Soporte para varios modelos de Claude:
  - `claude-haiku-4-5-20251001` (por defecto)
  - `claude-sonnet-4-5-20250929`
  - `claude-3-5-sonnet-20240620`
  - `claude-3-opus-20240229`
- **🎨 Interfaz TUI** - Interfaz de usuario en terminal con [huh](https://github.com/charmbracelet/huh)
- **📝 Preguntas Directas** - Ejecuta preguntas sin entrar en modo interactivo
- **🔧 Configuración Sencilla** - Archivo YAML en `~/.config/orgmai/`

## 📦 Instalación

### Instalación Rápida (Recomendada)

```bash
curl -fsSL https://raw.githubusercontent.com/osmargm1202/orgmai/main/orgmai.sh | bash
```

O con wget:

```bash
wget -qO- https://raw.githubusercontent.com/osmargm1202/orgmai/main/orgmai.sh | bash
```

### Descarga Directa del Binario

```bash
# Descargar binario
curl -L -o ~/.local/bin/orgmai https://custom.or-gm.com/orgmai

# Dar permisos de ejecución
chmod +x ~/.local/bin/orgmai

# Asegurar que ~/.local/bin está en PATH
export PATH="$HOME/.local/bin:$PATH"
```

### Instalación Manual (Compilar desde código fuente)

**Requisitos:**
- Go 1.21 o superior

```bash
# Clonar repositorio
git clone https://github.com/osmargm1202/orgmai.git
cd orgmai

# Compilar
go build -o orgmai ./cmd/orgmai

# Mover binario a PATH
mv orgmai ~/.local/bin/

# Opcional: asegurar que ~/.local/bin está en PATH
echo 'export PATH="$HOME/.local/bin:$PATH"' >> ~/.bashrc
source ~/.bashrc
```

## 🚀 Uso

### Configuración Inicial

```bash
# Configurar tu API key de Claude
orgmai apikey

# Seleccionar modelo (opcional)
orgmai config
```

### Comandos

| Comando | Descripción |
|---------|-------------|
| `orgmai apikey` | Configurar API key de Claude |
| `orgmai config` | Seleccionar modelo de Claude |
| `orgmai chat` | Iniciar chat interactivo |
| `orgmai chat [pregunta]` | Hacer una pregunta directa |
| `orgmai prev` | Continuar conversación anterior |
| `orgmai --debug` | Ejecutar con logs de debug |

### Ejemplos

```bash
# Pregunta directa
orgmai chat "¿Qué es Kubernetes?"

# Modo interactivo
orgmai chat

# Continuar conversación anterior
orgmai prev

# Con debug activado
orgmai chat --debug "Explica Docker"
```

## ⚙️ Configuración

La configuración se guarda en `~/.config/orgmai/config.yaml`:

```yaml
claude_api_key: "sk-ant-api..."
model: "claude-haiku-4-5-20251001"
```

Las conversaciones se almacenan en `~/.config/orgmai/conversations/` en formato Markdown.

## 🔗 Links

- **Repositorio:** [github.com/osmargm1202/orgmai](https://github.com/osmargm1202/orgmai)
- **Descarga Directa:** [custom.or-gm.com/orgmai](https://custom.or-gm.com/orgmai)
- **Anthropic API:** [console.anthropic.com](https://console.anthropic.com)


## ⚙️ Lanzar en desarrollo

Antes de lanzar

```
export GH_TOKEN=$(gh auth token) | source .env
```

## 📄 Licencia

Este proyecto está bajo la licencia **MIT**.

```
MIT License

Copyright (c) 2024 osmargm1202

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```

---

**Desarrollado con ❤️ usando Go y la API de Claude**
