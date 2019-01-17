package wallet

import (
	"testing"
	"time"
)

func GetTestDB() *BalanceCache {
	db := NewBalanceCache()

	db.Set("a", 0)
	db.Set("b", 1)
	db.Set("c", 9999999999999999999)
	return db
}

func TestBalanceCache_Get(t *testing.T) {
	db := GetTestDB()
	db.ttl = time.Millisecond * 500

	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 0, true},
		{"Case b", db, args{"b"}, 1, true},
		{"Case c", db, args{"c"}, 9999999999999999999, true},
		{"Case d", db, args{"d"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	time.Sleep(time.Millisecond * 600)

	db.Set("new", 2)

	tests = []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 0, false},
		{"Case b", db, args{"b"}, 0, false},
		{"Case c", db, args{"c"}, 0, false},
		{"Case d", db, args{"d"}, 0, false},
		{"Case new", db, args{"new"}, 2, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBalanceCache_Clear(t *testing.T) {
	db := GetTestDB()

	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 0, true},
		{"Case b", db, args{"b"}, 1, true},
		{"Case c", db, args{"c"}, 9999999999999999999, true},
		{"Case d", db, args{"d"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	db.Clear()

	tests = []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 0, false},
		{"Case b", db, args{"b"}, 0, false},
		{"Case c", db, args{"c"}, 0, false},
		{"Case d", db, args{"d"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func TestBalanceCache_Set(t *testing.T) {
	db := NewBalanceCache()

	type args struct {
		key string
	}
	tests := []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 0, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}

	db.Set("a", 1)

	tests2 := []struct {
		name  string
		c     *BalanceCache
		args  args
		want  uint64
		want1 bool
	}{
		{"Case a", db, args{"a"}, 1, true},
	}
	for _, tt := range tests2 {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := tt.c.Get(tt.args.key)
			if got != tt.want {
				t.Errorf("BalanceCache.Get() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("BalanceCache.Get() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
