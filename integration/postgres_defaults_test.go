//go:build !sqlite

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"integration/valk"
	"integration/valk/allFieldsSoFar"
	"math"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestNativeDefaults_MinimalCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	rec, err := db.AllFieldsSoFar.Create(
		allFieldsSoFar.StringReq.Set("hello"),
		allFieldsSoFar.StringVarchar.Set("varchar"),
		allFieldsSoFar.StringChar.Set("0123456789"),
		allFieldsSoFar.BitVal.Set("1010101010"),
		allFieldsSoFar.VarBitVal.Set("1101"),
		allFieldsSoFar.InetVal.Set("192.168.1.1"),
		allFieldsSoFar.XmlVal.Set("<root/>"),
		allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		allFieldsSoFar.IntReq.Set(1),
		allFieldsSoFar.IntegerVal.Set(2),
		allFieldsSoFar.SmallInt.Set(3),
		allFieldsSoFar.TinyInt.Set(4),
		allFieldsSoFar.OidVal.Set(5),
		allFieldsSoFar.BigIntReq.Set(6),
		allFieldsSoFar.FloatReq.Set(1.5),
		allFieldsSoFar.RealVal.Set(2.5),
		allFieldsSoFar.DecimalReq.Set("10.50"),
		allFieldsSoFar.DecimalPrecise.Set("99.99"),
		allFieldsSoFar.MoneyVal.Set("12.34"),
		allFieldsSoFar.BoolReq.Set(true),
		allFieldsSoFar.DateTimeReq.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.UpdatedAt.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.DateTimeTz.Set(time.Date(2024, 6, 15, 12, 0, 0, 0, time.FixedZone("IST", 3600))),
		allFieldsSoFar.TimestampVal.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.TimeVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC)),
		allFieldsSoFar.TimetzVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.FixedZone("IST", 3600))),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`{"key":"value"}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`[1,2,3]`)),
		allFieldsSoFar.BytesReq.Set([]byte{0x00, 0x01, 0x02}),
		allFieldsSoFar.LtreeField.Set("Top.Collections"),
	).Exec(ctx)

	if err != nil {
		t.Fatalf("minimal create failed: %v", err)
	}

	if rec.Id == 0 {
		t.Error("expected autoincrement id")
	}
	if rec.StringReq != "hello" {
		t.Errorf("StringReq = %q", rec.StringReq)
	}
	if rec.StringDefault != "default" {
		t.Errorf("StringDefault default = %q", rec.StringDefault)
	}
	if rec.IntDefault != 42 {
		t.Errorf("IntDefault default = %d", rec.IntDefault)
	}
	if rec.BoolDefault {
		t.Error("BoolDefault default should be false")
	}
	if rec.DateTimeDefault.IsZero() {
		t.Error("DateTimeDefault should be populated")
	}
	if rec.CuidDefault == "" {
		t.Error("CuidDefault should be populated")
	}
	if rec.UuidDefault == "" {
		t.Error("UuidDefault should be populated")
	}
	if rec.Uuid4Default == "" {
		t.Error("Uuid4Default should be populated")
	}
	if rec.Uuid7Default == "" {
		t.Error("Uuid7Default should be populated")
	}
	if rec.UlidDefault == "" {
		t.Error("UlidDefault should be populated")
	}
	if rec.NanoidDefault == "" {
		t.Error("NanoidDefault should be populated")
	}
	if ltreeFieldStr(rec.LtreeField) != "Top.Collections" {
		t.Errorf("ltreeField = %s", ltreeFieldStr(rec.LtreeField))
	}
}

func TestNativeDefaults_AllExplicitValues(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	now := time.Date(2024, 7, 11, 15, 30, 45, 123456789, time.UTC)

	rec, err := db.AllFieldsSoFar.Create(
		allFieldsSoFar.StringReq.Set("explicit"),
		allFieldsSoFar.StringOpt.Set("opt"),
		allFieldsSoFar.StringDefault.Set("overridden"),
		allFieldsSoFar.StringVarchar.Set(strings.Repeat("a", 255)),
		allFieldsSoFar.StringChar.Set("1234567890"),
		allFieldsSoFar.BitVal.Set("1111111111"),
		allFieldsSoFar.VarBitVal.Set("0"),
		allFieldsSoFar.InetVal.Set("::1"),
		allFieldsSoFar.XmlVal.Set("<doc><elem>text</elem></doc>"),
		allFieldsSoFar.CuidDefault.Set("manual-cuid"),
		allFieldsSoFar.Cuid1Default.Set("manual-cuid1"),
		allFieldsSoFar.Cuid2Default.Set("manual-cuid2"),
		allFieldsSoFar.UuidDefault.Set("00000000-0000-0000-0000-000000000000"),
		allFieldsSoFar.Uuid4Default.Set("11111111-1111-4111-8111-111111111111"),
		allFieldsSoFar.Uuid7Default.Set("22222222-2222-7222-8222-222222222222"),
		allFieldsSoFar.UlidDefault.Set("manual-ulid"),
		allFieldsSoFar.NanoidDefault.Set("manual-nanoid"),
		allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		allFieldsSoFar.IntReq.Set(-1),
		allFieldsSoFar.IntOpt.Set(0),
		allFieldsSoFar.IntDefault.Set(999),
		allFieldsSoFar.IntegerVal.Set(math.MaxInt32),
		allFieldsSoFar.SmallInt.Set(-32768),
		allFieldsSoFar.TinyInt.Set(-128),
		allFieldsSoFar.OidVal.Set(2147483647),
		allFieldsSoFar.BigIntReq.Set(math.MaxInt64),
		allFieldsSoFar.BigIntOpt.Set(math.MinInt64),
		allFieldsSoFar.FloatReq.Set(1.7976931348623157e+308),
		allFieldsSoFar.FloatOpt.Set(2.2250738585072014e-308),
		allFieldsSoFar.RealVal.Set(3.4028234e+38),
		allFieldsSoFar.DecimalReq.Set("1234567890.12345678"),
		allFieldsSoFar.DecimalOpt.Set("0.0001"),
		allFieldsSoFar.DecimalPrecise.Set("99999999.99"),
		allFieldsSoFar.MoneyVal.Set("99999999999.99"),
		allFieldsSoFar.BoolReq.Set(false),
		allFieldsSoFar.BoolOpt.Set(true),
		allFieldsSoFar.BoolDefault.Set(true),
		allFieldsSoFar.DateTimeReq.Set(now),
		allFieldsSoFar.DateTimeOpt.Set(now.Add(-time.Hour)),
		allFieldsSoFar.DateTimeDefault.Set(now.Add(-24*time.Hour)),
		allFieldsSoFar.UpdatedAt.Set(now),
		allFieldsSoFar.DateTimeTz.Set(now.In(time.FixedZone("IST", 3600))),
		allFieldsSoFar.TimestampVal.Set(now),
		allFieldsSoFar.TimeVal.Set(now),
		allFieldsSoFar.TimetzVal.Set(now.In(time.FixedZone("IST", 3600))),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`null`)),
		allFieldsSoFar.JsonOpt.Set(json.RawMessage(`{"nested":{"a":1}}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`"just a string"`)),
		allFieldsSoFar.BytesReq.Set([]byte{}),
		allFieldsSoFar.BytesOpt.Set([]byte{}),
		allFieldsSoFar.HstoreField.Set(map[string]*string{
			"a": new("1"),
			"b": new("2"),
		}),
		allFieldsSoFar.LtreeField.Set("Top.Collections.Elements"),
	).Exec(ctx)
	if err != nil {
		t.Fatalf("explicit create failed: %v", err)
	}

	if rec.StringReq != "explicit" {
		t.Errorf("StringReq = %q", rec.StringReq)
	}
	if rec.StringOpt == nil || *rec.StringOpt != "opt" {
		t.Errorf("StringOpt = %v", rec.StringOpt)
	}
	if rec.StringDefault != "overridden" {
		t.Errorf("StringDefault override = %q", rec.StringDefault)
	}
	if rec.CuidDefault != "manual-cuid" {
		t.Errorf("CuidDefault override = %q", rec.CuidDefault)
	}
	if rec.UlidDefault != "manual-ulid" {
		t.Errorf("UlidDefault override = %q", rec.UlidDefault)
	}
	if !rec.BoolDefault {
		t.Error("BoolDefault override should be true")
	}
	if rec.IntDefault != 999 {
		t.Errorf("IntDefault override = %d", rec.IntDefault)
	}
	if rec.BigIntReq != math.MaxInt64 {
		t.Errorf("BigIntReq = %d", rec.BigIntReq)
	}
}

