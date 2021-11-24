package localstorage_test

import (
	"math/rand"
	"strconv"
	"testing"

	"github.com/serge64/localstorage"
)

func TestLocalStorage_Get(t *testing.T) {
	db := localstorage.New(10)
	_ = db.Put("key", "value")

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
	db := localstorage.New(10)

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
			err := db.Put(tc.key, tc.value)
			if err != tc.expected {
				t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", tc.expected, err)
			}
		})
	}
}

func TestLocalStorage_Del(t *testing.T) {
	db := localstorage.New(10)
	_ = db.Put("key", "value")

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
	db := localstorage.New(10)

	keys := db.Keys()
	if len(keys) != 0 {
		t.Errorf("Values not equals:\n- expected: 0\n- actual: %d", len(keys))
	}

	_ = db.Put("key", "value")
	expected := "key"

	for i := 0; i < 2; i++ {
		keys = db.Keys()
		if expected != keys[0] {
			t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", expected, keys[0])
		}
	}
}

func TestLocalStorage_Values(t *testing.T) {
	db := localstorage.New(10)

	values := db.Values()
	if len(values) != 0 {
		t.Errorf("Values not equals:\n- expected: 0\n- actual: %d", len(values))
	}

	_ = db.Put("key", "value")
	expected := "value"

	for i := 0; i < 2; i++ {
		values = db.Values()
		if expected != values[0].(string) {
			t.Errorf("Values not equals:\n- expected: %s\n- actual: %s", expected, values[0].(string))
		}
	}
}

func BenchmarkLocalStorage_AsyncAll(b *testing.B) {
	db := localstorage.New(1024)
	keys := GenerateKeys(1024)
	count := len(keys)

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_, _ = db.Get(keys[rand.Intn(count)])
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(count)], struct{}{})
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(count)], struct{}{})
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Put(keys[rand.Intn(count)], struct{}{})
		}
	})

	b.RunParallel(func(p *testing.PB) {
		for p.Next() {
			_ = db.Del(keys[rand.Intn(count)])
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
	db := localstorage.New(128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{})
	}

	for i := 0; i < b.N; i++ {
		_, _ = db.Get(keys[rand.Intn(len(keys))])
	}
}

func BenchmarkLocalStorage_Put(b *testing.B) {
	db := localstorage.New(128)
	keys := GenerateKeys(128)

	for i := 0; i < b.N; i++ {
		_ = db.Put(keys[rand.Intn(len(keys))], struct{}{})
	}
}

func BenchmarkLocalStorage_Del(b *testing.B) {
	db := localstorage.New(128)
	keys := GenerateKeys(128)

	for _, v := range keys {
		_ = db.Put(v, struct{}{})
	}

	for i := 0; i < b.N; i++ {
		_ = db.Del(keys[rand.Intn(len(keys))])
	}
}

func BenchmarkLocalStorage_Keys(b *testing.B) {
	db := localstorage.New(128)
	keys := GenerateKeys(128)

	go func() {
		for _, v := range keys {
			_ = db.Put(v, struct{}{})
		}
	}()

	for i := 0; i < b.N; i++ {
		_ = db.Keys()
	}
}

func BenchmarkLocalStorage_Values(b *testing.B) {
	db := localstorage.New(128)
	keys := GenerateKeys(128)

	go func() {
		for _, v := range keys {
			_ = db.Put(v, struct{}{})
		}
	}()

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
