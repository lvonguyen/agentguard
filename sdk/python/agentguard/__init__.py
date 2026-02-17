"""
AgentGuard Python SDK

Provides middleware integration for LangChain, CrewAI, and other agent frameworks
to enable runtime security governance, policy enforcement, and observability.

Usage with LangChain:
    from agentguard import AgentGuard, LangChainMiddleware
    
    guard = AgentGuard(api_key="your-api-key")
    agent = guard.wrap(your_langchain_agent)
    result = await agent.invoke({"input": "user query"})
"""

import asyncio
import hashlib
import json
import time
import uuid
from abc import ABC, abstractmethod
from contextlib import asynccontextmanager
from dataclasses import dataclass, field
from datetime import datetime
from enum import Enum
from typing import Any, Callable, Dict, List, Optional, TypeVar

import httpx

__version__ = "0.1.0"

# Type variables
T = TypeVar("T")


class DecisionType(Enum):
    """Policy decision types."""
    ALLOW = "allow"
    DENY = "deny"
    WARN = "warn"
    REQUIRE_APPROVAL = "require_approval"


class SignalType(Enum):
    """Security signal types."""
    INJECTION_ATTEMPT = "injection_attempt"
    DATA_EXFILTRATION = "data_exfiltration"
    TOOL_ABUSE = "tool_abuse"
    PRIVILEGE_ESCALATION = "privilege_escalation"
    ANOMALOUS_BEHAVIOR = "anomalous_behavior"
    POLICY_VIOLATION = "policy_violation"
    RATE_LIMIT_EXCEEDED = "rate_limit_exceeded"


@dataclass
class PolicyDecision:
    """Result of a policy evaluation."""
    allow: bool
    decision: DecisionType
    reasons: List[str] = field(default_factory=list)
    violations: List[Dict[str, Any]] = field(default_factory=list)
    eval_time_us: int = 0
    metadata: Dict[str, Any] = field(default_factory=dict)


@dataclass
class SecuritySignal:
    """A security-relevant event detected during execution."""
    id: str
    type: SignalType
    severity: str  # low, medium, high, critical
    title: str
    description: str
    evidence: Dict[str, Any] = field(default_factory=dict)
    timestamp: datetime = field(default_factory=datetime.utcnow)
    mitigated: bool = False


@dataclass
class Span:
    """A single operation within a trace."""
    span_id: str
    name: str
    type: str  # llm, retrieval, tool, chain, agent, policy
    start_time: datetime
    end_time: Optional[datetime] = None
    parent_span_id: Optional[str] = None
    status: str = "running"
    attributes: Dict[str, Any] = field(default_factory=dict)
    events: List[Dict[str, Any]] = field(default_factory=list)
    
    @property
    def duration_ms(self) -> int:
        if self.end_time is None:
            return 0
        return int((self.end_time - self.start_time).total_seconds() * 1000)


@dataclass
class Trace:
    """A complete execution trace for an agent invocation."""
    trace_id: str
    agent_id: str
    session_id: str
    user_id: Optional[str] = None
    start_time: datetime = field(default_factory=datetime.utcnow)
    end_time: Optional[datetime] = None
    status: str = "running"
    spans: List[Span] = field(default_factory=list)
    security_signals: List[SecuritySignal] = field(default_factory=list)
    metadata: Dict[str, Any] = field(default_factory=dict)


