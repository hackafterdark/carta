package value

import (
	"database/sql"
	"math"
	"reflect"
	"testing"
	"time"
)

func TestCell_Scan(t *testing.T) {
	testCases := []struct {
		name     string
		src      interface{}
		expected Cell
	}{
		{
			name: "Scan int64",
			src:  int64(123),
			expected: Cell{
				kind:  reflect.Int64,
				bits:  123,
				valid: true,
			},
		},
		{
			name: "Scan float64",
			src:  float64(123.45),
			expected: Cell{
				kind:  reflect.Float64,
				bits:  math.Float64bits(123.45),
				valid: true,
			},
		},
		{
			name: "Scan bool",
			src:  true,
			expected: Cell{
				kind:  reflect.Bool,
				bits:  1,
				valid: true,
			},
		},
		{
			name: "Scan string",
			src:  "hello",
			expected: Cell{
				kind:  reflect.String,
				text:  "hello",
				valid: true,
			},
		},
		{
			name: "Scan bytes",
			src:  []byte("world"),
			expected: Cell{
				kind:  reflect.String,
				text:  "world",
				valid: true,
			},
		},
		{
			name: "Scan time",
			src:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			expected: Cell{
				kind:  reflect.Struct,
				time:  time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
				valid: true,
			},
		},
		{
			name:     "Scan nil",
			src:      nil,
			expected: Cell{valid: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cell := &Cell{}
			err := cell.Scan(tc.src)
			if err != nil {
				t.Fatalf("unexpected error during scan: %v", err)
			}
			if cell.kind != tc.expected.kind {
				t.Errorf("expected kind %v, got %v", tc.expected.kind, cell.kind)
			}
			if cell.bits != tc.expected.bits {
				t.Errorf("expected bits %v, got %v", tc.expected.bits, cell.bits)
			}
			if cell.text != tc.expected.text {
				t.Errorf("expected text %q, got %q", tc.expected.text, cell.text)
			}
			if !cell.time.Equal(tc.expected.time) {
				t.Errorf("expected time %v, got %v", tc.expected.time, cell.time)
			}
			if cell.valid != tc.expected.valid {
				t.Errorf("expected valid %v, got %v", tc.expected.valid, cell.valid)
			}
		})
	}
}

func TestCell_Getters(t *testing.T) {
	now := time.Now()

	testCases := []struct {
		name        string
		cell        Cell
		getBool     bool
		getInt64    int64
		getFloat64  float64
		getString   string
		getTime     time.Time
		getNullBool sql.NullBool
	}{
		{
			name:        "Int value",
			cell:        Cell{kind: reflect.Int64, bits: 123, valid: true},
			getBool:     true,
			getInt64:    123,
			getFloat64:  123,
			getNullBool: sql.NullBool{Bool: true, Valid: true},
		},
		{
			name:        "String value",
			cell:        Cell{kind: reflect.String, text: "hello", valid: true},
			getString:   "hello",
			getNullBool: sql.NullBool{Bool: false, Valid: true},
		},
		{
			name:        "Time value",
			cell:        Cell{kind: reflect.Struct, time: now, valid: true},
			getTime:     now,
			getNullBool: sql.NullBool{Bool: false, Valid: true},
		},
		{
			name:        "Null value",
			cell:        Cell{valid: false},
			getNullBool: sql.NullBool{Valid: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			b, _ := tc.cell.Bool()
			if b != tc.getBool {
				t.Errorf("Bool(): expected %v, got %v", tc.getBool, b)
			}

			i, _ := tc.cell.Int64()
			if i != tc.getInt64 {
				t.Errorf("Int64(): expected %v, got %v", tc.getInt64, i)
			}

			f, _ := tc.cell.Float64()
			if f != tc.getFloat64 {
				t.Errorf("Float64(): expected %v, got %v", tc.getFloat64, f)
			}

			s, _ := tc.cell.String()
			if s != tc.getString {
				t.Errorf("String(): expected %q, got %q", tc.getString, s)
			}

			tm, _ := tc.cell.Time()
			if !tm.Equal(tc.getTime) {
				t.Errorf("Time(): expected %v, got %v", tc.getTime, tm)
			}

			nb, _ := tc.cell.NullBool()
			if nb != tc.getNullBool {
				t.Errorf("NullBool(): expected %v, got %v", tc.getNullBool, nb)
			}
		})
	}
}

