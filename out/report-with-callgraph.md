# Performance Triage Report

Generated: 2026-02-28T19:16:03+01:00

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 20003.00 | 20003.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 20003.00 | 20003.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 20003.00 | 20003.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 20003.00 | 20003.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 20003.00 | 20003.00 |
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
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 1.00 | 1.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 1.00 | 1.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 1.00 | 1.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 1.00 | 1.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| runtime/pprof.StartCPUProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 855 | 1.00 | 1.00 |
| net/http/pprof.Profile | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 157 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 1.00 | 1.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 1.00 | 1.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 1.00 | 1.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 1.00 | 1.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| runtime/pprof.writeGoroutineStacks | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 762 | 1.00 | 1.00 |
| runtime/pprof.writeGoroutine | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 753 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 0.00 | 0.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 0.00 | 0.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 0.00 | 0.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 0.00 | 0.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 0.00 | 0.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 0.00 | 0.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 0.00 | 0.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 0.00 | 0.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 0.00 | 0.00 |

### Callgraph Analysis (Depth 3)

```
net.IP.String (32768.0% cum, 32768.0% flat)
  net.ipEmptyString (32768.0% cum, 32768.0% flat)
  net/http.(*conn).serve (32768.0% cum, 32768.0% flat)
main.allocHeavyHandler (20003.0% cum, 20003.0% flat)
  net/http.HandlerFunc.ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.(*ServeMux).ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.serverHandler.ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.(*conn).serve (20003.0% cum, 20003.0% flat)
syscall.anyToSockaddr (8192.0% cum, 8192.0% flat)
  syscall.Accept4 (8192.0% cum, 8192.0% flat)
  internal/poll.accept (8192.0% cum, 8192.0% flat)
  internal/poll.(*FD).Accept (8192.0% cum, 8192.0% flat)
  net.(*netFD).accept (8192.0% cum, 8192.0% flat)
  net.(*TCPListener).accept (8192.0% cum, 8192.0% flat)
  net.(*TCPListener).Accept (8192.0% cum, 8192.0% flat)
  net/http.(*Server).Serve (8192.0% cum, 8192.0% flat)
  net/http.(*Server).ListenAndServe (8192.0% cum, 8192.0% flat)
  net/http.ListenAndServe (8192.0% cum, 8192.0% flat)
  runtime.main (8192.0% cum, 8192.0% flat)
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
runtime/pprof.(*profileBuilder).emitLocation (293.0% cum, 293.0% flat)
  runtime/pprof.(*profileBuilder).appendLocsForStack (293.0% cum, 293.0% flat)
  runtime/pprof.writeHeapProto (293.0% cum, 293.0% flat)
  runtime/pprof.writeHeapInternal (293.0% cum, 293.0% flat)
  runtime/pprof.writeAlloc (293.0% cum, 293.0% flat)
  runtime/pprof.(*Profile).WriteTo (293.0% cum, 293.0% flat)
  net/http/pprof.handler.ServeHTTP (293.0% cum, 293.0% flat)
  net/http/pprof.Index (293.0% cum, 293.0% flat)
  net/http.HandlerFunc.ServeHTTP (293.0% cum, 293.0% flat)
  net/http.(*ServeMux).ServeHTTP (293.0% cum, 293.0% flat)
  net/http.serverHandler.ServeHTTP (293.0% cum, 293.0% flat)
  net/http.(*conn).serve (293.0% cum, 293.0% flat)
compress/flate.newDeflateFast (8.0% cum, 8.0% flat)
  compress/flate.NewWriter (8.0% cum, 8.0% flat)
  runtime/pprof.(*profileBuilder).build (8.0% cum, 8.0% flat)
  runtime/pprof.printCountCycleProfile (8.0% cum, 8.0% flat)
  runtime/pprof.writeProfileInternal (8.0% cum, 8.0% flat)
  runtime/pprof.writeBlock (8.0% cum, 8.0% flat)
  runtime/pprof.(*Profile).WriteTo (8.0% cum, 8.0% flat)
  net/http/pprof.handler.ServeHTTP (8.0% cum, 8.0% flat)
  net/http/pprof.Index (8.0% cum, 8.0% flat)
  net/http.HandlerFunc.ServeHTTP (8.0% cum, 8.0% flat)
  net/http.(*ServeMux).ServeHTTP (8.0% cum, 8.0% flat)
  net/http.serverHandler.ServeHTTP (8.0% cum, 8.0% flat)
  net/http.(*conn).serve (8.0% cum, 8.0% flat)
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
compress/flate.NewWriter (1.0% cum, 1.0% flat)
  runtime/pprof.(*profileBuilder).build (1.0% cum, 1.0% flat)
  runtime/pprof.writeHeapProto (1.0% cum, 1.0% flat)
  runtime/pprof.writeHeapInternal (1.0% cum, 1.0% flat)
  runtime/pprof.writeAlloc (1.0% cum, 1.0% flat)
  runtime/pprof.(*Profile).WriteTo (1.0% cum, 1.0% flat)
  net/http/pprof.handler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http/pprof.Index (1.0% cum, 1.0% flat)
  net/http.HandlerFunc.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*ServeMux).ServeHTTP (1.0% cum, 1.0% flat)
  net/http.serverHandler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*conn).serve (1.0% cum, 1.0% flat)
runtime/pprof.StartCPUProfile (1.0% cum, 1.0% flat)
  net/http/pprof.Profile (1.0% cum, 1.0% flat)
  net/http.HandlerFunc.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*ServeMux).ServeHTTP (1.0% cum, 1.0% flat)
  net/http.serverHandler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*conn).serve (1.0% cum, 1.0% flat)
runtime/pprof.writeGoroutineStacks (1.0% cum, 1.0% flat)
  runtime/pprof.writeGoroutine (1.0% cum, 1.0% flat)
  runtime/pprof.(*Profile).WriteTo (1.0% cum, 1.0% flat)
  net/http/pprof.handler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http/pprof.Index (1.0% cum, 1.0% flat)
  net/http.HandlerFunc.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*ServeMux).ServeHTTP (1.0% cum, 1.0% flat)
  net/http.serverHandler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*conn).serve (1.0% cum, 1.0% flat)
compress/flate.(*compressor).init (0.0% cum, 0.0% flat)
  compress/flate.NewWriter (0.0% cum, 0.0% flat)
  runtime/pprof.(*profileBuilder).build (0.0% cum, 0.0% flat)
  runtime/pprof.profileWriter (0.0% cum, 0.0% flat)
```

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 20003.00 | 20003.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 20003.00 | 20003.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 20003.00 | 20003.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 20003.00 | 20003.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 20003.00 | 20003.00 |
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
| runtime/pprof.StartCPUProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 855 | 2.00 | 2.00 |
| net/http/pprof.Profile | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 157 | 2.00 | 2.00 |
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
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 1.00 | 1.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 1.00 | 1.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 1.00 | 1.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 1.00 | 1.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 1.00 | 1.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 1.00 | 1.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 1.00 | 1.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 1.00 | 1.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 1.00 | 1.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |
| runtime/pprof.writeGoroutineStacks | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 762 | 1.00 | 1.00 |
| runtime/pprof.writeGoroutine | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 753 | 1.00 | 1.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 1.00 | 1.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 1.00 | 1.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 1.00 | 1.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 1.00 | 1.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 1.00 | 1.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 1.00 | 1.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 1.00 | 1.00 |

