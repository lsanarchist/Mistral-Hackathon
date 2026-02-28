# Performance Triage Report

Generated: 2026-02-28T19:12:26+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully
  - Callgraph analysis enabled (depth 3)
  - Regression analysis performed against baseline

## Cpu: Top cpu hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: cpu
  - Artifact: out/cpu.pb.gz

### Regression Analysis

- **Baseline Score**: 0
- **Current Score**: 0
- **Delta**: 0 (0.0%)
- **Severity**: None
- **Confidence**: 50%


## Heap: Top heap hotspots

- **Severity**: Critical
- **Score**: 90
- **Evidence**:
  - Profile: heap
  - Artifact: out/heap.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 10001.00 | 10001.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10001.00 | 10001.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10001.00 | 10001.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10001.00 | 10001.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10001.00 | 10001.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 2308.00 | 2308.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 2308.00 | 2308.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 2308.00 | 2308.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 2308.00 | 2308.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 2308.00 | 2308.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 2308.00 | 2308.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 2308.00 | 2308.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 2308.00 | 2308.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 2308.00 | 2308.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 513.00 | 513.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 513.00 | 513.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 513.00 | 513.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 513.00 | 513.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 513.00 | 513.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 513.00 | 513.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 513.00 | 513.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 513.00 | 513.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.goschedImpl | /usr/lib/go-1.24/src/runtime/proc.go | 4235 | 256.00 | 256.00 |
| runtime.gopreempt_m | /usr/lib/go-1.24/src/runtime/proc.go | 4252 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |

### Callgraph Analysis (Depth 3)

```
net.IP.String (32768.0% cum, 32768.0% flat)
  net.ipEmptyString (32768.0% cum, 32768.0% flat)
  net/http.(*conn).serve (32768.0% cum, 32768.0% flat)
main.allocHeavyHandler (10001.0% cum, 10001.0% flat)
  net/http.HandlerFunc.ServeHTTP (10001.0% cum, 10001.0% flat)
  net/http.(*ServeMux).ServeHTTP (10001.0% cum, 10001.0% flat)
  net/http.serverHandler.ServeHTTP (10001.0% cum, 10001.0% flat)
  net/http.(*conn).serve (10001.0% cum, 10001.0% flat)
runtime.allocm (2308.0% cum, 2308.0% flat)
  runtime.newm (2308.0% cum, 2308.0% flat)
  runtime.startm (2308.0% cum, 2308.0% flat)
  runtime.wakep (2308.0% cum, 2308.0% flat)
  runtime.resetspinning (2308.0% cum, 2308.0% flat)
  runtime.schedule (2308.0% cum, 2308.0% flat)
  runtime.mstart1 (2308.0% cum, 2308.0% flat)
  runtime.mstart0 (2308.0% cum, 2308.0% flat)
  runtime.mstart (2308.0% cum, 2308.0% flat)
net.open (8.0% cum, 8.0% flat)
  net.maxListenerBacklog (8.0% cum, 8.0% flat)
  net.listenerBacklog.func1 (8.0% cum, 8.0% flat)
  sync.(*Once).doSlow (8.0% cum, 8.0% flat)
  sync.(*Once).Do (8.0% cum, 8.0% flat)
  net.socket (8.0% cum, 8.0% flat)
  net.internetSocket (8.0% cum, 8.0% flat)
  net.(*sysListener).listenTCPProto (8.0% cum, 8.0% flat)
  net.(*sysListener).listenMPTCP (8.0% cum, 8.0% flat)
  net.(*ListenConfig).Listen (8.0% cum, 8.0% flat)
  net.Listen (8.0% cum, 8.0% flat)
  net/http.(*Server).ListenAndServe (8.0% cum, 8.0% flat)
  net/http.ListenAndServe (8.0% cum, 8.0% flat)
  runtime.main (8.0% cum, 8.0% flat)
compress/flate.newDeflateFast (0.0% cum, 0.0% flat)
  compress/flate.NewWriter (0.0% cum, 0.0% flat)
  runtime/pprof.(*profileBuilder).build (0.0% cum, 0.0% flat)
  runtime/pprof.profileWriter (0.0% cum, 0.0% flat)
```

### Regression Analysis

- **Baseline Score**: 80
- **Current Score**: 80
- **Delta**: 0 (0.0%)
- **Severity**: None
- **Confidence**: 50%


## Mutex: Top mutex hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: mutex
  - Artifact: out/mutex.pb.gz

### Regression Analysis

- **Baseline Score**: 0
- **Current Score**: 0
- **Delta**: 0 (0.0%)
- **Severity**: None
- **Confidence**: 50%


## Block: Top block hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: block
  - Artifact: out/block.pb.gz

### Regression Analysis

- **Baseline Score**: 0
- **Current Score**: 0
- **Delta**: 0 (0.0%)
- **Severity**: None
- **Confidence**: 50%


