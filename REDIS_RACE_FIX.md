# Redis Comment Floor Counter Race Condition Fix

## Problem Analysis

The original implementation in `internal/logic/comment.go:46-68` had a race condition:

```go
// BEFORE (Race Condition)
floor, err := redisClient.INCR(key)
if err == nil {
    return int(floor), nil
}

// If Redis fails, query database
tmpCmt := &model.Comment{}
_, err = MasterDB.Where("objid=? AND objtype=?", objid, objtype).OrderBy("floor DESC").Get(tmpCmt)
nextFloor := tmpCmt.Floor + 1
```

**Issues:**
1. Multiple goroutines could query the database simultaneously when Redis fails
2. All get the same `max_floor` and return `max_floor + 1`
3. Result: duplicate floor numbers

## Solution: Atomic Lua Script + Distributed Lock

### 1. Redis Lua Script (Primary Path)

```lua
local key = KEYS[1]
local ttl = ARGV[1]

-- Atomic INCR + EXPIRE
local floor = redis.call('INCR', key)
if floor == 1 then
    redis.call('EXPIRE', key, ttl)
end

return floor
```

**Benefits:**
- Atomic operation (no race condition)
- Automatic initialization with TTL
- Single Redis round-trip

### 2. Distributed Lock Fallback (Failure Path)

When Redis is unavailable:
1. Acquire distributed lock using `SET key value NX EX 10`
2. Query database for max floor
3. Initialize Redis counter
4. Release lock

**Code Implementation:**
```go
// Try lock acquisition (max 3 seconds)
for i := 0; i < 30; i++ {
    reply, err := redisClient.Do("SET", lockKey, lockValue, "NX", "EX", 10)
    if err == nil && reply != nil {
        locked = true
        break
    }
    time.Sleep(100 * time.Millisecond)
}
```

## Modified Files

### 1. `internal/logic/comment.go`

**Changes:**
- Added Lua script constant `getNextFloorScript`
- Rewrote `GetNextCommentFloor()` to use Lua script
- Added `getNextFloorWithLock()` for fallback with distributed lock
- Properly handle Redis key prefix

**Key Code:**
```go
// Execute Lua script (atomic)
result, err := redisClient.Do("EVAL", getNextFloorScript, 1, prefixedKey, ttl)
if err == nil {
    if floor, ok := result.(int64); ok {
        return int(floor), nil
    }
}

// Fallback to distributed lock
return getNextFloorWithLock(objid, objtype, key, ttl)
```

### 2. `internal/logic/comment_floor_test.go`

**Added:**
- `TestGetNextCommentFloorConcurrency`: 100-goroutine concurrent test
- Validates:
  - All floor numbers are unique
  - No missing floor numbers
  - Range: 1 to 100

## Testing

### Unit Test
```bash
# Run key format test (no Redis required)
go test -v -run TestCommentFloorKeyFormat ./internal/logic/
```

### Integration Test (Requires Redis + MySQL)
```bash
# Run full integration tests
go test -v -tags=integration -run TestGetNextCommentFloorIntegration ./internal/logic/

# Run concurrency test
go test -v -tags=integration -run TestGetNextCommentFloorConcurrency ./internal/logic/
```

### Verification Script
```bash
# Run mock concurrency test
go run verify_race_fix.go
```

Expected output:
```
Atomic Lua Script Version:
  Total goroutines: 100
  Unique floors: 100
  ✅ PASSED: All floors are unique
```

## Verification Checklist

- [x] Lua script atomic guarantee
- [x] Distributed lock fallback
- [x] Proper Redis key prefix handling
- [x] Code compiles successfully (`go build ./cmd/studygolang`)
- [x] No race conditions in mock test
- [x] Proper error handling and logging
- [x] Resource cleanup (defer Close(), defer DEL lock)

## Performance Impact

- **Normal path (Redis available):** Single Lua script execution (~1ms)
- **Fallback path (Redis down):** Distributed lock + DB query (~10-50ms)
- **Concurrency:** No contention in normal path, serialized in fallback

## Deployment Notes

1. **No database migration required** (logic change only)
2. **Backward compatible** with existing comments
3. **Redis requirement:** Supports EVAL command (Redis 2.6+)
4. **Monitoring:** Check logs for "Redis Lua script execution failed"

## Related Files

- `internal/logic/comment.go` - Core implementation
- `internal/logic/comment_floor_test.go` - Integration tests
- `internal/model/comment.go` - Comment model
- `internal/pkg/nosql/` - Redis client (external dependency)

## References

- [Redis EVAL command](https://redis.io/commands/eval/)
- [Redis SET NX EX](https://redis.io/commands/set/)
- [Distributed Lock Pattern](https://redis.io/topics/distlock)
