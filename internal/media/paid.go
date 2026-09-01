// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	PaidScopePreview = "preview"
	PaidScopeProof   = "proof"
	PaidScopeFinal   = "final"

	defaultPaidDisclosure = "This live Kie.ai action may consume credits. Exact cost is not known locally."
	paidConfirmationTTL   = 10 * time.Minute
)

type PaidConfirmation struct {
	ID             string     `json:"id"`
	BriefID        string     `json:"brief_id"`
	Scope          string     `json:"scope"`
	GenerationKind string     `json:"generation_kind"`
	Model          string     `json:"model"`
	PlanHash       string     `json:"plan_hash"`
	BriefHash      string     `json:"brief_hash"`
	Disclosure     string     `json:"disclosure"`
	PriceEstimate  string     `json:"price_estimate,omitempty"`
	ExactCostKnown bool       `json:"exact_cost_known"`
	ConfirmedBy    string     `json:"confirmed_by"`
	ConfirmedAt    time.Time  `json:"confirmed_at"`
	ExpiresAt      time.Time  `json:"expires_at"`
	UsedAt         *time.Time `json:"used_at,omitempty"`
}

type PaidConfirmationRequest struct {
	Scope          string
	GenerationKind string
	ConfirmedBy    string
	Disclosure     string
	PriceEstimate  string
	ExactCostKnown bool
	Acknowledged   bool
	Now            time.Time
}

func NewPaidConfirmation(brief *Brief, plan *Plan, request PaidConfirmationRequest) (*PaidConfirmation, error) {
	if brief == nil || strings.TrimSpace(brief.ID) == "" || plan == nil || strings.TrimSpace(plan.Model) == "" {
		return nil, fmt.Errorf("a saved brief and paid plan are required")
	}
	if !request.Acknowledged {
		return nil, fmt.Errorf("paid generation disclosure must be explicitly acknowledged")
	}
	expectedKind := generationKindForScope(request.Scope)
	if expectedKind == "" {
		return nil, fmt.Errorf("unsupported paid confirmation scope %q", request.Scope)
	}
	if strings.TrimSpace(request.GenerationKind) == "" {
		return nil, fmt.Errorf("paid generation kind is required")
	}
	if request.GenerationKind != expectedKind {
		return nil, fmt.Errorf("paid confirmation scope %q requires generation kind %q, got %q", request.Scope, expectedKind, request.GenerationKind)
	}
	if strings.TrimSpace(request.ConfirmedBy) == "" {
		return nil, fmt.Errorf("paid confirmation source is required")
	}
	now := request.Now.UTC()
	if now.IsZero() {
		now = time.Now().UTC()
	}
	disclosure := strings.TrimSpace(request.Disclosure)
	if disclosure == "" {
		disclosure = defaultPaidDisclosure
	}
	planHash, err := planFingerprint(plan)
	if err != nil {
		return nil, fmt.Errorf("fingerprinting paid plan: %w", err)
	}
	return &PaidConfirmation{
		ID: newID("confirm"), BriefID: brief.ID, Scope: request.Scope,
		GenerationKind: request.GenerationKind, Model: plan.Model,
		PlanHash: planHash, BriefHash: creativeBriefFingerprint(brief),
		Disclosure: disclosure, PriceEstimate: strings.TrimSpace(request.PriceEstimate), ExactCostKnown: request.ExactCostKnown,
		ConfirmedBy: request.ConfirmedBy, ConfirmedAt: now, ExpiresAt: now.Add(paidConfirmationTTL),
	}, nil
}

func ValidatePaidConfirmation(confirmation *PaidConfirmation, brief *Brief, plan *Plan, scope, generationKind string, now time.Time) error {
	if confirmation == nil {
		return fmt.Errorf("a fresh paid confirmation is required")
	}
	if confirmation.UsedAt != nil {
		return fmt.Errorf("paid confirmation %s was already used", confirmation.ID)
	}
	if now.IsZero() {
		now = time.Now().UTC()
	}
	if !now.UTC().Before(confirmation.ExpiresAt) {
		return fmt.Errorf("paid confirmation %s expired; confirm this paid action again", confirmation.ID)
	}
	if brief == nil || plan == nil || confirmation.BriefID != brief.ID {
		return fmt.Errorf("paid confirmation %s does not match this brief", confirmation.ID)
	}
	if confirmation.Scope != scope || confirmation.GenerationKind != generationKind {
		return fmt.Errorf("paid confirmation %s does not match the requested action", confirmation.ID)
	}
	planHash, err := planFingerprint(plan)
	if err != nil {
		return fmt.Errorf("fingerprinting paid plan: %w", err)
	}
	if confirmation.Model != plan.Model || confirmation.PlanHash != planHash || confirmation.BriefHash != creativeBriefFingerprint(brief) {
		return fmt.Errorf("paid confirmation %s is stale because the brief, model, or settings changed", confirmation.ID)
	}
	return nil
}

