# Repository Integration Tests

This directory contains integration tests for the repository layer, specifically testing the Phase 1 functions for the Poop Wrapped feature.

## Running Tests

### Run all tests
```bash
go test ./repository/... -v
```

### Run with coverage
```bash
go test ./repository/... -cover
```

### Run benchmarks
```bash
go test ./repository/... -bench=. -benchmem
```

### Using Makefile
```bash
make test              # Run all tests with verbose output
make test-coverage     # Run tests with coverage report
make test-benchmark    # Run benchmark tests
```

## Test Structure

The test suite (`repository_test.go`) includes:

### Setup Functions
- `setupTestDB()`: Creates an in-memory SQLite database for each test
- `insertTestData()`: Inserts realistic test data with multiple users and scenarios

### Test Data

The test data includes 5 users with different patterns:

1. **Alice** (user_id: 1001)
   - 5 poops in January 2025
   - Morning person (8-10 AM)
   - Also has 3 poops in 2024 (for year filtering tests)

2. **Bob** (user_id: 1002)
   - 10 poops in January 2025
   - Night owl (11 PM - 2 AM)
   - Has a 3-day consecutive streak (for Consistency King award)

3. **Charlie** (user_id: 1003)
   - 8 poops in January 2025
   - Weekend warrior (mostly Saturdays and Sundays)
   - Some days have multiple poops

4. **Machine Gun** (user_id: 1004)
   - 5 poops in a single day (Jan 15)
   - Should win "Machine Gun" award

5. **Early Bird** (user_id: 1005)
   - 6 poops in January 2025
   - All between 5-8 AM
   - Should win "Early Bird" award

## Test Cases

### TestGetYearlyPoopCount
Tests counting poops for specific users and years:
- ✅ Valid users with data
- ✅ Year filtering (2024 vs 2025)
- ✅ Non-existent users (returns 0)
- ✅ Users with no data in specific year

### TestGetPoopsByHour
Tests hour distribution:
- ✅ Returns all 24 hours (even with 0 counts)
- ✅ Hours are in correct order (0-23)
- ✅ Correctly counts poops per hour
- ✅ Handles missing hours (fills with 0)

### TestGetPoopsByDayOfWeek
Tests day-of-week distribution:
- ✅ Returns all 7 days
- ✅ Days in correct order (Monday-Sunday)
- ✅ Correctly identifies weekend vs weekday patterns

### TestGetYearlyRanking
Tests user ranking:
- ✅ Correct rank calculation
- ✅ Percentage calculation (users below this user)
- ✅ Ranking order (more poops = better rank)
- ✅ Edge case: single user (rank 1, percentage 0)

### TestGetGroupYearlyStats
Tests group statistics:
- ✅ Returns all users
- ✅ Sorted by poop count (descending)
- ✅ Correct counts for each user

### TestGetGroupAwards
Tests award calculations:
- ✅ Early Bird award (most 5-8 AM poops)
- ✅ Machine Gun award (most in single day)
- ✅ Consistency King award (longest streak)
- ✅ Weekend Warrior award (highest % on weekends)
- ✅ Night Owl award (most 11 PM - 4 AM poops)

### Edge Case Tests
- `TestGetYearlyRanking_EdgeCases`: Single user scenario
- `TestGetPoopsByHour_AllHoursPresent`: All 24 hours have data
- `TestGetPoopsByDayOfWeek_AllDaysPresent`: All 7 days have data

## Benchmark Tests

Performance benchmarks for:
- `BenchmarkGetYearlyPoopCount`
- `BenchmarkGetPoopsByHour`
- `BenchmarkGetGroupAwards`

## Expected Results

When running tests, you should see:
- ✅ All tests passing
- ✅ Award winners match expected users:
  - Early Bird: `early_bird`
  - Machine Gun: `machine_gun` (5 poops in one day)
  - Consistency King: `bob` (3-day streak)
  - Weekend Warrior: `charlie` (highest weekend %)

## Troubleshooting

### Tests fail with "database is locked"
- This shouldn't happen with in-memory databases, but if it does, ensure tests aren't running in parallel
- Use `go test -p 1` to run tests sequentially

### Award tests fail
- Check that test data matches expected patterns
- Verify SQL queries handle edge cases (no winners, ties, etc.)

### Year filtering issues
- Ensure timestamps are formatted correctly
- Check that `strftime('%Y', timestamp)` works as expected

## Adding New Tests

When adding new repository functions:

1. Add test data in `insertTestData()` if needed
2. Create a test function following the pattern: `TestFunctionName`
3. Use table-driven tests for multiple scenarios
4. Test edge cases (empty results, single user, etc.)
5. Add benchmarks for performance-critical functions

Example:
```go
func TestNewFunction(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()
    insertTestData(t, db)

    ctx := context.Background()
    result, err := NewFunction(ctx, db, params...)
    if err != nil {
        t.Fatalf("NewFunction() error = %v", err)
    }
    // Assertions...
}
```

