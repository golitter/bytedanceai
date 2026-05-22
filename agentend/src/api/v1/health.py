from fastapi import APIRouter

from src.app.config import settings

router = APIRouter(tags=["health"])


@router.get("/health")
async def health_check() -> dict:
    return {"status": "ok", "version": settings.app.version}