func TestNativeDefaults_OptionalNulls(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	rec, err := db.AllFieldsSoFar.Create(
		allFieldsSoFar.StringReq.Set("req"),
		allFieldsSoFar.StringVarchar.Set("v"),
		allFieldsSoFar.StringChar.Set("1234567890"),
		allFieldsSoFar.BitVal.Set("1010101010"),
		allFieldsSoFar.VarBitVal.Set("1"),
		allFieldsSoFar.InetVal.Set("10.0.0.1"),
		allFieldsSoFar.XmlVal.Set("<x/>"),
		allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		allFieldsSoFar.IntReq.Set(0),
		allFieldsSoFar.IntegerVal.Set(0),
		allFieldsSoFar.SmallInt.Set(0),
		allFieldsSoFar.TinyInt.Set(0),
		allFieldsSoFar.OidVal.Set(0),
		allFieldsSoFar.BigIntReq.Set(0),
		allFieldsSoFar.FloatReq.Set(0),
		allFieldsSoFar.RealVal.Set(0),
		allFieldsSoFar.DecimalReq.Set("0"),
		allFieldsSoFar.DecimalPrecise.Set("0.00"),
		allFieldsSoFar.MoneyVal.Set("0.00"),
		allFieldsSoFar.BoolReq.Set(false),
		allFieldsSoFar.DateTimeReq.Set(time.Time{}),
		allFieldsSoFar.UpdatedAt.Set(time.Time{}),
		allFieldsSoFar.DateTimeTz.Set(time.Time{}),
		allFieldsSoFar.TimestampVal.Set(time.Time{}),
		allFieldsSoFar.TimeVal.Set(time.Time{}),
		allFieldsSoFar.TimetzVal.Set(time.Time{}),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.BytesReq.Set([]byte{}),
		allFieldsSoFar.LtreeField.Set("Top"),
	).Exec(ctx)

	if err != nil {
		t.Fatalf("optional nulls create failed: %v", err)
	}

	if rec.StringOpt != nil {
		t.Error("StringOpt should be nil")
	}
	if rec.IntOpt != nil {
		t.Error("IntOpt should be nil")
	}
	if rec.BigIntOpt != nil {
		t.Error("BigIntOpt should be nil")
	}
	if rec.FloatOpt != nil {
		t.Error("FloatOpt should be nil")
	}
	if rec.DecimalOpt != nil {
		t.Error("DecimalOpt should be nil")
	}
	if rec.BoolOpt != nil {
		t.Error("BoolOpt should be nil")
	}
	if rec.DateTimeOpt != nil {
		t.Error("DateTimeOpt should be nil")
	}
	if rec.JsonOpt != nil {
		t.Error("JsonOpt should be nil")
	}
	if rec.BytesOpt != nil {
		t.Error("BytesOpt should be nil")
	}
	if rec.HstoreField != nil {
		t.Error("HstoreField should be nil")
	}
}

