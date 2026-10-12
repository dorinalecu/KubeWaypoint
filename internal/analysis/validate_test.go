package analysis

import (
	"errors"
	"strings"
	"testing"
)

const secretMarker = "SECRET-MARKER-7f3a"

func strPtr(s string) *string { return &s }

func TestInputValidateValid(t *testing.T) {
	// Each non-ASCII rune below is 2 bytes ("é"), 3 bytes ("日") or 4 bytes ("😀").
	tests := []struct {
		name string
		in   Input
	}{
		{"description only", Input{Description: "Node worker-1 never joined the cluster."}},
		{"nil optionals", Input{Description: "x", CaseID: nil, Metadata: nil}},
		{"empty metadata", Input{Description: "x", Metadata: map[string]string{}}},
		{"description surrounded by whitespace", Input{Description: " \t\n text \r\n"}},
		{"description at limit", Input{Description: strings.Repeat("a", MaxDescriptionBytes)}},
		{"description at limit in 2-byte runes", Input{Description: strings.Repeat("é", MaxDescriptionBytes/2)}},
		{"description at limit in 3-byte runes plus padding", Input{Description: strings.Repeat("日", MaxDescriptionBytes/3) + "a"}},
		{"description limit counts bytes before trimming", Input{Description: strings.Repeat(" ", MaxDescriptionBytes-1) + "a"}},
		{"case id at limit", Input{Description: "x", CaseID: strPtr(strings.Repeat("c", MaxCaseIDBytes))}},
		{"case id at limit in 4-byte runes", Input{Description: "x", CaseID: strPtr(strings.Repeat("😀", MaxCaseIDBytes/4))}},
		{"case id with surrounding whitespace", Input{Description: "x", CaseID: strPtr(" case-1 ")}},
		{"metadata empty value", Input{Description: "x", Metadata: map[string]string{"k": ""}}},
		{"metadata key at limit", Input{Description: "x", Metadata: map[string]string{strings.Repeat("k", MaxMetadataKeyBytes): "v"}}},
		{"metadata key at limit in 3-byte runes", Input{Description: "x", Metadata: map[string]string{strings.Repeat("日", MaxMetadataKeyBytes/3) + "k": "v"}}},
		{"metadata value at limit", Input{Description: "x", Metadata: map[string]string{"k": strings.Repeat("v", MaxMetadataValueBytes)}}},
		{"metadata value at limit in 2-byte runes", Input{Description: "x", Metadata: map[string]string{"k": strings.Repeat("é", MaxMetadataValueBytes/2)}}},
		{"metadata 16 entries", Input{Description: "x", Metadata: metadataWithEntries(MaxMetadataEntries)}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.in.Validate(); err != nil {
				t.Fatalf("Validate() = %v, want nil", err)
			}
		})
	}
}

func TestInputValidateInvalid(t *testing.T) {
	bad := "\xff\xfe"
	tests := []struct {
		name string
		in   Input
		want ValidationError
	}{
		{"empty description", Input{Description: ""}, ValidationError{FieldDescription, CodeBlank}},
		{"whitespace description", Input{Description: " \t\r\n "}, ValidationError{FieldDescription, CodeBlank}},
		{"unicode whitespace description", Input{Description: "\u00a0\u2003"}, ValidationError{FieldDescription, CodeBlank}},
		{"description invalid UTF-8", Input{Description: "abc" + bad + secretMarker}, ValidationError{FieldDescription, CodeInvalidUTF8}},
		{"description truncated rune", Input{Description: "abc\xe6\x97"}, ValidationError{FieldDescription, CodeInvalidUTF8}},
		{"description one byte over limit", Input{Description: strings.Repeat("a", MaxDescriptionBytes+1)}, ValidationError{FieldDescription, CodeTooLong}},
		{"description over limit in 2-byte runes", Input{Description: strings.Repeat("é", MaxDescriptionBytes/2) + "a"}, ValidationError{FieldDescription, CodeTooLong}},
		{"description over limit counts untrimmed bytes", Input{Description: strings.Repeat(" ", MaxDescriptionBytes) + "a"}, ValidationError{FieldDescription, CodeTooLong}},
		{"explicit empty case id", Input{Description: "x", CaseID: strPtr("")}, ValidationError{FieldCaseID, CodeBlank}},
		{"whitespace case id", Input{Description: "x", CaseID: strPtr(" \t ")}, ValidationError{FieldCaseID, CodeBlank}},
		{"case id invalid UTF-8", Input{Description: "x", CaseID: strPtr("id" + bad + secretMarker)}, ValidationError{FieldCaseID, CodeInvalidUTF8}},
		{"case id one byte over limit", Input{Description: "x", CaseID: strPtr(strings.Repeat("c", MaxCaseIDBytes+1))}, ValidationError{FieldCaseID, CodeTooLong}},
		{"case id over limit in 4-byte runes", Input{Description: "x", CaseID: strPtr(strings.Repeat("😀", MaxCaseIDBytes/4) + "a")}, ValidationError{FieldCaseID, CodeTooLong}},
		{"metadata 17 entries", Input{Description: "x", Metadata: metadataWithEntries(MaxMetadataEntries + 1)}, ValidationError{FieldMetadata, CodeTooMany}},
		{"metadata empty key", Input{Description: "x", Metadata: map[string]string{"": secretMarker}}, ValidationError{FieldMetadataKey, CodeEmpty}},
		{"metadata key invalid UTF-8", Input{Description: "x", Metadata: map[string]string{"k" + bad + secretMarker: "v"}}, ValidationError{FieldMetadataKey, CodeInvalidUTF8}},
		{"metadata key one byte over limit", Input{Description: "x", Metadata: map[string]string{strings.Repeat("k", MaxMetadataKeyBytes+1): "v"}}, ValidationError{FieldMetadataKey, CodeTooLong}},
		{"metadata key over limit in 3-byte runes", Input{Description: "x", Metadata: map[string]string{strings.Repeat("日", MaxMetadataKeyBytes/3+1): "v"}}, ValidationError{FieldMetadataKey, CodeTooLong}},
		{"metadata value invalid UTF-8", Input{Description: "x", Metadata: map[string]string{"k": "v" + bad + secretMarker}}, ValidationError{FieldMetadataValue, CodeInvalidUTF8}},
		{"metadata value one byte over limit", Input{Description: "x", Metadata: map[string]string{"k": strings.Repeat("v", MaxMetadataValueBytes+1)}}, ValidationError{FieldMetadataValue, CodeTooLong}},
		{"metadata value over limit in 2-byte runes", Input{Description: "x", Metadata: map[string]string{"k": strings.Repeat("é", MaxMetadataValueBytes/2) + "a"}}, ValidationError{FieldMetadataValue, CodeTooLong}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			var ve ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("Validate() = %v, want ValidationError", err)
			}
			if ve != tc.want {
				t.Fatalf("Validate() = %+v, want %+v", ve, tc.want)
			}
			assertSafeError(t, err)
		})
	}
}

