import os
import boto3
from pathlib import Path
from typing import List, Dict, Any
from datetime import datetime
from _logs_config import logger
from fastapi import FastAPI, HTTPException, Response
from openai import OpenAI
from dotenv import load_dotenv
import uvicorn
import json
from contextlib import asynccontextmanager
from proposal_services import (
    generate_proposal_html,
    download_from_r2,
    upload_to_r2,
)


# Import SQLModel and database utilities
from _models_config import (
    ProposalGenerationRequest,
    ProposalResponse,
    ProposalListItem,
    ProposalUpdate,
    TitleSubtitleUpdateRequest,
)
from _database_config import (
    create_db_and_tables,
    get_proposals_from_db,
    get_proposal_from_db,
    update_proposal_in_db,
    search_proposals_by_title_subtitle,
)


load_dotenv(override=True)


@asynccontextmanager
async def deploy_logger(app: FastAPI):
    try:
        from _logs_config import DeployLogger

        DeployLogger()
    except Exception as e:
        logger.error(f"Deploy log failed: {str(e)}")
    yield


# FastAPI app
app = FastAPI(
    title="Propuestas API",
    description="Microservicio para generar propuestas comerciales en HTML con OpenAI",
    version="2.0.0",
)


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


@app.post("/generate-proposal", response_model=ProposalResponse)
async def generate_proposal(request: ProposalGenerationRequest):
    """
    Generate a complete HTML proposal in a single call.
    This is the main endpoint for creating proposals.
    """
    try:
        params = {
            "title": request.title,
            "subtitle": request.subtitle,
            "prompt": request.prompt,
            "model": request.model or "gpt-5-chat-latest",
        }

        result = generate_proposal_html(params)

        return ProposalResponse(
            id=result["proposal_id"],
            html_url=result["html_url"],
            created_at=result["created_at"],
        )

    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "proposal_generation_error",
                "request": request.dict(),
                "error": str(e),
            }
        )
        raise HTTPException(
            status_code=500, detail=f"Error generating proposal: {str(e)}"
        )


@app.get("/proposals", response_model=List[ProposalListItem])
async def list_proposals():
    """List all proposals from database"""
    try:
        db_proposals = get_proposals_from_db()

        proposals = []
        for proposal in db_proposals:
            metadata = (
                json.loads(proposal.proposal_metadata)
                if proposal.proposal_metadata
                else {}
            )
            proposals.append(
                ProposalListItem(
                    id=proposal.id,
                    title=proposal.title,
                    subtitle=proposal.subtitle,
                    prompt=proposal.prompt,
                    html_url=f"/proposal/{proposal.id}/html"
                    if proposal.html_url
                    else None,
                    created_at=proposal.created_at.isoformat(),
                    size_html=metadata.get("file_size_html", 0),
                )
            )

        return proposals

    except Exception as e:
        logger.report_exc_info({"error_type": "list_proposals_error", "error": str(e)})
        raise HTTPException(
            status_code=500, detail=f"Error listing proposals: {str(e)}"
        )


@app.get("/proposals/search", response_model=List[ProposalListItem])
async def search_proposals(q: str, limit: int = 10):
    """Search proposals by title/subtitle tokens (case-insensitive, partial, order-agnostic).

    - Splits `q` on spaces into tokens.
    - Each token must be present in either title OR subtitle.
    - Results are unique and ordered by creation date desc.
    """
    try:
        results = search_proposals_by_title_subtitle(q, limit)
        proposals: list[ProposalListItem] = []
        for proposal in results:
            metadata = (
                json.loads(proposal.proposal_metadata)
                if proposal.proposal_metadata
                else {}
            )
            proposals.append(
                ProposalListItem(
                    id=proposal.id,
                    title=proposal.title,
                    subtitle=proposal.subtitle,
                    prompt=proposal.prompt,
                    html_url=f"/proposal/{proposal.id}/html"
                    if proposal.html_url
                    else None,
                    created_at=proposal.created_at.isoformat(),
                    size_html=metadata.get("file_size_html", 0),
                )
            )
        return proposals
    except Exception as e:
        logger.report_exc_info(
            {"error_type": "search_proposals_error", "query": q, "error": str(e)}
        )
        raise HTTPException(
            status_code=500, detail=f"Error searching proposals: {str(e)}"
        )


@app.get("/proposal/{proposal_id}/html")
async def download_proposal_html(proposal_id: str):
    """Download a proposal HTML file"""
    try:
        key: str = f"propuestas/html/{proposal_id}.html"
        file_content: bytes = download_from_r2(key)

        filename: str = f"{proposal_id}.html"

        return Response(
            content=file_content,
            media_type="text/html",
            headers={"Content-Disposition": f"attachment; filename={filename}"},
        )

    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "download_proposal_error",
                "proposal_id": proposal_id,
                "file_type": "html",
                "error": str(e),
            }
        )
        raise HTTPException(status_code=404, detail=f"Proposal not found: {str(e)}")


