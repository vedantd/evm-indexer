# EVM Indexer

**Go-based EVM chain indexer** — ingests blocks and receipts, decodes ERC-20/721 events, handles reorgs, and serves balances/ownership APIs with observability (Prometheus, OpenTelemetry, structured logging).

The project is built **incrementally** with an engineering-first approach: small commits, tested components, and extensibility in mind.

---

##  Scope

- [x] **Ingestion**: backfill planner, WebSocket newHeads, bounded HTTP fetcher
- [ ] **Transformation**: decode logs → ERC-20 transfers, ERC-721 transfers
- [ ] **State Management**: reorg detector, rollbacker, Postgres transactional store
- [ ] **API**: per-chain queries (height, balances, ownership)
- [ ] **Observability**: Prometheus metrics, health checks, tracing, logs
- [ ] **Extensibility**: multiple chains, receipts mode (generic, Erigon, QuickNode)

---

##  Progress Tracker

- [x] Project skeleton (CLI + version)
- [x] Architecture diagrams
- [ ] Config loader
- [ ] Logger
- [ ] Ingestion planner
- [ ] WebSocket client
- [ ] HTTP fetcher
- [ ] Transform layer
- [ ] State and reorg handling
- [ ] Database schema migrations
- [ ] REST API
- [ ] Metrics and health endpoints

---

##  System Architecture 

```mermaid
flowchart TB
    subgraph Config["Config and Registry"]
      CFG["chains.yaml and env vars per chain: name, chain_id, rpc_http, rpc_ws, start_block, batch_size, receipts_mode"]
    end

    subgraph Pipeline["Per Chain Pipeline"]
      direction TB
      subgraph Ingest["Ingestion"]
        direction LR
        WS["WS client newHeads"]
        BP["Backfill planner"]
        HF["HTTP fetcher"]
      end
      subgraph Transform["Transform"]
        direction TB
        LE["Log extractor topics and ABI decode"]
        DE["Domain events builder ERC20 Transfer and ERC721 Transfer"]
      end
      subgraph State["Index and State"]
        direction TB
        RD["Reorg detector parent hash check"]
        RB["Rollbacker reverse materialized state"]
        ST["Store Postgres transaction"]
      end
    end

    subgraph DB["Postgres keyed by chain_id"]
      direction TB
      BL["blocks"]
      T20["erc20_transfers"]
      T721["erc721_transfers"]
      B20["erc20_balances"]
      O721["erc721_ownership"]
      CUR["cursor"]
    end

    subgraph API["API and Observability"]
      direction TB
      APILBL["HTTP API"]
      API1["GET v1 chain height"]
      API2["GET v1 chain erc20 balance"]
      API3["GET v1 chain erc721 owner"]
      API4["POST v1 multi erc20 balances"]
      HL["GET health"]
      MET["GET metrics"]
      PROM["Prometheus"]
      OTEL["OpenTelemetry"]
      JAEGER["Jaeger UI"]
      GRAF["Grafana dashboards"]
      LOGS["Structured logs zerolog"]
    end

    subgraph Legend["Notes"]
      direction TB
      L1["Each chain runs as an independent pipeline for fault isolation"]
      L2["Receipts mode options: generic per tx, erigon eth_getBlockReceipts, quicknode qn_getBlockWithReceipts"]
      L3["All tables include chain_id so multiple chains can share one DB or be sharded later"]
      L4["API supports per chain routes and multi chain aggregation"]
    end

    CFG --> Pipeline
    WS --> HF
    BP --> HF
    HF --> LE
    LE --> DE
    DE --> RD
    RD --> ST
    RB --> ST
    ST --> BL
    ST --> T20
    ST --> T721
    ST --> B20
    ST --> O721
    ST --> CUR
    BL --> APILBL
    B20 --> APILBL
    O721 --> APILBL
    CUR --> APILBL
    APILBL --> API1
    APILBL --> API2
    APILBL --> API3
    APILBL --> API4
    APILBL --> HL
    APILBL --> MET
    MET --> PROM --> GRAF
    OTEL --> JAEGER
    ST --> LOGS



