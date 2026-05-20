from fastapi import APIRouter, Depends, HTTPException

from src.api.dependencies import get_session_manager
from src.session.manager import SessionManager

router = APIRouter(prefix="/v1/session", tags=["session"])


@router.get("")
async def list_sessions(
    mgr: SessionManager = Depends(get_session_manager),
) -> list[dict]:
    sessions = mgr.list()
    return [
        {
            "session_id": s.id,
            "agent_type": s.agent_type,
            "state": s.state.value,
            "created_at": s.created_at.isoformat(),
            "last_active": s.last_active.isoformat(),
        }
        for s in sessions
    ]


@router.get("/{session_id}")
async def get_session(
    session_id: str,
    mgr: SessionManager = Depends(get_session_manager),
) -> dict:
    session = mgr.get(session_id)
    if not session:
        raise HTTPException(status_code=404, detail="Session not found")
    return {
        "session_id": session.id,
        "agent_type": session.agent_type,
        "state": session.state.value,
        "workspace_path": session.workspace_path,
        "created_at": session.created_at.isoformat(),
        "last_active": session.last_active.isoformat(),
        "history_count": len(session.history),
    }


@router.post("/{session_id}/interrupt")
async def interrupt_session(
    session_id: str,
    mgr: SessionManager = Depends(get_session_manager),
) -> dict:
    session = mgr.get(session_id)
    if not session:
        raise HTTPException(status_code=404, detail="Session not found")
    if session.state.value != "running":
        return {"message": "session not running"}
    from src.session.models import SessionState

    mgr.update_state(session_id, SessionState.INTERRUPTED)
    return {"message": "session interrupted"}


@router.delete("/{session_id}")
async def delete_session(
    session_id: str,
    mgr: SessionManager = Depends(get_session_manager),
) -> dict:
    destroyed = await mgr.destroy(session_id)
    if not destroyed:
        raise HTTPException(status_code=404, detail="Session not found")
    return {"message": "session destroyed"}
