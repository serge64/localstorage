package localstorage_test

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/serge64/localstorage"
)

func TestLocalStorage_Get(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)
	_ = db.Put("key", "value", 0)

	testcases := []struct {
		name     string
		key      string
		expected string
	}{
		{
			name:     "valid",
			key:      "key",
			expected: "value",
		},
		{
			name:     "no valid",
			key:      "invalidkey",
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			value, _ := db.Get(tc.key)
			if value == nil {
				value = ""
			} else {
				value = value.(string)
			}
			if value != tc.expected {
				t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", tc.expected, value)
			}
		})
	}
}

func TestLocalStorage_Put(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)

	testcases := []struct {
		name     string
		key      string
		value    interface{}
		expected error
	}{
		{
			name:  "valid",
			key:   "key",
			value: "value",
		},
		{
			name:     "no valid",
			key:      "key",
			value:    "value",
			expected: localstorage.ErrNotUniqueKey,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Put(tc.key, tc.value, 0)
			if err != tc.expected {
				t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", tc.expected, err)
			}
		})
	}
}

func TestLocalStorage_PutTTL(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)

	testcases := []struct {
		name     string
		key      string
		value    interface{}
		ttl      time.Duration
		timeout  time.Duration
		expected string
	}{
		{
			name:     "valid",
			key:      "key1",
			value:    "value",
			expected: "value",
		},
		{
			name:     "no valid",
			key:      "key2",
			value:    "value",
			ttl:      time.Duration(10) * time.Millisecond,
			timeout:  time.Duration(11) * time.Millisecond,
			expected: "",
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Put(tc.key, tc.value, tc.ttl)
			if err != nil {
				t.Errorf("No error expected but got %s", err)
			}

			time.Sleep(tc.timeout)

			value, _ := db.Get(tc.key)
			if value == nil {
				value = ""
			} else {
				value = value.(string)
			}

			if value != tc.expected {
				t.Errorf("Values not equals:\n - expected: %s\n- actual: %s", tc.expected, value)
			}
		})
	}
}

func TestLocalStorage_Del(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)
	_ = db.Put("key", "value", 0)

	testcases := []struct {
		name     string
		key      string
		expected error
	}{
		{
			name: "valid",
			key:  "key",
		},
		{
			name:     "no valid",
			key:      "invalidkey",
			expected: localstorage.ErrNotFoundKey,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			err := db.Del(tc.key)
			if err != tc.expected {
				t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", tc.expected, err)
			}
		})
	}
}

func TestLocalStorage_Keys(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)

	keys := db.Keys()
	if len(keys) != 0 {
		t.Errorf("Values not equals:\n- expected: 0\n- actual: %d", len(keys))
	}

	_ = db.Put("key", "value", 0)
	expected := "key"
	keys = db.Keys()
	if expected != keys[0] {
		t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", expected, keys[0])
	}
}

func TestLocalStorage_Values(t *testing.T) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 10)

	values := db.Values()
	if len(values) != 0 {
		t.Errorf("Values not equals:\n- expected: 0\n- actual: %d", len(values))
	}

	_ = db.Put("key", "value", 0)
	expected := "value"
	values = db.Values()
	if expected != values[0].(string) {
		t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", expected, values[0].(string))
	}
}

func BenchmarkLocalStorage_AllInConcurency(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 1024*100)
	ttl := time.Millisecond
	ttl2 := time.Duration(100) * time.Microsecond

	keys := GenerateKeys(1024 * 100)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, _ = db.Get(keys[rand.Intn(len(keys))])
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(len(keys))], struct{}{}, 0)
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(len(keys))], struct{}{}, ttl)
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(len(keys))], struct{}{}, ttl2)
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Del(keys[rand.Intn(9)])
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Keys()
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Values()
		}
	})
}

func BenchmarkLocalStorage_Get(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{}, 0)
	}

	for i := 0; i < b.N; i++ {
		_, _ = db.Get(keys[rand.Intn(len(keys))])
	}
}

func BenchmarkLocalStorage_Put(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	keys := GenerateKeys(128)

	for i := 0; i < b.N; i++ {
		_ = db.Put(keys[rand.Intn(len(keys))], struct{}{}, 0)
	}
}

func BenchmarkLocalStorage_PutTTL(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	ttl := time.Duration(10) * time.Microsecond
	keys := GenerateKeys(128)

	for i := 0; i < b.N; i++ {
		_ = db.Put(keys[rand.Intn(len(keys))], struct{}{}, ttl)
	}
}

func BenchmarkLocalStorage_Del(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{}, 0)
	}

	for i := 0; i < b.N; i++ {
		_ = db.Del(keys[rand.Intn(len(keys))])
	}
}

func BenchmarkLocalStorage_Keys(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{}, 0)
	}

	for i := 0; i < b.N; i++ {
		_ = db.Keys()
	}
}

func BenchmarkLocalStorage_Values(b *testing.B) {
	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	db := localstorage.New(ctx, 128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{}, 0)
	}

	for i := 0; i < b.N; i++ {
		_ = db.Values()
	}
}

func GenerateKeys(size int) []string {
	keys := make([]string, 0, size)
	for i := 0; i < cap(keys); i++ {
		keys = append(keys, strconv.Itoa(i))
	}
	return keys
}
