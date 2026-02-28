# Performance Triage Report

Generated: 2026-02-28T19:23:33+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully
  - Callgraph analysis enabled (depth 3)

## Cpu: Top cpu hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: cpu
  - Artifact: out/cpu.pb.gz

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 30004.00 | 30004.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 30004.00 | 30004.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 30004.00 | 30004.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 30004.00 | 30004.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 30004.00 | 30004.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 21845.00 | 21845.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 21845.00 | 21845.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 21845.00 | 21845.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 21845.00 | 21845.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 8192.00 | 8192.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 8192.00 | 8192.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 8192.00 | 8192.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 8192.00 | 8192.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 8192.00 | 8192.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 8192.00 | 8192.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 8192.00 | 8192.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 8192.00 | 8192.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 8192.00 | 8192.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 8192.00 | 8192.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 8192.00 | 8192.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 8192.00 | 8192.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 5461.00 | 5461.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 5461.00 | 5461.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 5461.00 | 5461.00 |
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
| runtime/pprof.(*profileBuilder).emitLocation | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 622 | 293.00 | 293.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 467 | 293.00 | 293.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 47 | 293.00 | 293.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 293.00 | 293.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 293.00 | 293.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 293.00 | 293.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 293.00 | 293.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 293.00 | 293.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 293.00 | 293.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 293.00 | 293.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 293.00 | 293.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 293.00 | 293.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.goschedImpl | /usr/lib/go-1.24/src/runtime/proc.go | 4235 | 256.00 | 256.00 |
| runtime.gopreempt_m | /usr/lib/go-1.24/src/runtime/proc.go | 4252 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 8.00 | 8.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 8.00 | 8.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8.00 | 8.00 |
| net.open | /usr/lib/go-1.24/src/net/parse.go | 80 | 8.00 | 8.00 |
| net.maxListenerBacklog | /usr/lib/go-1.24/src/net/sock_linux.go | 35 | 8.00 | 8.00 |
| net.listenerBacklog.func1 | /usr/lib/go-1.24/src/net/net.go | 400 | 8.00 | 8.00 |
| sync.(*Once).doSlow | /usr/lib/go-1.24/src/sync/once.go | 78 | 8.00 | 8.00 |
| sync.(*Once).Do | /usr/lib/go-1.24/src/sync/once.go | 69 | 8.00 | 8.00 |
| net.listenerBacklog | /usr/lib/go-1.24/src/net/net.go | 400 | 8.00 | 8.00 |
| net.socket | /usr/lib/go-1.24/src/net/sock_posix.go | 57 | 8.00 | 8.00 |
| net.internetSocket | /usr/lib/go-1.24/src/net/ipsock_posix.go | 167 | 8.00 | 8.00 |
| net.(*sysListener).listenTCPProto | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 189 | 8.00 | 8.00 |
| net.(*sysListener).listenMPTCP | /usr/lib/go-1.24/src/net/mptcpsock_linux.go | 79 | 8.00 | 8.00 |
| net.(*ListenConfig).Listen | /usr/lib/go-1.24/src/net/dial.go | 819 | 8.00 | 8.00 |
| net.Listen | /usr/lib/go-1.24/src/net/dial.go | 898 | 8.00 | 8.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3346 | 8.00 | 8.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 8.00 | 8.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 8.00 | 8.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 4.00 | 4.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 4.00 | 4.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 4.00 | 4.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 4.00 | 4.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 4.00 | 4.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 4.00 | 4.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 4.00 | 4.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 4.00 | 4.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 4.00 | 4.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 4.00 | 4.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 4.00 | 4.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 4.00 | 4.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 4.00 | 4.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 4.00 | 4.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 4.00 | 4.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 4.00 | 4.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 4.00 | 4.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 4.00 | 4.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 2.00 | 2.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 2.00 | 2.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 2.00 | 2.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 2.00 | 2.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 2.00 | 2.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 2.00 | 2.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 2.00 | 2.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 2.00 | 2.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 2.00 | 2.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 2.00 | 2.00 |
| runtime/pprof.writeGoroutineStacks | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 762 | 2.00 | 2.00 |
| runtime/pprof.writeGoroutine | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 753 | 2.00 | 2.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 2.00 | 2.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 2.00 | 2.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 2.00 | 2.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 2.00 | 2.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 2.00 | 2.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 2.00 | 2.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 2.00 | 2.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 2.00 | 2.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 2.00 | 2.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 2.00 | 2.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 2.00 | 2.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 2.00 | 2.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 2.00 | 2.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 2.00 | 2.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 2.00 | 2.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 2.00 | 2.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 2.00 | 2.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 2.00 | 2.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 2.00 | 2.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 2.00 | 2.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 2.00 | 2.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 2.00 | 2.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 2.00 | 2.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 2.00 | 2.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 2.00 | 2.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 2.00 | 2.00 |

