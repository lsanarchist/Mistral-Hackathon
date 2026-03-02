# Performance Triage Report

Generated: 2026-03-02T01:44:57+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully

## Cpu: Top cpu hotspots

- **Severity**: Critical
- **Score**: 90
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 390 | 859.00 | 859.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 859.00 | 859.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 859.00 | 859.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 859.00 | 859.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 859.00 | 859.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 859.00 | 859.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 859.00 | 859.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 859.00 | 859.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 388 | 270.00 | 270.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 270.00 | 270.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 270.00 | 270.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 270.00 | 270.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 270.00 | 270.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 270.00 | 270.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 270.00 | 270.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 270.00 | 270.00 |
| runtime.madvise | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 544 | 202.00 | 202.00 |
| runtime.sysUnusedOS | /usr/lib/go-1.24/src/runtime/mem_linux.go | 63 | 202.00 | 202.00 |
| runtime.sysUnused | /usr/lib/go-1.24/src/runtime/mem.go | 62 | 202.00 | 202.00 |
| runtime.(*pageAlloc).scavengeOne | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 778 | 202.00 | 202.00 |
| runtime.(*pageAlloc).scavenge.func1 | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 683 | 202.00 | 202.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 202.00 | 202.00 |
| runtime.(*pageAlloc).scavenge | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 682 | 202.00 | 202.00 |
| runtime.(*scavengerState).init.func2 | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 395 | 202.00 | 202.00 |
| runtime.(*scavengerState).run | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 602 | 202.00 | 202.00 |
| runtime.bgscavenge | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 656 | 202.00 | 202.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 183.00 | 183.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 183.00 | 183.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 183.00 | 183.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 183.00 | 183.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 183.00 | 183.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 183.00 | 183.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 386 | 174.00 | 174.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 174.00 | 174.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 174.00 | 174.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 174.00 | 174.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 174.00 | 174.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 174.00 | 174.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 174.00 | 174.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 174.00 | 174.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 170.00 | 170.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 170.00 | 170.00 |
| runtime.gcDrainMarkWorkerIdle | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1100 | 170.00 | 170.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1520 | 170.00 | 170.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 170.00 | 170.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 170.00 | 170.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 96 | 158.00 | 158.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 158.00 | 158.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 158.00 | 158.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 158.00 | 158.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 158.00 | 158.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 158.00 | 158.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 158.00 | 158.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 158.00 | 158.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 158.00 | 158.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 158.00 | 158.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 93 | 136.00 | 136.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 136.00 | 136.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 136.00 | 136.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 136.00 | 136.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 136.00 | 136.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 136.00 | 136.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 136.00 | 136.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 136.00 | 136.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 136.00 | 136.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 136.00 | 136.00 |
| runtime.(*unwinder).resolveInternal | /usr/lib/go-1.24/src/runtime/traceback.go | 378 | 130.00 | 130.00 |
| runtime.(*unwinder).initAt | /usr/lib/go-1.24/src/runtime/traceback.go | 224 | 130.00 | 130.00 |
| runtime.(*unwinder).init | /usr/lib/go-1.24/src/runtime/traceback.go | 129 | 130.00 | 130.00 |
| runtime.scanstack | /usr/lib/go-1.24/src/runtime/mgcmark.go | 904 | 130.00 | 130.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 240 | 130.00 | 130.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 130.00 | 130.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 130.00 | 130.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 130.00 | 130.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 130.00 | 130.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 130.00 | 130.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 130.00 | 130.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 384 | 128.00 | 128.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 128.00 | 128.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 128.00 | 128.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 128.00 | 128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 128.00 | 128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 128.00 | 128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 128.00 | 128.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 112.00 | 112.00 |
| runtime.gcDrainN | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1317 | 112.00 | 112.00 |
| runtime.gcAssistAlloc1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 670 | 112.00 | 112.00 |
| runtime.gcAssistAlloc.func2 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 561 | 112.00 | 112.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 112.00 | 112.00 |
| runtime.gcAssistAlloc | /usr/lib/go-1.24/src/runtime/mgcmark.go | 560 | 112.00 | 112.00 |
| runtime.deductAssistCredit | /usr/lib/go-1.24/src/runtime/malloc.go | 1691 | 112.00 | 112.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1044 | 112.00 | 112.00 |
| runtime.rawstring | /usr/lib/go-1.24/src/runtime/string.go | 311 | 112.00 | 112.00 |
| runtime.rawstringtmp | /usr/lib/go-1.24/src/runtime/string.go | 175 | 112.00 | 112.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 52 | 112.00 | 112.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 112.00 | 112.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 112.00 | 112.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 112.00 | 112.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 112.00 | 112.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 112.00 | 112.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 112.00 | 112.00 |
| runtime.suspendG | /usr/lib/go-1.24/src/runtime/preempt.go | 176 | 104.00 | 104.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 232 | 104.00 | 104.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 104.00 | 104.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 104.00 | 104.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 104.00 | 104.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 104.00 | 104.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 104.00 | 104.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 104.00 | 104.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 96 | 102.00 | 102.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 102.00 | 102.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 102.00 | 102.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 102.00 | 102.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 102.00 | 102.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 102.00 | 102.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 102.00 | 102.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 102.00 | 102.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 102.00 | 102.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 102.00 | 102.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 93 | 100.00 | 100.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 100.00 | 100.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 100.00 | 100.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 100.00 | 100.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 100.00 | 100.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 100.00 | 100.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 100.00 | 100.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 100.00 | 100.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 100.00 | 100.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 100.00 | 100.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 100 | 98.00 | 98.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 98.00 | 98.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 98.00 | 98.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 98.00 | 98.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 98.00 | 98.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 98.00 | 98.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 98.00 | 98.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 98.00 | 98.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 98.00 | 98.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 98.00 | 98.00 |
| runtime.memclrNoHeapPointers | /usr/lib/go-1.24/src/runtime/memclr_amd64.s | 100 | 74.00 | 74.00 |
| runtime.memclrNoHeapPointersChunked | /usr/lib/go-1.24/src/runtime/malloc.go | 1718 | 74.00 | 74.00 |
| runtime.mallocgcLarge | /usr/lib/go-1.24/src/runtime/malloc.go | 1600 | 74.00 | 74.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1063 | 74.00 | 74.00 |
| runtime.makeslice | /usr/lib/go-1.24/src/runtime/slice.go | 116 | 74.00 | 74.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 74.00 | 74.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 74.00 | 74.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 74.00 | 74.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 74.00 | 74.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 74.00 | 74.00 |
| runtime.(*gcBits).bitp | /usr/lib/go-1.24/src/runtime/mheap.go | 2420 | 69.00 | 69.00 |
| runtime.(*mspan).markBitsForIndex | /usr/lib/go-1.24/src/runtime/mbitmap.go | 1209 | 69.00 | 69.00 |
| runtime.greyobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1590 | 69.00 | 69.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1463 | 69.00 | 69.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 69.00 | 69.00 |
| runtime.gcDrainMarkWorkerIdle | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1100 | 69.00 | 69.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1520 | 69.00 | 69.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 69.00 | 69.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 69.00 | 69.00 |
| runtime.readgstatus | /usr/lib/go-1.24/src/runtime/proc.go | 1150 | 58.00 | 58.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 207 | 58.00 | 58.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 58.00 | 58.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 58.00 | 58.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 58.00 | 58.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 58.00 | 58.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 58.00 | 58.00 |
| runtime.(*gcBits).bitp | /usr/lib/go-1.24/src/runtime/mheap.go | 2420 | 57.00 | 57.00 |
| runtime.(*mspan).markBitsForIndex | /usr/lib/go-1.24/src/runtime/mbitmap.go | 1209 | 57.00 | 57.00 |
| runtime.greyobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1590 | 57.00 | 57.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1463 | 57.00 | 57.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 57.00 | 57.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 57.00 | 57.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 57.00 | 57.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 57.00 | 57.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 57.00 | 57.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 391 | 56.00 | 56.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 56.00 | 56.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 56.00 | 56.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 320 | 56.00 | 56.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 56.00 | 56.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 56.00 | 56.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 56.00 | 56.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 56.00 | 56.00 |