class AgentGuardClient:
    """HTTP client for AgentGuard API."""
    
    def __init__(
        self,
        api_key: str,
        base_url: str = "http://localhost:8080",
        timeout: float = 30.0,
    ):
        self.api_key = api_key
        self.base_url = base_url.rstrip("/")
        self.timeout = timeout
        self._client: Optional[httpx.AsyncClient] = None
    
    async def _get_client(self) -> httpx.AsyncClient:
        if self._client is None:
            self._client = httpx.AsyncClient(
                base_url=self.base_url,
                timeout=self.timeout,
                headers={
                    "Authorization": f"Bearer {self.api_key}",
                    "Content-Type": "application/json",
                    "X-AgentGuard-SDK-Version": __version__,
                },
            )
        return self._client
    
    async def close(self) -> None:
        if self._client is not None:
            await self._client.aclose()
            self._client = None
    
    async def evaluate_policy(
        self,
        agent_id: str,
        tool_name: Optional[str] = None,
        tool_params: Optional[Dict[str, Any]] = None,
        data_classification: Optional[str] = None,
    ) -> PolicyDecision:
        """Evaluate policies before an action."""
        client = await self._get_client()
        
        payload = {
            "agent": {"id": agent_id},
        }
        if tool_name:
            payload["tool"] = {
                "name": tool_name,
                "parameters": tool_params or {},
            }
        if data_classification:
            payload["data"] = {"classification": data_classification}
        
        response = await client.post("/api/v1/policies/evaluate", json=payload)
        response.raise_for_status()
        
        data = response.json()
        return PolicyDecision(
            allow=data.get("allow", True),
            decision=DecisionType(data.get("decision", "allow")),
            reasons=data.get("reasons", []),
            violations=data.get("violations", []),
            eval_time_us=data.get("eval_time_us", 0),
        )
    
    async def ingest_trace(self, trace: Trace) -> None:
        """Send a trace to AgentGuard for storage and analysis."""
        client = await self._get_client()
        
        payload = {
            "trace_id": trace.trace_id,
            "agent_id": trace.agent_id,
            "session_id": trace.session_id,
            "user_id": trace.user_id,
            "start_time": trace.start_time.isoformat(),
            "end_time": trace.end_time.isoformat() if trace.end_time else None,
            "status": trace.status,
            "spans": [self._span_to_dict(s) for s in trace.spans],
            "security_signals": [self._signal_to_dict(s) for s in trace.security_signals],
            "metadata": trace.metadata,
        }
        
        response = await client.post("/api/v1/observe/traces", json=payload)
        response.raise_for_status()
    
    async def pre_invoke(
        self,
        agent_id: str,
        tool_name: str,
        tool_params: Dict[str, Any],
        session_id: str,
    ) -> PolicyDecision:
        """Pre-invocation hook for policy evaluation."""
        client = await self._get_client()
        
        payload = {
            "agent_id": agent_id,
            "tool": {
                "name": tool_name,
                "parameters": tool_params,
            },
            "session_id": session_id,
            "timestamp": datetime.utcnow().isoformat(),
        }
        
        response = await client.post("/api/v1/sdk/pre-invoke", json=payload)
        response.raise_for_status()
        
        data = response.json()
        return PolicyDecision(
            allow=data.get("allow", True),
            decision=DecisionType(data.get("decision", "allow")),
            reasons=data.get("reasons", []),
        )
    
    async def post_invoke(
        self,
        agent_id: str,
        tool_name: str,
        result: Any,
        duration_ms: int,
        session_id: str,
    ) -> None:
        """Post-invocation hook for observability."""
        client = await self._get_client()
        
        payload = {
            "agent_id": agent_id,
            "tool": {"name": tool_name},
            "result_hash": self._hash_content(result),
            "duration_ms": duration_ms,
            "session_id": session_id,
            "timestamp": datetime.utcnow().isoformat(),
        }
        
        response = await client.post("/api/v1/sdk/post-invoke", json=payload)
        response.raise_for_status()
    
    def _span_to_dict(self, span: Span) -> Dict[str, Any]:
        return {
            "span_id": span.span_id,
            "parent_span_id": span.parent_span_id,
            "name": span.name,
            "type": span.type,
            "start_time": span.start_time.isoformat(),
            "end_time": span.end_time.isoformat() if span.end_time else None,
            "duration_ms": span.duration_ms,
            "status": span.status,
            "attributes": span.attributes,
            "events": span.events,
        }
    
    def _signal_to_dict(self, signal: SecuritySignal) -> Dict[str, Any]:
        return {
            "id": signal.id,
            "type": signal.type.value,
            "severity": signal.severity,
            "title": signal.title,
            "description": signal.description,
            "evidence": signal.evidence,
            "timestamp": signal.timestamp.isoformat(),
            "mitigated": signal.mitigated,
        }
    
    @staticmethod
    def _hash_content(content: Any) -> str:
        """Create a hash of content for logging without storing actual data."""
        content_str = json.dumps(content, sort_keys=True, default=str)
        return hashlib.sha256(content_str.encode()).hexdigest()[:16]


