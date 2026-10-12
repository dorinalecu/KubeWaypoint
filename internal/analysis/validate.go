package analysis

import (
	"sort"
	"strings"
	"unicode/utf8"
)

// Input limits, measured in raw UTF-8 bytes.
const (
	MaxDescriptionBytes   = 16 * 1024
	MaxCaseIDBytes        = 128
	MaxMetadataEntries    = 16
	MaxMetadataKeyBytes   = 64
	MaxMetadataValueBytes = 256
)

// Fixed values for ValidationError.Field.
const (
	FieldDescription   = "description"
	FieldCaseID        = "case_id"
	FieldMetadata      = "metadata"
	FieldMetadataKey   = "metadata.key"
	FieldMetadataValue = "metadata.value"
)

// Fixed values for ValidationError.Code.
const (
	CodeInvalidUTF8 = "invalid_utf8"
	CodeBlank       = "blank"
	CodeEmpty       = "empty"
	CodeTooLong     = "too_long"
	CodeTooMany     = "too_many_entries"
)

// ValidationError reports the first failed input constraint. It never carries
// supplied values, so its text is safe to log or return.
type ValidationError struct {
	Field string
	Code  string
}

func (e ValidationError) Error() string {
	return "invalid input: field=" + e.Field + " code=" + e.Code
}

// Validate checks the input against the domain constraints. It returns nil
// only for a valid input and does not modify it. Checks run in a fixed order:
// description, case_id, then metadata entry count and entries by sorted key.
func (in Input) Validate() error {
	if len(in.Description) > MaxDescriptionBytes {
		return ValidationError{FieldDescription, CodeTooLong}
	}
	if !utf8.ValidString(in.Description) {
		return ValidationError{FieldDescription, CodeInvalidUTF8}
	}
	if strings.TrimSpace(in.Description) == "" {
		return ValidationError{FieldDescription, CodeBlank}
	}

	if in.CaseID != nil {
		id := *in.CaseID
		if len(id) > MaxCaseIDBytes {
			return ValidationError{FieldCaseID, CodeTooLong}
		}
		if !utf8.ValidString(id) {
			return ValidationError{FieldCaseID, CodeInvalidUTF8}
		}
		if strings.TrimSpace(id) == "" {
			return ValidationError{FieldCaseID, CodeBlank}
		}
	}

	if len(in.Metadata) > MaxMetadataEntries {
		return ValidationError{FieldMetadata, CodeTooMany}
	}
	keys := make([]string, 0, len(in.Metadata))
	for k := range in.Metadata {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		switch {
		case k == "":
			return ValidationError{FieldMetadataKey, CodeEmpty}
		case len(k) > MaxMetadataKeyBytes:
			return ValidationError{FieldMetadataKey, CodeTooLong}
		case !utf8.ValidString(k):
			return ValidationError{FieldMetadataKey, CodeInvalidUTF8}
		}
		v := in.Metadata[k]
		switch {
		case len(v) > MaxMetadataValueBytes:
			return ValidationError{FieldMetadataValue, CodeTooLong}
		case !utf8.ValidString(v):
			return ValidationError{FieldMetadataValue, CodeInvalidUTF8}
		}
	}
	return nil
}
