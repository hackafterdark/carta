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
			name: "Scan int",
			src:  int(123),
			expected: Cell{
				kind:  reflect.Int64,
				bits:  123,
				valid: true,
			},
		},
		{
			name: "Scan int8",
			src:  int8(123),
			expected: Cell{
				kind:  reflect.Int64,
				bits:  123,
				valid: true,
			},
		},
		{
			name: "Scan int16",
			src:  int16(123),
			expected: Cell{
				kind:  reflect.Int64,
				bits:  123,
				valid: true,
			},
		},
		{
			name: "Scan int32",
			src:  int32(123),
			expected: Cell{
				kind:  reflect.Int64,
				bits:  123,
				valid: true,
			},
		},
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
		name          string
		cell          Cell
		getBool       bool
		getInt64      int64
		getFloat64    float64
		getString     string
		getTime       time.Time
		getNullBool   sql.NullBool
		getNullInt64  sql.NullInt64
		getNullString sql.NullString
		getNullTime   sql.NullTime
	}{
		{
			name:          "Int value",
			cell:          Cell{kind: reflect.Int64, bits: 123, valid: true},
			getBool:       true,
			getInt64:      123,
			getFloat64:    123,
			getNullBool:   sql.NullBool{Bool: true, Valid: true},
			getNullInt64:  sql.NullInt64{Int64: 123, Valid: true},
			getNullString: sql.NullString{String: "", Valid: true},
			getNullTime:   sql.NullTime{Time: time.Time{}, Valid: true},
		},
		{
			name:          "String value",
			cell:          Cell{kind: reflect.String, text: "hello", valid: true},
			getString:     "hello",
			getNullBool:   sql.NullBool{Bool: false, Valid: true},
			getNullInt64:  sql.NullInt64{Int64: 0, Valid: true},
			getNullString: sql.NullString{String: "hello", Valid: true},
			getNullTime:   sql.NullTime{Time: time.Time{}, Valid: true},
		},
		{
			name:          "Time value",
			cell:          Cell{kind: reflect.Struct, time: now, valid: true},
			getTime:       now,
			getNullBool:   sql.NullBool{Bool: false, Valid: true},
			getNullInt64:  sql.NullInt64{Int64: 0, Valid: true},
			getNullString: sql.NullString{String: "", Valid: true},
			getNullTime:   sql.NullTime{Time: now, Valid: true},
		},
		{
			name:          "Null value",
			cell:          Cell{valid: false},
			getNullBool:   sql.NullBool{Valid: false},
			getNullInt64:  sql.NullInt64{Valid: false},
			getNullString: sql.NullString{Valid: false},
			getNullTime:   sql.NullTime{Valid: false},
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

			ni, _ := tc.cell.NullInt64()
			if ni != tc.getNullInt64 {
				t.Errorf("NullInt64(): expected %v, got %v", tc.getNullInt64, ni)
			}

			ns, _ := tc.cell.NullString()
			if ns != tc.getNullString {
				t.Errorf("NullString(): expected %v, got %v", tc.getNullString, ns)
			}

			nt, _ := tc.cell.NullTime()
			if nt != tc.getNullTime {
				t.Errorf("NullTime(): expected %v, got %v", tc.getNullTime, nt)
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
	testCases := []struct {
		name     string
		cell     func() *Cell
		expected interface{}
	}{
		{
			name: "AsInterface bool",
			cell: func() *Cell {
				c := NewCell("BOOL")
				c.SetBool(true)
				return c
			},
			expected: true,
		},
		{
			name: "AsInterface int32",
			cell: func() *Cell {
				c := NewCell("INT")
				c.SetInt64(123)
				c.kind = reflect.Int32
				return c
			},
			expected: int32(123),
		},
		{
			name: "AsInterface int64",
			cell: func() *Cell {
				c := NewCell("BIGINT")
				c.SetInt64(123)
				return c
			},
			expected: int64(123),
		},
		{
			name: "AsInterface uint32",
			cell: func() *Cell {
				c := NewCell("INT UNSIGNED")
				c.SetInt64(123)
				c.kind = reflect.Uint32
				return c
			},
			expected: uint32(123),
		},
		{
			name: "AsInterface uint64",
			cell: func() *Cell {
				c := NewCell("BIGINT UNSIGNED")
				c.SetInt64(123)
				c.kind = reflect.Uint64
				return c
			},
			expected: uint64(123),
		},
		{
			name: "AsInterface float32",
			cell: func() *Cell {
				c := NewCell("FLOAT")
				c.bits = uint64(math.Float32bits(123.45))
				c.kind = reflect.Float32
				c.valid = true
				return c
			},
			expected: float32(123.45),
		},
		{
			name: "AsInterface float64",
			cell: func() *Cell {
				c := NewCell("DOUBLE")
				c.SetFloat64(123.45)
				return c
			},
			expected: 123.45,
		},
		{
			name: "AsInterface string",
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
			val, err := tc.cell().AsInterface()
			if err != nil {
				t.Fatalf("AsInterface() failed: %v", err)
			}
			if val != tc.expected {
				t.Errorf("AsInterface() failed: expected %v, got %v", tc.expected, val)
			}
		})
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
		{
			name: "Int64 UID",
			cell: func() *Cell {
				c := NewCell("BIGINT")
				c.SetInt64(123456789)
				return c
			},
			expected: "21i3v9",
		},
		{
			name: "Float64 UID",
			cell: func() *Cell {
				c := NewCell("DOUBLE")
				c.SetFloat64(123.45)
				return c
			},
			expected: "z8nf4cjfzqod",
		},
		{
			name: "Time UID",
			cell: func() *Cell {
				c := NewCell("TIMESTAMP")
				c.SetTime(time.Unix(123456789, 0))
				return c
			},
			expected: "21i3v9",
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

func TestCell_Setters(t *testing.T) {
	t.Run("SetBool", func(t *testing.T) {
		c := NewCell("BOOL")
		c.SetBool(true)
		if c.Kind() != reflect.Bool || !c.IsValid() || c.bits != 1 {
			t.Error("SetBool(true) failed")
		}
		c.SetBool(false)
		if c.Kind() != reflect.Bool || !c.IsValid() || c.bits != 0 {
			t.Error("SetBool(false) failed")
		}
	})

	t.Run("SetFloat64", func(t *testing.T) {
		c := NewCell("DOUBLE")
		c.SetFloat64(123.45)
		if c.Kind() != reflect.Float64 || !c.IsValid() || math.Float64frombits(c.bits) != 123.45 {
			t.Error("SetFloat64() failed")
		}
	})

	t.Run("SetInt64", func(t *testing.T) {
		c := NewCell("BIGINT")
		c.SetInt64(123)
		if c.Kind() != reflect.Int64 || !c.IsValid() || c.bits != 123 {
			t.Error("SetInt64() failed")
		}
	})

	t.Run("SetString", func(t *testing.T) {
		c := NewCell("TEXT")
		c.SetString("hello")
		if c.Kind() != reflect.String || !c.IsValid() || c.text != "hello" {
			t.Error("SetString() failed")
		}
	})

	t.Run("SetTime", func(t *testing.T) {
		now := time.Now()
		c := NewCell("TIMESTAMP")
		c.SetTime(now)
		if c.Kind() != reflect.Struct || !c.IsValid() || !c.time.Equal(now) {
			t.Error("SetTime() failed")
		}
	})

	t.Run("SetNull", func(t *testing.T) {
		c := NewCell("TEXT")
		c.SetNull()
		if c.IsValid() {
			t.Error("SetNull() failed")
		}
	})
}

func TestCell_Remaining(t *testing.T) {
	t.Run("NewCellWithData", func(t *testing.T) {
		c := NewCellWithData("INT", 123)
		if c.Kind() != reflect.Int64 || !c.IsValid() || c.bits != 123 {
			t.Error("NewCellWithData() failed")
		}
	})

	t.Run("Timestamp", func(t *testing.T) {
		now := time.Now()
		c := NewCell("TIMESTAMP")
		c.SetTime(now)
		ts, err := c.Timestamp()
		if err != nil {
			t.Fatalf("Timestamp() failed: %v", err)
		}
		if ts.Seconds != now.Unix() {
			t.Errorf("Timestamp() failed: expected seconds %v, got %v", now.Unix(), ts.Seconds)
		}
	})

	t.Run("NullFloat64", func(t *testing.T) {
		c := NewCell("DOUBLE")
		c.SetFloat64(123.45)
		nf, err := c.NullFloat64()
		if err != nil {
			t.Fatalf("NullFloat64() failed: %v", err)
		}
		if !nf.Valid || nf.Float64 != 123.45 {
			t.Errorf("NullFloat64() failed: expected %v, got %v", 123.45, nf.Float64)
		}

		c.SetNull()
		nf, err = c.NullFloat64()
		if err != nil {
			t.Fatalf("NullFloat64() failed: %v", err)
		}
		if nf.Valid {
			t.Error("NullFloat64() failed: expected invalid, got valid")
		}
	})

	t.Run("NullInt32", func(t *testing.T) {
		c := NewCell("INT")
		c.SetInt64(123)
		ni, err := c.NullInt32()
		if err != nil {
			t.Fatalf("NullInt32() failed: %v", err)
		}
		if !ni.Valid || ni.Int32 != 123 {
			t.Errorf("NullInt32() failed: expected %v, got %v", 123, ni.Int32)
		}

		c.SetNull()
		ni, err = c.NullInt32()
		if err != nil {
			t.Fatalf("NullInt32() failed: %v", err)
		}
		if ni.Valid {
			t.Error("NullInt32() failed: expected invalid, got valid")
		}
	})

	t.Run("Int64 from float64", func(t *testing.T) {
		c := NewCell("DOUBLE")
		c.SetFloat64(123.45)
		i, err := c.Int64()
		if err != nil {
			t.Fatalf("Int64() from float64 failed: %v", err)
		}
		if i != 123 {
			t.Errorf("Int64() from float64 failed: expected 123, got %v", i)
		}
	})

	t.Run("Uid default", func(t *testing.T) {
		c := NewCell("OTHER")
		c.kind = reflect.Map
		c.valid = true
		uid := c.Uid()
		if uid != "" {
			t.Errorf("Uid() default failed: expected empty string, got %q", uid)
		}
	})
}