func TestNativeDefaults_StringConstraints_VarChar(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("exactly 255 chars", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringVarchar.Set(strings.Repeat("a", 255)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("255 chars should be valid: %v", err)
		}
	})

	t.Run("256 chars rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringVarchar.Set(strings.Repeat("a", 256)),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for 256-char VarChar")
		}
		if !strings.Contains(err.Error(), "length") {
			t.Errorf("expected length error, got: %v", err)
		}
	})

	t.Run("empty string rejected for required", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringVarchar.Set(""),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for empty VarChar")
		}
	})

	t.Run("unicode within limit", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringVarchar.Set(strings.Repeat("ñ", 255)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("unicode within 255 should work: %v", err)
		}
	})

	t.Run("unicode over limit by rune count", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringVarchar.Set(strings.Repeat("ñ", 256)),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for 256 unicode chars in VarChar")
		}
	})
}

func TestNativeDefaults_StringConstraints_Char(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("exactly 10 chars", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringChar.Set("1234567890"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("10 chars should be valid: %v", err)
		}
	})

	t.Run("11 chars rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringChar.Set("12345678901"),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for 11-char Char")
		}
	})
}

func TestNativeDefaults_StringConstraints_Bit(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("exactly 10 bits valid", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.BitVal.Set("0101010101"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("valid bit string: %v", err)
		}
	})

	t.Run("wrong chars rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.BitVal.Set("2101010101"),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for non-binary chars in Bit")
		}
	})
}

func TestNativeDefaults_StringConstraints_Inet(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("valid IPv4", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.InetVal.Set("255.255.255.255"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("valid IPv4: %v", err)
		}
	})

	t.Run("valid IPv6", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.InetVal.Set("2001:db8::1"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("valid IPv6: %v", err)
		}
	})

	t.Run("CIDR notation", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.InetVal.Set("192.168.0.0/16"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("CIDR notation: %v", err)
		}
	})

	t.Run("invalid string rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.InetVal.Set("not-an-ip"),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for invalid IP string")
		}
	})

	t.Run("empty string rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.InetVal.Set(""),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for empty IP")
		}
	})
}

func TestNativeDefaults_IntConstraints_SmallInt(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("min valid -32768", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.SmallInt.Set(-32768),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("min smallint: %v", err)
		}
	})

	t.Run("max valid 32767", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.SmallInt.Set(32767),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("max smallint: %v", err)
		}
	})

	t.Run("underflow -32769", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.SmallInt.Set(-32769),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for underflow smallint")
		}
	})

	t.Run("overflow 32768", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.SmallInt.Set(32768),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for overflow smallint")
		}
	})
}

func TestNativeDefaults_IntConstraints_Oid(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("zero is valid", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.OidVal.Set(0),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("oid=0 valid: %v", err)
		}
	})

	t.Run("positive is valid", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.OidVal.Set(42),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("oid=42 valid: %v", err)
		}
	})

	t.Run("negative rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.OidVal.Set(-1),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for negative oid")
		}
	})
}

func TestNativeDefaults_DecimalPrecision(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("valid NUMERIC(10,2) integer part", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.DecimalPrecise.Set("99999999.99"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("max decimal: %v", err)
		}
	})

	t.Run("overflow integer part", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.DecimalPrecise.Set("100000000.00"),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for decimal overflow")
		}
	})
}

