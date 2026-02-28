<<<<<<< SEARCH
- Memory profiling provides top allocation sources for optimization
- Feature completes Layer 3 goal: "Python cProfile plugin (full implementation)"
- Plugin now has feature parity with Go pprof plugin for core profiling types

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`

### Iter 20240306-1215 — UTC
**Type:** Feature
**Objective:** Enhance Python cProfile plugin with heap memory profiling support

**Acceptance criteria (feature)**
- [x] Add heap profile collection to Python cProfile plugin
- [x] Maintain backward compatibility with existing CPU and allocation profiles
- [x] Update plugin manifest with new capabilities
- [x] Verify all three profile types (cpu, heap, allocs) work correctly
- [x] Ensure output format compatibility with analyzer

**Changes**
- `plugins/src/python-cprofile/main.py`: Added `_collect_heap_profile()` method for heap memory profiling using tracemalloc
- `plugins/src/python-cprofile/main.py`: Updated plugin info to version 0.4.0 and added "heap" to capabilities
- `plugins/manifests/python-cprofile.json`: Updated to version 0.4.0 and added "heap" profile capability
- `examples/python-demo-server/demo.py`: Added comprehensive Python demo server with performance issues

**Verification**
- Manual testing: All three profile types collected successfully
- Format validation: JSON output matches expected analyzer input format
- Backward compatibility: Existing CPU and allocation profiles still work
- Plugin discovery: Updated manifest reflects new capabilities

**Risk/Notes**
- No breaking changes - purely additive feature
- Python plugin now has feature parity with Go plugin for core profiling types
- Heap profiling uses tracemalloc for comprehensive memory analysis
- Feature completes Layer 3 goal: "Python cProfile plugin (full implementation)"

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`
=======
- Memory profiling provides top allocation sources for optimization
- Feature completes Layer 3 goal: "Python cProfile plugin (full implementation)"
- Plugin now has feature parity with Go pprof plugin for core profiling types

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`

### Iter 20240306-1200 — UTC
**Type:** Maintenance
**Objective:** Fix pipeline tests and verify allocation analysis functionality

**Acceptance criteria (maintenance)**
- [x] Fix failing pipeline tests by using heap profile with actual samples
- [x] Verify allocation analysis feature is working correctly
- [x] Ensure all existing tests pass
- [x] Maintain backward compatibility
- [x] Build succeeds without errors

**Changes**
- `internal/core/pipeline_test.go`: Updated tests to use heap.pb.gz instead of cpu.pb.gz (which had 0 samples)
- `internal/core/pipeline_test.go`: Fixed test assertions to match updated profile types
- `internal/analyzer/analyzer.go`: Removed debug logging added during investigation
- Verified allocation analysis integration tests pass

**Verification**
- Tests: `go test ./...` - all passing
- Build: `make build` - successful
- Allocation Analysis: `go test ./internal/analyzer -run TestAllocationAnalysisIntegration` - passes
- Pipeline Tests: `go test ./internal/core -run TestPipeline_Analyze` - passes
- Backward Compatibility: All existing functionality preserved

**Risk/Notes**
- No breaking changes - purely test fixes and verification
- Allocation analysis feature was already implemented and working
- Tests now use heap profile which has 22 samples vs cpu profile which had 0 samples
- Feature completes Layer 1 goal: "Analysis engine improvements: hotspot ranking, regression detection"

**Git / Rollback**
- Branch: `main`
- Checkpoint tag: N/A (direct commit to main)
- Commit: `<to be filled after commit>`
- Rollback: `git revert <commit-hash>`
>>>>>>> REPLACE