### Callgraph Analysis

```
net.IP.String (cum: 32768.0, flat: 32768.0, depth: 0)
  ├── net.ipEmptyString (cum: 32768.0, flat: 32768.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 32768.0, flat: 32768.0, depth: 1)
main.allocHeavyHandler (cum: 30004.0, flat: 30004.0, depth: 0)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 30004.0, flat: 30004.0, depth: 1)
internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All (cum: 21845.0, flat: 21845.0, depth: 0)
  ├── unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 (cum: 21845.0, flat: 21845.0, depth: 1)
  ├── unique.registerCleanup.func1 (cum: 21845.0, flat: 21845.0, depth: 1)
  ├── runtime.unique_runtime_registerUniqueMapCleanup.func2 (cum: 21845.0, flat: 21845.0, depth: 1)
syscall.anyToSockaddr (cum: 8192.0, flat: 8192.0, depth: 0)
  ├── syscall.Accept4 (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── internal/poll.accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── internal/poll.(*FD).Accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*netFD).accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*TCPListener).accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*TCPListener).Accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.(*Server).Serve (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.(*Server).ListenAndServe (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.ListenAndServe (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── runtime.main (cum: 8192.0, flat: 8192.0, depth: 1)
internal/sync.runtime_SemacquireMutex (cum: 5461.0, flat: 5461.0, depth: 0)
  ├── internal/sync.(*Mutex).lockSlow (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── internal/sync.(*Mutex).Lock (cum: 5461.0, flat: 5461.0, depth: 1)
runtime/pprof.allFrames (cum: 5461.0, flat: 5461.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).appendLocsForStack (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeap (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http/pprof.Index (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 5461.0, flat: 5461.0, depth: 1)
runtime.allocm (cum: 2308.0, flat: 2308.0, depth: 0)
  ├── runtime.newm (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.startm (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.wakep (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.resetspinning (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.schedule (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart1 (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart0 (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart (cum: 2308.0, flat: 2308.0, depth: 1)
runtime/pprof.(*profileBuilder).emitLocation (cum: 293.0, flat: 293.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).appendLocsForStack (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeAlloc (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http/pprof.Index (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 293.0, flat: 293.0, depth: 1)
compress/flate.newDeflateFast (cum: 8.0, flat: 8.0, depth: 0)
  ├── compress/flate.NewWriter (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.(*profileBuilder).build (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.printCountCycleProfile (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.writeProfileInternal (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.writeBlock (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http/pprof.Index (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 8.0, flat: 8.0, depth: 1)
net.open (cum: 8.0, flat: 8.0, depth: 0)
  ├── net.maxListenerBacklog (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.listenerBacklog.func1 (cum: 8.0, flat: 8.0, depth: 1)
  ├── sync.(*Once).doSlow (cum: 8.0, flat: 8.0, depth: 1)
  ├── sync.(*Once).Do (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.socket (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.internetSocket (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*sysListener).listenTCPProto (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*sysListener).listenMPTCP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*ListenConfig).Listen (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.Listen (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*Server).ListenAndServe (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.ListenAndServe (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime.main (cum: 8.0, flat: 8.0, depth: 1)
compress/flate.(*compressor).init (cum: 2.0, flat: 2.0, depth: 0)
  ├── compress/flate.NewWriter (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.(*profileBuilder).build (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.printCountCycleProfile (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.writeProfileInternal (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.writeMutex (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.Index (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 2.0, flat: 2.0, depth: 1)
runtime/pprof.writeGoroutineStacks (cum: 2.0, flat: 2.0, depth: 0)
  ├── runtime/pprof.writeGoroutine (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.Index (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 2.0, flat: 2.0, depth: 1)
runtime/pprof.StartCPUProfile (cum: 2.0, flat: 2.0, depth: 0)
  ├── net/http/pprof.Profile (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 2.0, flat: 2.0, depth: 1)
compress/flate.NewWriter (cum: 2.0, flat: 2.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).build (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.writeAlloc (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.Index (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 2.0, flat: 2.0, depth: 1)
```

