package main

import (
	"cache/pkg"
	"testing"
)

func TestCache(t *testing.T) {
	// Initialize cache with LRU eviction policy.
	cache := db.NewCache(2, db.NewLRU())

	t.Run("Add and Get entry", func(t *testing.T) {
		cache.AddEntry("key1", "value1")
		value, exists := cache.Get("key1")
		if !exists || value != "value1" {
			t.Errorf("Expected 'value1', got '%s'", value)
		}
	})

	t.Run("Get non-existent key", func(t *testing.T) {
		_, exists := cache.Get("key_non_existent")
		if exists {
			t.Error("Expected key to be non-existent")
		}
	})

	t.Run("Eviction after exceeding capacity (LRU policy)", func(t *testing.T) {
		// Add entries to exceed capacity.
		cache.AddEntry("key1", "value1")
		cache.AddEntry("key2", "value2")
		cache.AddEntry("key3", "value3") // Should evict "key1"

		// Verify eviction has occurred (LRU).
		_, exists := cache.Get("key1") // "key1" should be evicted.
		if exists {
			t.Error("Expected 'key1' to be evicted, but it was found")
		}

		// Verify remaining keys.
		value, exists := cache.Get("key2")
		if !exists || value != "value2" {
			t.Errorf("Expected 'value2' for 'key2', got '%s'", value)
		}

		value, exists = cache.Get("key3")
		if !exists || value != "value3" {
			t.Errorf("Expected 'value3' for 'key3', got '%s'", value)
		}
	})

	t.Run("Change eviction policy dynamically", func(t *testing.T) {
		// Clear cache and switch eviction policy.
		cache.Clear()
		cache.ChangePolicy(db.NewLRU()) // Switching to a new instance of LRU for test clarity.

		// Add new entries.
		cache.AddEntry("key1", "value1")
		cache.AddEntry("key2", "value2")
		cache.AddEntry("key3", "value3") // Should evict "key1" again.

		// Verify new policy behavior (LRU).
		_, exists := cache.Get("key1")
		if exists {
			t.Error("Expected 'key1' to be evicted after policy change")
		}

		value, exists := cache.Get("key2")
		if !exists || value != "value2" {
			t.Errorf("Expected 'value2' for 'key2', got '%s'", value)
		}
	})
}
