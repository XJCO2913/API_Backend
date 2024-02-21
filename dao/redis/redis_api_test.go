package redis

import (
    "context"
    "testing"
    "time"
)

func TestSetAndGetKeyValue(t *testing.T) {
    ctx := context.Background()
    key := "testKey"
    valueExpected := "Hello Redis!"

    if err := SetKeyValue(ctx, key, valueExpected, 10*time.Second); err != nil {
        t.Fatalf("Failed to set key: %v", err)
    }

    valueGot, err := GetKeyValue(ctx, key)
    if err != nil {
        t.Fatalf("Failed to get key: %v", err)
    }

    if valueGot != valueExpected {
        t.Fatalf("Expected value '%s', got '%s'", valueExpected, valueGot)
    }
}
