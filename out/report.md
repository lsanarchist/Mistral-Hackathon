# Performance Triage Report

Generated: 2026-02-28T23:29:35+01:00

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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 50003.00 | 50003.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 50003.00 | 50003.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 50003.00 | 50003.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 50003.00 | 50003.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 50003.00 | 50003.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| net/http.init.func14 | /usr/lib/go-1.24/src/net/http/header.go | 161 | 21845.00 | 21845.00 |
| sync.(*Pool).Get | /usr/lib/go-1.24/src/sync/pool.go | 155 | 21845.00 | 21845.00 |
| net/http.Header.sortedKeyValues | /usr/lib/go-1.24/src/net/http/header.go | 168 | 21845.00 | 21845.00 |
| net/http.Header.writeSubset | /usr/lib/go-1.24/src/net/http/header.go | 195 | 21845.00 | 21845.00 |
| net/http.Header.WriteSubset | /usr/lib/go-1.24/src/net/http/header.go | 187 | 21845.00 | 21845.00 |
| net/http.(*chunkWriter).writeHeader | /usr/lib/go-1.24/src/net/http/server.go | 1577 | 21845.00 | 21845.00 |
| net/http.(*chunkWriter).Write | /usr/lib/go-1.24/src/net/http/server.go | 376 | 21845.00 | 21845.00 |
| bufio.(*Writer).Flush | /usr/lib/go-1.24/src/bufio/bufio.go | 643 | 21845.00 | 21845.00 |
| net/http.(*response).finishRequest | /usr/lib/go-1.24/src/net/http/server.go | 1715 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2108 | 21845.00 | 21845.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 5461.00 | 5461.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 5461.00 | 5461.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 5461.00 | 5461.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 2052.00 | 2052.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 2052.00 | 2052.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 2052.00 | 2052.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 2052.00 | 2052.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 2052.00 | 2052.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 2052.00 | 2052.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 2052.00 | 2052.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 2052.00 | 2052.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 2052.00 | 2052.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 513.00 | 513.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 513.00 | 513.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 513.00 | 513.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 513.00 | 513.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 513.00 | 513.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 513.00 | 513.00 |
| runtime.goschedImpl | /usr/lib/go-1.24/src/runtime/proc.go | 4235 | 513.00 | 513.00 |
| runtime.gopreempt_m | /usr/lib/go-1.24/src/runtime/proc.go | 4252 | 513.00 | 513.00 |
| runtime.newstack | /usr/lib/go-1.24/src/runtime/stack.go | 1074 | 513.00 | 513.00 |
| runtime.morestack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 621 | 513.00 | 513.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |
| compress/flate.(*huffmanEncoder).generate | /usr/lib/go-1.24/src/compress/flate/huffman_code.go | 277 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).writeBlockDynamic | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 509 | 228.00 | 228.00 |
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
| bufio.NewReaderSize | /usr/lib/go-1.24/src/bufio/bufio.go | 57 | 128.00 | 128.00 |
| bufio.NewReader | /usr/lib/go-1.24/src/bufio/bufio.go | 63 | 128.00 | 128.00 |
| net/http.newBufioReader | /usr/lib/go-1.24/src/net/http/server.go | 859 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2013 | 128.00 | 128.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 17.00 | 17.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 17.00 | 17.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 17.00 | 17.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 17.00 | 17.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 17.00 | 17.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 17.00 | 17.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 10.00 | 10.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 10.00 | 10.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 10.00 | 10.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 10.00 | 10.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 10.00 | 10.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 10.00 | 10.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 10.00 | 10.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 10.00 | 10.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 10.00 | 10.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 10.00 | 10.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10.00 | 10.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10.00 | 10.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10.00 | 10.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10.00 | 10.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 583 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 8.00 | 8.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 8.00 | 8.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 8.00 | 8.00 |
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
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
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
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 583 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 8.00 | 8.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 8.00 | 8.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
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
| runtime/pprof.profileWriter | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 885 | 8.00 | 8.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 8.00 | 8.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 8.00 | 8.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 8.00 | 8.00 |
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
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 50003.00 | 50003.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 50003.00 | 50003.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 50003.00 | 50003.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 50003.00 | 50003.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 50003.00 | 50003.00 |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691.00 | 43691.00 |
| unique.registerCleanup.func1 | /usr/lib/go-1.24/src/unique/handle.go | 162 | 43691.00 | 43691.00 |
| runtime.unique_runtime_registerUniqueMapCleanup.func2 | /usr/lib/go-1.24/src/runtime/mgc.go | 1798 | 43691.00 | 43691.00 |
| net/http.init.func14 | /usr/lib/go-1.24/src/net/http/header.go | 161 | 21845.00 | 21845.00 |
| sync.(*Pool).Get | /usr/lib/go-1.24/src/sync/pool.go | 155 | 21845.00 | 21845.00 |
| net/http.Header.sortedKeyValues | /usr/lib/go-1.24/src/net/http/header.go | 168 | 21845.00 | 21845.00 |
| net/http.Header.writeSubset | /usr/lib/go-1.24/src/net/http/header.go | 195 | 21845.00 | 21845.00 |
| net/http.Header.WriteSubset | /usr/lib/go-1.24/src/net/http/header.go | 187 | 21845.00 | 21845.00 |
| net/http.(*chunkWriter).writeHeader | /usr/lib/go-1.24/src/net/http/server.go | 1577 | 21845.00 | 21845.00 |
| net/http.(*chunkWriter).Write | /usr/lib/go-1.24/src/net/http/server.go | 376 | 21845.00 | 21845.00 |
| bufio.(*Writer).Flush | /usr/lib/go-1.24/src/bufio/bufio.go | 643 | 21845.00 | 21845.00 |
| net/http.(*response).finishRequest | /usr/lib/go-1.24/src/net/http/server.go | 1715 | 21845.00 | 21845.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2108 | 21845.00 | 21845.00 |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).lockSlow | /usr/lib/go-1.24/src/internal/sync/mutex.go | 149 | 5461.00 | 5461.00 |
| internal/sync.(*Mutex).Lock | /usr/lib/go-1.24/src/internal/sync/mutex.go | 70 | 5461.00 | 5461.00 |
| sync.(*Mutex).Lock | /usr/lib/go-1.24/src/sync/mutex.go | 46 | 5461.00 | 5461.00 |
| main.mutexContentionHandler.func1 | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 57 | 5461.00 | 5461.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 2052.00 | 2052.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 2052.00 | 2052.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 2052.00 | 2052.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 2052.00 | 2052.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 2052.00 | 2052.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 2052.00 | 2052.00 |
| runtime.mstart1 | /usr/lib/go-1.24/src/runtime/proc.go | 1894 | 2052.00 | 2052.00 |
| runtime.mstart0 | /usr/lib/go-1.24/src/runtime/proc.go | 1840 | 2052.00 | 2052.00 |
| runtime.mstart | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 395 | 2052.00 | 2052.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 513.00 | 513.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 513.00 | 513.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 513.00 | 513.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 513.00 | 513.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 513.00 | 513.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 513.00 | 513.00 |
| runtime.goschedImpl | /usr/lib/go-1.24/src/runtime/proc.go | 4235 | 513.00 | 513.00 |
| runtime.gopreempt_m | /usr/lib/go-1.24/src/runtime/proc.go | 4252 | 513.00 | 513.00 |
| runtime.newstack | /usr/lib/go-1.24/src/runtime/stack.go | 1074 | 513.00 | 513.00 |
| runtime.morestack | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 621 | 513.00 | 513.00 |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 256.00 | 256.00 |
| runtime.newm | /usr/lib/go-1.24/src/runtime/proc.go | 2812 | 256.00 | 256.00 |
| runtime.startm | /usr/lib/go-1.24/src/runtime/proc.go | 3038 | 256.00 | 256.00 |
| runtime.wakep | /usr/lib/go-1.24/src/runtime/proc.go | 3185 | 256.00 | 256.00 |
| runtime.resetspinning | /usr/lib/go-1.24/src/runtime/proc.go | 3937 | 256.00 | 256.00 |
| runtime.schedule | /usr/lib/go-1.24/src/runtime/proc.go | 4095 | 256.00 | 256.00 |
| runtime.park_m | /usr/lib/go-1.24/src/runtime/proc.go | 4201 | 256.00 | 256.00 |
| runtime.mcall | /usr/lib/go-1.24/src/runtime/asm_amd64.s | 459 | 256.00 | 256.00 |
| compress/flate.(*huffmanEncoder).generate | /usr/lib/go-1.24/src/compress/flate/huffman_code.go | 277 | 228.00 | 228.00 |
| compress/flate.(*huffmanBitWriter).writeBlockDynamic | /usr/lib/go-1.24/src/compress/flate/huffman_bit_writer.go | 509 | 228.00 | 228.00 |
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
| os.readFileContents | /usr/lib/go-1.24/src/os/file.go | 838 | 171.00 | 171.00 |
| os.ReadFile | /usr/lib/go-1.24/src/os/file.go | 805 | 171.00 | 171.00 |
| runtime/pprof.(*profileBuilder).readMapping | /usr/lib/go-1.24/src/runtime/pprof/proto_other.go | 18 | 171.00 | 171.00 |
| runtime/pprof.newProfileBuilder | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 270 | 171.00 | 171.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 435 | 171.00 | 171.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 171.00 | 171.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 171.00 | 171.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 171.00 | 171.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 171.00 | 171.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 171.00 | 171.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 171.00 | 171.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 171.00 | 171.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 171.00 | 171.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 171.00 | 171.00 |
| bufio.NewReaderSize | /usr/lib/go-1.24/src/bufio/bufio.go | 57 | 128.00 | 128.00 |
| bufio.NewReader | /usr/lib/go-1.24/src/bufio/bufio.go | 63 | 128.00 | 128.00 |
| net/http.newBufioReader | /usr/lib/go-1.24/src/net/http/server.go | 859 | 128.00 | 128.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2013 | 128.00 | 128.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 17.00 | 17.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 17.00 | 17.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 17.00 | 17.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 17.00 | 17.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 17.00 | 17.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 17.00 | 17.00 |
| runtime/pprof.writeHeap | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 593 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.newDeflateFast | /usr/lib/go-1.24/src/compress/flate/deflatefast.go | 64 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 586 | 17.00 | 17.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 17.00 | 17.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 17.00 | 17.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 17.00 | 17.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 17.00 | 17.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 17.00 | 17.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 17.00 | 17.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 17.00 | 17.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 17.00 | 17.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 17.00 | 17.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 17.00 | 17.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 17.00 | 17.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 17.00 | 17.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 17.00 | 17.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 587 | 10.00 | 10.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 10.00 | 10.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 10.00 | 10.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 10.00 | 10.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 10.00 | 10.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 10.00 | 10.00 |
| runtime/pprof.writeBlock | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 918 | 10.00 | 10.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 10.00 | 10.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 10.00 | 10.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 10.00 | 10.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 10.00 | 10.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 10.00 | 10.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 10.00 | 10.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 10.00 | 10.00 |
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
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 583 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.writeHeapProto | /usr/lib/go-1.24/src/runtime/pprof/protomem.go | 66 | 8.00 | 8.00 |
| runtime/pprof.writeHeapInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 634 | 8.00 | 8.00 |
| runtime/pprof.writeAlloc | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 599 | 8.00 | 8.00 |
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
| compress/flate.(*compressor).init | /usr/lib/go-1.24/src/compress/flate/deflate.go | 583 | 8.00 | 8.00 |
| compress/flate.NewWriter | /usr/lib/go-1.24/src/compress/flate/deflate.go | 664 | 8.00 | 8.00 |
| compress/gzip.(*Writer).Write | /usr/lib/go-1.24/src/compress/gzip/gzip.go | 191 | 8.00 | 8.00 |
| runtime/pprof.(*profileBuilder).build | /usr/lib/go-1.24/src/runtime/pprof/proto.go | 390 | 8.00 | 8.00 |
| runtime/pprof.printCountCycleProfile | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 455 | 8.00 | 8.00 |
| runtime/pprof.writeProfileInternal | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 944 | 8.00 | 8.00 |
| runtime/pprof.writeMutex | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 923 | 8.00 | 8.00 |
| runtime/pprof.(*Profile).WriteTo | /usr/lib/go-1.24/src/runtime/pprof/pprof.go | 377 | 8.00 | 8.00 |
| net/http/pprof.handler.ServeHTTP | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 272 | 8.00 | 8.00 |
| net/http/pprof.Index | /usr/lib/go-1.24/src/net/http/pprof/pprof.go | 389 | 8.00 | 8.00 |
| net/http.HandlerFunc.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2294 | 8.00 | 8.00 |
| net/http.(*ServeMux).ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 2822 | 8.00 | 8.00 |
| net/http.serverHandler.ServeHTTP | /usr/lib/go-1.24/src/net/http/server.go | 3301 | 8.00 | 8.00 |
| net/http.(*conn).serve | /usr/lib/go-1.24/src/net/http/server.go | 2102 | 8.00 | 8.00 |