### Callgraph Analysis (Depth 3)

```
net.IP.String (32768.0% cum, 32768.0% flat)
  net.ipEmptyString (32768.0% cum, 32768.0% flat)
  net/http.(*conn).serve (32768.0% cum, 32768.0% flat)
main.allocHeavyHandler (20003.0% cum, 20003.0% flat)
  net/http.HandlerFunc.ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.(*ServeMux).ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.serverHandler.ServeHTTP (20003.0% cum, 20003.0% flat)
  net/http.(*conn).serve (20003.0% cum, 20003.0% flat)
syscall.anyToSockaddr (8192.0% cum, 8192.0% flat)
  syscall.Accept4 (8192.0% cum, 8192.0% flat)
  internal/poll.accept (8192.0% cum, 8192.0% flat)
  internal/poll.(*FD).Accept (8192.0% cum, 8192.0% flat)
  net.(*netFD).accept (8192.0% cum, 8192.0% flat)
  net.(*TCPListener).accept (8192.0% cum, 8192.0% flat)
  net.(*TCPListener).Accept (8192.0% cum, 8192.0% flat)
  net/http.(*Server).Serve (8192.0% cum, 8192.0% flat)
  net/http.(*Server).ListenAndServe (8192.0% cum, 8192.0% flat)
  net/http.ListenAndServe (8192.0% cum, 8192.0% flat)
  runtime.main (8192.0% cum, 8192.0% flat)
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
runtime/pprof.(*profileBuilder).emitLocation (293.0% cum, 293.0% flat)
  runtime/pprof.(*profileBuilder).appendLocsForStack (293.0% cum, 293.0% flat)
  runtime/pprof.writeHeapProto (293.0% cum, 293.0% flat)
  runtime/pprof.writeHeapInternal (293.0% cum, 293.0% flat)
  runtime/pprof.writeAlloc (293.0% cum, 293.0% flat)
  runtime/pprof.(*Profile).WriteTo (293.0% cum, 293.0% flat)
  net/http/pprof.handler.ServeHTTP (293.0% cum, 293.0% flat)
  net/http/pprof.Index (293.0% cum, 293.0% flat)
  net/http.HandlerFunc.ServeHTTP (293.0% cum, 293.0% flat)
  net/http.(*ServeMux).ServeHTTP (293.0% cum, 293.0% flat)
  net/http.serverHandler.ServeHTTP (293.0% cum, 293.0% flat)
  net/http.(*conn).serve (293.0% cum, 293.0% flat)
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
compress/flate.newDeflateFast (8.0% cum, 8.0% flat)
  compress/flate.NewWriter (8.0% cum, 8.0% flat)
  runtime/pprof.(*profileBuilder).build (8.0% cum, 8.0% flat)
  runtime/pprof.printCountCycleProfile (8.0% cum, 8.0% flat)
  runtime/pprof.writeProfileInternal (8.0% cum, 8.0% flat)
  runtime/pprof.writeBlock (8.0% cum, 8.0% flat)
  runtime/pprof.(*Profile).WriteTo (8.0% cum, 8.0% flat)
  net/http/pprof.handler.ServeHTTP (8.0% cum, 8.0% flat)
  net/http/pprof.Index (8.0% cum, 8.0% flat)
  net/http.HandlerFunc.ServeHTTP (8.0% cum, 8.0% flat)
  net/http.(*ServeMux).ServeHTTP (8.0% cum, 8.0% flat)
  net/http.serverHandler.ServeHTTP (8.0% cum, 8.0% flat)
  net/http.(*conn).serve (8.0% cum, 8.0% flat)
runtime/pprof.StartCPUProfile (2.0% cum, 2.0% flat)
  net/http/pprof.Profile (2.0% cum, 2.0% flat)
  net/http.HandlerFunc.ServeHTTP (2.0% cum, 2.0% flat)
  net/http.(*ServeMux).ServeHTTP (2.0% cum, 2.0% flat)
  net/http.serverHandler.ServeHTTP (2.0% cum, 2.0% flat)
  net/http.(*conn).serve (2.0% cum, 2.0% flat)
compress/flate.(*compressor).init (2.0% cum, 2.0% flat)
  compress/flate.NewWriter (2.0% cum, 2.0% flat)
  runtime/pprof.(*profileBuilder).build (2.0% cum, 2.0% flat)
  runtime/pprof.printCountCycleProfile (2.0% cum, 2.0% flat)
  runtime/pprof.writeProfileInternal (2.0% cum, 2.0% flat)
  runtime/pprof.writeMutex (2.0% cum, 2.0% flat)
  runtime/pprof.(*Profile).WriteTo (2.0% cum, 2.0% flat)
  net/http/pprof.handler.ServeHTTP (2.0% cum, 2.0% flat)
  net/http/pprof.Index (2.0% cum, 2.0% flat)
  net/http.HandlerFunc.ServeHTTP (2.0% cum, 2.0% flat)
  net/http.(*ServeMux).ServeHTTP (2.0% cum, 2.0% flat)
  net/http.serverHandler.ServeHTTP (2.0% cum, 2.0% flat)
  net/http.(*conn).serve (2.0% cum, 2.0% flat)
compress/flate.NewWriter (1.0% cum, 1.0% flat)
  runtime/pprof.(*profileBuilder).build (1.0% cum, 1.0% flat)
  runtime/pprof.profileWriter (1.0% cum, 1.0% flat)
runtime/pprof.writeGoroutineStacks (1.0% cum, 1.0% flat)
  runtime/pprof.writeGoroutine (1.0% cum, 1.0% flat)
  runtime/pprof.(*Profile).WriteTo (1.0% cum, 1.0% flat)
  net/http/pprof.handler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http/pprof.Index (1.0% cum, 1.0% flat)
  net/http.HandlerFunc.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*ServeMux).ServeHTTP (1.0% cum, 1.0% flat)
  net/http.serverHandler.ServeHTTP (1.0% cum, 1.0% flat)
  net/http.(*conn).serve (1.0% cum, 1.0% flat)
```

---

*Generated by triageprof*