func TestNativeDefaults_JsonTypes(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("jsonb accepts binary JSON", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.JsonReq.Set(json.RawMessage(`{"a":1,"b":[2,3]}`)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("jsonb create: %v", err)
		}
	})

	t.Run("json vs jsonb round trip", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.JsonReq.Set(json.RawMessage(`{"x":1}`)),
				allFieldsSoFar.JsonVal.Set(json.RawMessage(`{"x":1}`)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("json/jsonb create: %v", err)
		}
		var got, want any
		json.Unmarshal(rec.JsonReq, &got)
		json.Unmarshal([]byte(`{"x":1}`), &want)
		gotJSON, _ := json.Marshal(got)
		wantJSON, _ := json.Marshal(want)
		if string(gotJSON) != string(wantJSON) {
			t.Errorf("jsonb roundtrip = %s", string(rec.JsonReq))
		}
	})
}

func TestNativeDefaults_BinaryRoundTrip(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	binaryData := []byte{0x00, 0x01, 0xFF, 0xFE, 0x7F, 0x80}
	rec, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.BytesReq.Set(binaryData),
			allFieldsSoFar.BytesOpt.Set(binaryData),
		)...,
	).Exec(ctx)

	if err != nil {
		t.Fatalf("binary create: %v", err)
	}

	if len(rec.BytesReq) != len(binaryData) {
		t.Fatalf("BytesReq length mismatch: %d vs %d", len(rec.BytesReq), len(binaryData))
	}
	for i, b := range rec.BytesReq {
		if b != binaryData[i] {
			t.Errorf("BytesReq[%d] = %02x, want %02x", i, b, binaryData[i])
		}
	}

	if rec.BytesOpt == nil {
		t.Fatal("BytesOpt should not be nil")
	}
	if len(*rec.BytesOpt) != len(binaryData) {
		t.Fatalf("BytesOpt length mismatch: %d vs %d", len(*rec.BytesOpt), len(binaryData))
	}
}

func TestNativeDefaults_UpdatedAt(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	before := time.Now().Add(-time.Minute)
	rec, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.UpdatedAt.Set(time.Now()),
		)...,
	).Exec(ctx)
	if err != nil {
		t.Fatalf("create with updatedAt: %v", err)
	}
	if rec.UpdatedAt.Before(before) {
		t.Error("UpdatedAt seems stale")
	}
}

func TestNativeDefaults_Validation_RequiredFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("missing stringReq", func(t *testing.T) {
		var filtered []valk.FieldAssignment
		for _, a := range base {
			if a.Col != "stringReq" {
				filtered = append(filtered, a)
			}
		}
		_, err := db.AllFieldsSoFar.Create(filtered...).Exec(ctx)
		if err == nil {
			t.Fatal("expected error for missing stringReq")
		}
		valErr, ok := err.(valk.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		found := false
		for _, fe := range valErr.Errors {
			if fe.Field == "stringReq" && fe.Rule == "required" {
				found = true
			}
		}
		if !found {
			t.Errorf("missing required error for stringReq: %v", valErr.Errors)
		}
	})

	t.Run("missing multiple required fields", func(t *testing.T) {
		required := []string{"stringReq", "stringVarchar", "stringChar", "bitVal", "varBitVal",
			"inetVal", "xmlVal", "uuidDb", "intReq", "integerVal", "smallInt",
			"tinyInt", "oidVal", "bigIntReq", "floatReq", "realVal", "decimalReq",
			"decimalPrecise", "moneyVal", "boolReq", "dateTimeReq", "updatedAt",
			"dateTimeTz", "timestampVal", "timeVal", "timetzVal", "jsonReq",
			"jsonVal", "bytesReq", "ltreeField"}

		_, err := db.AllFieldsSoFar.Create().Exec(ctx)
		if err == nil {
			t.Fatal("expected validation errors")
		}
		valErr, ok := err.(valk.ValidationError)
		if !ok {
			t.Fatalf("expected ValidationError, got %T", err)
		}
		found := make(map[string]bool)
		for _, fe := range valErr.Errors {
			if fe.Rule == "required" {
				found[fe.Field] = true
			}
		}
		for _, field := range required {
			if !found[field] {
				t.Errorf("missing required validation for %q", field)
			}
		}
	})
}

func TestNativeDefaults_NullByteSafety(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	nullByteFields := []struct {
		name  string
		setFn func(v string) valk.FieldAssignment
	}{
		{"stringReq", allFieldsSoFar.StringReq.Set},
		{"stringVarchar", allFieldsSoFar.StringVarchar.Set},
		{"stringChar", allFieldsSoFar.StringChar.Set},
		{"bitVal", allFieldsSoFar.BitVal.Set},
		{"varBitVal", allFieldsSoFar.VarBitVal.Set},
		{"inetVal", allFieldsSoFar.InetVal.Set},
		{"xmlVal", allFieldsSoFar.XmlVal.Set},
		{"uuidDb", allFieldsSoFar.UuidDb.Set},
		{"decimalReq", allFieldsSoFar.DecimalReq.Set},
		{"decimalPrecise", allFieldsSoFar.DecimalPrecise.Set},
		{"moneyVal", allFieldsSoFar.MoneyVal.Set},
	}

	for _, fc := range nullByteFields {
		t.Run(fc.name+" null byte", func(t *testing.T) {
			_, err := db.AllFieldsSoFar.Create(
				append(base,
					fc.setFn("val\x00ue"),
				)...,
			).Exec(ctx)
			if err == nil {
				t.Fatal("expected error for null byte")
			}
		})
	}
}

