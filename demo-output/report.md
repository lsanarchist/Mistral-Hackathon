# Performance Triage Report

Generated: 2026-03-02T03:04:47+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully

### 🤖 LLM Insights

**Overview**: The profiling data reveals critical performance bottlenecks primarily in JSON encoding and reflection operations, with secondary concerns around runtime scheduling and memory allocation patterns. The system shows excessive heap allocations and CPU overhead from reflection-based operations, particularly in JSON serialization paths.
**Overall Severity**: high (Confidence: 90%)
**Key Themes**: Excessive reflection usage in serialization, High memory allocation pressure, Runtime scheduling overhead, Inefficient data structure handling

#### 📊 Performance Categories

   - **cpu**: 25 findings
   - **blocking**: 15 findings
   - **mutex**: 10 findings
   - **io**: 5 findings
   - **memory**: 45 findings


#### 🚨 Top Risks

**1. Excessive reflection in JSON serialization causing high allocation pressure and CPU overhead**
   - Severity: high
   - Impact: Primary bottleneck affecting both CPU and memory performance. Causes cascading effects on GC, blocking operations, and overall throughput.
   - Likelihood: 
**2. Inefficient goroutine scheduling and potential blocking operations**
   - Severity: medium
   - Impact: Secondary bottleneck that may limit scalability under load. Causes increased latency and reduced concurrency efficiency.
   - Likelihood: 
**3. Potential mutex contention in shared resources**
   - Severity: low
   - Impact: Currently minor but could become significant under higher load. May limit horizontal scalability.
   - Likelihood: 

#### 🎯 Top Action Items

**1. Replace reflection-based JSON encoding with code-generated marshalers for all critical types**
   - Priority: high
   - Estimated Effort: 8-16 hours
**2. Implement object pooling for JSON encoding buffers and temporary objects**
   - Priority: high
   - Estimated Effort: 4-8 hours
**3. Analyze and optimize goroutine usage patterns to reduce scheduling overhead**
   - Priority: medium
   - Estimated Effort: 6-12 hours

## Cpu: Top cpu hotspots

