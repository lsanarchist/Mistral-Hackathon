# Performance Triage Report

Generated: 2026-03-01T00:32:37+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully

## Cpu: Top cpu hotspots

- **Severity**: Critical
- **Score**: 90
- **Evidence**:
  - Profile: cpu
  - Artifact: demo-output/cpu.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 390 | 756.00 | 756.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 756.00 | 756.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 756.00 | 756.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 756.00 | 756.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 756.00 | 756.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 756.00 | 756.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 756.00 | 756.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 756.00 | 756.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 388 | 207.00 | 207.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 207.00 | 207.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 207.00 | 207.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 207.00 | 207.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 207.00 | 207.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 207.00 | 207.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 207.00 | 207.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 207.00 | 207.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 386 | 137.00 | 137.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 137.00 | 137.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 137.00 | 137.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 137.00 | 137.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 137.00 | 137.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 137.00 | 137.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 137.00 | 137.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 137.00 | 137.00 |
| runtime.madvise | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 544 | 125.00 | 125.00 |
| runtime.sysUnusedOS | /usr/lib/go-1.24/src/runtime/mem_linux.go | 63 | 125.00 | 125.00 |
| runtime.sysUnused | /usr/lib/go-1.24/src/runtime/mem.go | 62 | 125.00 | 125.00 |
| runtime.(*pageAlloc).scavengeOne | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 778 | 125.00 | 125.00 |
| runtime.(*pageAlloc).scavenge.func1 | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 683 | 125.00 | 125.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 125.00 | 125.00 |
| runtime.(*pageAlloc).scavenge | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 682 | 125.00 | 125.00 |
| runtime.(*scavengerState).init.func2 | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 395 | 125.00 | 125.00 |
| runtime.(*scavengerState).run | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 602 | 125.00 | 125.00 |
| runtime.bgscavenge | /usr/lib/go-1.24/src/runtime/mgcscavenge.go | 656 | 125.00 | 125.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 384 | 103.00 | 103.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 103.00 | 103.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 103.00 | 103.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 103.00 | 103.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 103.00 | 103.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 103.00 | 103.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 103.00 | 103.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 103.00 | 103.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 84.00 | 84.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 84.00 | 84.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 84.00 | 84.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 84.00 | 84.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 84.00 | 84.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 84.00 | 84.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 76.00 | 76.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 76.00 | 76.00 |
| runtime.gcDrainMarkWorkerIdle | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1100 | 76.00 | 76.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1520 | 76.00 | 76.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 76.00 | 76.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 76.00 | 76.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 391 | 47.00 | 47.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 47.00 | 47.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 47.00 | 47.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 47.00 | 47.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 47.00 | 47.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 47.00 | 47.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 47.00 | 47.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 47.00 | 47.00 |
| runtime.futex | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 558 | 38.00 | 38.00 |
| runtime.futexsleep | /usr/lib/go-1.24/src/runtime/os_linux.go | 75 | 38.00 | 38.00 |
| runtime.notesleep | /usr/lib/go-1.24/src/runtime/lock_futex.go | 47 | 38.00 | 38.00 |
| runtime.mPark | /usr/lib/go-1.24/src/runtime/proc.go | 1919 | 38.00 | 38.00 |
| runtime.stopm | /usr/lib/go-1.24/src/runtime/proc.go | 2950 | 38.00 | 38.00 |
| runtime.findRunnable | /usr/lib/go-1.24/src/runtime/proc.go | 3699 | 38.00 | 38.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4072 | 38.00 | 38.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 38.00 | 38.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 38.00 | 38.00 |
| runtime.(*gcBits).bitp | /usr/lib/go-1.24/src/runtime/mheap.go | 2420 | 29.00 | 29.00 |
| runtime.(*mspan).markBitsForIndex | /usr/lib/go-1.24/src/runtime/mbitmap.go | 1209 | 29.00 | 29.00 |
| runtime.greyobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1590 | 29.00 | 29.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1463 | 29.00 | 29.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 29.00 | 29.00 |
| runtime.gcDrainMarkWorkerIdle | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1100 | 29.00 | 29.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1520 | 29.00 | 29.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 29.00 | 29.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 29.00 | 29.00 |
| runtime.memmove | /usr/lib/go-1.24/src/runtime/memmove_amd64.s | 389 | 27.00 | 27.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 54 | 27.00 | 27.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 27.00 | 27.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 27.00 | 27.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 27.00 | 27.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 27.00 | 27.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 27.00 | 27.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 27.00 | 27.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 25.00 | 25.00 |
| runtime.gcDrainN | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1317 | 25.00 | 25.00 |
| runtime.gcAssistAlloc1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 670 | 25.00 | 25.00 |
| runtime.gcAssistAlloc.func2 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 561 | 25.00 | 25.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 25.00 | 25.00 |
| runtime.gcAssistAlloc | /usr/lib/go-1.24/src/runtime/mgcmark.go | 560 | 25.00 | 25.00 |
| runtime.deductAssistCredit | /usr/lib/go-1.24/src/runtime/malloc.go | 1691 | 25.00 | 25.00 |
| runtime.mallocgc | /usr/lib/go-1.24/src/runtime/malloc.go | 1044 | 25.00 | 25.00 |
| runtime.rawstring | /usr/lib/go-1.24/src/runtime/string.go | 311 | 25.00 | 25.00 |
| runtime.rawstringtmp | /usr/lib/go-1.24/src/runtime/string.go | 175 | 25.00 | 25.00 |
| runtime.concatstrings | /usr/lib/go-1.24/src/runtime/string.go | 52 | 25.00 | 25.00 |
| runtime.concatstring4 | /usr/lib/go-1.24/src/runtime/string.go | 69 | 25.00 | 25.00 |
| main.inefficientStringHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 314 | 25.00 | 25.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 25.00 | 25.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 25.00 | 25.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 25.00 | 25.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 25.00 | 25.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1444 | 24.00 | 24.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 24.00 | 24.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 24.00 | 24.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1516 | 24.00 | 24.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 24.00 | 24.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 24.00 | 24.00 |
| runtime.(*unwinder).resolveInternal | /usr/lib/go-1.24/src/runtime/traceback.go | 378 | 18.00 | 18.00 |
| runtime.(*unwinder).initAt | /usr/lib/go-1.24/src/runtime/traceback.go | 224 | 18.00 | 18.00 |
| runtime.(*unwinder).init | /usr/lib/go-1.24/src/runtime/traceback.go | 129 | 18.00 | 18.00 |
| runtime.scanstack | /usr/lib/go-1.24/src/runtime/mgcmark.go | 904 | 18.00 | 18.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 240 | 18.00 | 18.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 18.00 | 18.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 18.00 | 18.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 18.00 | 18.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 18.00 | 18.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 18.00 | 18.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 18.00 | 18.00 |
| runtime.futex | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 558 | 17.00 | 17.00 |
| runtime.futexwakeup | /usr/lib/go-1.24/src/runtime/os_linux.go | 88 | 17.00 | 17.00 |
| runtime.notewakeup | /usr/lib/go-1.24/src/runtime/lock_futex.go | 32 | 17.00 | 17.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3063 | 17.00 | 17.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 17.00 | 17.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 17.00 | 17.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 17.00 | 17.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 17.00 | 17.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 17.00 | 17.00 |
| runtime.tgkill | /usr/lib/go-1.24/src/runtime/sys_linux_amd64.s | 177 | 17.00 | 17.00 |
| runtime.signalM | /usr/lib/go-1.24/src/runtime/os_linux.go | 569 | 17.00 | 17.00 |
| runtime.preemptM | /usr/lib/go-1.24/src/runtime/signal_unix.go | 385 | 17.00 | 17.00 |
| runtime.suspendG | /usr/lib/go-1.24/src/runtime/preempt.go | 232 | 17.00 | 17.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 232 | 17.00 | 17.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 17.00 | 17.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 17.00 | 17.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 17.00 | 17.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 17.00 | 17.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 17.00 | 17.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 17.00 | 17.00 |
| runtime.(*mspan).base | /usr/lib/go-1.24/src/runtime/mheap.go | 499 | 16.00 | 16.00 |
| runtime.greyobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1613 | 16.00 | 16.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1463 | 16.00 | 16.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 16.00 | 16.00 |
| runtime.gcDrainMarkWorkerIdle | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1100 | 16.00 | 16.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1520 | 16.00 | 16.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 16.00 | 16.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 16.00 | 16.00 |
| runtime.nanotime | /usr/lib/go-1.24/src/runtime/time_nofake.go | 33 | 14.00 | 14.00 |
| runtime.suspendG | /usr/lib/go-1.24/src/runtime/preempt.go | 247 | 14.00 | 14.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 232 | 14.00 | 14.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 14.00 | 14.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 14.00 | 14.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 14.00 | 14.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 14.00 | 14.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 14.00 | 14.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 14.00 | 14.00 |
| runtime.(*gcBits).bitp | /usr/lib/go-1.24/src/runtime/mheap.go | 2420 | 12.00 | 12.00 |
| runtime.(*mspan).markBitsForIndex | /usr/lib/go-1.24/src/runtime/mbitmap.go | 1209 | 12.00 | 12.00 |
| runtime.greyobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1590 | 12.00 | 12.00 |
| runtime.scanobject | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1463 | 12.00 | 12.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1228 | 12.00 | 12.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 12.00 | 12.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 12.00 | 12.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 12.00 | 12.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 12.00 | 12.00 |
| runtime.suspendG | /usr/lib/go-1.24/src/runtime/preempt.go | 176 | 12.00 | 12.00 |
| runtime.markroot.func1 | /usr/lib/go-1.24/src/runtime/mgcmark.go | 232 | 12.00 | 12.00 |
| runtime.markroot | /usr/lib/go-1.24/src/runtime/mgcmark.go | 214 | 12.00 | 12.00 |
| runtime.gcDrain | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1186 | 12.00 | 12.00 |
| runtime.gcDrainMarkWorkerDedicated | /usr/lib/go-1.24/src/runtime/mgcmark.go | 1110 | 12.00 | 12.00 |
| runtime.gcBgMarkWorker.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1501 | 12.00 | 12.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 12.00 | 12.00 |
| runtime.gcBgMarkWorker | /usr/lib/go-1.24/src/runtime/mgc.go | 1483 | 12.00 | 12.00 |

