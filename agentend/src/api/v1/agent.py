import asyncio

from fastapi import APIRouter, Depends, HTTPException
from sse_starlette.sse import EventSourceResponse

from src.adapters.base import BaseAgentAdapter
from src.adapters.registry import AdapterRegistry
from src.api.dependencies import get_adapter_registry, get_rule_engine, get_session_manager
from src.app.config import settings
from src.rules.engine import RuleEngine
from src.schemas.request import AgentRequest
from src.schemas.response import AgentResponse
from src.session.manager import SessionManager
from src.session.models import SessionState

router = APIRouter(prefix="/v1/agent", tags=["agent"])


async def _get_or_create_session(request: AgentRequest, session_mgr: SessionManager) -> str:
    if request.session_id:
        session = session_mgr.get(request.session_id)
        if session:
            return session.id
    session = session_mgr.create(
        agent_type=request.agent_type,
        workspace_path=request.workspace_path,
    )
    return session.id


def _build_rule_context(request: AgentRequest) -> dict:
    return {
        "message": request.message,
        "agent_type": request.agent_type,
        "workspace_path": request.workspace_path,
        "allowed_tools": request.config.get("allowed_tools", []) if request.config else [],
    }


async def _execute_stream(
    request: AgentRequest,
    adapter: BaseAgentAdapter,
    session_id: str,
    rule_result: dict,
    session_mgr: SessionManager,
):
    session_mgr.update_state(session_id, SessionState.RUNNING)
    session_mgr.record_history(session_id, {"role": "user", "content": request.message})

    try:
        async for event in adapter.stream_chat(
            session_id,
            request.message,
            cli_session_id=request.session_id,
            system_prompt_append="\n".join(rule_result.get("system_prompt_append", [])) or None,
            allowed_tools=rule_result.get("allowed_tools") or None,
            max_turns=rule_result.get("max_turns"),
        ):
            yield {
                "event": event.type,
                "data": event.model_dump_json(),
            }
    finally:
        session_mgr.update_state(session_id, SessionState.COMPLETED)


@router.post("/stream")
async def agent_stream(
    request: AgentRequest,
    adapter_registry: AdapterRegistry = Depends(get_adapter_registry),
    rule_engine: RuleEngine = Depends(get_rule_engine),
    session_mgr: SessionManager = Depends(get_session_manager),
) -> EventSourceResponse:
    passed, rule_result = rule_engine.evaluate(_build_rule_context(request))
    if not passed:
        raise HTTPException(status_code=400, detail=rule_result)

    adapter_cls = adapter_registry.get(request.agent_type)
    adapter = adapter_cls()
    session_id = await _get_or_create_session(request, session_mgr)

    return EventSourceResponse(_execute_stream(request, adapter, session_id, rule_result, session_mgr))


@router.post("/execute", response_model=AgentResponse)
async def agent_execute(
    request: AgentRequest,
    adapter_registry: AdapterRegistry = Depends(get_adapter_registry),
    rule_engine: RuleEngine = Depends(get_rule_engine),
    session_mgr: SessionManager = Depends(get_session_manager),
) -> AgentResponse:
    passed, rule_result = rule_engine.evaluate(_build_rule_context(request))
    if not passed:
        raise HTTPException(status_code=400, detail=rule_result)

    adapter_cls = adapter_registry.get(request.agent_type)
    adapter = adapter_cls()
    session_id = await _get_or_create_session(request, session_mgr)
    session_mgr.update_state(session_id, SessionState.RUNNING)
    session_mgr.record_history(session_id, {"role": "user", "content": request.message})

    try:
        response = await asyncio.wait_for(
            adapter.chat(
                session_id,
                request.message,
                cli_session_id=request.session_id,
                system_prompt_append="\n".join(rule_result.get("system_prompt_append", [])) or None,
                allowed_tools=rule_result.get("allowed_tools") or None,
                max_turns=rule_result.get("max_turns"),
            ),
            timeout=settings.EXECUTION_TIMEOUT,
        )
        session_mgr.update_state(session_id, SessionState.COMPLETED)
        session_mgr.record_history(session_id, {"role": "assistant", "content": response.content})
        return response
    except asyncio.TimeoutError:
        session_mgr.update_state(session_id, SessionState.ERROR)
        await adapter.interrupt(session_id)
        raise HTTPException(status_code=408, detail="execution timeout")