- **Severity**: Low
- **Score**: 50
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| runtime.execute | /usr/lib/go-1.24/src/runtime/proc.go | 3286 | 3.00 | 3.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4127 | 3.00 | 3.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 3.00 | 3.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 3.00 | 3.00 |
| runtime.futex | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 558 | 3.00 | 3.00 |
| runtime.futexwakeup | /usr/lib/go-1.24/src/runtime/os_linux.go | 88 | 3.00 | 3.00 |
| runtime.notewakeup | /usr/lib/go-1.24/src/runtime/lock_futex.go | 32 | 3.00 | 3.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3063 | 3.00 | 3.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 3.00 | 3.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 3.00 | 3.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 3.00 | 3.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 3.00 | 3.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 3.00 | 3.00 |
| runtime.write1 | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 99 | 2.00 | 2.00 |
| runtime.write | /usr/lib/go-1.24/src/runtime/time_nofake.go | 57 | 2.00 | 2.00 |
| runtime.netpollBreak | /usr/lib/go-1.24/src/runtime/netpoll_epoll.go | 76 | 2.00 | 2.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3696 | 2.00 | 2.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 2.00 | 2.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 2.00 | 2.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 2.00 | 2.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1026 | 2.00 | 2.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 2.00 | 2.00 |
| runtime.stealWork | /usr/lib/go-1.24/src/runtime/proc.go | 3764 | 2.00 | 2.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3434 | 2.00 | 2.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 2.00 | 2.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 2.00 | 2.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 2.00 | 2.00 |
| runtime.ready | /usr/lib/go-1.24/src/runtime/proc.go | 1067 | 2.00 | 2.00 |
| runtime.goroutineReady.goready.func1 | /usr/lib/go-1.24/src/runtime/proc.go | 456 | 2.00 | 2.00 |
| runtime.goready | /usr/lib/go-1.24/src/runtime/proc.go | 455 | 2.00 | 2.00 |
| runtime.goroutineReady | /usr/lib/go-1.24/src/runtime/time.go | 417 | 2.00 | 2.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1176 | 2.00 | 2.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 2.00 | 2.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 2.00 | 2.00 |
| runtime.stealWork | /usr/lib/go-1.24/src/runtime/proc.go | 3764 | 2.00 | 2.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3434 | 2.00 | 2.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 2.00 | 2.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 2.00 | 2.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 2.00 | 2.00 |
| runtime.casgstatus | /usr/lib/go-1.24/src/runtime/proc.go | 1254 | 1.00 | 1.00 |
| runtime.ready | /usr/lib/go-1.24/src/runtime/proc.go | 1074 | 1.00 | 1.00 |
| runtime.goroutineReady.goready.func1 | /usr/lib/go-1.24/src/runtime/proc.go | 456 | 1.00 | 1.00 |
| runtime.goready | /usr/lib/go-1.24/src/runtime/proc.go | 455 | 1.00 | 1.00 |
| runtime.goroutineReady | /usr/lib/go-1.24/src/runtime/time.go | 417 | 1.00 | 1.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1176 | 1.00 | 1.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 1.00 | 1.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 1.00 | 1.00 |
| runtime.stealWork | /usr/lib/go-1.24/src/runtime/proc.go | 3764 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3434 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*timers).addHeap | /usr/lib/go-1.24/src/runtime/time.go | 437 | 1.00 | 1.00 |
| runtime.(*timer).maybeAdd | /usr/lib/go-1.24/src/runtime/time.go | 692 | 1.00 | 1.00 |
| runtime.(*timer).modify | /usr/lib/go-1.24/src/runtime/time.go | 620 | 1.00 | 1.00 |
| runtime.(*timer).reset | /usr/lib/go-1.24/src/runtime/time.go | 706 | 1.00 | 1.00 |
| runtime.resetForSleep | /usr/lib/go-1.24/src/runtime/time.go | 347 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4180 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*timers).siftDown | /usr/lib/go-1.24/src/runtime/time.go | 1329 | 1.00 | 1.00 |
| runtime.(*timers).deleteMin | /usr/lib/go-1.24/src/runtime/time.go | 534 | 1.00 | 1.00 |
| runtime.(*timer).updateHeap | /usr/lib/go-1.24/src/runtime/time.go | 274 | 1.00 | 1.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1105 | 1.00 | 1.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 1.00 | 1.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 1.00 | 1.00 |
| runtime.stealWork | /usr/lib/go-1.24/src/runtime/proc.go | 3764 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3434 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1176 | 1.00 | 1.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 1.00 | 1.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3340 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4153 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*timers).siftUp | /usr/lib/go-1.24/src/runtime/time.go | 1288 | 1.00 | 1.00 |
| runtime.(*timers).addHeap | /usr/lib/go-1.24/src/runtime/time.go | 438 | 1.00 | 1.00 |
| runtime.(*timer).maybeAdd | /usr/lib/go-1.24/src/runtime/time.go | 692 | 1.00 | 1.00 |
| runtime.(*timer).modify | /usr/lib/go-1.24/src/runtime/time.go | 620 | 1.00 | 1.00 |
| runtime.(*timer).reset | /usr/lib/go-1.24/src/runtime/time.go | 706 | 1.00 | 1.00 |
| runtime.resetForSleep | /usr/lib/go-1.24/src/runtime/time.go | 347 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4180 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| main.goroutineLeakHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 506 | 1.00 | 1.00 |
| runtime.ready | /usr/lib/go-1.24/src/runtime/proc.go | 1067 | 1.00 | 1.00 |
| runtime.goroutineReady.goready.func1 | /usr/lib/go-1.24/src/runtime/proc.go | 456 | 1.00 | 1.00 |
| runtime.goready | /usr/lib/go-1.24/src/runtime/proc.go | 455 | 1.00 | 1.00 |
| runtime.goroutineReady | /usr/lib/go-1.24/src/runtime/time.go | 417 | 1.00 | 1.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1176 | 1.00 | 1.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 1.00 | 1.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3340 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.nanotime | /usr/lib/go-1.24/src/runtime/time_nofake.go | 33 | 1.00 | 1.00 |
| time.Sleep | /usr/lib/go-1.24/src/runtime/time.go | 324 | 1.00 | 1.00 |
| main.goroutineLeakHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 503 | 1.00 | 1.00 |
| runtime.gopark | /usr/lib/go-1.24/src/runtime/proc.go | 436 | 1.00 | 1.00 |
| time.Sleep | /usr/lib/go-1.24/src/runtime/time.go | 338 | 1.00 | 1.00 |
| main.goroutineLeakHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 503 | 1.00 | 1.00 |
| runtime.(*timers).siftDown | /usr/lib/go-1.24/src/runtime/time.go | 1329 | 1.00 | 1.00 |
| runtime.(*timers).deleteMin | /usr/lib/go-1.24/src/runtime/time.go | 534 | 1.00 | 1.00 |
| runtime.(*timer).updateHeap | /usr/lib/go-1.24/src/runtime/time.go | 274 | 1.00 | 1.00 |
| runtime.(*timer).unlockAndRun | /usr/lib/go-1.24/src/runtime/time.go | 1105 | 1.00 | 1.00 |
| runtime.(*timers).run | /usr/lib/go-1.24/src/runtime/time.go | 1051 | 1.00 | 1.00 |
| runtime.(*timers).check | /usr/lib/go-1.24/src/runtime/time.go | 985 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3340 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*timers).unlock | /usr/lib/go-1.24/src/runtime/time.go | 171 | 1.00 | 1.00 |
| runtime.(*timer).maybeAdd | /usr/lib/go-1.24/src/runtime/time.go | 695 | 1.00 | 1.00 |
| runtime.(*timer).modify | /usr/lib/go-1.24/src/runtime/time.go | 620 | 1.00 | 1.00 |
| runtime.(*timer).reset | /usr/lib/go-1.24/src/runtime/time.go | 706 | 1.00 | 1.00 |
| runtime.resetForSleep | /usr/lib/go-1.24/src/runtime/time.go | 347 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4180 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| runtime.(*randomEnum).next | /usr/lib/go-1.24/src/runtime/proc.go | 7342 | 1.00 | 1.00 |
| runtime.stealWork | /usr/lib/go-1.24/src/runtime/proc.go | 3740 | 1.00 | 1.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3434 | 1.00 | 1.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 1.00 | 1.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 1.00 | 1.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 1.00 | 1.00 |
| math/rand.(*Rand).Int31n | /usr/lib/go-1.24/src/math/rand/rand.go | 149 | 1.00 | 1.00 |
| math/rand.(*Rand).Intn | /usr/lib/go-1.24/src/math/rand/rand.go | 183 | 1.00 | 1.00 |
| math/rand.Intn | /usr/lib/go-1.24/src/math/rand/rand.go | 464 | 1.00 | 1.00 |
| main.goroutineLeakHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 504 | 1.00 | 1.00 |