**Callgraph Statistics**: 126 nodes, max depth 1

## Mutex: Top mutex hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: mutex
  - Artifact: out/mutex.pb.gz

## Block: Top block hotspots

- **Severity**: Low
- **Score**: 50
- **Evidence**:
  - Profile: block
  - Artifact: out/block.pb.gz

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 30004.00 | 30004.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 30004.00 | 30004.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 30004.00 | 30004.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 30004.00 | 30004.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 30004.00 | 30004.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 21845.00 | 21845.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 21845.00 | 21845.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 21845.00 | 21845.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 21845.00 | 21845.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 8192.00 | 8192.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 8192.00 | 8192.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 8192.00 | 8192.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 8192.00 | 8192.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 8192.00 | 8192.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 8192.00 | 8192.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 8192.00 | 8192.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 8192.00 | 8192.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 8192.00 | 8192.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 8192.00 | 8192.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 8192.00 | 8192.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 8192.00 | 8192.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 5461.00 | 5461.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 5461.00 | 5461.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 5461.00 | 5461.00 |
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
| runtime/pprof.(*profileBuilder).emitLocation | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 622 | 293.00 | 293.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 467 | 293.00 | 293.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 47 | 293.00 | 293.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 293.00 | 293.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 293.00 | 293.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 293.00 | 293.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 293.00 | 293.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 293.00 | 293.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 293.00 | 293.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 293.00 | 293.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 293.00 | 293.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 293.00 | 293.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.goschedImpl | /usr/lib/go-1.24/src/runtime/proc.go | 4235 | 256.00 | 256.00 |
| runtime.gopreempt_m | /usr/lib/go-1.24/src/runtime/proc.go | 4252 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8.00 | 8.00 |
| net.open | /usr/lib/go-1.24/src/net/parse.go | 80 | 8.00 | 8.00 |
| net.maxListenerBacklog | /usr/lib/go-1.24/src/net/sock_linux.go | 35 | 8.00 | 8.00 |
| net.listenerBacklog.func1 | /usr/lib/go-1.24/src/net/net.go | 400 | 8.00 | 8.00 |
| sync.(*Once).doSlow | /usr/lib/go-1.24/src/sync/once.go | 78 | 8.00 | 8.00 |
| sync.(*Once).Do | /usr/lib/go-1.24/src/sync/once.go | 69 | 8.00 | 8.00 |
| net.listenerBacklog | /usr/lib/go-1.24/src/net/net.go | 400 | 8.00 | 8.00 |
| net.socket | /usr/lib/go-1.24/src/net/sock_posix.go | 57 | 8.00 | 8.00 |
| net.internetSocket | /usr/lib/go-1.24/src/net/ipsock_posix.go | 167 | 8.00 | 8.00 |
| net.(*sysListener).listenTCPProto | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 189 | 8.00 | 8.00 |
| net.(*sysListener).listenMPTCP | /usr/lib/go-1.24/src/net/mptcpsock_linux.go | 79 | 8.00 | 8.00 |
| net.(*ListenConfig).Listen | /usr/lib/go-1.24/src/net/dial.go | 819 | 8.00 | 8.00 |
| net.Listen | /usr/lib/go-1.24/src/net/dial.go | 898 | 8.00 | 8.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3346 | 8.00 | 8.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 8.00 | 8.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 8.00 | 8.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 583 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 8.00 | 8.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 8.00 | 8.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 8.00 | 8.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 8.00 | 8.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 4.00 | 4.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 4.00 | 4.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 4.00 | 4.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 4.00 | 4.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 4.00 | 4.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 4.00 | 4.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 4.00 | 4.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 4.00 | 4.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 4.00 | 4.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 4.00 | 4.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 4.00 | 4.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 4.00 | 4.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 4.00 | 4.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 4.00 | 4.00 |
| runtime/pprof.StartCPUProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 855 | 3.00 | 3.00 |
| net/http/pprof.Profile | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 157 | 3.00 | 3.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 3.00 | 3.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 3.00 | 3.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 3.00 | 3.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 3.00 | 3.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 2.00 | 2.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 2.00 | 2.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 2.00 | 2.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 2.00 | 2.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 2.00 | 2.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 2.00 | 2.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 2.00 | 2.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 2.00 | 2.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 2.00 | 2.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 2.00 | 2.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 2.00 | 2.00 |