class AgentGuard:
    """Main AgentGuard SDK class for wrapping agents."""
    
    def __init__(
        self,
        api_key: str,
        agent_id: str,
        base_url: str = "http://localhost:8080",
        enabled: bool = True,
        fail_open: bool = False,
    ):
        """
        Initialize AgentGuard.
        
        Args:
            api_key: API key for authentication
            agent_id: Unique identifier for this agent
            base_url: AgentGuard API base URL
            enabled: Whether to enable policy enforcement
            fail_open: If True, allow actions when AgentGuard is unavailable
        """
        self.agent_id = agent_id
        self.enabled = enabled
        self.fail_open = fail_open
        self.client = AgentGuardClient(api_key, base_url)
        self._current_trace: Optional[Trace] = None
        self._current_span_stack: List[Span] = []
    
    async def close(self) -> None:
        """Close the client connection."""
        await self.client.close()
    
    @asynccontextmanager
    async def trace(
        self,
        session_id: str,
        user_id: Optional[str] = None,
        metadata: Optional[Dict[str, Any]] = None,
    ):
        """Context manager for tracing an agent execution."""
        trace = Trace(
            trace_id=str(uuid.uuid4()),
            agent_id=self.agent_id,
            session_id=session_id,
            user_id=user_id,
            metadata=metadata or {},
        )
        self._current_trace = trace
        
        try:
            yield trace
            trace.status = "completed"
        except Exception as e:
            trace.status = "failed"
            trace.metadata["error"] = str(e)
            raise
        finally:
            trace.end_time = datetime.utcnow()
            self._current_trace = None
            
            # Send trace to AgentGuard
            if self.enabled:
                try:
                    await self.client.ingest_trace(trace)
                except Exception:
                    # Log error but don't fail
                    pass
    
    @asynccontextmanager
    async def span(
        self,
        name: str,
        span_type: str,
        attributes: Optional[Dict[str, Any]] = None,
    ):
        """Context manager for tracing a span within an execution."""
        parent_span_id = None
        if self._current_span_stack:
            parent_span_id = self._current_span_stack[-1].span_id
        
        span = Span(
            span_id=str(uuid.uuid4()),
            name=name,
            type=span_type,
            start_time=datetime.utcnow(),
            parent_span_id=parent_span_id,
            attributes=attributes or {},
        )
        
        self._current_span_stack.append(span)
        if self._current_trace:
            self._current_trace.spans.append(span)
        
        try:
            yield span
            span.status = "completed"
        except Exception as e:
            span.status = "error"
            span.attributes["error"] = str(e)
            raise
        finally:
            span.end_time = datetime.utcnow()
            self._current_span_stack.pop()
    
    async def check_tool_access(
        self,
        tool_name: str,
        tool_params: Optional[Dict[str, Any]] = None,
        session_id: Optional[str] = None,
    ) -> PolicyDecision:
        """
        Check if a tool invocation is allowed by policy.
        
        Args:
            tool_name: Name of the tool to invoke
            tool_params: Parameters being passed to the tool
            session_id: Current session ID
        
        Returns:
            PolicyDecision with allow/deny and reasons
        """
        if not self.enabled:
            return PolicyDecision(allow=True, decision=DecisionType.ALLOW)
        
        try:
            return await self.client.pre_invoke(
                agent_id=self.agent_id,
                tool_name=tool_name,
                tool_params=tool_params or {},
                session_id=session_id or str(uuid.uuid4()),
            )
        except Exception as e:
            if self.fail_open:
                return PolicyDecision(
                    allow=True,
                    decision=DecisionType.WARN,
                    reasons=[f"AgentGuard unavailable: {e}. Proceeding with fail-open."],
                )
            else:
                return PolicyDecision(
                    allow=False,
                    decision=DecisionType.DENY,
                    reasons=[f"AgentGuard unavailable: {e}. Blocking due to fail-closed policy."],
                )
    
    def add_security_signal(
        self,
        signal_type: SignalType,
        severity: str,
        title: str,
        description: str,
        evidence: Optional[Dict[str, Any]] = None,
    ) -> None:
        """Add a security signal to the current trace."""
        if not self._current_trace:
            return
        
        signal = SecuritySignal(
            id=str(uuid.uuid4()),
            type=signal_type,
            severity=severity,
            title=title,
            description=description,
            evidence=evidence or {},
        )
        self._current_trace.security_signals.append(signal)