## Allocs: Top allocs hotspots

- **Severity**: Critical
- **Score**: 90
- **Evidence**:
  - Profile: allocs
  - Artifact: out/allocs.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 10002.00 | 10002.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10002.00 | 10002.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10002.00 | 10002.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10002.00 | 10002.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10002.00 | 10002.00 |
| runtime/pprof.allFrames | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 232 | 5461.00 | 5461.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 451 | 5461.00 | 5461.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 47 | 5461.00 | 5461.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 5461.00 | 5461.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 5461.00 | 5461.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 5461.00 | 5461.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 5461.00 | 5461.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 5461.00 | 5461.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 5461.00 | 5461.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 5461.00 | 5461.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 5461.00 | 5461.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 5461.00 | 5461.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 2308.00 | 2308.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 2308.00 | 2308.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 2308.00 | 2308.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 2308.00 | 2308.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 2308.00 | 2308.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 2308.00 | 2308.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 2308.00 | 2308.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 2308.00 | 2308.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 2308.00 | 2308.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 513.00 | 513.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 513.00 | 513.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 513.00 | 513.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 513.00 | 513.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 513.00 | 513.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 513.00 | 513.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 513.00 | 513.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 513.00 | 513.00 |

### Callgraph Analysis (Depth 3)

```
net.IP.String (32768.0% cum, 32768.0% flat)
  net.ipEmptyString (32768.0% cum, 32768.0% flat)
  net/http.(*conn).serve (32768.0% cum, 32768.0% flat)
main.allocHeavyHandler (10002.0% cum, 10002.0% flat)
  net/http.HandlerFunc.ServeHTTP (10002.0% cum, 10002.0% flat)
  net/http.(*ServeMux).ServeHTTP (10002.0% cum, 10002.0% flat)
  net/http.serverHandler.ServeHTTP (10002.0% cum, 10002.0% flat)
  net/http.(*conn).serve (10002.0% cum, 10002.0% flat)
runtime/pprof.allFrames (5461.0% cum, 5461.0% flat)
  runtime/pprof.(*profileBuilder).appendLocsForStack (5461.0% cum, 5461.0% flat)
  runtime/pprof.writeHeapProto (5461.0% cum, 5461.0% flat)
  runtime/pprof.writeHeapInternal (5461.0% cum, 5461.0% flat)
  runtime/pprof.writeHeap (5461.0% cum, 5461.0% flat)
  runtime/pprof.(*Profile).WriteTo (5461.0% cum, 5461.0% flat)
  net/http/pprof.handler.ServeHTTP (5461.0% cum, 5461.0% flat)
  net/http/pprof.Index (5461.0% cum, 5461.0% flat)
  net/http.HandlerFunc.ServeHTTP (5461.0% cum, 5461.0% flat)
  net/http.(*ServeMux).ServeHTTP (5461.0% cum, 5461.0% flat)
  net/http.serverHandler.ServeHTTP (5461.0% cum, 5461.0% flat)
  net/http.(*conn).serve (5461.0% cum, 5461.0% flat)
runtime.allocm (2308.0% cum, 2308.0% flat)
  runtime.newm (2308.0% cum, 2308.0% flat)
  runtime.startm (2308.0% cum, 2308.0% flat)
  runtime.wakep (2308.0% cum, 2308.0% flat)
  runtime.resetspinning (2308.0% cum, 2308.0% flat)
  runtime.schedule (2308.0% cum, 2308.0% flat)
  runtime.mstart1 (2308.0% cum, 2308.0% flat)
  runtime.mstart0 (2308.0% cum, 2308.0% flat)
  runtime.mstart (2308.0% cum, 2308.0% flat)
net.open (8.0% cum, 8.0% flat)
  net.maxListenerBacklog (8.0% cum, 8.0% flat)
  net.listenerBacklog.func1 (8.0% cum, 8.0% flat)
  sync.(*Once).doSlow (8.0% cum, 8.0% flat)
  sync.(*Once).Do (8.0% cum, 8.0% flat)
  net.socket (8.0% cum, 8.0% flat)
  net.internetSocket (8.0% cum, 8.0% flat)
  net.(*sysListener).listenTCPProto (8.0% cum, 8.0% flat)
  net.(*sysListener).listenMPTCP (8.0% cum, 8.0% flat)
  net.(*ListenConfig).Listen (8.0% cum, 8.0% flat)
  net.Listen (8.0% cum, 8.0% flat)
  net/http.(*Server).ListenAndServe (8.0% cum, 8.0% flat)
  net/http.ListenAndServe (8.0% cum, 8.0% flat)
  runtime.main (8.0% cum, 8.0% flat)
```

### Regression Analysis

- **Baseline Score**: 80
- **Current Score**: 80
- **Delta**: 0 (0.0%)
- **Severity**: None
- **Confidence**: 50%


---

*Generated by triageprof*