@app.put("/proposal/{proposal_id}", response_model=ProposalResponse)
async def modify_proposal(proposal_id: str, request: ProposalGenerationRequest):
    """Modify an existing proposal by regenerating HTML with new content"""
    try:
        # Check if proposal exists in database
        existing_proposal = get_proposal_from_db(proposal_id)
        if not existing_proposal:
            logger.error("Proposal not found")
            raise HTTPException(status_code=404, detail="Proposal not found")

        modelo: str = request.model if request.model else "gpt-5-chat-latest"

        # Download both YAML files
        prompt_instructions: str = download_from_r2(
            "templates/prompt/propuesta.yaml"
        ).decode("utf-8")
        
        html_instructions: str = download_from_r2(
            "templates/html/propuesta.yaml"
        ).decode("utf-8")

        # Download CSS template
        css_content: str = download_from_r2(
            "templates/css/template.css"
        ).decode("utf-8")

        # Replace CSS placeholder
        html_instructions_with_css = html_instructions.replace(
            "[AQUÍ VA TODO EL CSS DE template.css]", css_content
        )

        # Combine instructions
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

        # Get existing HTML to provide context
        html_key: str = f"propuestas/html/{proposal_id}.html"
        try:
            existing_html: str = download_from_r2(html_key).decode("utf-8")
            enhanced_prompt: str = f"""
Contenido actual de la propuesta:
{existing_html}

Modifica la propuesta según lo siguiente:
Título: {request.title}
Subtítulo: {request.subtitle}

Instrucciones de modificación:
{request.prompt}
"""
        except:
            # If HTML doesn't exist, generate from scratch
            enhanced_prompt = f"""
Título: {request.title}
Subtítulo: {request.subtitle}

Prompt del usuario:
{request.prompt}
"""

        # Generate new HTML
        html_content: str = OpenAiApiResponse(
            modelo=modelo, texto=enhanced_prompt, instructions=combined_instructions
        )

        # Clean up HTML
        html_content = html_content.strip()
        if html_content.startswith("```html"):
            html_content = html_content[7:]
        if html_content.startswith("```"):
            html_content = html_content[3:]
        if html_content.endswith("```"):
            html_content = html_content[:-3]
        html_content = html_content.strip()

        # Save HTML temporarily
        temp_dir: Path = Path(__file__).parent / "temp"
        temp_dir.mkdir(exist_ok=True)

        html_file: Path = temp_dir / f"{proposal_id}.html"
        with open(html_file, "w", encoding="utf-8") as f:
            f.write(html_content)

        # Upload HTML to R2
        html_key: str = f"propuestas/html/{proposal_id}.html"
        html_url: str = upload_to_r2(html_file, html_key)

        # Update database
        metadata = {
            "file_size_html": html_file.stat().st_size,
            "generation_method": "openai_html_modification",
            "modified": True,
            "original_prompt": existing_proposal.prompt,
        }

        update_data = ProposalUpdate(
            title=request.title,
            subtitle=request.subtitle,
            prompt=request.prompt,
            html_url=html_url,
            proposal_metadata=metadata,
        )

        update_proposal_in_db(proposal_id, update_data)

        # Cleanup
        html_file.unlink()

        timestamp = (
            existing_proposal.created_at.isoformat()
            if hasattr(existing_proposal, "created_at")
            else datetime.now().isoformat()
        )

        return ProposalResponse(
            id=proposal_id, html_url=html_url, created_at=timestamp
        )

    except HTTPException:
        raise
    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "modify_proposal_error",
                "proposal_id": proposal_id,
                "prompt": request.prompt,
                "error": str(e),
            }
        )
        raise HTTPException(
            status_code=500, detail=f"Error modifying proposal: {str(e)}"
        )


