"""
Servicios centralizados para generación de propuestas
Simplificado para generar solo HTML directamente
"""

import os
import uuid
import boto3
from pathlib import Path
from typing import Dict, Any
from datetime import datetime
from openai import OpenAI
from dotenv import load_dotenv
import json

# Import database utilities
from _database_config import (
    save_proposal_to_db,
    get_proposal_from_db,
    update_proposal_in_db,
)
from _models_config import ProposalCreate, ProposalUpdate
from _logs_config import logger

load_dotenv(override=True)


# R2 Configuration
R2_ACCESS_KEY: str = os.getenv("CLOUDFLARE_R2_ACCESS_KEY_ID", "")
R2_SECRET_KEY: str = os.getenv("CLOUDFLARE_R2_SECRET_ACCESS_KEY", "")
R2_ENDPOINT_URL: str = os.getenv("CLOUDFLARE_R2_ENDPOINT_URL", "")
R2_BUCKET: str = os.getenv("CLOUDFLARE_R2_BUCKET", "propuestas-bucket")

# Initialize S3 client for R2
s3_client = boto3.client(
    "s3",
    endpoint_url=R2_ENDPOINT_URL,
    aws_access_key_id=R2_ACCESS_KEY,
    aws_secret_access_key=R2_SECRET_KEY,
    region_name="auto",
)


def OpenAiApiResponse(modelo: str, texto: str, instructions: str):
    """Generate text using OpenAI API"""
    api_key: str | None = os.getenv("OPENAI_API_KEY")
    if not api_key:
        logger.error("OPENAI_API_KEY no configurada")
        raise ValueError("OPENAI_API_KEY no está configurada en settings.")

    client = OpenAI(api_key=api_key)
    try:
        response = client.chat.completions.create(
            model=modelo,
            messages=[
                {"role": "system", "content": instructions},
                {"role": "user", "content": texto},
            ],
            temperature=0.2,
        )
    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "openai_api_error",
                "model": modelo,
                "prompt_length": len(texto),
                "error": str(e),
            }
        )
        raise
    return response.choices[0].message.content


def download_from_r2(key: str) -> bytes:
    """Download file from R2"""
    try:
        response = s3_client.get_object(Bucket=R2_BUCKET, Key=key)
        return response["Body"].read()
    except Exception as e:
        logger.report_exc_info(
            {"error_type": "r2_download_error", "key": key, "error": str(e)}
        )
        raise


def upload_to_r2(file_path: Path, key: str) -> str:
    """Upload file to Cloudflare R2 and return public URL"""
    try:
        with open(file_path, "rb") as file:
            s3_client.upload_fileobj(file, R2_BUCKET, key)

        # Return the R2 URL
        public_url = f"{R2_ENDPOINT_URL}/{R2_BUCKET}/{key}"
        return public_url
    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "r2_upload_error",
                "file_path": str(file_path),
                "key": key,
                "error": str(e),
            }
        )
        raise


def generate_proposal_html(params: Dict[str, Any]) -> Dict[str, Any]:
    """
    Generate HTML proposal directly from prompt - combines both YAML instructions.
    This is the main function that generates proposals in a single call.
    """
    try:
        proposal_id = str(uuid.uuid4())[:8]
        timestamp = datetime.now().isoformat()
        modelo = params.get("model", "gpt-5-chat-latest")

        # Download both YAML files
        prompt_instructions: str = download_from_r2(
            "templates/prompt/propuesta.yaml"
        ).decode("utf-8")
        
        html_instructions: str = download_from_r2(
            "templates/html/propuesta.yaml"
        ).decode("utf-8")

        # Download CSS template to include in instructions
        css_content: str = download_from_r2(
            "templates/css/template.css"
        ).decode("utf-8")

        # Merge both instructions: prompt YAML defines content structure,
        # html YAML defines HTML formatting
        # Replace placeholder in html_instructions with actual CSS
        html_instructions_with_css = html_instructions.replace(
            "[AQUÍ VA TODO EL CSS DE template.css]", css_content
        )

        # Combine both instructions
        combined_instructions = f"""
{prompt_instructions}

---

{html_instructions_with_css}

---

IMPORTANTE: 
- Genera el HTML completo directamente desde el prompt del usuario.
- Usa el contenido generado según las reglas de prompt/propuesta.yaml.
- Formatea ese contenido en HTML usando la estructura de html/propuesta.yaml.
- Incluye TODO el CSS dentro de <style> tag.
- El HTML debe ser completo y listo para usar.
"""

        # Create user prompt with title and subtitle
        user_prompt = f"""
Título: {params["title"]}
Subtítulo: {params["subtitle"]}

Prompt del usuario:
{params["prompt"]}
"""

        # Generate HTML directly
        html_content = OpenAiApiResponse(
            modelo=modelo, texto=user_prompt, instructions=combined_instructions
        )

        # Clean up HTML (remove markdown code blocks if present)
        html_content = html_content.strip()
        if html_content.startswith("```html"):
            html_content = html_content[7:]
        if html_content.startswith("```"):
            html_content = html_content[3:]
        if html_content.endswith("```"):
            html_content = html_content[:-3]
        html_content = html_content.strip()

        # Save HTML temporarily
        temp_dir = Path(__file__).parent / "temp"
        temp_dir.mkdir(exist_ok=True)

        html_file = temp_dir / f"{proposal_id}.html"
        with open(html_file, "w", encoding="utf-8") as f:
            f.write(html_content)

        # Upload HTML to R2
        html_key = f"propuestas/html/{proposal_id}.html"
        html_url = upload_to_r2(html_file, html_key)

        # Save to database
        metadata = {
            "file_size_html": html_file.stat().st_size,
            "generation_method": "openai_html_generation_direct",
        }

        proposal_data: ProposalCreate = ProposalCreate(
            title=params["title"],
            subtitle=params["subtitle"],
            prompt=params["prompt"],
            html_url=html_url,
            proposal_metadata=json.dumps(metadata),
        )

        save_proposal_to_db(proposal_data, proposal_id)

        # Clean up temporary file
        html_file.unlink()

        return {
            "proposal_id": proposal_id,
            "html_url": html_url,
            "created_at": timestamp,
            "status": "html_generated",
        }

    except Exception as e:
        logger.report_exc_info(
            {"error_type": "proposal_generation_error", "params": params, "error": str(e)}
        )
        raise