## Heap: Top heap hotspots

- **Severity**: Critical
- **Score**: 90
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 163842.00 | 163842.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 163842.00 | 163842.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 163842.00 | 163842.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 163842.00 | 163842.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 163842.00 | 163842.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 163842.00 | 163842.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 163842.00 | 163842.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 163842.00 | 163842.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 163842.00 | 163842.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 163842.00 | 163842.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 131074.00 | 131074.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 131074.00 | 131074.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 131074.00 | 131074.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 131074.00 | 131074.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 131074.00 | 131074.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 131074.00 | 131074.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 131074.00 | 131074.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 131074.00 | 131074.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 131074.00 | 131074.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 131074.00 | 131074.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 131074.00 | 131074.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 131074.00 | 131074.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 131074.00 | 131074.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 131074.00 | 131074.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 131074.00 | 131074.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 98305.00 | 98305.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 98305.00 | 98305.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 98305.00 | 98305.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 98305.00 | 98305.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 98305.00 | 98305.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 98305.00 | 98305.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 98305.00 | 98305.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 98305.00 | 98305.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 98305.00 | 98305.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 98305.00 | 98305.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 75096.00 | 75096.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 75096.00 | 75096.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 75096.00 | 75096.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 75096.00 | 75096.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 75096.00 | 75096.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 276 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 43691.00 | 43691.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| net/textproto.MIMEHeader.Set | /usr/lib/go-1.24/src/net/textproto/header.go | 22 | 32768.00 | 32768.00 |
| net/http.Header.Set | /usr/lib/go-1.24/src/net/http/header.go | 40 | 32768.00 | 32768.00 |
| main.noCacheHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 347 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
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
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 32768.00 | 32768.00 |
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
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 270 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 21846.00 | 21846.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21846.00 | 21846.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21846.00 | 21846.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21846.00 | 21846.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21846.00 | 21846.00 |
| main.goroutineLeakHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 499 | 21845.00 | 21845.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21845.00 | 21845.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21845.00 | 21845.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21845.00 | 21845.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 20030.00 | 20030.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 20030.00 | 20030.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 20030.00 | 20030.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 20030.00 | 20030.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 20030.00 | 20030.00 |
| main.generateTags | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 139 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 106 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| time.Time.Format | /usr/lib/go-1.24/src/time/format.go | 650 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 99 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 14995.00 | 14995.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 14995.00 | 14995.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 14995.00 | 14995.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 14995.00 | 14995.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 14995.00 | 14995.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 269 | 12746.00 | 12746.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 12746.00 | 12746.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 12746.00 | 12746.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 12746.00 | 12746.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 12746.00 | 12746.00 |

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

