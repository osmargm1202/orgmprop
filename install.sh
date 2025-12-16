#!/bin/bash
# ORGMPROP Installer
# Usage: curl -fsSL https://custom.or-gm.com/orgmprop.sh | bash

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="osmargm1202/propuestas"
BINARY_NAME="orgmprop"
INSTALL_DIR="$HOME/.local/bin"
CONFIG_DIR="$HOME/.config/orgmprop"
DOTFILES_REPO="https://github.com/osmargm1202/dotfiles.git"

echo -e "${BLUE}"
echo "╔══════════════════════════════════════════════════════╗"
echo "║           ORGMPROP Installer                         ║"
echo "║       Generador de Propuestas CLI                    ║"
echo "╚══════════════════════════════════════════════════════╝"
echo -e "${NC}"

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)
    
    case "$ARCH" in
        x86_64)
            ARCH="amd64"
            ;;
        aarch64|arm64)
            ARCH="arm64"
            ;;
        *)
            echo -e "${RED}Arquitectura no soportada: $ARCH${NC}"
            exit 1
            ;;
    esac
    
    case "$OS" in
        linux)
            PLATFORM="linux-$ARCH"
            ;;
        darwin)
            PLATFORM="darwin-$ARCH"
            ;;
        *)
            echo -e "${RED}Sistema operativo no soportado: $OS${NC}"
            exit 1
            ;;
    esac
    
    echo -e "${GREEN}Plataforma detectada: $PLATFORM${NC}"
}

# Download binary
download_binary() {
    echo -e "${BLUE}Descargando $BINARY_NAME...${NC}"
    
    # Try to get the latest release URL
    BINARY_URL="https://github.com/$REPO/releases/latest/download/$BINARY_NAME-$PLATFORM"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download binary
    if command -v curl &> /dev/null; then
        curl -fsSL -o "$INSTALL_DIR/$BINARY_NAME" "$BINARY_URL"
    elif command -v wget &> /dev/null; then
        wget -q -O "$INSTALL_DIR/$BINARY_NAME" "$BINARY_URL"
    else
        echo -e "${RED}Error: curl o wget requerido${NC}"
        exit 1
    fi
    
    # Make executable
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
    
    echo -e "${GREEN}Binario instalado en: $INSTALL_DIR/$BINARY_NAME${NC}"
}

# Setup configuration
setup_config() {
    echo -e "${BLUE}Configurando...${NC}"
    
    # Create config directory
    mkdir -p "$CONFIG_DIR"
    
    # Try to download dotfiles for config
    if command -v git &> /dev/null; then
        TEMP_DIR=$(mktemp -d)
        if git clone --depth 1 "$DOTFILES_REPO" "$TEMP_DIR" 2>/dev/null; then
            if [ -d "$TEMP_DIR/.config/orgmprop" ]; then
                cp -n "$TEMP_DIR/.config/orgmprop"/* "$CONFIG_DIR/" 2>/dev/null || true
                echo -e "${GREEN}Configuración copiada desde dotfiles${NC}"
            fi
        fi
        rm -rf "$TEMP_DIR"
    fi
    
    echo -e "${GREEN}Directorio de configuración: $CONFIG_DIR${NC}"
}

# Update PATH
update_path() {
    # Check if path already includes install dir
    if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
        echo -e "${YELLOW}Añadiendo $INSTALL_DIR a PATH...${NC}"
        
        # Detect shell and update rc file
        SHELL_NAME=$(basename "$SHELL")
        case "$SHELL_NAME" in
            bash)
                RC_FILE="$HOME/.bashrc"
                ;;
            zsh)
                RC_FILE="$HOME/.zshrc"
                ;;
            fish)
                RC_FILE="$HOME/.config/fish/config.fish"
                FISH_EXPORT="set -gx PATH \$HOME/.local/bin \$PATH"
                ;;
            *)
                RC_FILE="$HOME/.profile"
                ;;
        esac
        
        if [ "$SHELL_NAME" = "fish" ]; then
            if ! grep -q ".local/bin" "$RC_FILE" 2>/dev/null; then
                echo "$FISH_EXPORT" >> "$RC_FILE"
            fi
        else
            if ! grep -q ".local/bin" "$RC_FILE" 2>/dev/null; then
                echo 'export PATH="$HOME/.local/bin:$PATH"' >> "$RC_FILE"
            fi
        fi
        
        echo -e "${GREEN}PATH actualizado en $RC_FILE${NC}"
        echo -e "${YELLOW}Ejecuta: source $RC_FILE${NC}"
    fi
}

# Main installation
main() {
    detect_platform
    download_binary
    setup_config
    update_path
    
    echo ""
    echo -e "${GREEN}╔══════════════════════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║         ¡Instalación completada!                     ║${NC}"
    echo -e "${GREEN}╚══════════════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${BLUE}Próximos pasos:${NC}"
    echo "  1. Recarga tu shell: source ~/.bashrc (o ~/.zshrc)"
    echo "  2. Configura tu API key: orgmprop config apikey"
    echo "  3. Configura la carpeta base: orgmprop config folder"
    echo "  4. ¡Listo! Ejecuta: orgmprop menu"
    echo ""
    echo -e "${BLUE}Más información: https://github.com/$REPO${NC}"
}

main