class LangChainMiddleware:
    """
    Middleware for LangChain agents.
    
    Usage:
        from langchain.agents import AgentExecutor
        from agentguard import AgentGuard, LangChainMiddleware
        
        guard = AgentGuard(api_key="...", agent_id="my-agent")
        middleware = LangChainMiddleware(guard)
        
        # Wrap the agent executor
        wrapped_agent = middleware.wrap(agent_executor)
        result = await wrapped_agent.ainvoke({"input": "query"})
    """
    
    def __init__(self, guard: AgentGuard):
        self.guard = guard
    
    def wrap(self, agent):
        """Wrap a LangChain AgentExecutor with security middleware."""
        # Import here to avoid hard dependency
        try:
            from langchain.callbacks.base import BaseCallbackHandler
        except ImportError:
            raise ImportError("langchain is required for LangChainMiddleware")
        
        class AgentGuardCallback(BaseCallbackHandler):
            def __init__(self, guard: AgentGuard):
                self.guard = guard
                self._tool_start_times: Dict[str, float] = {}
            
            async def on_tool_start(
                self,
                serialized: Dict[str, Any],
                input_str: str,
                **kwargs,
            ) -> None:
                tool_name = serialized.get("name", "unknown")
                
                # Check policy before tool execution
                decision = await self.guard.check_tool_access(
                    tool_name=tool_name,
                    tool_params={"input": input_str},
                )
                
                if not decision.allow:
                    raise PermissionError(
                        f"Tool '{tool_name}' blocked by policy: {', '.join(decision.reasons)}"
                    )
                
                self._tool_start_times[tool_name] = time.time()
            
            async def on_tool_end(
                self,
                output: str,
                **kwargs,
            ) -> None:
                # Record tool completion
                pass
            
            async def on_llm_start(
                self,
                serialized: Dict[str, Any],
                prompts: List[str],
                **kwargs,
            ) -> None:
                # Could check for prompt injection here
                pass
        
        # Add our callback to the agent
        callback = AgentGuardCallback(self.guard)
        if hasattr(agent, "callbacks") and agent.callbacks:
            agent.callbacks.append(callback)
        else:
            agent.callbacks = [callback]
        
        return agent


class CrewAIMiddleware:
    """
    Middleware for CrewAI agents.

    Usage:
        from crewai import Crew, Agent, Task
        from agentguard import AgentGuard, CrewAIMiddleware

        guard = AgentGuard(api_key="...", agent_id="my-crew")
        middleware = CrewAIMiddleware(guard)

        crew = Crew(agents=[...], tasks=[...])
        wrapped_crew = middleware.wrap(crew)
        result = wrapped_crew.kickoff()
    """

    def __init__(self, guard: AgentGuard):
        self.guard = guard

    def wrap(self, crew):
        """Wrap a CrewAI Crew with security middleware."""
        original_kickoff = crew.kickoff

        async def wrapped_kickoff(*args, **kwargs):
            session_id = str(uuid.uuid4())

            async with self.guard.trace(session_id=session_id):
                # Pre-execution policy check
                for agent in getattr(crew, 'agents', []):
                    agent_name = getattr(agent, 'role', 'unknown')
                    tools = getattr(agent, 'tools', [])

                    for tool in tools:
                        tool_name = getattr(tool, 'name', str(tool))
                        decision = await self.guard.check_tool_access(
                            tool_name=tool_name,
                            session_id=session_id,
                        )
                        if not decision.allow:
                            raise PermissionError(
                                f"Tool '{tool_name}' for agent '{agent_name}' blocked: "
                                f"{', '.join(decision.reasons)}"
                            )

                # Execute with tracing
                async with self.guard.span("crew_execution", "agent"):
                    result = await original_kickoff(*args, **kwargs)

                return result

        # Handle both sync and async
        if asyncio.iscoroutinefunction(original_kickoff):
            crew.kickoff = wrapped_kickoff
        else:
            def sync_wrapped(*args, **kwargs):
                return asyncio.run(wrapped_kickoff(*args, **kwargs))
            crew.kickoff = sync_wrapped

        return crew