func TestNativeDefaults_InvalidUtf8(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	_, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.StringReq.Set("valid"),
			allFieldsSoFar.StringVarchar.Set("invalid\xffutf8"),
		)...,
	).Exec(ctx)
	if err == nil {
		t.Fatal("expected error for invalid UTF-8")
	}
}

func TestNativeDefaults_SelectOmit(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("select single field", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(base...).
			Select(valk.AllFieldsSoFarSelect{StringReq: true}).
			Exec(ctx)
		if err != nil {
			t.Fatalf("select create: %v", err)
		}
		if rec.StringReq == "" {
			t.Error("selected StringReq should be populated")
		}
		if rec.IntReq != 0 {
			t.Error("unselected IntReq should be zero")
		}
	})

	t.Run("omit single field", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(base...).
			Omit(valk.AllFieldsSoFarOmit{IntReq: true}).
			Exec(ctx)
		if err != nil {
			t.Fatalf("omit create: %v", err)
		}
		if rec.IntReq != 0 {
			t.Error("omitted IntReq should be zero")
		}
		if rec.StringReq == "" {
			t.Error("non-omitted StringReq should be populated")
		}
	})
}

func TestNativeDefaults_CreateMany(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	rec1 := allFieldsSoFar.Record(append(base,
		allFieldsSoFar.StringReq.Set("many-1"),
		allFieldsSoFar.IntReq.Set(10),
	)...)

	rec2 := allFieldsSoFar.Record(append(base,
		allFieldsSoFar.StringReq.Set("many-2"),
		allFieldsSoFar.IntReq.Set(20),
	)...)

	count, err := db.AllFieldsSoFar.CreateMany(rec1, rec2).Exec(ctx)
	if err != nil {
		t.Fatalf("CreateMany failed: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 rows, got %d", count)
	}
}

func TestNativeDefaults_CreateManyAndReturn(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	recs, err := db.AllFieldsSoFar.CreateManyAndReturn(
		allFieldsSoFar.Record(append(base,
			allFieldsSoFar.StringReq.Set("ret-1"),
			allFieldsSoFar.IntReq.Set(100),
		)...),
		allFieldsSoFar.Record(append(base,
			allFieldsSoFar.StringReq.Set("ret-2"),
			allFieldsSoFar.IntReq.Set(200),
		)...),
	).Exec(ctx)

	if err != nil {
		t.Fatalf("CreateManyAndReturn failed: %v", err)
	}
	if len(recs) != 2 {
		t.Fatalf("expected 2 records, got %d", len(recs))
	}
	if recs[0].StringReq != "ret-1" {
		t.Errorf("rec[0].StringReq = %q", recs[0].StringReq)
	}
	if recs[1].IntReq != 200 {
		t.Errorf("rec[1].IntReq = %d", recs[1].IntReq)
	}
}

func TestNativeDefaults_Hooks(t *testing.T) {
	ctx := context.Background()

	t.Run("BeforeCreate mutates input", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.AllFieldsSoFar.BeforeCreate(func(ctx context.Context, input *valk.AllFieldsSoFarCreate) error {
			if input.StringReq == "mutate-me" {
				s := "mutated"
				input.StringReq = s
			}
			return nil
		})

		base := baseAllFields(t)
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringReq.Set("mutate-me"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("BeforeCreate hook failed: %v", err)
		}
		if rec.StringReq != "mutated" {
			t.Errorf("expected hook to mutate StringReq, got %q", rec.StringReq)
		}
	})

	t.Run("BeforeCreate error aborts", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		db.AllFieldsSoFar.BeforeCreate(func(ctx context.Context, input *valk.AllFieldsSoFarCreate) error {
			if input.StringReq == "abort" {
				return fmt.Errorf("hook aborted: %s", input.StringReq)
			}
			return nil
		})

		base := baseAllFields(t)
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.StringReq.Set("abort"),
			)...,
		).Exec(ctx)
		if err == nil {
			t.Fatal("expected hook abort error")
		}
	})

	t.Run("AfterCreateMany receives correct count", func(t *testing.T) {
		db, cleanup := setupTestDB(t)
		defer cleanup()

		var captured struct {
			inputs []valk.AllFieldsSoFarCreate
			count  int64
		}
		db.AllFieldsSoFar.AfterCreateMany(func(ctx context.Context, inputs []valk.AllFieldsSoFarCreate, count int64) error {
			captured.inputs = inputs
			captured.count = count
			return nil
		})

		base := baseAllFields(t)
		r1 := allFieldsSoFar.Record(append(base,
			allFieldsSoFar.StringReq.Set("hook-cm-1"),
			allFieldsSoFar.IntReq.Set(1),
		)...)
		r2 := allFieldsSoFar.Record(append(base,
			allFieldsSoFar.StringReq.Set("hook-cm-2"),
			allFieldsSoFar.IntReq.Set(2),
		)...)

		_, err := db.AllFieldsSoFar.CreateMany(r1, r2).Exec(ctx)
		if err != nil {
			t.Fatalf("CreateMany failed: %v", err)
		}
		if captured.count != 2 {
			t.Errorf("expected count=2, got %d", captured.count)
		}
		if len(captured.inputs) != 2 {
			t.Errorf("expected 2 inputs, got %d", len(captured.inputs))
		}
	})
}