func planFingerprint(plan *Plan) (string, error) {
	if plan == nil {
		return "", fmt.Errorf("paid plan is required")
	}
	data, err := json.Marshal(plan)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}

func PlanFingerprint(plan *Plan) (string, error) { return planFingerprint(plan) }

func BriefFingerprint(brief *Brief) string { return creativeBriefFingerprint(brief) }

func paidActionReview(brief *Brief, nextAction string) *PaidActionReview {
	var (
		plan       *Plan
		scope      string
		kind       string
		disclosure = defaultPaidDisclosure
	)
	switch nextAction {
	case "generate_preview", "regenerate_preview":
		plan, scope, kind = BuildPreviewPlan(brief), PaidScopePreview, GenerationKindPreview
	case "offer_proof", "regenerate_proof":
		var option ProofOption
		var err error
		plan, option, err = BuildProofPlan(brief)
		if err != nil {
			return nil
		}
		disclosure = option.Disclosure
		scope, kind = PaidScopeProof, GenerationKindProof
	case "review_then_submit":
		if brief == nil || brief.Plan == nil {
			return nil
		}
		plan, scope, kind = BuildPlan(brief), PaidScopeFinal, GenerationKindFinal
	default:
		return nil
	}
	review := &PaidActionReview{
		Scope: scope, GenerationKind: kind, Model: plan.Model,
		CostStatus: plan.CostStatus, Disclosure: disclosure,
	}
	planHash, err := PlanFingerprint(plan)
	if err != nil {
		review.BlockedReason = "the current paid plan could not be fingerprinted; revise the brief before requesting confirmation"
		return review
	}
	review.PlanHash = planHash
	return review
}

func PaidPlanForScope(brief *Brief, scope string) (*Plan, string, error) {
	switch scope {
	case PaidScopePreview:
		if brief == nil || brief.MediaType != "video" {
			return nil, "", fmt.Errorf("preview confirmation applies only to video briefs")
		}
		return BuildPreviewPlan(brief), GenerationKindPreview, nil
	case PaidScopeProof:
		plan, _, err := BuildProofPlan(brief)
		return plan, GenerationKindProof, err
	case PaidScopeFinal:
		if brief == nil || brief.Plan == nil {
			return nil, "", fmt.Errorf("ready brief plan is required")
		}
		return BuildPlan(brief), GenerationKindFinal, nil
	default:
		return nil, "", fmt.Errorf("unsupported paid confirmation scope %q", scope)
	}
}

func generationKindForScope(scope string) string {
	switch scope {
	case PaidScopePreview:
		return GenerationKindPreview
	case PaidScopeProof:
		return GenerationKindProof
	case PaidScopeFinal:
		return GenerationKindFinal
	default:
		return ""
	}
}

func (s *Store) SavePaidConfirmation(confirmation *PaidConfirmation) error {
	if confirmation == nil || strings.TrimSpace(confirmation.ID) == "" {
		return fmt.Errorf("paid confirmation id is required")
	}
	return s.writeJSON(s.confirmationPath(confirmation.ID), confirmation)
}

func (s *Store) GetPaidConfirmation(id string) (*PaidConfirmation, error) {
	var confirmation PaidConfirmation
	if err := s.readJSON(s.confirmationPath(id), &confirmation); err != nil {
		return nil, err
	}
	return &confirmation, nil
}

func (s *Store) UsePaidConfirmation(id string, brief *Brief, plan *Plan, scope, generationKind string, now time.Time) (*PaidConfirmation, error) {
	if strings.TrimSpace(id) == "" {
		return nil, fmt.Errorf("a fresh paid confirmation is required")
	}
	confirmation, err := s.GetPaidConfirmation(id)
	if err != nil {
		return nil, err
	}
	if err := ValidatePaidConfirmation(confirmation, brief, plan, scope, generationKind, now); err != nil {
		return nil, err
	}
	if err := s.claimPaidConfirmation(confirmation.ID); err != nil {
		return nil, err
	}
	usedAt := now.UTC()
	if usedAt.IsZero() {
		usedAt = time.Now().UTC()
	}
	confirmation.UsedAt = &usedAt
	if err := s.SavePaidConfirmation(confirmation); err != nil {
		return nil, err
	}
	return confirmation, nil
}

func (s *Store) confirmationPath(id string) string {
	return filepath.Join(s.root, "paid-confirmations", safeID(id)+".json")
}

func (s *Store) claimPaidConfirmation(id string) error {
	dir := filepath.Join(s.root, "paid-confirmation-claims")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}
	path := filepath.Join(dir, safeID(id)+".claim")
	file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- private app-derived claim path.
	if errors.Is(err, os.ErrExist) {
		return fmt.Errorf("paid confirmation %s was already used", id)
	}
	if err != nil {
		return err
	}
	if _, err := fmt.Fprintf(file, "%s\n", time.Now().UTC().Format(time.RFC3339Nano)); err != nil {
		_ = file.Close()
		return err
	}
	return file.Close()
}