class SecurityEnricher:
    """
    Security signal enrichment for AI execution traces.
    Detects security-relevant patterns in prompts, outputs, and tool usage.
    """

    # Injection patterns
    INJECTION_PATTERNS = [
        r"ignore\s+(previous|all|above)\s+(instructions?|prompts?)",
        r"disregard\s+(your|all|previous)\s+(rules?|instructions?)",
        r"you\s+are\s+now\s+(a|an|in)",
        r"pretend\s+(you're?|to\s+be)",
        r"roleplay\s+as",
        r"jailbreak",
        r"DAN\s+mode",
        r"\[system\]|\[SYSTEM\]",
        r"```\s*(system|assistant)",
    ]

    # PII patterns
    PII_PATTERNS = {
        "ssn": r"\b\d{3}-\d{2}-\d{4}\b",
        "credit_card": r"\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b",
        "email": r"\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Z|a-z]{2,}\b",
        "phone": r"\b\d{3}[-.]?\d{3}[-.]?\d{4}\b",
        "api_key": r"\b(sk-|api[_-]?key|bearer\s+)[a-zA-Z0-9]{20,}\b",
    }

    # Dangerous tool patterns
    DANGEROUS_TOOLS = [
        "execute_code", "run_shell", "file_write", "delete",
        "database_query", "http_request", "send_email",
    ]

    def __init__(self):
        import re
        self._injection_regex = [re.compile(p, re.IGNORECASE) for p in self.INJECTION_PATTERNS]
        self._pii_regex = {k: re.compile(v, re.IGNORECASE) for k, v in self.PII_PATTERNS.items()}

    def analyze_prompt(self, prompt: str) -> List[SecuritySignal]:
        """Analyze a prompt for security signals."""
        signals = []

        # Check for injection attempts
        for pattern in self._injection_regex:
            if pattern.search(prompt):
                signals.append(SecuritySignal(
                    id=str(uuid.uuid4()),
                    type=SignalType.INJECTION_ATTEMPT,
                    severity="high",
                    title="Potential Prompt Injection Detected",
                    description=f"Pattern matched: {pattern.pattern}",
                    evidence={"prompt_snippet": prompt[:200], "pattern": pattern.pattern},
                ))
                break

        # Check for PII in prompt
        for pii_type, pattern in self._pii_regex.items():
            if pattern.search(prompt):
                signals.append(SecuritySignal(
                    id=str(uuid.uuid4()),
                    type=SignalType.DATA_EXFILTRATION,
                    severity="medium",
                    title=f"Potential {pii_type.upper()} in Prompt",
                    description=f"Detected {pii_type} pattern in user input",
                    evidence={"pii_type": pii_type},
                ))

        return signals

    def analyze_output(self, output: str) -> List[SecuritySignal]:
        """Analyze an output for security signals."""
        signals = []

        # Check for PII exposure in output
        for pii_type, pattern in self._pii_regex.items():
            matches = pattern.findall(output)
            if matches:
                signals.append(SecuritySignal(
                    id=str(uuid.uuid4()),
                    type=SignalType.DATA_EXFILTRATION,
                    severity="high",
                    title=f"{pii_type.upper()} Exposure in Output",
                    description=f"Detected {len(matches)} instance(s) of {pii_type} in output",
                    evidence={"pii_type": pii_type, "count": len(matches)},
                ))

        return signals

    def analyze_tool_use(
        self,
        tool_name: str,
        tool_params: Dict[str, Any],
        output: Any,
    ) -> List[SecuritySignal]:
        """Analyze tool usage for security signals."""
        signals = []

        # Check for dangerous tool usage
        tool_lower = tool_name.lower()
        for dangerous in self.DANGEROUS_TOOLS:
            if dangerous in tool_lower:
                signals.append(SecuritySignal(
                    id=str(uuid.uuid4()),
                    type=SignalType.TOOL_ABUSE,
                    severity="medium",
                    title=f"Sensitive Tool Usage: {tool_name}",
                    description=f"Tool '{tool_name}' is classified as sensitive",
                    evidence={
                        "tool_name": tool_name,
                        "param_keys": list(tool_params.keys()),
                    },
                ))
                break

        # Check for privilege escalation patterns
        if "sudo" in str(tool_params).lower() or "admin" in str(tool_params).lower():
            signals.append(SecuritySignal(
                id=str(uuid.uuid4()),
                type=SignalType.PRIVILEGE_ESCALATION,
                severity="high",
                title="Potential Privilege Escalation",
                description="Tool parameters contain privilege escalation indicators",
                evidence={"tool_name": tool_name},
            ))

        return signals


