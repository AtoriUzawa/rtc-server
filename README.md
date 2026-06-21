# RTC Server

A real-time communication server built on the [Cira](https://github.com/AtoriUzawa/cira) WebSocket framework, providing WebRTC signaling relay, P2P call management, live streaming rooms, and meeting rooms.
It works together with the RTC Android client: [rtc-android](https://github.com/AtoriUzawa/rtc-android)

## Introduction

RTC Server is a lightweight, single-process signaling server that enables WebRTC-capable clients to discover peers, exchange SDP offers/answers, relay ICE candidates, and manage room-based group communication scenarios. It does **not** process media — it only routes signaling messages between connected clients.

## Features

- **WebRTC Signaling Relay** — Offer, Answer, and ICE Candidate forwarding between peers
- **P2P Call Management** — Initiate and hang up calls with synchronous request-response and configurable timeout
- **Live Streaming Rooms** — Owner/viewer role model with heat-based ranking and cursor pagination
- **Meeting Rooms** — Host/member role model with automatic room teardown when the host leaves
- **Connection Lifecycle** — Automatic cleanup on disconnect via close callbacks
- **Concurrency Safe** — All managers use `sync.RWMutex` for concurrent access

## Architecture

```mermaid
graph TB
    subgraph Clients
        A[Client A]
        B[Client B]
        C[Client C]
    end

    subgraph "RTC Server"
        HTTP[HTTP Server<br/>gin :18888]
        WS[WebSocket Server<br/>cira :18887]

        subgraph Modules
            SIG[signal.Manager<br/>conns map]
            P2P[p2p.Manager<br/>peers map + skiplist]
            LIVE[live.Manager<br/>rooms map + skiplist]
            MEET[meeting.Manager<br/>rooms map]
        end

        HTTP -->|POST /api/*/list| P2P
        HTTP -->|POST /api/*/list| LIVE
        HTTP -->|POST /api/*/list| MEET

        WS -->|signal.rtc.*| SIG
        WS -->|p2p.*| P2P
        WS -->|live.*| LIVE
        WS -->|meeting.*| MEET

        P2P --> SIG
        LIVE --> SIG
        MEET --> SIG
    end

    A <-->|WebSocket| WS
    B <-->|WebSocket| WS
    C <-->|WebSocket| WS
```

The server exposes two ports: an HTTP API server (gin) for resource listing, and a WebSocket server (cira) for all real-time signaling. Four internal modules share a single `signal.Manager` that maps client IDs to their WebSocket connections.

## Project Structure

```
rtc-server/
├── cmd/main.go                       # Entry point
├── internal/
│   ├── app/app.go                    # Module assembly and server startup
│   ├── signal/                       # WebRTC signaling relay
│   │   ├── manager.go                # Connection registry (map[id]*cira.Conn)
│   │   ├── handler_ws.go             # WS: register, unregister, forward
│   │   ├── model.go                  # RegisterReq, UnRegisterReq
│   │   ├── protocol.go               # SignalMessage, Type (offer/answer/candidate)
│   │   ├── module.go                 # Module composition
│   │   └── router.go                 # Route: signal.rtc.*
│   ├── p2p/                          # P2P call management
│   │   ├── manager.go                # Peer registry (map + skiplist)
│   │   ├── service.go                # Call / Hangup logic
│   │   ├── handler_ws.go / handler_http.go
│   │   ├── model.go                  # Peer
│   │   ├── dto.go                    # CallReq, HangupReq, ListReq/Resp, PeerDTO
│   │   ├── module.go / router.go
│   ├── live/                         # Live streaming rooms
│   │   ├── manager.go                # Room registry (map + skiplist by heat)
│   │   ├── service.go                # Create / Join / Leave + broadcast
│   │   ├── model.go                  # Room, RoomMember, RoomItem
│   │   ├── dto.go                    # CreateReq, JoinReq, LeaveReq, RoomDTO
│   │   ├── handler_ws.go / handler_http.go
│   │   ├── module.go / router.go
│   ├── meeting/                      # Meeting rooms
│   │   ├── manager.go                # Room registry (map only)
│   │   ├── service.go                # Create / Join / Leave + broadcast
│   │   ├── model.go                  # Room, RoomMember (Host/Member roles)
│   │   ├── dto.go                    # CreateReq, JoinReq, LeaveReq, RoomDTO
│   │   ├── handler_ws.go / handler_http.go
│   │   ├── module.go / router.go
│   ├── router/router.go              # Aggregated route registration
│   └── transport/
│       ├── httpx/                    # HTTP JSON bind + response helpers
│       └── wsx/                      # WS JSON bind + response helpers
└── pkg/
    ├── skiplist/                     # Generic skip list
    ├── xerror/                       # Error types with code/message
    ├── heap/                         # Binary heap (unused)
    ├── xlog/                         # zap logger (unused)
    ├── jwt/                          # JWT utilities (unused)
    ├── redis/                        # Redis client (unused)
    └── idgen/                        # ID generator (unused)
```

## P2P Call

P2P calling uses cira's synchronous `Call()` primitive — the server sends a request to the callee and blocks until a response arrives or the timeout fires.

```mermaid
sequenceDiagram
    participant A as Client A (caller)
    participant S as RTC Server
    participant B as Client B (callee)

    Note over A,B: Peer registration (before call)
    A->>S: WS: p2p.register {id:"A"}
    B->>S: WS: p2p.register {id:"B"}

    Note over A,B: P2P Call
    A->>S: WS: p2p.call {from:"A", to:"B", call_id:"c1"}
    S->>S: lookup peer B in p2p.Manager
    S->>S: lookup conn(B) in signal.Manager
    S->>B: cira.Call("p2p.call", req, 30s timeout)
    B-->>S: response {status:"ok"}
    S-->>A: resp {status:"ok"}

    Note over A,B: Hangup
    A->>S: WS: p2p.hangup {from:"A", to:"B"}
    S->>B: push p2p.hangup {from:"A", to:"B"}
```

**Key implementation details:**
- `Call()` blocks with a 30-second timeout (`ctx.Timeout = 30 * time.Second`)
- Timeout returns `{"status": "timeout"}` to the caller
- Offline peer returns `{"status": "offline"}`
- Peer lookup: `p2p.Manager.Peer(id)` then `signal.Manager.Conn(id)`

**Peer list (HTTP):**
```
POST /api/p2p/list  {cursor, limit}
→ {list: [{id, nickname}], next_cursor}
```
Paginated via skip list, limit clamped to [1, 10].

## Live Streaming

Live rooms use a **skiplist ranked by heat** (`member count × 10`) with cursor-based pagination, and broadcast member changes to all participants.

```mermaid
graph LR
    subgraph "Live Room (id=room1)"
        OWNER[Owner<br/>RoleOwner]
        V1[Viewer A<br/>RoleViewer]
        V2[Viewer B<br/>RoleViewer]
    end

    subgraph "RTC Server"
        LM[live.Manager]
        SM[signal.Manager]
    end

    OWNER -->|conn.Do → Push live.update| SM
    LM -->|rooms map| OWNER
    LM -->|skiplist by heat| OWNER
    SM -->|conn lookup| V1
    SM -->|conn lookup| V2
```

**Room lifecycle:**
```
Owner creates → room joins RoleOwner → OnClose deletes room
Viewer joins  → room joins RoleViewer → OnClose leaves room
Owner leaves  → Manager.Delete + room.Close (all members cleared)
```

**Heat ranking:**
- `Heat = len(members) * 10`
- Skiplist comparator: `higher heat first, then lower ID first`
- Cursor format: `"{id}|{heat}"`

**Broadcast on member change:**
```
live.update pushed to all room members via signal.Manager.Conn(memberID).Do()
```

**Room list (HTTP):**
```
POST /api/live/list  {cursor, limit}
→ {list: [{id, title, owner_id, count}], next_cursor}
```

## Meeting

Meeting rooms use a simpler flat-map model without ranking. The key difference from live: when the **host leaves**, the entire room is destroyed.

```mermaid
graph TB
    subgraph "Meeting Room (id=room1)"
        HOST[Host<br/>Role=Host]
        M1[Member A<br/>Role=Member]
        M2[Member B<br/>Role=Member]
    end

    subgraph "RTC Server"
        MM[meeting.Manager<br/>rooms map]
        SM[signal.Manager]
    end

    HOST --- MM
    M1 --- MM
    M2 --- MM
    MM --- SM

    HOST -.->|"host leaves →<br/>Manager.Delete + Close"| MM
```

**Key differences from Live:**

| | Live | Meeting |
|---|---|---|
| Ranking | Skiplist by heat | None (map) |
| Pagination | Cursor-based | Full list |
| Roles | Owner / Viewer | Host / Member |
| Room destroy | Owner leaves | Host leaves |
| Create params | ID + Title | ID only |

**Meeting list (HTTP):**
```
POST /api/meeting/list  (no pagination)
→ {id: {id, host_id, members: {id: {id, role}}}}
```

## Signaling Flow

WebRTC signaling messages (SDP offers, answers, ICE candidates) are relayed through the `signal` module. The server acts as a transparent forwarder — it never inspects or modifies the SDP/ICE payload.

```mermaid
sequenceDiagram
    participant A as Client A
    participant S as RTC Server
    participant B as Client B

    Note over A,B: Registration
    A->>S: WS: signal.rtc.register {id:"A"}
    B->>S: WS: signal.rtc.register {id:"B"}

    Note over A,B: WebRTC Signaling
    A->>S: WS: signal.rtc.forward<br/>{type:offer, from:"A", to:"B", call_id, payload:sdp}
    S->>S: json.Unmarshal → SignalMessage
    S->>S: h.m.Conn(req.To) → find B's connection
    S->>B: conn.Do → Push signal.rtc.forward<br/>{type:offer, from:"A", to:"B", call_id, payload:sdp}

    B->>S: WS: signal.rtc.forward<br/>{type:answer, from:"B", to:"A", call_id, payload:sdp}
    S->>A: Push signal.rtc.forward<br/>{type:answer, from:"B", to:"A", call_id, payload:sdp}

    B->>S: WS: signal.rtc.forward<br/>{type:candidate, from:"B", to:"A", call_id, payload:ice}
    S->>A: Push signal.rtc.forward<br/>{type:candidate, from:"B", to:"A", call_id, payload:ice}
```

**SignalMessage protocol:**
```json
{
    "type": "offer|answer|candidate",
    "from": "sender-id",
    "to": "target-id",
    "call_id": "correlation-uuid",
    "payload": { }
}
```

The `payload` field carries the raw SDP or ICE candidate data as JSON. The server never deserializes it — only the `type`, `from`, and `to` fields are parsed for routing.

## Quick Start

### Prerequisites

- Go 1.25+

### Install and Run

```bash
git clone https://github.com/AtoriUzawa/rtc-server
cd rtc-server
go run ./cmd
```

The server starts:
- HTTP API on `:18888`
- WebSocket on `:18887`

### Client Registration

Connect to `ws://localhost:18887/ws` and send:

```json
{"route": "signal.rtc.register", "type": "request", "id": "1", "data": {"id": "alice"}}
```

### Initiate a P2P Call

```json
{"route": "p2p.call", "type": "request", "id": "2", "data": {"from": "alice", "to": "bob", "call_id": "c1"}}
```

### Create a Live Room

```json
{"route": "live.create", "type": "request", "id": "3", "data": {"id": "room1", "title": "My Stream"}}
```

## License

MIT
