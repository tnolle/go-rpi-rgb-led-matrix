# HTTP API

The Go server listens on port `8085`. Catalog endpoints return JSON arrays under `items`, while display commands return the renderer that was started.

## Display stream

```text
GET /display/stream
Upgrade: websocket
```

The WebSocket streams the pixels currently shown by the server. When a frame
is available, the first message is JSON metadata:

```json
{"type":"display","version":1,"width":128,"height":64,"pixel_format":"rgb24"}
```

Every following message is a binary, complete RGB24 frame. Pixels are in
row-major order, with three bytes per pixel:

```text
R G B R G B ...
```

The binary message length is always `width * height * 3`. A newly connected
client receives the latest rendered frame immediately. Slow clients may skip
intermediate frames, but always receive the newest available frame; streaming
never blocks the physical or terminal display renderer.

The metadata and dimensions remain valid for the lifetime of the connection.
Protocol version `1` does not support changing matrix geometry without
reconnecting.

## Catalogs

```text
GET /images
GET /gifs
GET /dashboards
GET /animations
```

```json
{"items":["plasma","ripple"]}
```

## Display commands

```text
PUT /image?name=autodarts
PUT /gif?name=celebration
PUT /gif-once?name=success
PUT /dashboard?name=clock
PUT /animation?name=plasma
```

A successful response means the content was validated and its renderer started:

```json
{"display":{"type":"animation","name":"plasma","temporary":false}}
```

GIF-once playback is temporary and restores the previous non-temporary renderer when it finishes.

## Errors

Errors have a stable code and a human-readable message:

```json
{"error":{"code":"animation_not_found","message":"animation \"unknown\" does not exist"}}
```

| Status | Meaning |
| --- | --- |
| `400 Bad Request` | The required `name` parameter is missing or blank. |
| `404 Not Found` | The requested content is not in the server catalog. |
| `408 Request Timeout` | The caller cancelled the request. |
| `422 Unprocessable Entity` | An image or GIF exists but cannot be decoded. |
| `500 Internal Server Error` | Catalog access or renderer startup failed. |
| `503 Service Unavailable` | A dashboard dependency is unavailable. |
| `504 Gateway Timeout` | Renderer preparation exceeded five seconds. |

Rejected commands do not replace the active renderer. Internal failure details are logged by the server and are not exposed in API responses.

## Health

```text
GET /healthz
```

Returns `200 OK` with the plain-text body `OK`.
