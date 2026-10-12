// Package analysis holds the transport-independent domain contract for
// KubeWaypoint's deterministic incident analysis.
package analysis

// RulesVersion identifies the version of the classification rules.
const RulesVersion = "1.0"

// Category is the classification outcome of an analysis.
type Category string

const (
	CategoryNodeRegistration        Category = "node_registration"
	CategoryDNS                     Category = "dns"
	CategoryInsufficientInformation Category = "insufficient_information"
	CategoryUnknown                 Category = "unknown"
)

// Confidence is a fixed rubric value, not a probability or calibrated score.
// Active evidence for both supported categories yields ConfidenceUnresolved
// regardless of strength, and repeated evidence never raises confidence.
type Confidence float64

const (
	// ConfidenceUnresolved means no category is supported.
	ConfidenceUnresolved Confidence = 0
	// ConfidenceTentative means the category is tentatively supported.
	ConfidenceTentative Confidence = 0.5
	// ConfidenceSupported means the category is explicitly supported.
	ConfidenceSupported Confidence = 0.8
)

// Input is the free-text incident description to analyze.
type Input struct {
	Description string            `json:"description"`
	CaseID      *string           `json:"case_id,omitempty"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// Result is the domain analysis outcome. Collections are never omitted; a
// successfully generated result must carry empty non-nil slices rather than
// nil ones so they encode as [].
type Result struct {
	RulesVersion           string                   `json:"rules_version"`
	RequiresEngineerReview bool                     `json:"requires_engineer_review"`
	Summary                string                   `json:"summary"`
	Category               Category                 `json:"category"`
	Confidence             Confidence               `json:"confidence"`
	ClassificationEvidence []ClassificationEvidence `json:"classification_evidence"`
	ReportedFacts          []ReportedFact           `json:"reported_facts"`
	MissingInformation     []InformationGap         `json:"missing_information"`
	ClarifyingQuestions    []ClarifyingQuestion     `json:"clarifying_questions"`
	InitialChecks          []InitialCheck           `json:"initial_checks"`
}

// ClassificationEvidence links a rule to the input text that activated it.
type ClassificationEvidence struct {
	RuleID        string `json:"rule_id"`
	SourceExcerpt string `json:"source_excerpt"`
}

// ReportedFact is a statement taken from the input, not a verified finding.
type ReportedFact struct {
	Statement     string `json:"statement"`
	SourceExcerpt string `json:"source_excerpt"`
}

// InformationGap names information missing from the input.
type InformationGap struct {
	ID          string `json:"id"`
	Description string `json:"description"`
}

// ClarifyingQuestion asks for the information described by the gap whose ID
// equals InformationID.
type ClarifyingQuestion struct {
	ID            string `json:"id"`
	InformationID string `json:"information_id"`
	Question      string `json:"question"`
}

// InitialCheck is a read-only starting check for the engineer.
type InitialCheck struct {
	ID               string `json:"id"`
	Description      string `json:"description"`
	Purpose          string `json:"purpose"`
	ExpectedEvidence string `json:"expected_evidence"`
}