## Heap: Top heap hotspots

- **Severity**: Critical
- **Score**: 90
- **Evidence**:
  - Profile: heap
  - Artifact: demo-output/heap.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 43691.00 | 43691.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 43691.00 | 43691.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 43691.00 | 43691.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 43691.00 | 43691.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 43691.00 | 43691.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
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
| net/textproto.MIMEHeader.Set | /usr/lib/go-1.24/src/net/textproto/header.go | 22 | 32768.00 | 32768.00 |
| net/http.Header.Set | /usr/lib/go-1.24/src/net/http/header.go | 40 | 32768.00 | 32768.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 159 | 32768.00 | 32768.00 |
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
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 32768.00 | 32768.00 |
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
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 21846.00 | 21846.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21846.00 | 21846.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21846.00 | 21846.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21846.00 | 21846.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21846.00 | 21846.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 12289.00 | 12289.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 12289.00 | 12289.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 12289.00 | 12289.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 12289.00 | 12289.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 12289.00 | 12289.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 12289.00 | 12289.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 12289.00 | 12289.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 12289.00 | 12289.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 12289.00 | 12289.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 12289.00 | 12289.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 12289.00 | 12289.00 |
| time.Sleep | /usr/lib/go-1.24/src/runtime/time.go | 313 | 10923.00 | 10923.00 |
| main.processHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 407 | 10923.00 | 10923.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 278 | 10923.00 | 10923.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10923.00 | 10923.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10923.00 | 10923.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10923.00 | 10923.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10923.00 | 10923.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 10084.00 | 10084.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 10084.00 | 10084.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 10084.00 | 10084.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 10084.00 | 10084.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 10084.00 | 10084.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 10084.00 | 10084.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 10084.00 | 10084.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10084.00 | 10084.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10084.00 | 10084.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10084.00 | 10084.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10084.00 | 10084.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 9104.00 | 9104.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 9104.00 | 9104.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 9104.00 | 9104.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 9104.00 | 9104.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 9104.00 | 9104.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 6993.00 | 6993.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 6993.00 | 6993.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 6993.00 | 6993.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 6993.00 | 6993.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 6993.00 | 6993.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 269 | 5462.00 | 5462.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 5462.00 | 5462.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 5462.00 | 5462.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 5462.00 | 5462.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 5462.00 | 5462.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 5461.00 | 5461.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 5461.00 | 5461.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 5461.00 | 5461.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 5461.00 | 5461.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 5461.00 | 5461.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 5461.00 | 5461.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 5461.00 | 5461.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 5461.00 | 5461.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 5461.00 | 5461.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 5461.00 | 5461.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 88 | 5461.00 | 5461.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 5461.00 | 5461.00 |
| runtime.acquireSudog | /usr/lib/go-1.24/src/runtime/proc.go | 484 | 5461.00 | 5461.00 |
| runtime.chanrecv | /usr/lib/go-1.24/src/runtime/chan.go | 635 | 5461.00 | 5461.00 |
| runtime.chanrecv1 | /usr/lib/go-1.24/src/runtime/chan.go | 506 | 5461.00 | 5461.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1797 | 5461.00 | 5461.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 4096.00 | 4096.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 4096.00 | 4096.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 4096.00 | 4096.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 4096.00 | 4096.00 |
| main.generateTags | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 139 | 4096.00 | 4096.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 106 | 4096.00 | 4096.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 4096.00 | 4096.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 4096.00 | 4096.00 |