### Allocation Analysis

- **Total Allocations**: 124532
- **Top 10% Concentration**: 92.8%
- **Allocation Severity**: Critical
- **Allocation Score**: 90/100

⚠️ **High Allocation Concentration Detected**
Top functions account for 92.8% of all allocations.
This indicates potential memory allocation hotspots that may benefit from optimization.

#### Top Allocation Hotspots

| Function | File | Line | Count | Percentage |
|----------|------|------|-------|------------|
| main.allocHeavyHandler | /home/doomguy/Documents/hackaton/Mistral-Hackathon/examples/demo-server/main.go | 43 | 50003 | 40.2% |
| unique.addUniqueMap[go.shape.struct { net/netip.isV6 bool; net/netip.zoneV6 string }].func1 | /usr/lib/go-1.24/src/unique/handle.go | 134 | 43691 | 35.1% |
| net/http.init.func14 | /usr/lib/go-1.24/src/net/http/header.go | 161 | 21845 | 17.5% |
| internal/sync.runtime_SemacquireMutex | /usr/lib/go-1.24/src/runtime/sema.go | 95 | 5461 | 4.4% |
| runtime.allocm | /usr/lib/go-1.24/src/runtime/proc.go | 2276 | 2052 | 1.6% |

---

*Generated by triageprof*