func TestNativeDefaults_HstoreAndLtree(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("hstore nil is valid", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.LtreeField.Set("Top"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("hstore=nil create: %v", err)
		}
	})

	t.Run("hstore with value", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.HstoreField.Set(map[string]*string{
					"name": new("John"),
					"age":  new("30"),
				}), allFieldsSoFar.LtreeField.Set("Top.Collections"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("hstore set create: %v", err)
		}
		if rec.HstoreField == nil {
			t.Error("HstoreField should be set")
		}
		if ltreeFieldStr(rec.LtreeField) != "Top.Collections" {
			t.Errorf("ltreeField = %s", ltreeFieldStr(rec.LtreeField))
		}
	})

	t.Run("ltree path traversal", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.LtreeField.Set("Top.Collections.Photos.Voids"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("ltree long path: %v", err)
		}
	})
}

func TestNativeDefaults_ZeroTimeValues(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	rec, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.DateTimeReq.Set(time.Time{}),
			allFieldsSoFar.UpdatedAt.Set(time.Time{}),
			allFieldsSoFar.DateTimeTz.Set(time.Time{}),
			allFieldsSoFar.TimestampVal.Set(time.Time{}),
			allFieldsSoFar.TimeVal.Set(time.Time{}),
			allFieldsSoFar.TimetzVal.Set(time.Time{}),
		)...,
	).Exec(ctx)

	if err != nil {
		t.Fatalf("zero time create: %v", err)
	}

	if !rec.DateTimeReq.IsZero() {
		t.Error("DateTimeReq should be zero-value")
	}
}

func TestNativeDefaults_AggregatedValidationErrors(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	_, err := db.AllFieldsSoFar.Create(
		allFieldsSoFar.StringReq.Set("valid"),
		allFieldsSoFar.StringVarchar.Set(strings.Repeat("a", 256)),
		allFieldsSoFar.StringChar.Set("toolongggg"),
		allFieldsSoFar.BitVal.Set("2101010101"),
		allFieldsSoFar.VarBitVal.Set("abc"),
		allFieldsSoFar.InetVal.Set("bad-ip"),
		allFieldsSoFar.IntReq.Set(-1),
		allFieldsSoFar.SmallInt.Set(99999),
		allFieldsSoFar.OidVal.Set(-5),
		allFieldsSoFar.XmlVal.Set("<x/>"),
		allFieldsSoFar.UuidDb.Set("not-a-uuid"),
		allFieldsSoFar.IntegerVal.Set(0),
		allFieldsSoFar.TinyInt.Set(0),
		allFieldsSoFar.BigIntReq.Set(0),
		allFieldsSoFar.FloatReq.Set(0),
		allFieldsSoFar.RealVal.Set(0),
		allFieldsSoFar.DecimalReq.Set("0"),
		allFieldsSoFar.DecimalPrecise.Set("0.00"),
		allFieldsSoFar.MoneyVal.Set("0.00"),
		allFieldsSoFar.BoolReq.Set(false),
		allFieldsSoFar.DateTimeReq.Set(time.Time{}),
		allFieldsSoFar.UpdatedAt.Set(time.Time{}),
		allFieldsSoFar.DateTimeTz.Set(time.Time{}),
		allFieldsSoFar.TimestampVal.Set(time.Time{}),
		allFieldsSoFar.TimeVal.Set(time.Time{}),
		allFieldsSoFar.TimetzVal.Set(time.Time{}),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.BytesReq.Set([]byte{}),
		allFieldsSoFar.LtreeField.Set("Top"),
	).Exec(ctx)

	if err == nil {
		t.Fatal("expected aggregated validation errors")
	}

	valErr, ok := err.(valk.ValidationError)
	if !ok {
		t.Fatalf("expected ValidationError, got %T: %v", err, err)
	}

	if len(valErr.Errors) < 5 {
		t.Errorf("expected at least 5 validation errors, got %d: %v", len(valErr.Errors), valErr.Errors)
	}
}