@app.post("/proposal/{proposal_id}/regenerate", response_model=ProposalResponse)
async def regenerate_proposal(proposal_id: str, request: ProposalGenerationRequest):
    """Regenerate a proposal from scratch using new title, subtitle, and prompt.

    - Keeps the same proposal_id
    - Replaces the HTML content in R2
    - Updates title/subtitle/prompt in DB
    """
    try:
        existing = get_proposal_from_db(proposal_id)
        if not existing:
            logger.error("Proposal not found")
            raise HTTPException(status_code=404, detail="Proposal not found")

        # Use the main generation function
        params = {
            "title": request.title,
            "subtitle": request.subtitle,
            "prompt": request.prompt,
            "model": request.model or "gpt-5-chat-latest",
        }

        # Generate new HTML (this will create a new proposal_id, so we need to handle it differently)
        # Instead, we'll generate directly and update the existing record
        modelo: str = request.model or "gpt-5-chat-latest"

        # Download both YAML files
        prompt_instructions: str = download_from_r2(
            "templates/prompt/propuesta.yaml"
        ).decode("utf-8")
        
        html_instructions: str = download_from_r2(
            "templates/html/propuesta.yaml"
        ).decode("utf-8")

        # Download CSS template
        css_content: str = download_from_r2(
            "templates/css/template.css"
        ).decode("utf-8")

        # Replace CSS placeholder
        html_instructions_with_css = html_instructions.replace(
            "[AQUÍ VA TODO EL CSS DE template.css]", css_content
        )

        # Combine instructions
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

        # Generate new HTML from provided prompt (fresh start)
        user_prompt = f"""
Título: {request.title}
Subtítulo: {request.subtitle}

Prompt del usuario:
{request.prompt}
"""

        html_content: str = OpenAiApiResponse(
            modelo=modelo, texto=user_prompt, instructions=combined_instructions
        )

        # Clean up HTML
        html_content = html_content.strip()
        if html_content.startswith("```html"):
            html_content = html_content[7:]
        if html_content.startswith("```"):
            html_content = html_content[3:]
        if html_content.endswith("```"):
            html_content = html_content[:-3]
        html_content = html_content.strip()

        # Save and upload HTML replacing the existing one
        temp_dir: Path = Path(__file__).parent / "temp"
        temp_dir.mkdir(exist_ok=True)
        html_file: Path = temp_dir / f"{proposal_id}.html"
        with open(html_file, "w", encoding="utf-8") as f:
            f.write(html_content)

        html_key: str = f"propuestas/html/{proposal_id}.html"
        html_url: str = upload_to_r2(html_file, html_key)

        # Prepare metadata and DB update
        metadata = {
            "file_size_html": html_file.stat().st_size,
            "generation_method": "openai_html_regeneration",
            "regenerated": True,
            "previous_prompt": existing.prompt,
        }

        update_data = ProposalUpdate(
            title=request.title,
            subtitle=request.subtitle,
            prompt=request.prompt,
            html_url=html_url,
            proposal_metadata=metadata,
        )

        update_proposal_in_db(proposal_id, update_data)

        # Cleanup
        html_file.unlink(missing_ok=True)

        created_ts = (
            existing.created_at.isoformat()
            if hasattr(existing, "created_at")
            else datetime.now().isoformat()
        )

        return ProposalResponse(id=proposal_id, html_url=html_url, created_at=created_ts)
    except HTTPException:
        raise
    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "regenerate_proposal_error",
                "proposal_id": proposal_id,
                "error": str(e),
            }
        )
        raise HTTPException(
            status_code=500, detail=f"Error regenerating proposal: {str(e)}"
        )


@app.get("/")
async def health_check():
    """Health check endpoint"""
    return {"status": "healthy", "service": "propuestas-api", "version": "2.0.0"}


# Partial update: update only title and subtitle
@app.patch("/proposal/{proposal_id}/title-subtitle", response_model=ProposalResponse)
async def update_title_subtitle(proposal_id: str, request: TitleSubtitleUpdateRequest):
    """Update only title and subtitle of a proposal without regenerating content."""
    try:
        existing = get_proposal_from_db(proposal_id)
        if not existing:
            logger.error("Proposal not found")
            raise HTTPException(status_code=404, detail="Proposal not found")

        update_data = ProposalUpdate(
            title=request.title,
            subtitle=request.subtitle,
        )
        updated = update_proposal_in_db(proposal_id, update_data)
        timestamp = (
            updated.created_at.isoformat()
            if hasattr(updated, "created_at")
            else datetime.now().isoformat()
        )
        return ProposalResponse(
            id=proposal_id,
            html_url=updated.html_url,
            created_at=timestamp,
        )
    except HTTPException:
        raise
    except Exception as e:
        logger.report_exc_info(
            {
                "error_type": "update_title_subtitle_error",
                "proposal_id": proposal_id,
                "error": str(e),
            }
        )
        raise HTTPException(
            status_code=500, detail=f"Error updating title/subtitle: {str(e)}"
        )


# Error handlers
@app.exception_handler(Exception)
async def global_exception_handler(request, exc) -> HTTPException:
    logger.report_exc_info(
        {
            "error_type": "unhandled_exception",
            "request_url": str(request.url),
            "request_method": request.method,
        }
    )
    return HTTPException(status_code=500, detail="Internal server error")


if __name__ == "__main__":
    # Verify required environment variables
    required_vars = [
        "OPENAI_API_KEY",
        "ROLLBAR_ACCESS_TOKEN",
        "ROLLBAR_ENVIRONMENT",
        "CLOUDFLARE_R2_ACCESS_KEY_ID",
        "CLOUDFLARE_R2_SECRET_ACCESS_KEY",
        "CLOUDFLARE_R2_ENDPOINT_URL",
        "CLOUDFLARE_R2_BUCKET",
        "DATABASE_URL",
    ]

    missing_vars: list[str] = [var for var in required_vars if not os.getenv(var)]
    if missing_vars:
        logger.report_message(
            f"Missing environment variables: {missing_vars}",
            "error",
            extra_data={"missing_vars": missing_vars},
        )
        exit(1)

    # Send deploy log once
    try:
        from _logs_config import DeployLogger

        DeployLogger()
    except Exception as e:
        logger.error(f"Deploy log failed: {str(e)}")

    # Initialize database
    try:
        create_db_and_tables()
    except Exception as e:
        logger.error(f"Failed to initialize database: {e}")
        logger.report_exc_info({"error_type": "database_init_error"})
        exit(1)

    uvicorn.run("main:app", host="0.0.0.0", port=8000, reload=True)
