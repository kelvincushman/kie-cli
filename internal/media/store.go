// Copyright 2026 Kelvin Cushman and contributors. Licensed under Apache-2.0. See LICENSE.

package media

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"kie-pp-cli/internal/cliutil"
)

const (
	maxImageReferenceBytes int64 = 30 << 20
	maxAudioReferenceBytes int64 = 15 << 20
	maxVideoReferenceBytes int64 = 200 << 20
)

type Store struct {
	root string
}

func DefaultStore() (*Store, error) {
	dir, err := cliutil.DataDir()
	if err != nil {
		return nil, err
	}
	return NewStore(filepath.Join(dir, "media")), nil
}

func NewStore(root string) *Store {
	return &Store{root: filepath.Clean(root)}
}

func (s *Store) Root() string { return s.root }

func (s *Store) SaveBrief(b *Brief) error {
	if b == nil || strings.TrimSpace(b.ID) == "" {
		return fmt.Errorf("brief id is required")
	}
	return s.writeJSON(filepath.Join(s.root, "briefs", b.ID+".json"), b)
}

func (s *Store) GetBrief(id string) (*Brief, error) {
	var b Brief
	if err := s.readJSON(filepath.Join(s.root, "briefs", safeID(id)+".json"), &b); err != nil {
		return nil, err
	}
	return &b, nil
}