class LangfuseExporter:
    """
    Export traces to Langfuse for observability.

    Usage:
        from agentguard import LangfuseExporter

        exporter = LangfuseExporter(
            public_key="pk-...",
            secret_key="sk-...",
            host="https://cloud.langfuse.com"
        )

        await exporter.export_trace(trace)
    """

    def __init__(
        self,
        public_key: str,
        secret_key: str,
        host: str = "https://cloud.langfuse.com",
    ):
        self.public_key = public_key
        self.secret_key = secret_key
        self.host = host.rstrip("/")
        self._client: Optional[httpx.AsyncClient] = None

    async def _get_client(self) -> httpx.AsyncClient:
        if self._client is None:
            import base64
            auth = base64.b64encode(
                f"{self.public_key}:{self.secret_key}".encode()
            ).decode()

            self._client = httpx.AsyncClient(
                base_url=self.host,
                headers={
                    "Authorization": f"Basic {auth}",
                    "Content-Type": "application/json",
                },
            )
        return self._client

    async def close(self) -> None:
        if self._client:
            await self._client.aclose()
            self._client = None

    async def export_trace(self, trace: Trace) -> None:
        """Export an AgentGuard trace to Langfuse."""
        client = await self._get_client()

        # Convert to Langfuse trace format
        langfuse_trace = {
            "id": trace.trace_id,
            "name": f"agent-{trace.agent_id}",
            "userId": trace.user_id,
            "sessionId": trace.session_id,
            "metadata": {
                **trace.metadata,
                "agentguard_version": __version__,
                "security_signals_count": len(trace.security_signals),
            },
            "input": trace.metadata.get("input"),
            "output": trace.metadata.get("output"),
        }

        response = await client.post("/api/public/traces", json=langfuse_trace)
        response.raise_for_status()

        # Export spans as generations/observations
        for span in trace.spans:
            await self._export_span(trace.trace_id, span)

        # Export security signals as events
        for signal in trace.security_signals:
            await self._export_security_signal(trace.trace_id, signal)

    async def _export_span(self, trace_id: str, span: Span) -> None:
        """Export a span to Langfuse."""
        client = await self._get_client()

        observation_type = "generation" if span.type == "llm" else "span"

        payload = {
            "id": span.span_id,
            "traceId": trace_id,
            "type": observation_type,
            "name": span.name,
            "parentObservationId": span.parent_span_id,
            "startTime": span.start_time.isoformat(),
            "endTime": span.end_time.isoformat() if span.end_time else None,
            "metadata": span.attributes,
            "level": "ERROR" if span.status == "error" else "DEFAULT",
        }

        # Add LLM-specific fields
        if span.type == "llm" and span.attributes.get("model"):
            payload["model"] = span.attributes.get("model")
            payload["promptTokens"] = span.attributes.get("prompt_tokens")
            payload["completionTokens"] = span.attributes.get("completion_tokens")

        response = await client.post("/api/public/observations", json=payload)
        response.raise_for_status()

    async def _export_security_signal(
        self,
        trace_id: str,
        signal: SecuritySignal,
    ) -> None:
        """Export a security signal as a Langfuse event."""
        client = await self._get_client()

        payload = {
            "id": signal.id,
            "traceId": trace_id,
            "type": "event",
            "name": f"security:{signal.type.value}",
            "metadata": {
                "severity": signal.severity,
                "title": signal.title,
                "description": signal.description,
                "evidence": signal.evidence,
                "mitigated": signal.mitigated,
            },
            "level": "ERROR" if signal.severity in ["critical", "high"] else "WARNING",
        }

        response = await client.post("/api/public/observations", json=payload)
        response.raise_for_status()