## Mutex: Top mutex hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: mutex
  - Artifact: demo-output/mutex.pb.gz

## Block: Top block hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: block
  - Artifact: demo-output/block.pb.gz

## Allocs: Top allocs hotspots

- **Severity**: High
- **Score**: 70
- **Evidence**:
  - Profile: allocs
  - Artifact: demo-output/allocs.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Key | /usr/lib/go-1.24/src/reflect/map_swiss.go | 267 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 769 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537.00 | 65537.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 65537.00 | 65537.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 65537.00 | 65537.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 65537.00 | 65537.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 65537.00 | 65537.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 65537.00 | 65537.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 65537.00 | 65537.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 65537.00 | 65537.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 65537.00 | 65537.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 65537.00 | 65537.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 65537.00 | 65537.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 65537.00 | 65537.00 |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 43691.00 | 43691.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 43691.00 | 43691.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 43691.00 | 43691.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 43691.00 | 43691.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 43691.00 | 43691.00 |
| net/textproto.MIMEHeader.Set | /usr/lib/go-1.24/src/net/textproto/header.go | 22 | 32768.00 | 32768.00 |
| net/http.Header.Set | /usr/lib/go-1.24/src/net/http/header.go | 40 | 32768.00 | 32768.00 |
| main.usersHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 159 | 32768.00 | 32768.00 |
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
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 32768.00 | 32768.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 32768.00 | 32768.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 32768.00 | 32768.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 32768.00 | 32768.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 32768.00 | 32768.00 |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768.00 | 32768.00 |
| reflect.(*MapIter).Value | /usr/lib/go-1.24/src/reflect/map_swiss.go | 311 | 32768.00 | 32768.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 772 | 32768.00 | 32768.00 |
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
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 268 | 21846.00 | 21846.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 21846.00 | 21846.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 21846.00 | 21846.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 21846.00 | 21846.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 21846.00 | 21846.00 |
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 16384.00 | 16384.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 16384.00 | 16384.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 16384.00 | 16384.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16384.00 | 16384.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 12289.00 | 12289.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 12289.00 | 12289.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 12289.00 | 12289.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 12289.00 | 12289.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 12289.00 | 12289.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 12289.00 | 12289.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 12289.00 | 12289.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 12289.00 | 12289.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 12289.00 | 12289.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 12289.00 | 12289.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 12289.00 | 12289.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 12289.00 | 12289.00 |
| time.Sleep | /usr/lib/go-1.24/src/runtime/time.go | 313 | 10923.00 | 10923.00 |
| main.processHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 407 | 10923.00 | 10923.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 278 | 10923.00 | 10923.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10923.00 | 10923.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10923.00 | 10923.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10923.00 | 10923.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10923.00 | 10923.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 764 | 10084.00 | 10084.00 |
| encoding/json.arrayEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 870 | 10084.00 | 10084.00 |
| encoding/json.sliceEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 843 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 10084.00 | 10084.00 |
| encoding/json.interfaceEncoder | /usr/lib/go-1.24/src/encoding/json/encode.go | 680 | 10084.00 | 10084.00 |
| encoding/json.mapEncoder.encode | /usr/lib/go-1.24/src/encoding/json/encode.go | 784 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).reflectValue | /usr/lib/go-1.24/src/encoding/json/encode.go | 333 | 10084.00 | 10084.00 |
| encoding/json.(*encodeState).marshal | /usr/lib/go-1.24/src/encoding/json/encode.go | 309 | 10084.00 | 10084.00 |
| encoding/json.(*Encoder).Encode | /usr/lib/go-1.24/src/encoding/json/stream.go | 210 | 10084.00 | 10084.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 292 | 10084.00 | 10084.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10084.00 | 10084.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10084.00 | 10084.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10084.00 | 10084.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10084.00 | 10084.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 280 | 9104.00 | 9104.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 9104.00 | 9104.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 9104.00 | 9104.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 9104.00 | 9104.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 9104.00 | 9104.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 287 | 6993.00 | 6993.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 6993.00 | 6993.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 6993.00 | 6993.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 6993.00 | 6993.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 6993.00 | 6993.00 |
| main.exportHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 269 | 5462.00 | 5462.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 5462.00 | 5462.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 5462.00 | 5462.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 5462.00 | 5462.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 5462.00 | 5462.00 |
| runtime.acquireSudog | /usr/lib/go-1.24/src/runtime/proc.go | 484 | 5461.00 | 5461.00 |
| runtime.chanrecv | /usr/lib/go-1.24/src/runtime/chan.go | 635 | 5461.00 | 5461.00 |
| runtime.chanrecv1 | /usr/lib/go-1.24/src/runtime/chan.go | 506 | 5461.00 | 5461.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1797 | 5461.00 | 5461.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 5461.00 | 5461.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 5461.00 | 5461.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 5461.00 | 5461.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 5461.00 | 5461.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 5461.00 | 5461.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 5461.00 | 5461.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 5461.00 | 5461.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 5461.00 | 5461.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 5461.00 | 5461.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 5461.00 | 5461.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 88 | 5461.00 | 5461.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 5461.00 | 5461.00 |
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
| main.generateFriends | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 115 | 4096.00 | 4096.00 |
| main.initializeDemoData | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 98 | 4096.00 | 4096.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 41 | 4096.00 | 4096.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 4096.00 | 4096.00 |

### Allocation Analysis

- **Total Allocations**: 739760
- **Top 10% Concentration**: 58.2%
- **Allocation Severity**: High
- **Allocation Score**: 70/100

⚠️ **High Allocation Concentration Detected**
Top functions account for 58.2% of all allocations.
This indicates potential memory allocation hotspots that may benefit from optimization.

#### Top Allocation Hotspots

| Function | File | Line | Count | Percentage |
|----------|------|------|-------|------------|
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537 | 8.9% |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 65537 | 8.9% |
| fmt.Sprintf | /usr/lib/go-1.24/src/fmt/print.go | 240 | 43691 | 5.9% |
| net/textproto.MIMEHeader.Set | /usr/lib/go-1.24/src/net/textproto/header.go | 22 | 32768 | 4.4% |
| reflect.copyVal | /usr/lib/go-1.24/src/reflect/value.go | 1791 | 32768 | 4.4% |

---

*Generated by triageprof*
