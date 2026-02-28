# Performance Triage Report

Generated: 2026-02-28T14:37:14+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully

## Cpu: Top cpu hotspots

- **Severity**: Medium
- **Score**: 50
- **Evidence**:
  - Profile: cpu
  - Artifact: out-test/cpu.pb.gz

## Heap: Top heap hotspots

- **Severity**: Medium
- **Score**: 90
- **Evidence**:
  - Profile: heap
  - Artifact: out-test/heap.pb.gz

### Top Hotspots

| Function | File | Line | Cumulative | Flat |
|----------|------|------|------------|------|
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
| net/textproto.readMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 591 | 1489.00 | 1489.00 |
| net/textproto.(*Reader).ReadMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 507 | 1489.00 | 1489.00 |
| net/http.readRequest | /usr/lib/go-1.24/src/net/http/request.go | 1133 | 1489.00 | 1489.00 |
| net/http.(*conn).readRequest | /usr/lib/go-1.24/src/net/http/server.go | 1048 | 1489.00 | 1489.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2027 | 1489.00 | 1489.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 513.00 | 513.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 513.00 | 513.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 513.00 | 513.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 513.00 | 513.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 513.00 | 513.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 513.00 | 513.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 513.00 | 513.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 513.00 | 513.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 513.00 | 513.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 4.00 | 4.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 4.00 | 4.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 4.00 | 4.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 4.00 | 4.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 4.00 | 4.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 4.00 | 4.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 4.00 | 4.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 4.00 | 4.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 4.00 | 4.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 4.00 | 4.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 4.00 | 4.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 4.00 | 4.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 4.00 | 4.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 4.00 | 4.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 2.00 | 2.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 2.00 | 2.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 2.00 | 2.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 2.00 | 2.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 2.00 | 2.00 |
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
| strconv.syntaxError | /usr/lib/go-1.24/src/strconv/atoi.go | 48 | 0.00 | 0.00 |
| strconv.ParseInt | /usr/lib/go-1.24/src/strconv/atoi.go | 201 | 0.00 | 0.00 |
| strconv.Atoi | /usr/lib/go-1.24/src/strconv/atoi.go | 272 | 0.00 | 0.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 261 | 0.00 | 0.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 0.00 | 0.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 0.00 | 0.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 0.00 | 0.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 0.00 | 0.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 0.00 | 0.00 |
| runtime/pprof.writeGoroutineStacks | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 762 | 0.00 | 0.00 |
| runtime/pprof.writeGoroutine | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 753 | 0.00 | 0.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 0.00 | 0.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 0.00 | 0.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 0.00 | 0.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 0.00 | 0.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 0.00 | 0.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 0.00 | 0.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 0.00 | 0.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 0.00 | 0.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 0.00 | 0.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 0.00 | 0.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 0.00 | 0.00 |

## Mutex: Top mutex hotspots

- **Severity**: Medium
- **Score**: 50
- **Evidence**:
  - Profile: mutex
  - Artifact: out-test/mutex.pb.gz

## Block: Top block hotspots

- **Severity**: Medium
- **Score**: 50
- **Evidence**:
  - Profile: block
  - Artifact: out-test/block.pb.gz

---

*Generated by triageprof*