func TestInputValidateFieldPrecedence(t *testing.T) {
	in := Input{
		Description: " ",
		CaseID:      strPtr(""),
		Metadata:    metadataWithEntries(MaxMetadataEntries + 1),
	}
	want := ValidationError{FieldDescription, CodeBlank}
	if err := in.Validate(); err != want {
		t.Fatalf("Validate() = %v, want %v", err, want)
	}
	in.Description = "x"
	want = ValidationError{FieldCaseID, CodeBlank}
	if err := in.Validate(); err != want {
		t.Fatalf("Validate() = %v, want %v", err, want)
	}
	in.CaseID = nil
	want = ValidationError{FieldMetadata, CodeTooMany}
	if err := in.Validate(); err != want {
		t.Fatalf("Validate() = %v, want %v", err, want)
	}
}

func TestInputValidateMetadataDeterministic(t *testing.T) {
	// The first failing entry in sorted key order is reported, regardless of
	// Go's randomized map iteration order.
	long := strings.Repeat("k", MaxMetadataKeyBytes+1)
	tests := []struct {
		name string
		md   map[string]string
		want ValidationError
	}{
		{"empty key first", map[string]string{"": "v", long: "v", "b": strings.Repeat("v", MaxMetadataValueBytes+1)}, ValidationError{FieldMetadataKey, CodeEmpty}},
		{"long key before later bad value", map[string]string{long: "v", "z": "\xff"}, ValidationError{FieldMetadataKey, CodeTooLong}},
		{"first sorted key's value wins", map[string]string{"a": "\xff", "b": strings.Repeat("v", MaxMetadataValueBytes+1)}, ValidationError{FieldMetadataValue, CodeInvalidUTF8}},
		{"key error of later key after valid earlier entry", map[string]string{"a": "ok", "b\xff": "v"}, ValidationError{FieldMetadataKey, CodeInvalidUTF8}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := Input{Description: "x", Metadata: tc.md}
			for i := 0; i < 50; i++ {
				if err := in.Validate(); err != tc.want {
					t.Fatalf("iteration %d: Validate() = %v, want %v", i, err, tc.want)
				}
			}
		})
	}
}

func TestInputValidateDoesNotMutate(t *testing.T) {
	caseID := " case-1 "
	md := map[string]string{"b": "2", "a": "1", "c\xff": "3"}
	in := Input{Description: "  text  ", CaseID: &caseID, Metadata: md}
	if err := in.Validate(); err == nil {
		t.Fatal("Validate() = nil, want error")
	}
	if in.Description != "  text  " || caseID != " case-1 " || in.CaseID != &caseID {
		t.Fatal("Validate mutated description or case id")
	}
	if len(md) != 3 || md["a"] != "1" || md["b"] != "2" || md["c\xff"] != "3" {
		t.Fatalf("Validate mutated metadata: %v", md)
	}
}

func TestValidationErrorText(t *testing.T) {
	err := ValidationError{FieldMetadataKey, CodeTooLong}
	if got, want := err.Error(), "invalid input: field=metadata.key code=too_long"; got != want {
		t.Fatalf("Error() = %q, want %q", got, want)
	}
}

func assertSafeError(t *testing.T, err error) {
	t.Helper()
	var ve ValidationError
	if !errors.As(err, &ve) {
		t.Fatalf("not a ValidationError: %v", err)
	}
	for _, s := range []string{err.Error(), ve.Field, ve.Code} {
		if strings.Contains(s, secretMarker) || strings.Contains(s, "\xff") {
			t.Fatalf("error leaks input: %q", s)
		}
	}
}

func metadataWithEntries(n int) map[string]string {
	md := make(map[string]string, n)
	for i := 0; i < n; i++ {
		md[strings.Repeat("k", i+1)] = "v"
	}
	return md
}