func TestCell_Conversions(t *testing.T) {
	// Test successful string to numeric conversions
	t.Run("Successful conversions", func(t *testing.T) {
		cell := NewCell("TEXT")
		cell.SetString("123")

		i32, err := cell.Int32()
		if err != nil || i32 != 123 {
			t.Errorf("Int32() failed: expected 123, got %d, err: %v", i32, err)
		}

		i64, err := cell.Int64()
		if err != nil || i64 != 123 {
			t.Errorf("Int64() failed: expected 123, got %d, err: %v", i64, err)
		}

		u32, err := cell.Uint32()
		if err != nil || u32 != 123 {
			t.Errorf("Uint32() failed: expected 123, got %d, err: %v", u32, err)
		}

		u64, err := cell.Uint64()
		if err != nil || u64 != 123 {
			t.Errorf("Uint64() failed: expected 123, got %d, err: %v", u64, err)
		}

		cell.SetString("123.45")
		f32, err := cell.Float32()
		if err != nil || f32 != 123.45 {
			t.Errorf("Float32() failed: expected 123.45, got %f, err: %v", f32, err)
		}

		f64, err := cell.Float64()
		if err != nil || f64 != 123.45 {
			t.Errorf("Float64() failed: expected 123.45, got %f, err: %v", f64, err)
		}
	})

	// Test conversion errors
	t.Run("Conversion errors", func(t *testing.T) {
		cell := NewCell("TEXT")
		cell.SetString("not a number")

		if _, err := cell.Int32(); err == nil {
			t.Error("Int32() expected an error, but got nil")
		}
		if _, err := cell.Int64(); err == nil {
			t.Error("Int64() expected an error, but got nil")
		}
		if _, err := cell.Uint32(); err == nil {
			t.Error("Uint32() expected an error, but got nil")
		}
		if _, err := cell.Uint64(); err == nil {
			t.Error("Uint64() expected an error, but got nil")
		}
		if _, err := cell.Float32(); err == nil {
			t.Error("Float32() expected an error, but got nil")
		}
		if _, err := cell.Float64(); err == nil {
			t.Error("Float64() expected an error, but got nil")
		}
	})
}

func TestCell_AsInterface(t *testing.T) {
	cell := NewCell("INT")
	cell.SetInt64(123)
	val, err := cell.AsInterface()
	if err != nil {
		t.Fatalf("AsInterface() failed: %v", err)
	}
	if val.(int64) != 123 {
		t.Errorf("AsInterface() failed: expected 123, got %v", val)
	}
}

func TestCell_Uid(t *testing.T) {
	testCases := []struct {
		name     string
		cell     func() *Cell
		expected string
	}{
		{
			name: "Null UID",
			cell: func() *Cell {
				c := NewCell("")
				c.SetNull()
				return c
			},
			expected: "cnull",
		},
		{
			name: "Bool true UID",
			cell: func() *Cell {
				c := NewCell("BOOL")
				c.SetBool(true)
				return c
			},
			expected: "ctrue",
		},
		{
			name: "Bool false UID",
			cell: func() *Cell {
				c := NewCell("BOOL")
				c.SetBool(false)
				return c
			},
			expected: "cfalse",
		},
		{
			name: "String UID",
			cell: func() *Cell {
				c := NewCell("TEXT")
				c.SetString("hello")
				return c
			},
			expected: "hello",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			uid := tc.cell().Uid()
			if uid != tc.expected {
				t.Errorf("expected UID %q, got %q", tc.expected, uid)
			}
		})
	}
}
