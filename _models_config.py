from datetime import datetime
from typing import Optional, Dict, Any
from sqlmodel import SQLModel, Field
from sqlalchemy import JSON
from pydantic import BaseModel


class ProposalBase(SQLModel):
    """Base model for proposals"""

    title: str
    subtitle: str
    prompt: str
    md_url: Optional[str] = None
    html_url: Optional[str] = None
    pdf_url: Optional[str] = None
    proposal_metadata: Optional[str] = Field(default="{}", sa_column=JSON)


class Proposal(ProposalBase, table=True):
    """Proposal model for database table"""

    __tablename__ = "proposals"

    id: str = Field(primary_key=True, max_length=8)
    created_at: datetime = Field(default_factory=datetime.utcnow)
    updated_at: datetime = Field(default_factory=datetime.utcnow)


class ProposalCreate(ProposalBase):
    """Model for creating proposals"""

    pass


class ProposalRead(ProposalBase):
    """Model for reading proposals"""

    id: str
    created_at: datetime
    updated_at: datetime


class ProposalUpdate(SQLModel):
    """Model for updating proposals"""

    title: Optional[str] = None
    subtitle: Optional[str] = None
    prompt: Optional[str] = None
    md_url: Optional[str] = None
    html_url: Optional[str] = None
    pdf_url: Optional[str] = None
    proposal_metadata: Optional[Dict[str, Any]] = None


class ProposalGenerationRequest(SQLModel):
    """Request model for generating a proposal (HTML only)"""

    title: str
    subtitle: str
    prompt: str
    model: str = "gpt-5-chat-latest"


class ProposalResponse(BaseModel):
    """Response model for proposal generation"""

    id: str
    html_url: Optional[str] = None
    created_at: str


class TitleSubtitleUpdateRequest(SQLModel):
    """Request model to update only title and subtitle"""

    title: str
    subtitle: str


class ProposalListItem(BaseModel):
    """Model for listing proposals"""

    id: str
    title: str
    subtitle: str
    prompt: str
    html_url: Optional[str] = None
    created_at: str
    size_html: int