## Heap: Top heap hotspots

- **Severity**: Critical
- **Score**: 90
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 43692.00 | 43692.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 43692.00 | 43692.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 43692.00 | 43692.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 43692.00 | 43692.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 43692.00 | 43692.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 96 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 97 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 32768.00 | 32768.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 32768.00 | 32768.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| runtime.(*timers).addHeap | /usr/lib/go-1.24/src/runtime/time.go | 437 | 32768.00 | 32768.00 |
| runtime.(*timer).maybeAdd | /usr/lib/go-1.24/src/runtime/time.go | 692 | 32768.00 | 32768.00 |
| runtime.(*timer).modify | /usr/lib/go-1.24/src/runtime/time.go | 620 | 32768.00 | 32768.00 |
| runtime.(*timer).reset | /usr/lib/go-1.24/src/runtime/time.go | 706 | 32768.00 | 32768.00 |
| runtime.resetForSleep | /usr/lib/go-1.24/src/runtime/time.go | 347 | 32768.00 | 32768.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4180 | 32768.00 | 32768.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 29950.00 | 29950.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 29950.00 | 29950.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 29950.00 | 29950.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 29950.00 | 29950.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 29950.00 | 29950.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 21845.00 | 21845.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 21845.00 | 21845.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 21845.00 | 21845.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 21845.00 | 21845.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 21845.00 | 21845.00 |
| main.goroutineLeakHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 499 | 21845.00 | 21845.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21845.00 | 21845.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21845.00 | 21845.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21845.00 | 21845.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 10925.00 | 10925.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10925.00 | 10925.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10925.00 | 10925.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10925.00 | 10925.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10925.00 | 10925.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 101 | 10923.00 | 10923.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 10923.00 | 10923.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 10923.00 | 10923.00 |
| runtime/pprof.allFrames | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 214 | 8740.00 | 8740.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 451 | 8740.00 | 8740.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 376 | 8740.00 | 8740.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8740.00 | 8740.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 8193.00 | 8193.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 8193.00 | 8193.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 8193.00 | 8193.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 8193.00 | 8193.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 8193.00 | 8193.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 8193.00 | 8193.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 8193.00 | 8193.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 8193.00 | 8193.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 8193.00 | 8193.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8193.00 | 8193.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8193.00 | 8193.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8193.00 | 8193.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8193.00 | 8193.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 39 | 8192.00 | 8192.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 8192.00 | 8192.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 8192.00 | 8192.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8192.00 | 8192.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8192.00 | 8192.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8192.00 | 8192.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8192.00 | 8192.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8192.00 | 8192.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8192.00 | 8192.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8192.00 | 8192.00 |