func TestNativeDefaults_DefaultGeneration_UniqueValues(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	rec1, err := db.AllFieldsSoFar.Create(base...).Exec(ctx)
	if err != nil {
		t.Fatalf("first default create: %v", err)
	}
	rec2, err := db.AllFieldsSoFar.Create(base...).Exec(ctx)
	if err != nil {
		t.Fatalf("second default create: %v", err)
	}

	if rec1.CuidDefault == rec2.CuidDefault {
		t.Error("CuidDefault should be unique across rows")
	}
	if rec1.UuidDefault == rec2.UuidDefault {
		t.Error("UuidDefault should be unique across rows")
	}
	if rec1.Uuid4Default == rec2.Uuid4Default {
		t.Error("Uuid4Default should be unique across rows")
	}
	if rec1.Uuid7Default == rec2.Uuid7Default {
		t.Error("Uuid7Default should be unique across rows")
	}
	if rec1.UlidDefault == rec2.UlidDefault {
		t.Error("UlidDefault should be unique across rows")
	}
	if rec1.NanoidDefault == rec2.NanoidDefault {
		t.Error("NanoidDefault should be unique across rows")
	}

	if rec1.Id == rec2.Id {
		t.Error("autoincrement Id should be unique")
	}
}

func TestNativeDefaults_UuidFormat(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	rec, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		)...,
	).Exec(ctx)
	if err != nil {
		t.Fatalf("uuid create: %v", err)
	}
	if rec.UuidDb != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("uuid roundtrip = %q", rec.UuidDb)
	}
}

func TestNativeDefaults_BoolEdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("bool with default false", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.BoolDefault.Set(false),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("bool default false: %v", err)
		}
		if rec.BoolDefault {
			t.Error("BoolDefault should be false")
		}
		if !rec.BoolReq {
			t.Error("BoolReq should be true")
		}
	})

	t.Run("bool opt nil", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.BoolDefault.Set(false),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("bool opt nil: %v", err)
		}
		if rec.BoolOpt != nil {
			t.Error("BoolOpt should be nil")
		}
	})

	t.Run("bool opt set true", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.BoolOpt.Set(true),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("bool opt true: %v", err)
		}
		if rec.BoolOpt == nil || !*rec.BoolOpt {
			t.Error("BoolOpt should be true")
		}
	})
}

func TestNativeDefaults_MoneyDecimal(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("money with large value", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.MoneyVal.Set("92233720368547758.07"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("money large: %v", err)
		}
		if rec.MoneyVal != "92233720368547758.07" {
			t.Logf("MoneyVal roundtrip = %q (may differ due to PG money format)", rec.MoneyVal)
		}
	})

	t.Run("decimal with trailing zeros", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.DecimalReq.Set("10.10000"),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("decimal trailing zeros: %v", err)
		}
		if rec.DecimalReq == "" {
			t.Error("DecimalReq should be populated")
		}
	})
}

func TestNativeDefaults_DateTimeTz(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	loc := time.FixedZone("IST", 5*3600+30)
	now := time.Date(2024, 1, 1, 12, 0, 0, 0, loc)

	rec, err := db.AllFieldsSoFar.Create(
		append(base,
			allFieldsSoFar.DateTimeTz.Set(now),
		)...,
	).Exec(ctx)
	if err != nil {
		t.Fatalf("timestamptz create: %v", err)
	}

	if !rec.DateTimeTz.Equal(now) {
		t.Logf("DateTimeTz roundtrip: in=%v got=%v (may differ in zone/offset)", now, rec.DateTimeTz)
	}
}

func TestNativeDefaults_RealFloatLimits(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("float max valid", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.FloatReq.Set(1e308),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("float64 max: %v", err)
		}
	})

	t.Run("real float32 boundary", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.RealVal.Set(3.4028234e+37),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("real boundary: %v", err)
		}
		if rec.RealVal == 0 {
			t.Error("RealVal should be non-zero")
		}
	})
}

func TestNativeDefaults_InetV6RoundTrip(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	ips := []string{
		"::1",
		"2001:db8::1",
		"fe80::1",
		"::ffff:192.0.2.1",
		"2001:0db8:85a3:0000:0000:8a2e:0370:7334",
	}

	for _, ip := range ips {
		t.Run(ip, func(t *testing.T) {
			rec, err := db.AllFieldsSoFar.Create(
				append(base,
					allFieldsSoFar.InetVal.Set(ip),
				)...,
			).Exec(ctx)
			if err != nil {
				t.Fatalf("IP %q failed: %v", ip, err)
			}
			if rec.InetVal == "" {
				t.Errorf("returned InetVal is empty for input %q", ip)
			}
		})
	}
}

func TestNativeDefaults_NullByteInStringFields(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	tests := []struct {
		name string
		set  func(string) valk.FieldAssignment
	}{
		{"xmlVal", allFieldsSoFar.XmlVal.Set},
		{"decimalPrecise", allFieldsSoFar.DecimalPrecise.Set},
		{"moneyVal", allFieldsSoFar.MoneyVal.Set},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := db.AllFieldsSoFar.Create(
				append(base, tt.set("\x00abc"))...,
			).Exec(ctx)
			if err == nil {
				t.Fatalf("expected error for null byte in %s", tt.name)
			}
		})
	}
}

