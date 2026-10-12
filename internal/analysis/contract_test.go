package analysis

import (
	"encoding/json"
	"reflect"
	"sort"
	"strings"
	"testing"
)

// These fixtures pin the JSON shape only. They are static samples, not output
// of an analyzer.
func populatedResultFixture() Result {
	return Result{
		RulesVersion:           RulesVersion,
		RequiresEngineerReview: true,
		Summary:                "Synthetic fixture summary.",
		Category:               CategoryDNS,
		Confidence:             ConfidenceSupported,
		ClassificationEvidence: []ClassificationEvidence{{RuleID: "fixture.rule", SourceExcerpt: "synthetic excerpt"}},
		ReportedFacts:          []ReportedFact{{Statement: "synthetic statement", SourceExcerpt: "synthetic excerpt"}},
		MissingInformation:     []InformationGap{{ID: "gap-1", Description: "synthetic gap"}},
		ClarifyingQuestions:    []ClarifyingQuestion{{ID: "q-1", InformationID: "gap-1", Question: "synthetic question?"}},
		InitialChecks:          []InitialCheck{{ID: "c-1", Description: "synthetic check", Purpose: "synthetic purpose", ExpectedEvidence: "synthetic evidence"}},
	}
}

func emptyCollectionsResultFixture() Result {
	return Result{
		RulesVersion:           RulesVersion,
		RequiresEngineerReview: true,
		Summary:                "Synthetic unresolved fixture.",
		Category:               CategoryUnknown,
		Confidence:             ConfidenceUnresolved,
		ClassificationEvidence: []ClassificationEvidence{},
		ReportedFacts:          []ReportedFact{},
		MissingInformation:     []InformationGap{},
		ClarifyingQuestions:    []ClarifyingQuestion{},
		InitialChecks:          []InitialCheck{},
	}
}

func sortedKeys(m map[string]json.RawMessage) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func objectKeys(t *testing.T, raw json.RawMessage) []string {
	t.Helper()
	var m map[string]json.RawMessage
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatalf("not a JSON object: %v: %s", err, raw)
	}
	return sortedKeys(m)
}

func TestResultJSONFields(t *testing.T) {
	data, err := json.Marshal(populatedResultFixture())
	if err != nil {
		t.Fatal(err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatal(err)
	}

	wantTop := []string{
		"category", "clarifying_questions", "classification_evidence", "confidence",
		"initial_checks", "missing_information", "reported_facts",
		"requires_engineer_review", "rules_version", "summary",
	}
	if got := sortedKeys(top); !reflect.DeepEqual(got, wantTop) {
		t.Fatalf("top-level fields = %v, want %v", got, wantTop)
	}

	nested := map[string][]string{
		"classification_evidence": {"rule_id", "source_excerpt"},
		"reported_facts":          {"source_excerpt", "statement"},
		"missing_information":     {"description", "id"},
		"clarifying_questions":    {"id", "information_id", "question"},
		"initial_checks":          {"description", "expected_evidence", "id", "purpose"},
	}
	for field, want := range nested {
		var items []json.RawMessage
		if err := json.Unmarshal(top[field], &items); err != nil || len(items) != 1 {
			t.Fatalf("%s: want array of 1 object, got %s (%v)", field, top[field], err)
		}
		if got := objectKeys(t, items[0]); !reflect.DeepEqual(got, want) {
			t.Errorf("%s item fields = %v, want %v", field, got, want)
		}
	}

	// Types of the scalar fields.
	var scalars struct {
		RulesVersion           string  `json:"rules_version"`
		RequiresEngineerReview bool    `json:"requires_engineer_review"`
		Summary                string  `json:"summary"`
		Category               string  `json:"category"`
		Confidence             float64 `json:"confidence"`
	}
	if err := json.Unmarshal(data, &scalars); err != nil {
		t.Fatal(err)
	}
	if scalars.RulesVersion != "1.0" || !scalars.RequiresEngineerReview || scalars.Category != "dns" || scalars.Confidence != 0.8 {
		t.Fatalf("unexpected scalar values: %+v", scalars)
	}
	if !strings.Contains(string(top["confidence"]), "0.8") || strings.Contains(string(top["confidence"]), `"`) {
		t.Fatalf("confidence must be a JSON number, got %s", top["confidence"])
	}
}

func TestResultJSONEmptyCollectionsEncodeAsArrays(t *testing.T) {
	data, err := json.Marshal(emptyCollectionsResultFixture())
	if err != nil {
		t.Fatal(err)
	}
	var top map[string]json.RawMessage
	if err := json.Unmarshal(data, &top); err != nil {
		t.Fatal(err)
	}
	if len(top) != 10 {
		t.Fatalf("top-level field count = %d, want 10: %s", len(top), data)
	}
	for _, f := range []string{
		"classification_evidence", "reported_facts", "missing_information",
		"clarifying_questions", "initial_checks",
	} {
		if got := string(top[f]); got != "[]" {
			t.Errorf("%s = %s, want []", f, got)
		}
	}
	if got := string(top["confidence"]); got != "0" {
		t.Errorf("confidence = %s, want 0", got)
	}
}

func TestResultJSONRoundTrip(t *testing.T) {
	for name, want := range map[string]Result{
		"populated": populatedResultFixture(),
		"empty":     emptyCollectionsResultFixture(),
	} {
		t.Run(name, func(t *testing.T) {
			data, err := json.Marshal(want)
			if err != nil {
				t.Fatal(err)
			}
			var got Result
			if err := json.Unmarshal(data, &got); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Fatalf("round trip = %+v, want %+v", got, want)
			}
		})
	}
}

func TestContractConstants(t *testing.T) {
	if RulesVersion != "1.0" {
		t.Errorf("RulesVersion = %q", RulesVersion)
	}
	categories := map[Category]string{
		CategoryNodeRegistration:        "node_registration",
		CategoryDNS:                     "dns",
		CategoryInsufficientInformation: "insufficient_information",
		CategoryUnknown:                 "unknown",
	}
	for c, want := range categories {
		if string(c) != want {
			t.Errorf("category %q, want %q", c, want)
		}
	}
	confidences := map[Confidence]float64{
		ConfidenceUnresolved: 0,
		ConfidenceTentative:  0.5,
		ConfidenceSupported:  0.8,
	}
	for c, want := range confidences {
		if float64(c) != want {
			t.Errorf("confidence %v, want %v", c, want)
		}
	}
}

func TestInputJSONFields(t *testing.T) {
	tests := []struct {
		name string
		in   Input
		want string
	}{
		{"nil optionals omitted", Input{Description: "d"}, `{"description":"d"}`},
		{"empty metadata omitted", Input{Description: "d", Metadata: map[string]string{}}, `{"description":"d"}`},
		{"explicit empty case id kept", Input{Description: "d", CaseID: strPtr("")}, `{"description":"d","case_id":""}`},
		{"all fields", Input{Description: "d", CaseID: strPtr("c"), Metadata: map[string]string{"k": "v"}}, `{"description":"d","case_id":"c","metadata":{"k":"v"}}`},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.in)
			if err != nil {
				t.Fatal(err)
			}
			if string(data) != tc.want {
				t.Fatalf("got %s, want %s", data, tc.want)
			}
		})
	}
}