## Mutex: Top mutex hotspots

- **Severity**: Low
- **Score**: 50
### Evidence

- **profile**: Profile evidence (100.0% weight)

## Block: Top block hotspots

- **Severity**: Low
- **Score**: 50
### Evidence

- **profile**: Profile evidence (100.0% weight)

## Allocs: Top allocs hotspots

- **Severity**: Critical
- **Score**: 90
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 43692.00 | 43692.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 43692.00 | 43692.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 43692.00 | 43692.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 43692.00 | 43692.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 43692.00 | 43692.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 97 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 96 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| runtime.(*timers).addHeap | /usr/lib/go-1.24/src/runtime/time.go | 437 | 32768.00 | 32768.00 |
| runtime.(*timer).maybeAdd | /usr/lib/go-1.24/src/runtime/time.go | 692 | 32768.00 | 32768.00 |
| runtime.(*timer).modify | /usr/lib/go-1.24/src/runtime/time.go | 620 | 32768.00 | 32768.00 |
| runtime.(*timer).reset | /usr/lib/go-1.24/src/runtime/time.go | 706 | 32768.00 | 32768.00 |
| runtime.resetForSleep | /usr/lib/go-1.24/src/runtime/time.go | 347 | 32768.00 | 32768.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4180 | 32768.00 | 32768.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 32768.00 | 32768.00 |
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 32768.00 | 32768.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 32768.00 | 32768.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 32768.00 | 32768.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 32768.00 | 32768.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 32768.00 | 32768.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 32768.00 | 32768.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 32768.00 | 32768.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 32768.00 | 32768.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 32768.00 | 32768.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 32768.00 | 32768.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 29950.00 | 29950.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 29950.00 | 29950.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 29950.00 | 29950.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 29950.00 | 29950.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 29950.00 | 29950.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 21845.00 | 21845.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 21845.00 | 21845.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 21845.00 | 21845.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 21845.00 | 21845.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 21845.00 | 21845.00 |
| main.goroutineLeakHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 499 | 21845.00 | 21845.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21845.00 | 21845.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21845.00 | 21845.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21845.00 | 21845.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 10925.00 | 10925.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10925.00 | 10925.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10925.00 | 10925.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10925.00 | 10925.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10925.00 | 10925.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 101 | 10923.00 | 10923.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 10923.00 | 10923.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 10923.00 | 10923.00 |
| runtime/pprof.allFrames | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 214 | 8740.00 | 8740.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 451 | 8740.00 | 8740.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 376 | 8740.00 | 8740.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8740.00 | 8740.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 8193.00 | 8193.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 8193.00 | 8193.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 8193.00 | 8193.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 8193.00 | 8193.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 8193.00 | 8193.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 8193.00 | 8193.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 8193.00 | 8193.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 8193.00 | 8193.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 8193.00 | 8193.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 8193.00 | 8193.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8193.00 | 8193.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8193.00 | 8193.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8193.00 | 8193.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8193.00 | 8193.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 39 | 8192.00 | 8192.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 8192.00 | 8192.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 8192.00 | 8192.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8192.00 | 8192.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8192.00 | 8192.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8192.00 | 8192.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8192.00 | 8192.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8192.00 | 8192.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8192.00 | 8192.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8192.00 | 8192.00 |

### Allocation Analysis

- **Total Allocations**: 742343
- **Top 10% Concentration**: 79.0%
- **Allocation Severity**: Critical
- **Allocation Score**: 90/100

⚠️ **High Allocation Concentration Detected**
Top functions account for 79.0% of all allocations.
This indicates potential memory allocation hotspots that may benefit from optimization.

#### Top Allocation Hotspots

| Function | File | Line | Count | Percentage |
|----------|------|------|-------|------------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537 | 8.8% |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 43692 | 5.9% |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691 | 5.9% |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 32768 | 4.4% |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768 | 4.4% |

---

*Generated by triageprof*
