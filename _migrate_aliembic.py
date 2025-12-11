#!/usr/bin/env python3
"""
Script para ejecutar migraciones de Alembic
"""

import os
import subprocess
from pathlib import Path

def run_command(command, description):
    """Ejecutar comando y mostrar resultado"""
    print(f"\n🔄 {description}...")
    print(f"Comando: {command}")
    print("-" * 50)
    
    try:
        result = subprocess.run(command, shell=True, check=True, capture_output=True, text=True)
        print("✅ Éxito!")
        if result.stdout:
            print("Salida:", result.stdout)
        return True
    except subprocess.CalledProcessError as e:
        print("❌ Error!")
        print("Código de salida:", e.returncode)
        if e.stdout:
            print("Salida:", e.stdout)
        if e.stderr:
            print("Error:", e.stderr)
        return False

def check_env():
    """Verificar configuración de entorno"""
    print("🔍 Verificando configuración...")
    
    # Verificar si existe .env
    env_file = Path(".env")
    if not env_file.exists():
        print("⚠️  Archivo .env no encontrado")
        print("   Crea un archivo .env con tu DATABASE_URL")
        print("   Puedes usar: cp env_example.txt .env")
        return False
    
    # Verificar DATABASE_URL
    from dotenv import load_dotenv
    load_dotenv()
    
    database_url = os.getenv("DATABASE_URL")
    if not database_url:
        print("⚠️  DATABASE_URL no configurada en .env")
        return False
    
    print(f"✅ DATABASE_URL configurada: {database_url[:50]}...")
    return True

def main():
    """Función principal"""
    print("🚀 Script de Migraciones Alembic")
    print("=" * 50)
    
    # Verificar configuración
    if not check_env():
        print("\n❌ Configuración incompleta. Por favor:")
        print("1. Crea un archivo .env con tu DATABASE_URL")
        print("2. Ejecuta este script nuevamente")
        return
    
    # Mostrar opciones
    print("\n📋 Opciones disponibles:")
    print("1. Ver estado actual de migraciones")
    print("2. Aplicar migraciones (upgrade)")
    print("3. Revertir última migración (downgrade)")
    print("4. Ver historial de migraciones")
    print("5. Generar nueva migración")
    
    choice = input("\nSelecciona una opción (1-5): ").strip()
    
    if choice == "1":
        run_command("uv run alembic current", "Ver estado actual")
    elif choice == "2":
        run_command("uv run alembic upgrade head", "Aplicar migraciones")
    elif choice == "3":
        run_command("uv run alembic downgrade -1", "Revertir última migración")
    elif choice == "4":
        run_command("uv run alembic history", "Ver historial")
    elif choice == "5":
        message = input("Mensaje para la migración: ").strip()
        if message:
            run_command(f"uv run alembic revision -m '{message}'", f"Generar migración: {message}")
        else:
            print("❌ Mensaje requerido")
    else:
        print("❌ Opción inválida")

if __name__ == "__main__":
    main()