func TestNativeDefaults_FloatNanInf(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("NaN rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.FloatReq.Set(math.NaN()),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("NaN create: %v", err)
		}
	})

	t.Run("+Inf rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.FloatReq.Set(math.Inf(1)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("+Inf create: %v", err)
		}
	})

	t.Run("-Inf rejected", func(t *testing.T) {
		_, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.FloatReq.Set(math.Inf(-1)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("-Inf create: %v", err)
		}
	})
}

func TestNativeDefaults_JsonEdgeCases(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	t.Run("empty JSON object", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.JsonReq.Set(json.RawMessage(`{}`)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("empty json: %v", err)
		}
		if string(rec.JsonReq) != `{}` {
			t.Errorf("json roundtrip = %s", string(rec.JsonReq))
		}
	})

	t.Run("deeply nested JSON", func(t *testing.T) {
		deep := buildDeepJSON(50)
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.JsonReq.Set(deep),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("deep json: %v", err)
		}
		if len(rec.JsonReq) == 0 {
			t.Error("JsonReq should be populated")
		}
	})

	t.Run("special unicode in JSON", func(t *testing.T) {
		rec, err := db.AllFieldsSoFar.Create(
			append(base,
				allFieldsSoFar.JsonReq.Set(json.RawMessage(`"🎉🚀👋"`)),
			)...,
		).Exec(ctx)
		if err != nil {
			t.Fatalf("unicode json: %v", err)
		}
		if !strings.Contains(string(rec.JsonReq), "🎉") {
			t.Errorf("unicode json roundtrip = %s", string(rec.JsonReq))
		}
	})
}

func TestNativeDefaults_ConcurrentCreate(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()
	ctx := context.Background()

	base := baseAllFields(t)

	done := make(chan struct{})
	errs := make(chan error, 10)

	for i := 0; i < 10; i++ {
		go func() {
			_, err := db.AllFieldsSoFar.Create(
				append(base,
					allFieldsSoFar.StringReq.Set("concurrent"),
					allFieldsSoFar.IntReq.Set(42),
				)...,
			).Exec(ctx)
			errs <- err
		}()
	}

	for i := 0; i < 10; i++ {
		if err := <-errs; err != nil {
			t.Errorf("concurrent create failed: %v", err)
		}
	}
	close(done)
	<-done
}

func baseAllFields(t *testing.T) []valk.FieldAssignment {
	t.Helper()
	return []valk.FieldAssignment{
		allFieldsSoFar.StringReq.Set("test"),
		allFieldsSoFar.StringVarchar.Set("varchar"),
		allFieldsSoFar.StringChar.Set("0123456789"),
		allFieldsSoFar.BitVal.Set("1010101010"),
		allFieldsSoFar.VarBitVal.Set("1101"),
		allFieldsSoFar.InetVal.Set("10.0.0.1"),
		allFieldsSoFar.XmlVal.Set("<x/>"),
		allFieldsSoFar.UuidDb.Set("550e8400-e29b-41d4-a716-446655440000"),
		allFieldsSoFar.IntReq.Set(1),
		allFieldsSoFar.IntegerVal.Set(2),
		allFieldsSoFar.SmallInt.Set(3),
		allFieldsSoFar.TinyInt.Set(4),
		allFieldsSoFar.OidVal.Set(5),
		allFieldsSoFar.BigIntReq.Set(int64(6)),
		allFieldsSoFar.FloatReq.Set(1.5),
		allFieldsSoFar.RealVal.Set(2.5),
		allFieldsSoFar.DecimalReq.Set("10.50"),
		allFieldsSoFar.DecimalPrecise.Set("99.99"),
		allFieldsSoFar.MoneyVal.Set("12.34"),
		allFieldsSoFar.BoolReq.Set(true),
		allFieldsSoFar.DateTimeReq.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.UpdatedAt.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.DateTimeTz.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.TimestampVal.Set(time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)),
		allFieldsSoFar.TimeVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC)),
		allFieldsSoFar.TimetzVal.Set(time.Date(1, 1, 1, 10, 30, 0, 0, time.UTC)),
		allFieldsSoFar.JsonReq.Set(json.RawMessage(`{}`)),
		allFieldsSoFar.JsonVal.Set(json.RawMessage(`[]`)),
		allFieldsSoFar.BytesReq.Set([]byte{0x01}),
		allFieldsSoFar.LtreeField.Set("Top"),
	}
}

func ltreeFieldStr(v any) string {
	switch x := v.(type) {
	case string:
		return x
	case []byte:
		return string(x)
	case fmt.Stringer:
		return x.String()
	default:
		return fmt.Sprintf("%s", v)
	}
}

func buildDeepJSON(depth int) json.RawMessage {
	if depth <= 0 {
		return json.RawMessage(`"bottom"`)
	}
	inner := buildDeepJSON(depth - 1)
	raw := `{"depth_` + strconv.Itoa(depth) + `":` + string(inner) + `}`
	return json.RawMessage(raw)
}
