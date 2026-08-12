// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"fmt"
	"sort"
	"strings"
)

// Workflow is compact runtime metadata for one Kie-native skill port. Agents
// can discover this over CLI or MCP instead of loading every skill document.
type Workflow struct {
	Name        string   `json:"name"`
	Skill       string   `json:"skill"`
	Description string   `json:"description"`
	MediaTypes  []string `json:"media_types"`
	Stages      []string `json:"stages"`
	Unsupported []string `json:"unsupported,omitempty"`
}

var workflowCatalog = map[string]Workflow{
	"academy": {
		Name: "academy", Skill: "kie-lesson", MediaTypes: []string{"image", "video"},
		Description: "Select a source-linked public production lesson, adapt its method into an original Kie brief, then enforce script, storyboard, keyframe, and human approval gates.",
		Stages:      []string{"select-lesson", "qualify", "script", "lock-assets", "storyboard", "prompt-each-shot", "preview-each-shot", "approve-anchor", "generate", "continuity-review", "local-assembly"},
		Unsupported: []string{"copied-course-prompts", "trained-cross-model-soul", "automatic-quality-acceptance"},
	},
	"generate": {
		Name: "generate", Skill: "kie-generate", MediaTypes: []string{"image", "video"},
		Description: "Route general image and video generation through the shared director; the skill uses dedicated CLI commands for audio and music.",
		Stages:      []string{"qualify", "select-model", "preview-video", "approve-anchor", "generate", "inspect", "deliver"},
		Unsupported: []string{"true-3d-mesh", "proprietary-virality-score"},
	},
	"brandkit": {
		Name: "brandkit", Skill: "kie-brandkit", MediaTypes: []string{"image"},
		Description: "Build an approval-led local brand system with Kie concept media.",
		Stages:      []string{"strategy", "palette", "logo-directions", "typography", "imagery", "applications", "handoff"},
		Unsupported: []string{"kie-vector-logo-generation"},
	},
	"marketplace-cards": {
		Name: "marketplace-cards", Skill: "kie-marketplace-cards", MediaTypes: []string{"image"},
		Description: "Create truthful marketplace listing visuals and local exact-copy layouts.",
		Stages:      []string{"evidence", "asset-plan", "generate", "compose-copy", "compliance-review", "deliver"},
		Unsupported: []string{"marketplace-compliance-certification"},
	},
	"product-photoshoot": {
		Name: "product-photoshoot", Skill: "kie-product-photoshoot", MediaTypes: []string{"image"},
		Description: "Create consistent reference-led product campaign imagery.",
		Stages:      []string{"product-truth", "shoot-lock", "shot-list", "generate", "fidelity-review", "deliver"},
	},
	"identity": {
		Name: "identity", Skill: "kie-identity", MediaTypes: []string{"image", "video"},
		Description: "Use consented local likeness-reference bundles without biometric training.",
		Stages:      []string{"consent", "photo-review", "save-local-bundle", "preview-video", "approve-anchor", "generate", "identity-review"},
		Unsupported: []string{"cross-model-trained-soul"},
	},
	"video-explainer": {
		Name: "video-explainer", Skill: "kie-video-explainer", MediaTypes: []string{"image", "video"},
		Description: "Create a reusable style key and direct SeedDance visual blocks while the skill creates narration separately for local assembly.",
		Stages:      []string{"settings", "research", "script", "narration", "preview-each-block", "approve-anchor", "visual-blocks", "local-assembly", "qa"},
		Unsupported: []string{"kie-explainer-assembly-endpoint"},
	},
	"websites": {
		Name: "websites", Skill: "kie-websites", MediaTypes: []string{"image", "video"},
		Description: "Build locally hosted sites, apps, or games with Kie-generated media assets.",
		Stages:      []string{"product-scope", "local-build", "asset-plan", "preview-video", "approve-anchor", "generate", "integrate", "verify", "authorized-deploy"},
		Unsupported: []string{"kie-website-hosting", "true-3d-mesh"},
	},
	"youtube-thumbnail": {
		Name: "youtube-thumbnail", Skill: "kie-youtube-thumbnail", MediaTypes: []string{"image"},
		Description: "Create truthful 16:9 thumbnail variants with local exact-text composition.",
		Stages:      []string{"truthful-promise", "concepts", "generate-variants", "compose-text", "small-size-review", "deliver"},
	},
}

func ListWorkflows() []Workflow {
	result := make([]Workflow, 0, len(workflowCatalog))
	for _, workflow := range workflowCatalog {
		result = append(result, workflow)
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}

func GetWorkflow(name string) (Workflow, error) {
	name = normalizeWorkflow(name)
	workflow, ok := workflowCatalog[name]
	if !ok {
		return Workflow{}, fmt.Errorf("unknown media workflow %q", name)
	}
	return workflow, nil
}

func normalizeWorkflow(name string) string {
	name = strings.ToLower(strings.TrimSpace(name))
	name = strings.TrimPrefix(name, "kie-")
	switch name {
	case "lesson", "academy-film", "academy-director":
		return "academy"
	case "brand-kit":
		return "brandkit"
	case "marketplace", "cards":
		return "marketplace-cards"
	case "photoshoot", "product":
		return "product-photoshoot"
	case "explainer":
		return "video-explainer"
	case "website", "site", "game":
		return "websites"
	case "thumbnail", "youtube":
		return "youtube-thumbnail"
	default:
		return name
	}
}