### Callgraph Analysis

```
net.IP.String (cum: 32768.0, flat: 32768.0, depth: 0)
  ├── net.ipEmptyString (cum: 32768.0, flat: 32768.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 32768.0, flat: 32768.0, depth: 1)
main.allocHeavyHandler (cum: 30004.0, flat: 30004.0, depth: 0)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 30004.0, flat: 30004.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 30004.0, flat: 30004.0, depth: 1)
internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All (cum: 21845.0, flat: 21845.0, depth: 0)
  ├── unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 (cum: 21845.0, flat: 21845.0, depth: 1)
  ├── unique.registerCleanup.func1 (cum: 21845.0, flat: 21845.0, depth: 1)
  ├── runtime.unique_runtime_registerUniqueMapCleanup.func2 (cum: 21845.0, flat: 21845.0, depth: 1)
syscall.anyToSockaddr (cum: 8192.0, flat: 8192.0, depth: 0)
  ├── syscall.Accept4 (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── internal/poll.accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── internal/poll.(*FD).Accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*netFD).accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*TCPListener).accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net.(*TCPListener).Accept (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.(*Server).Serve (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.(*Server).ListenAndServe (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── net/http.ListenAndServe (cum: 8192.0, flat: 8192.0, depth: 1)
  ├── runtime.main (cum: 8192.0, flat: 8192.0, depth: 1)
runtime/pprof.allFrames (cum: 5461.0, flat: 5461.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).appendLocsForStack (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.writeHeap (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http/pprof.Index (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 5461.0, flat: 5461.0, depth: 1)
internal/sync.runtime_SemacquireMutex (cum: 5461.0, flat: 5461.0, depth: 0)
  ├── internal/sync.(*Mutex).lockSlow (cum: 5461.0, flat: 5461.0, depth: 1)
  ├── internal/sync.(*Mutex).Lock (cum: 5461.0, flat: 5461.0, depth: 1)
runtime.allocm (cum: 2308.0, flat: 2308.0, depth: 0)
  ├── runtime.newm (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.startm (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.wakep (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.resetspinning (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.schedule (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart1 (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart0 (cum: 2308.0, flat: 2308.0, depth: 1)
  ├── runtime.mstart (cum: 2308.0, flat: 2308.0, depth: 1)
runtime/pprof.(*profileBuilder).emitLocation (cum: 293.0, flat: 293.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).appendLocsForStack (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.writeAlloc (cum: 293.0, flat: 293.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http/pprof.Index (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 293.0, flat: 293.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 293.0, flat: 293.0, depth: 1)
net.open (cum: 8.0, flat: 8.0, depth: 0)
  ├── net.maxListenerBacklog (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.listenerBacklog.func1 (cum: 8.0, flat: 8.0, depth: 1)
  ├── sync.(*Once).doSlow (cum: 8.0, flat: 8.0, depth: 1)
  ├── sync.(*Once).Do (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.socket (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.internetSocket (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*sysListener).listenTCPProto (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*sysListener).listenMPTCP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.(*ListenConfig).Listen (cum: 8.0, flat: 8.0, depth: 1)
  ├── net.Listen (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*Server).ListenAndServe (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.ListenAndServe (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime.main (cum: 8.0, flat: 8.0, depth: 1)
compress/flate.(*compressor).init (cum: 8.0, flat: 8.0, depth: 0)
  ├── compress/flate.NewWriter (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.(*profileBuilder).build (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.writeHeap (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http/pprof.Index (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 8.0, flat: 8.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 8.0, flat: 8.0, depth: 1)
compress/flate.newDeflateFast (cum: 8.0, flat: 8.0, depth: 0)
  ├── compress/flate.NewWriter (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.(*profileBuilder).build (cum: 8.0, flat: 8.0, depth: 1)
  ├── runtime/pprof.profileWriter (cum: 8.0, flat: 8.0, depth: 1)
runtime/pprof.StartCPUProfile (cum: 3.0, flat: 3.0, depth: 0)
  ├── net/http/pprof.Profile (cum: 3.0, flat: 3.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 3.0, flat: 3.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 3.0, flat: 3.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 3.0, flat: 3.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 3.0, flat: 3.0, depth: 1)
runtime/pprof.writeGoroutineStacks (cum: 2.0, flat: 2.0, depth: 0)
  ├── runtime/pprof.writeGoroutine (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http/pprof.Index (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 2.0, flat: 2.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 2.0, flat: 2.0, depth: 1)
compress/flate.NewWriter (cum: 2.0, flat: 2.0, depth: 0)
  ├── runtime/pprof.(*profileBuilder).build (cum: 2.0, flat: 2.0, depth: 1)
  ├── runtime/pprof.profileWriter (cum: 2.0, flat: 2.0, depth: 1)
sync.(*Pool).pinSlow (cum: 0.0, flat: 0.0, depth: 0)
  ├── sync.(*Pool).pin (cum: 0.0, flat: 0.0, depth: 1)
  ├── sync.(*Pool).Get (cum: 0.0, flat: 0.0, depth: 1)
  ├── fmt.newPrinter (cum: 0.0, flat: 0.0, depth: 1)
  ├── fmt.Fprintf (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.(*chunkWriter).Write (cum: 0.0, flat: 0.0, depth: 1)
  ├── bufio.(*Writer).Flush (cum: 0.0, flat: 0.0, depth: 1)
  ├── bufio.(*Writer).Write (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.(*response).write (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.(*response).Write (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*huffmanBitWriter).write (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*huffmanBitWriter).writeTokens (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*huffmanBitWriter).writeBlockDynamic (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*compressor).encSpeed (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*compressor).close (cum: 0.0, flat: 0.0, depth: 1)
  ├── compress/flate.(*Writer).Close (cum: 0.0, flat: 0.0, depth: 1)
  ├── runtime/pprof.(*profileBuilder).build (cum: 0.0, flat: 0.0, depth: 1)
  ├── runtime/pprof.writeHeapProto (cum: 0.0, flat: 0.0, depth: 1)
  ├── runtime/pprof.writeHeapInternal (cum: 0.0, flat: 0.0, depth: 1)
  ├── runtime/pprof.writeHeap (cum: 0.0, flat: 0.0, depth: 1)
  ├── runtime/pprof.(*Profile).WriteTo (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http/pprof.handler.ServeHTTP (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http/pprof.Index (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.HandlerFunc.ServeHTTP (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.(*ServeMux).ServeHTTP (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.serverHandler.ServeHTTP (cum: 0.0, flat: 0.0, depth: 1)
  ├── net/http.(*conn).serve (cum: 0.0, flat: 0.0, depth: 1)
```

**Callgraph Statistics**: 135 nodes, max depth 1

---

*Generated by triageprof*
