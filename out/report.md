# Performance Triage Report

Generated: 2026-02-28T18:34:57+01:00

## Executive Summary

- **Overall Score**: 75/100
- **Top Issues**: performance
- **Notes**:
  - Analysis completed successfully

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 50126.00 | 50126.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 50126.00 | 50126.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 50126.00 | 50126.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 50126.00 | 50126.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 50126.00 | 50126.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 21845.00 | 21845.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 21845.00 | 21845.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 21845.00 | 21845.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 21845.00 | 21845.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 16385.00 | 16385.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 16385.00 | 16385.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 16385.00 | 16385.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 16385.00 | 16385.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 16385.00 | 16385.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 16385.00 | 16385.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 16385.00 | 16385.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 16385.00 | 16385.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 16385.00 | 16385.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 16385.00 | 16385.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 16385.00 | 16385.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16385.00 | 16385.00 |
| strconv.syntaxError | /usr/lib/go-1.24/src/strconv/atoi.go | 48 | 10923.00 | 10923.00 |
| strconv.ParseInt | /usr/lib/go-1.24/src/strconv/atoi.go | 201 | 10923.00 | 10923.00 |
| strconv.Atoi | /usr/lib/go-1.24/src/strconv/atoi.go | 272 | 10923.00 | 10923.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 261 | 10923.00 | 10923.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 10923.00 | 10923.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10923.00 | 10923.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10923.00 | 10923.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10923.00 | 10923.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10923.00 | 10923.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 10923.00 | 10923.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 10923.00 | 10923.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 10923.00 | 10923.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 10923.00 | 10923.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 10923.00 | 10923.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 1539.00 | 1539.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 1539.00 | 1539.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 1539.00 | 1539.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 1539.00 | 1539.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 1539.00 | 1539.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 1539.00 | 1539.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 1539.00 | 1539.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 1539.00 | 1539.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 1539.00 | 1539.00 |
| net/textproto.readMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 591 | 1489.00 | 1489.00 |
| net/textproto.(*Reader).ReadMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 507 | 1489.00 | 1489.00 |
| net/http.readRequest | /usr/lib/go-1.24/src/net/http/request.go | 1133 | 1489.00 | 1489.00 |
| net/http.(*conn).readRequest | /usr/lib/go-1.24/src/net/http/server.go | 1048 | 1489.00 | 1489.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2027 | 1489.00 | 1489.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 769.00 | 769.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 769.00 | 769.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 769.00 | 769.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 769.00 | 769.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 769.00 | 769.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 769.00 | 769.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 769.00 | 769.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 769.00 | 769.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.ready | /usr/lib/go-1.24/src/runtime/proc.go | 1080 | 256.00 | 256.00 |
| runtime.readyWithTime.goready.func1 | /usr/lib/go-1.24/src/runtime/proc.go | 456 | 256.00 | 256.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 256.00 | 256.00 |
| compress/flate.(*huffmanEncoder).generate | /usr/lib/go-1.24/src/compress/flate/huffman_code.go | 277 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).indexTokens | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 561 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).writeBlockDynamic | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 504 | 228.00 | 228.00 |
| compress/flate.(*compressor).encSpeed | /usr/lib/go-1.24/src/compress/flate/deflate.go | 363 | 228.00 | 228.00 |
| compress/flate.(*compressor).close | /usr/lib/go-1.24/src/compress/flate/deflate.go | 635 | 228.00 | 228.00 |
| compress/flate.(*Writer).Close | /usr/lib/go-1.24/src/compress/flate/deflate.go | 727 | 228.00 | 228.00 |
| compress/gzip.(*Writer).Close | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 242 | 228.00 | 228.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 391 | 228.00 | 228.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 228.00 | 228.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 228.00 | 228.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 228.00 | 228.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 228.00 | 228.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 228.00 | 228.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 228.00 | 228.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 228.00 | 228.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 228.00 | 228.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 228.00 | 228.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 228.00 | 228.00 |
| bufio.NewWriterSize | /usr/lib/go-1.24/src/bufio/bufio.go | 600 | 128.00 | 128.00 |
| net/http.newBufioWriterSize | /usr/lib/go-1.24/src/net/http/server.go | 894 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2014 | 128.00 | 128.00 |
| os.readFileContents | /usr/lib/go-1.24/src/os/file.go | 838 | 128.00 | 128.00 |
| os.ReadFile | /usr/lib/go-1.24/src/os/file.go | 805 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).readMapping | /usr/lib/go-1.24/src/runtime/pprof/proto_other.go | 18 | 128.00 | 128.00 |
| runtime/pprof.newProfileBuilder | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 270 | 128.00 | 128.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 435 | 128.00 | 128.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 128.00 | 128.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 128.00 | 128.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 128.00 | 128.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 128.00 | 128.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 128.00 | 128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 128.00 | 128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 128.00 | 128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 128.00 | 128.00 |
| os.readFileContents | /usr/lib/go-1.24/src/os/file.go | 838 | 128.00 | 128.00 |
| os.ReadFile | /usr/lib/go-1.24/src/os/file.go | 805 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).readMapping | /usr/lib/go-1.24/src/runtime/pprof/proto_other.go | 18 | 128.00 | 128.00 |
| runtime/pprof.newProfileBuilder | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 270 | 128.00 | 128.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 17 | 128.00 | 128.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 128.00 | 128.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 128.00 | 128.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 128.00 | 128.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 128.00 | 128.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 128.00 | 128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 128.00 | 128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 128.00 | 128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).emitLocation | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 606 | 108.00 | 108.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 467 | 108.00 | 108.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 47 | 108.00 | 108.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 108.00 | 108.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 108.00 | 108.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 108.00 | 108.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 108.00 | 108.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 108.00 | 108.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 108.00 | 108.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 108.00 | 108.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 108.00 | 108.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 108.00 | 108.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 12.00 | 12.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 12.00 | 12.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 12.00 | 12.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 12.00 | 12.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 12.00 | 12.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 10.00 | 10.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 10.00 | 10.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 10.00 | 10.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 10.00 | 10.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 10.00 | 10.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 10.00 | 10.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 10.00 | 10.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 10.00 | 10.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 10.00 | 10.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 10.00 | 10.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10.00 | 10.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10.00 | 10.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10.00 | 10.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10.00 | 10.00 |
| runtime/pprof.StartCPUProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 855 | 9.00 | 9.00 |
| net/http/pprof.Profile | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 157 | 9.00 | 9.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 9.00 | 9.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 9.00 | 9.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 9.00 | 9.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 9.00 | 9.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 8.00 | 8.00 |
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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 50128.00 | 50128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 50128.00 | 50128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 50128.00 | 50128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 50128.00 | 50128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 50128.00 | 50128.00 |
| internal/sync.(*HashTrieMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string },go.shape.struct { weak._ [0]*go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }; weak.u unsafe.Pointer }]).All | /usr/lib/go-1.24/src/internal/sync/hashtriemap.go | 483 | 21845.00 | 21845.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 21845.00 | 21845.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 21845.00 | 21845.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 21845.00 | 21845.00 |
| syscall.anyToSockaddr | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 675 | 16385.00 | 16385.00 |
| syscall.Accept4 | /usr/lib/go-1.24/src/syscall/syscall_linux.go | 699 | 16385.00 | 16385.00 |
| internal/poll.accept | /usr/lib/go-1.24/src/internal/poll/sock_cloexec.go | 17 | 16385.00 | 16385.00 |
| internal/poll.(*FD).Accept | /usr/lib/go-1.24/src/internal/poll/fd_unix.go | 611 | 16385.00 | 16385.00 |
| net.(*netFD).accept | /usr/lib/go-1.24/src/net/fd_unix.go | 172 | 16385.00 | 16385.00 |
| net.(*TCPListener).accept | /usr/lib/go-1.24/src/net/tcpsock_posix.go | 159 | 16385.00 | 16385.00 |
| net.(*TCPListener).Accept | /usr/lib/go-1.24/src/net/tcpsock.go | 380 | 16385.00 | 16385.00 |
| net/http.(*Server).Serve | /usr/lib/go-1.24/src/net/http/server.go | 3424 | 16385.00 | 16385.00 |
| net/http.(*Server).ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3350 | 16385.00 | 16385.00 |
| net/http.ListenAndServe | /usr/lib/go-1.24/src/net/http/server.go | 3665 | 16385.00 | 16385.00 |
| main.main | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 28 | 16385.00 | 16385.00 |
| runtime.main | /usr/lib/go-1.24/src/runtime/proc.go | 283 | 16385.00 | 16385.00 |
| strconv.syntaxError | /usr/lib/go-1.24/src/strconv/atoi.go | 48 | 10923.00 | 10923.00 |
| strconv.ParseInt | /usr/lib/go-1.24/src/strconv/atoi.go | 201 | 10923.00 | 10923.00 |
| strconv.Atoi | /usr/lib/go-1.24/src/strconv/atoi.go | 272 | 10923.00 | 10923.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 261 | 10923.00 | 10923.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 10923.00 | 10923.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10923.00 | 10923.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10923.00 | 10923.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10923.00 | 10923.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10923.00 | 10923.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 10923.00 | 10923.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 10923.00 | 10923.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 10923.00 | 10923.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 10923.00 | 10923.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 10923.00 | 10923.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 1539.00 | 1539.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 1539.00 | 1539.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 1539.00 | 1539.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 1539.00 | 1539.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 1539.00 | 1539.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 1539.00 | 1539.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 1539.00 | 1539.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 1539.00 | 1539.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 1539.00 | 1539.00 |
| net/textproto.readMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 591 | 1489.00 | 1489.00 |
| net/textproto.(*Reader).ReadMIMEHeader | /usr/lib/go-1.24/src/net/textproto/reader.go | 507 | 1489.00 | 1489.00 |
| net/http.readRequest | /usr/lib/go-1.24/src/net/http/request.go | 1133 | 1489.00 | 1489.00 |
| net/http.(*conn).readRequest | /usr/lib/go-1.24/src/net/http/server.go | 1048 | 1489.00 | 1489.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2027 | 1489.00 | 1489.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 769.00 | 769.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 769.00 | 769.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 769.00 | 769.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 769.00 | 769.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 769.00 | 769.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 769.00 | 769.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 769.00 | 769.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 769.00 | 769.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.ready | /usr/lib/go-1.24/src/runtime/proc.go | 1080 | 256.00 | 256.00 |
| runtime.readyWithTime.goready.func1 | /usr/lib/go-1.24/src/runtime/proc.go | 456 | 256.00 | 256.00 |
| runtime.systemstack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 514 | 256.00 | 256.00 |
| compress/flate.(*huffmanEncoder).generate | /usr/lib/go-1.24/src/compress/flate/huffman_code.go | 277 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).indexTokens | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 561 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).writeBlockDynamic | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 504 | 228.00 | 228.00 |
| compress/flate.(*compressor).encSpeed | /usr/lib/go-1.24/src/compress/flate/deflate.go | 363 | 228.00 | 228.00 |
| compress/flate.(*compressor).close | /usr/lib/go-1.24/src/compress/flate/deflate.go | 635 | 228.00 | 228.00 |
| compress/flate.(*Writer).Close | /usr/lib/go-1.24/src/compress/flate/deflate.go | 727 | 228.00 | 228.00 |
| compress/gzip.(*Writer).Close | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 242 | 228.00 | 228.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 391 | 228.00 | 228.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 228.00 | 228.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 228.00 | 228.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 228.00 | 228.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 228.00 | 228.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 228.00 | 228.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 228.00 | 228.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 228.00 | 228.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 228.00 | 228.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 228.00 | 228.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 228.00 | 228.00 |
| bufio.NewWriterSize | /usr/lib/go-1.24/src/bufio/bufio.go | 600 | 128.00 | 128.00 |
| net/http.newBufioWriterSize | /usr/lib/go-1.24/src/net/http/server.go | 894 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2014 | 128.00 | 128.00 |
| os.readFileContents | /usr/lib/go-1.24/src/os/file.go | 838 | 128.00 | 128.00 |
| os.ReadFile | /usr/lib/go-1.24/src/os/file.go | 805 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).readMapping | /usr/lib/go-1.24/src/runtime/pprof/proto_other.go | 18 | 128.00 | 128.00 |
| runtime/pprof.newProfileBuilder | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 270 | 128.00 | 128.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 435 | 128.00 | 128.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 128.00 | 128.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 128.00 | 128.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 128.00 | 128.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 128.00 | 128.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 128.00 | 128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 128.00 | 128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 128.00 | 128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 128.00 | 128.00 |
| os.readFileContents | /usr/lib/go-1.24/src/os/file.go | 838 | 128.00 | 128.00 |
| os.ReadFile | /usr/lib/go-1.24/src/os/file.go | 805 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).readMapping | /usr/lib/go-1.24/src/runtime/pprof/proto_other.go | 18 | 128.00 | 128.00 |
| runtime/pprof.newProfileBuilder | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 270 | 128.00 | 128.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 17 | 128.00 | 128.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 128.00 | 128.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 128.00 | 128.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 128.00 | 128.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 128.00 | 128.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 128.00 | 128.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 128.00 | 128.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 128.00 | 128.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 128.00 | 128.00 |
| runtime/pprof.(*profileBuilder).emitLocation | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 606 | 108.00 | 108.00 |
| runtime/pprof.(*profileBuilder).appendLocsForStack | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 467 | 108.00 | 108.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 47 | 108.00 | 108.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 108.00 | 108.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 108.00 | 108.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 108.00 | 108.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 108.00 | 108.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 108.00 | 108.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 108.00 | 108.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 108.00 | 108.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 108.00 | 108.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 108.00 | 108.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 12.00 | 12.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 12.00 | 12.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 12.00 | 12.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 12.00 | 12.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 12.00 | 12.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 12.00 | 12.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 12.00 | 12.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 12.00 | 12.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 12.00 | 12.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 12.00 | 12.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 12.00 | 12.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 12.00 | 12.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 12.00 | 12.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 12.00 | 12.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 12.00 | 12.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 12.00 | 12.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 12.00 | 12.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 12.00 | 12.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 12.00 | 12.00 |
| runtime/pprof.StartCPUProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 855 | 10.00 | 10.00 |
| net/http/pprof.Profile | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 157 | 10.00 | 10.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10.00 | 10.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10.00 | 10.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10.00 | 10.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10.00 | 10.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 9.00 | 9.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 9.00 | 9.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 9.00 | 9.00 |
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 9.00 | 9.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 9.00 | 9.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 9.00 | 9.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 9.00 | 9.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 9.00 | 9.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 9.00 | 9.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 9.00 | 9.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 9.00 | 9.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 9.00 | 9.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 9.00 | 9.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 9.00 | 9.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 9.00 | 9.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 9.00 | 9.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 9.00 | 9.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 663 | 9.00 | 9.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 9.00 | 9.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 9.00 | 9.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 9.00 | 9.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 9.00 | 9.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 9.00 | 9.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 9.00 | 9.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 9.00 | 9.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 9.00 | 9.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 9.00 | 9.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 9.00 | 9.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 9.00 | 9.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 9.00 | 9.00 |

---

*Generated by triageprof*
