from fastapi import FastAPI, Request
import uvicorn
from opentelemetry.trace.propagation.tracecontext import TraceContextTextMapPropagator
from opentelemetry.trace import get_current_span

app = FastAPI()

# 
@app.get("/health",)
async def read_root(
    request: Request,
):
    print(request.headers)
    return {"Hello": "World"}

@app.get("/info")
async def read_info(
    request: Request,
):
    print(request.headers)
    carrier ={'traceparent': request.headers['Traceparent']}
    ctx = TraceContextTextMapPropagator().extract(carrier=carrier)
    print(f"Received context: {ctx}")

    # Get the current span from the context
    span = get_current_span(ctx)
    span_context = span.get_span_context()

    trace_id = span_context.trace_id
    span_id = span_context.span_id
    
    print(f"Span Context: {span_context}")

    print(f"TraceID: {trace_id}, SpanID: {span_id}")

   
    return {"message": "Hello, this is an info endpoint"}

# ---------------for dockerfile--------------------

if __name__ == "__main__":

    uvicorn.run("app:app", host='0.0.0.0', port=8080, reload=True, log_level="info", workers=1)