func (s *Store) ListBriefs() ([]Brief, error) {
	var briefs []Brief
	if err := s.listJSON(filepath.Join(s.root, "briefs"), func(data []byte) error {
		var b Brief
		if err := json.Unmarshal(data, &b); err != nil {
			return err
		}
		briefs = append(briefs, b)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(briefs, func(i, j int) bool { return briefs[i].UpdatedAt.After(briefs[j].UpdatedAt) })
	return briefs, nil
}

func (s *Store) AddReference(source, name string) (*Reference, error) {
	return s.AddReferenceTyped(source, name, "")
}

func (s *Store) AddReferenceTyped(source, name, mediaType string) (*Reference, error) {
	source = normalizeDroppedReference(source)
	if source == "" {
		return nil, fmt.Errorf("reference source is required")
	}
	mediaType = normalizeReferenceMediaType(mediaType)
	if mediaType != "" && !supportedReferenceMediaType(mediaType) {
		return nil, fmt.Errorf("reference type must be image, video, or audio")
	}
	ref := &Reference{ID: newID("ref"), Name: strings.TrimSpace(name), MediaType: mediaType, Source: source, CreatedAt: time.Now().UTC()}
	if parsed, err := url.Parse(source); err == nil && (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Host != "" {
		ref.Kind = "url"
		ref.URL = source
		if ref.MediaType == "" {
			ref.MediaType = inferReferenceMediaType(parsed.Path)
		}
		if ref.MediaType == "" {
			ref.MediaType = "image"
		}
		if ref.Name == "" {
			ref.Name = source
		}
	} else {
		detected, err := inspectLocalReference(source, mediaType)
		if err != nil {
			return nil, err
		}
		ref.MediaType = detected
		ref.Kind = "file"
		if ref.Name == "" {
			ref.Name = filepath.Base(source)
		}
		dest := filepath.Join(s.root, "reference-files", ref.ID, filepath.Base(source))
		if err := copyPrivateFile(source, dest); err != nil {
			return nil, err
		}
		ref.StoredPath = dest
	}
	if err := s.writeJSON(filepath.Join(s.root, "references", ref.ID+".json"), ref); err != nil {
		return nil, err
	}
	return ref, nil
}

// VaultBriefReferences replaces local paths with private ref:<id> handles.
// This prevents absolute paths from being persisted in briefs or returned to
// agents while retaining explicit, reusable user-selected reference images.
func (s *Store) VaultBriefReferences(b *Brief) error {
	if b == nil {
		return nil
	}
	changed := false
	var err error
	if b.References, changed, err = s.vaultReferences(b.References, "image", changed); err != nil {
		return err
	}
	if b.ReferenceVideos, changed, err = s.vaultReferences(b.ReferenceVideos, "video", changed); err != nil {
		return err
	}
	if b.ReferenceAudio, changed, err = s.vaultReferences(b.ReferenceAudio, "audio", changed); err != nil {
		return err
	}
	if b.FirstFrame, changed, err = s.vaultSingleReference(b.FirstFrame, "image", changed); err != nil {
		return err
	}
	if b.LastFrame, changed, err = s.vaultSingleReference(b.LastFrame, "image", changed); err != nil {
		return err
	}
	if changed {
		Refresh(b)
	}
	return nil
}

func (s *Store) vaultReferences(sources []string, mediaType string, changed bool) ([]string, bool, error) {
	result := make([]string, 0, len(sources))
	for _, source := range sources {
		resolved, nextChanged, err := s.vaultSingleReference(source, mediaType, changed)
		if err != nil {
			return nil, changed, err
		}
		changed = nextChanged
		if resolved != "" {
			result = append(result, resolved)
		}
	}
	return result, changed, nil
}

func (s *Store) vaultSingleReference(source, mediaType string, changed bool) (string, bool, error) {
	source = normalizeDroppedReference(source)
	if source == "" || isHTTPReference(source) {
		return source, changed, nil
	}
	if strings.HasPrefix(source, "ref:") {
		ref, err := s.GetReference(source)
		if err != nil {
			return "", changed, err
		}
		if ref.MediaType != "" && ref.MediaType != mediaType {
			return "", changed, fmt.Errorf("reference %s is %s, not %s", source, ref.MediaType, mediaType)
		}
		return source, changed, nil
	}
	ref, err := s.AddReferenceTyped(source, "", mediaType)
	if err != nil {
		return "", changed, err
	}
	return "ref:" + ref.ID, true, nil
}

func PublicReferences(refs []Reference) []PublicReference {
	public := make([]PublicReference, 0, len(refs))
	for _, ref := range refs {
		public = append(public, ref.Public())
	}
	return public
}

func isHTTPReference(source string) bool {
	parsed, err := url.Parse(strings.TrimSpace(source))
	return err == nil && (parsed.Scheme == "https" || parsed.Scheme == "http") && parsed.Host != ""
}

func inspectLocalReference(source, expectedType string) (string, error) {
	source = normalizeDroppedReference(source)
	info, err := os.Lstat(source)
	if err != nil {
		return "", fmt.Errorf("reading reference %q: %w", filepath.Base(source), err)
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return "", fmt.Errorf("reference %q must not be a symbolic link", filepath.Base(source))
	}
	if !info.Mode().IsRegular() {
		return "", fmt.Errorf("reference %q is not a regular file", filepath.Base(source))
	}
	f, err := os.Open(source) // #nosec G304 -- the user explicitly selected this local reference.
	if err != nil {
		return "", fmt.Errorf("opening reference %q: %w", filepath.Base(source), err)
	}
	defer f.Close()
	header := make([]byte, 512)
	n, readErr := f.Read(header)
	if readErr != nil && !errors.Is(readErr, io.EOF) {
		return "", fmt.Errorf("reading reference %q: %w", filepath.Base(source), readErr)
	}
	mimeType := http.DetectContentType(header[:n])
	mediaType := mediaTypeForMIME(mimeType, filepath.Ext(source))
	if mediaType == "" {
		return "", fmt.Errorf("reference %q is %s; supported local media are JPEG, PNG, GIF, WebP, BMP, TIFF, MP4, MOV, MKV, MP3, WAV, AAC, M4A, or OGG", filepath.Base(source), mimeType)
	}
	if expectedType != "" && mediaType != expectedType {
		return "", fmt.Errorf("reference %q is %s, not %s", filepath.Base(source), mediaType, expectedType)
	}
	limit := maxImageReferenceBytes
	switch mediaType {
	case "audio":
		limit = maxAudioReferenceBytes
	case "video":
		limit = maxVideoReferenceBytes
	}
	if info.Size() > limit {
		return "", fmt.Errorf("%s reference %q exceeds the %d MiB local safety limit", mediaType, filepath.Base(source), limit>>20)
	}
	return mediaType, nil
}

func mediaTypeForMIME(mimeType, extension string) string {
	extensionType := inferReferenceMediaType(extension)
	switch strings.ToLower(mimeType) {
	case "image/jpeg", "image/png", "image/gif", "image/webp", "image/bmp", "image/tiff":
		return "image"
	case "video/mp4":
		// Some M4A containers are sniffed as video/mp4. The selected file
		// extension disambiguates the media role before upload.
		if extensionType == "audio" {
			return "audio"
		}
		return "video"
	case "video/quicktime", "video/x-matroska":
		return "video"
	case "audio/mpeg", "audio/mp3", "audio/wav", "audio/x-wav", "audio/vnd.wave", "audio/aac", "audio/mp4", "audio/ogg":
		return "audio"
	case "application/ogg":
		return extensionType
	}
	return extensionType
}

func inferReferenceMediaType(path string) string {
	extension := strings.ToLower(filepath.Ext(path))
	if strings.HasPrefix(path, ".") && !strings.Contains(path[1:], ".") {
		extension = strings.ToLower(path)
	}
	switch extension {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp", ".bmp", ".tif", ".tiff":
		return "image"
	case ".mp4", ".mov", ".mkv":
		return "video"
	case ".mp3", ".wav", ".aac", ".m4a", ".ogg":
		return "audio"
	default:
		return ""
	}
}

func normalizeReferenceMediaType(value string) string {
	value = strings.ToLower(strings.TrimSpace(value))
	switch value {
	case "photo", "picture", "still":
		return "image"
	case "movie", "clip":
		return "video"
	case "sound", "voice", "music":
		return "audio"
	default:
		return value
	}
}

func supportedReferenceMediaType(value string) bool {
	return value == "image" || value == "video" || value == "audio"
}

func (s *Store) GetReference(id string) (*Reference, error) {
	id = strings.TrimPrefix(strings.TrimSpace(id), "ref:")
	var ref Reference
	if err := s.readJSON(filepath.Join(s.root, "references", safeID(id)+".json"), &ref); err != nil {
		return nil, err
	}
	return &ref, nil
}

func (s *Store) ListReferences() ([]Reference, error) {
	var refs []Reference
	if err := s.listJSON(filepath.Join(s.root, "references"), func(data []byte) error {
		var ref Reference
		if err := json.Unmarshal(data, &ref); err != nil {
			return err
		}
		refs = append(refs, ref)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(refs, func(i, j int) bool { return refs[i].CreatedAt.After(refs[j].CreatedAt) })
	return refs, nil
}

func (s *Store) CreateIdentity(name string, sources []string, consent bool) (*IdentityProfile, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, fmt.Errorf("identity name is required")
	}
	if !consent {
		return nil, fmt.Errorf("explicit consent is required before saving likeness references")
	}
	if len(sources) == 0 || len(sources) > 20 {
		return nil, fmt.Errorf("identity profiles require 1 to 20 image references")
	}
	handles := make([]string, 0, len(sources))
	for _, source := range sources {
		source = normalizeDroppedReference(source)
		if strings.HasPrefix(source, "ref:") {
			ref, err := s.GetReference(source)
			if err != nil {
				return nil, err
			}
			if ref.MediaType != "" && ref.MediaType != "image" {
				return nil, fmt.Errorf("identity reference %s is %s, not image", source, ref.MediaType)
			}
			handles = append(handles, "ref:"+ref.ID)
			continue
		}
		ref, err := s.AddReferenceTyped(source, name, "image")
		if err != nil {
			return nil, err
		}
		handles = append(handles, "ref:"+ref.ID)
	}
	now := time.Now().UTC()
	identity := &IdentityProfile{
		ID: newID("identity"), Name: name, ImageReferences: handles,
		ConsentConfirmedAt: now, CreatedAt: now, UpdatedAt: now,
	}
	if err := s.writeJSON(filepath.Join(s.root, "identities", identity.ID+".json"), identity); err != nil {
		return nil, err
	}
	return identity, nil
}

func (s *Store) GetIdentity(id string) (*IdentityProfile, error) {
	id = strings.TrimPrefix(strings.TrimSpace(id), "identity:")
	var identity IdentityProfile
	if err := s.readJSON(filepath.Join(s.root, "identities", safeID(id)+".json"), &identity); err != nil {
		return nil, err
	}
	return &identity, nil
}

func (s *Store) ListIdentities() ([]IdentityProfile, error) {
	var identities []IdentityProfile
	if err := s.listJSON(filepath.Join(s.root, "identities"), func(data []byte) error {
		var identity IdentityProfile
		if err := json.Unmarshal(data, &identity); err != nil {
			return err
		}
		identities = append(identities, identity)
		return nil
	}); err != nil {
		return nil, err
	}
	sort.Slice(identities, func(i, j int) bool { return identities[i].UpdatedAt.After(identities[j].UpdatedAt) })
	return identities, nil
}

func (s *Store) SaveGeneration(g *Generation) error {
	if g == nil || strings.TrimSpace(g.ID) == "" {
		return fmt.Errorf("generation id is required")
	}
	return s.writeJSON(filepath.Join(s.root, "generations", g.ID+".json"), g)
}

func (s *Store) GetGeneration(id string) (*Generation, error) {
	var g Generation
	if err := s.readJSON(filepath.Join(s.root, "generations", safeID(id)+".json"), &g); err != nil {
		return nil, err
	}
	return &g, nil
}

func (s *Store) findGenerationByBriefIDAndKind(briefID, kind string) (*Generation, error) {
	var found *Generation
	err := s.listJSON(filepath.Join(s.root, "generations"), func(data []byte) error {
		var generation Generation
		if err := json.Unmarshal(data, &generation); err != nil {
			return err
		}
		generationKind := generation.Kind
		if generationKind == "" {
			generationKind = GenerationKindFinal
		}
		if generation.BriefID == briefID && generationKind == kind && (found == nil || generation.CreatedAt.After(found.CreatedAt)) {
			copy := generation
			found = &copy
		}
		return nil
	})
	return found, err
}

func (s *Store) findActiveGenerationByFingerprint(briefID, kind, fingerprint string) (*Generation, error) {
	var found *Generation
	err := s.listJSON(filepath.Join(s.root, "generations"), func(data []byte) error {
		var generation Generation
		if err := json.Unmarshal(data, &generation); err != nil {
			return err
		}
		generationKind := generation.Kind
		if generationKind == "" {
			generationKind = GenerationKindFinal
		}
		if generation.BriefID == briefID && generationKind == kind && generation.Fingerprint == fingerprint &&
			!generationFailed(generation.Status) && (found == nil || generation.CreatedAt.After(found.CreatedAt)) {
			copy := generation
			found = &copy
		}
		return nil
	})
	return found, err
}

func (s *Store) acquireSubmission(briefID string) (func(), error) {
	return s.acquireBriefLock(briefID, "submit")
}

func (s *Store) acquirePreviewSubmission(briefID string) (func(), error) {
	return s.acquireBriefLock(briefID, "preview")
}

func (s *Store) acquireBriefLock(briefID, operation string) (func(), error) {
	lockDir := filepath.Join(s.root, "locks")
	if err := os.MkdirAll(lockDir, 0o700); err != nil {
		return nil, err
	}
	path := filepath.Join(lockDir, safeID(briefID)+"."+safeID(operation)+".lock")
	for attempt := 0; attempt < 2; attempt++ {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600) // #nosec G304 -- private app-derived lock path.
		if err == nil {
			_, _ = fmt.Fprintf(file, "%s\n", time.Now().UTC().Format(time.RFC3339Nano))
			if closeErr := file.Close(); closeErr != nil {
				_ = os.Remove(path)
				return nil, closeErr
			}
			return func() { _ = os.Remove(path) }, nil
		}
		if !errors.Is(err, os.ErrExist) {
			return nil, err
		}
		info, statErr := os.Stat(path)
		if statErr != nil || time.Since(info.ModTime()) <= 15*time.Minute {
			return nil, fmt.Errorf("brief %s %s is already in progress", briefID, operation)
		}
		if err := os.Remove(path); err != nil {
			return nil, fmt.Errorf("removing stale submission lock: %w", err)
		}
	}
	return nil, fmt.Errorf("brief %s %s is already in progress", briefID, operation)
}

func (s *Store) writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return cliutil.AtomicWritePrivateFile(path, data, 0o600, 0o700)
}

func (s *Store) readJSON(path string, value any) error {
	data, err := os.ReadFile(path) // #nosec G304 -- path is rooted in the app data directory and ID is sanitized.
	if err != nil {
		return err
	}
	return json.Unmarshal(data, value)
}

func (s *Store) listJSON(dir string, consume func([]byte) error) error {
	entries, err := os.ReadDir(dir)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name())) // #nosec G304 -- names come from the private app directory.
		if err != nil {
			return err
		}
		if err := consume(data); err != nil {
			return fmt.Errorf("reading %s: %w", entry.Name(), err)
		}
	}
	return nil
}

func copyPrivateFile(source, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0o700); err != nil {
		return err
	}
	in, err := os.Open(source) // #nosec G304 -- the caller explicitly selected this reference file.
	if err != nil {
		return err
	}
	defer in.Close()
	out, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_EXCL, 0o600) // #nosec G304 -- destination is private and app-derived.
	if err != nil {
		return err
	}
	if _, err := io.Copy(out, in); err != nil {
		_ = out.Close()
		_ = os.Remove(dest)
		return err
	}
	return out.Close()
}

func newID(prefix string) string {
	buf := make([]byte, 8)
	if _, err := rand.Read(buf); err != nil {
		return fmt.Sprintf("%s_%d", prefix, time.Now().UnixNano())
	}
	return prefix + "_" + hex.EncodeToString(buf)
}

func safeID(id string) string {
	id = strings.TrimSpace(id)
	id = filepath.Base(id)
	id = strings.TrimSuffix(id, filepath.Ext(id))
	return id
}