class OTELExporter:
    """
    Export traces to OpenTelemetry-compatible backends.
    Integrates with the AgentGuard Go server's OTEL collector.
    """

    def __init__(self, endpoint: str = "http://localhost:4318"):
        self.endpoint = endpoint.rstrip("/")
        self._client: Optional[httpx.AsyncClient] = None

    async def _get_client(self) -> httpx.AsyncClient:
        if self._client is None:
            self._client = httpx.AsyncClient(
                base_url=self.endpoint,
                headers={"Content-Type": "application/json"},
            )
        return self._client

    async def close(self) -> None:
        if self._client:
            await self._client.aclose()
            self._client = None

    async def export_trace(self, trace: Trace) -> None:
        """Export trace to OTEL collector via OTLP/HTTP."""
        client = await self._get_client()

        # Convert to OTLP format
        resource_spans = {
            "resourceSpans": [{
                "resource": {
                    "attributes": [
                        {"key": "service.name", "value": {"stringValue": "agentguard-sdk"}},
                        {"key": "agent.id", "value": {"stringValue": trace.agent_id}},
                    ]
                },
                "scopeSpans": [{
                    "scope": {"name": "agentguard.python.sdk"},
                    "spans": [self._convert_span(trace.trace_id, s) for s in trace.spans],
                }],
            }]
        }

        response = await client.post("/v1/traces", json=resource_spans)
        response.raise_for_status()

    def _convert_span(self, trace_id: str, span: Span) -> Dict[str, Any]:
        """Convert AgentGuard span to OTLP span format."""
        return {
            "traceId": trace_id.replace("-", "")[:32].ljust(32, "0"),
            "spanId": span.span_id.replace("-", "")[:16].ljust(16, "0"),
            "parentSpanId": span.parent_span_id.replace("-", "")[:16].ljust(16, "0") if span.parent_span_id else None,
            "name": span.name,
            "kind": self._span_kind(span.type),
            "startTimeUnixNano": int(span.start_time.timestamp() * 1e9),
            "endTimeUnixNano": int(span.end_time.timestamp() * 1e9) if span.end_time else None,
            "attributes": [
                {"key": k, "value": {"stringValue": str(v)}}
                for k, v in span.attributes.items()
            ],
            "status": {
                "code": 2 if span.status == "error" else 1,
                "message": span.attributes.get("error", ""),
            },
        }

    @staticmethod
    def _span_kind(span_type: str) -> int:
        """Map AgentGuard span type to OTLP span kind."""
        mapping = {
            "llm": 3,  # CLIENT
            "tool": 3,  # CLIENT
            "retrieval": 3,  # CLIENT
            "agent": 2,  # SERVER
            "chain": 0,  # INTERNAL
            "policy": 0,  # INTERNAL
        }
        return mapping.get(span_type, 0)


# Convenience exports
__all__ = [
    "AgentGuard",
    "AgentGuardClient",
    "LangChainMiddleware",
    "CrewAIMiddleware",
    "SecurityEnricher",
    "LangfuseExporter",
    "OTELExporter",
    "PolicyDecision",
    "DecisionType",
    "SecuritySignal",
    "SignalType",
    "Span",
    "Trace",
]
