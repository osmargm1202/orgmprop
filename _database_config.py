"""
Database configuration and utilities for Neon PostgreSQL

FEATURES:
- Connection pooling with pre-ping verification
- TCP keepalives to maintain connection health
- Automatic retry with exponential backoff on connection failures
- Basic connectivity verification before data operations
"""

import os
import json
from datetime import datetime
from sqlmodel import SQLModel, create_engine, Session, select
from sqlalchemy import text
from _models_config import Proposal, ProposalCreate, ProposalUpdate
from dotenv import load_dotenv
from _logs_config import logger

# Load environment variables
load_dotenv(override=True)

# Database configuration
DATABASE_URL = os.getenv("DATABASE_URL")

if not DATABASE_URL:
    logger.report_message("DATABASE_URL no configurada", "error")
    raise ValueError("DATABASE_URL no está configurada en variables de entorno")

# Create engine with connection pooling configuration
engine = create_engine(
    DATABASE_URL,
    echo=False,
    pool_pre_ping=True,  # Verify connections before use
    pool_recycle=1800,   # Recycle connections every 30 minutes
    pool_size=5,         # Number of connections to maintain
    max_overflow=10,     # Additional connections when needed
    connect_args={
        "connect_timeout": 15, # Connection timeout
        "application_name": "propuestas-api",
        "keepalives": 1,       # Enable TCP keepalives
        "keepalives_idle": 30, # Send keepalive after 30 seconds of inactivity
        "keepalives_interval": 10, # Send keepalive every 10 seconds
        "keepalives_count": 3  # Number of keepalive packets to send before considering connection dead
    }
)

def create_db_and_tables():
    """Create database tables"""
    try:
        SQLModel.metadata.create_all(engine)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_creation_error"
        })
        raise

# Legacy function - kept for compatibility


def verify_connection():
    """Verify basic database connectivity"""
    try:
        with Session(engine) as session:
            # Test basic connectivity
            result = session.execute(text("SELECT 1")).fetchone()
            if result[0] != 1:
                raise ValueError("Database connectivity test failed")
        return True
    except Exception as e:
        logger.error(f"Database connectivity verification failed: {str(e)}")
        raise ValueError(f"Database connection verification failed: {str(e)}")

def get_session():
    """Get database session with automatic reconnection"""
    # First verify connection
    verify_connection()

    max_retries = 3
    retry_delay = 1

    for attempt in range(max_retries):
        try:
            with Session(engine) as session:
                # Test the connection
                session.execute(text("SELECT 1"))
                yield session
                return
        except Exception as e:
            if attempt < max_retries - 1:
                logger.warning(f"Database connection failed (attempt {attempt + 1}/{max_retries}): {str(e)}")
                import time
                time.sleep(retry_delay * (2 ** attempt))  # Exponential backoff
            else:
                logger.error(f"Database connection failed after {max_retries} attempts: {str(e)}")
                raise

def execute_with_retry(func, *args, **kwargs):
    """Execute database function with automatic retry on connection failure"""
    # Verify connection before executing
    verify_connection()

    max_retries = 3
    retry_delay = 1

    for attempt in range(max_retries):
        try:
            return func(*args, **kwargs)
        except Exception as e:
            if "SSL connection has been closed" in str(e) or "connection" in str(e).lower():
                if attempt < max_retries - 1:
                    logger.warning(f"Database connection failed (attempt {attempt + 1}/{max_retries}): {str(e)}")
                    import time
                    time.sleep(retry_delay * (2 ** attempt))  # Exponential backoff
                else:
                    logger.error(f"Database connection failed after {max_retries} attempts: {str(e)}")
                    raise
            else:
                raise

def save_proposal_to_db(proposal_data: ProposalCreate, proposal_id: str) -> Proposal:
    """Save proposal to database"""
    def _save():
        with Session(engine) as session:
            proposal = Proposal(
                id=proposal_id,
                title=proposal_data.title,
                subtitle=proposal_data.subtitle,
                prompt=proposal_data.prompt,
                md_url=proposal_data.md_url,
                html_url=proposal_data.html_url,
                pdf_url=proposal_data.pdf_url,
                proposal_metadata=json.dumps(proposal_data.proposal_metadata) if isinstance(proposal_data.proposal_metadata, dict) else (proposal_data.proposal_metadata or "{}")
            )
            session.add(proposal)
            session.commit()
            session.refresh(proposal)
            return proposal
    
    try:
        return execute_with_retry(_save)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_save_error",
            "proposal_id": proposal_id
        })
        raise

def get_proposals_from_db() -> list[Proposal]:
    """Get all proposals from database"""
    def _get_all():
        with Session(engine) as session:
            statement = select(Proposal).order_by(Proposal.created_at.desc())
            proposals = session.exec(statement).all()
            return proposals
    
    try:
        return execute_with_retry(_get_all)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_get_all_error"
        })
        raise

def get_proposal_from_db(proposal_id: str) -> Proposal | None:
    """Get specific proposal from database"""
    def _get_single():
        with Session(engine) as session:
            statement = select(Proposal).where(Proposal.id == proposal_id)
            proposal = session.exec(statement).first()
            return proposal
    
    try:
        return execute_with_retry(_get_single)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_get_single_error",
            "proposal_id": proposal_id
        })
        raise

def update_proposal_in_db(proposal_id: str, proposal_data: ProposalUpdate) -> Proposal | None:
    """Update proposal in database"""
    def _update():
        with Session(engine) as session:
            statement = select(Proposal).where(Proposal.id == proposal_id)
            proposal = session.exec(statement).first()
            
            if not proposal:
                return None
            
            # Update fields
            update_data = proposal_data.model_dump(exclude_unset=True)
            for field, value in update_data.items():
                if field == "proposal_metadata":
                    if isinstance(value, dict):
                        setattr(proposal, field, json.dumps(value))
                    else:
                        setattr(proposal, field, value or "{}")
                else:
                    setattr(proposal, field, value)
            
            proposal.updated_at = datetime.utcnow()
            session.add(proposal)
            session.commit()
            session.refresh(proposal)
            return proposal
    
    try:
        return execute_with_retry(_update)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_update_error",
            "proposal_id": proposal_id
        })
        raise

def search_proposals_by_title_subtitle(query: str, limit: int = 10) -> list[Proposal]:
    """Search proposals where all query tokens are present either in title or subtitle.

    Matching is case-insensitive and partial, order-agnostic. Tokens are split by whitespace.
    Devuelve las últimas 10 propuestas de la búsqueda.
    """
    def _search():
        tokens = [t.strip() for t in (query or "").split() if t.strip()]
        if not tokens:
            return []

        with Session(engine) as session:
            stmt = select(Proposal)

            # For each token require it to be found in title OR subtitle (AND across tokens)
            from sqlalchemy import or_, and_
            conditions = []
            for token in tokens:
                like = f"%{token}%"
                conditions.append(
                    or_(Proposal.title.ilike(like), Proposal.subtitle.ilike(like))
                )
            stmt = stmt.where(and_(*conditions)).order_by(Proposal.created_at.desc())

            results = session.exec(stmt).all()
            # Devuelve solo las últimas 10 propuestas (más recientes)
            return results[:10]
    
    try:
        return execute_with_retry(_search)
    except Exception:
        logger.report_exc_info({
            "error_type": "database_search_error",
            "query": query
        })
        raise