- **Severity**: High
- **Score**: 70
### Evidence

- **profile**: Profile evidence (100.0% weight)

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 163842.00 | 163842.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 163842.00 | 163842.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 163842.00 | 163842.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 163842.00 | 163842.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 163842.00 | 163842.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 163842.00 | 163842.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 163842.00 | 163842.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 163842.00 | 163842.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 163842.00 | 163842.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 163842.00 | 163842.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 163842.00 | 163842.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 163842.00 | 163842.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 131074.00 | 131074.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 131074.00 | 131074.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 131074.00 | 131074.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 131074.00 | 131074.00 |
| encoding/json.structEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 727 | 131074.00 | 131074.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 131074.00 | 131074.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 131074.00 | 131074.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 131074.00 | 131074.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 131074.00 | 131074.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 131074.00 | 131074.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 131074.00 | 131074.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 160 | 131074.00 | 131074.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 131074.00 | 131074.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 131074.00 | 131074.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 131074.00 | 131074.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 131074.00 | 131074.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 98305.00 | 98305.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 98305.00 | 98305.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 98305.00 | 98305.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 98305.00 | 98305.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 98305.00 | 98305.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 98305.00 | 98305.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 98305.00 | 98305.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 298 | 98305.00 | 98305.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 98305.00 | 98305.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 98305.00 | 98305.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 98305.00 | 98305.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 98305.00 | 98305.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 75096.00 | 75096.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 75096.00 | 75096.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 75096.00 | 75096.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 75096.00 | 75096.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 75096.00 | 75096.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 276 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 43691.00 | 43691.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 32768.00 | 32768.00 |
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
| net.IP.String | /usr/lib/go-1.24/src/net/ip.go | 315 | 32768.00 | 32768.00 |
| net.ipEmptyString | /usr/lib/go-1.24/src/net/ip.go | 332 | 32768.00 | 32768.00 |
| net.(*TCPAddr).String | /usr/lib/go-1.24/src/net/tcpsock.go | 48 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 1939 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 270 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
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
| net/textproto.MIMEHeader.Set | /usr/lib/go-1.24/src/net/textproto/header.go | 22 | 32768.00 | 32768.00 |
| net/http.Header.Set | /usr/lib/go-1.24/src/net/http/header.go | 40 | 32768.00 | 32768.00 |
| main.noCacheHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 347 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 21846.00 | 21846.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21846.00 | 21846.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21846.00 | 21846.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21846.00 | 21846.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21846.00 | 21846.00 |
| main.goroutineLeakHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 499 | 21845.00 | 21845.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21845.00 | 21845.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21845.00 | 21845.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21845.00 | 21845.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 20030.00 | 20030.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 20030.00 | 20030.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 20030.00 | 20030.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 20030.00 | 20030.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 20030.00 | 20030.00 |
| main.generateTags | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 139 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 106 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| time.Time.Format | /usr/lib/go-1.24/src/time/format.go | 650 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 99 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 14995.00 | 14995.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 14995.00 | 14995.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 14995.00 | 14995.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 14995.00 | 14995.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 14995.00 | 14995.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 269 | 12746.00 | 12746.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 12746.00 | 12746.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 12746.00 | 12746.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 12746.00 | 12746.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 12746.00 | 12746.00 |

### Allocation Analysis

- **Total Allocations**: 1402005
- **Top 10% Concentration**: 68.1%
- **Allocation Severity**: High
- **Allocation Score**: 70/100

⚠️ **High Allocation Concentration Detected**
Top functions account for 68.1% of all allocations.
This indicates potential memory allocation hotspots that may benefit from optimization.

#### Top Allocation Hotspots

| Function | File | Line | Count | Percentage |
|----------|------|------|-------|------------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 163842 | 11.7% |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 131074 | 9.3% |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 98305 | 7.0% |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 75096 | 5.4% |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 276 | 65537 | 4.7% |

---

*Generated by triageprof*
