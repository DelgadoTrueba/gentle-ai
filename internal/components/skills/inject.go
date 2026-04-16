package skills

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/gentleman-programming/gentle-ai/internal/agents"
	"github.com/gentleman-programming/gentle-ai/internal/assets"
	"github.com/gentleman-programming/gentle-ai/internal/components/filemerge"
	"github.com/gentleman-programming/gentle-ai/internal/model"
)

// isSDDSkill reports whether a skill ID belongs to the SDD orchestrator suite.
// SDD skills are installed by the SDD component; the skills component skips
// them to prevent duplicate writes when both components are selected.
func isSDDSkill(id model.SkillID) bool {
	return strings.HasPrefix(string(id), "sdd-")
}

type InjectionResult struct {
	Changed bool
	Files   []string
	Skipped []model.SkillID
}

// Inject writes the embedded SKILL.md files for each requested skill
// to the correct directory for the given agent adapter.
//
// The skills directory is determined by adapter.SkillsDir(), removing
// the need for any agent-specific switch statements.
//
// SDD skills (those whose IDs begin with "sdd-") are intentionally skipped
// here because the SDD component installs them as part of its own injection.
// This prevents a write conflict when both components are selected together.
//
// Individual skill failures (e.g., missing embedded asset) are logged
// and skipped rather than aborting the entire operation.
// TODO add markdown archivos
func Inject(homeDir string, adapter agents.Adapter, skillIDs []model.SkillID) (InjectionResult, error) {
	if !adapter.SupportsSkills() {
		return InjectionResult{Skipped: skillIDs}, nil
	}

	skillDir := adapter.SkillsDir(homeDir)
	if skillDir == "" {
		return InjectionResult{Skipped: skillIDs}, nil
	}

	paths := make([]string, 0, len(skillIDs))
	skipped := make([]model.SkillID, 0)
	changed := false

	for _, id := range skillIDs {
		// SDD skills are written by the SDD component — skip to avoid conflicts.
		if isSDDSkill(id) {
			continue
		}

		assetPath := "skills/" + string(id) + "/SKILL.md"
		content, readErr := assets.Read(assetPath)
		if readErr != nil {
			log.Printf("skills: skipping %q — embedded asset not found: %v", id, readErr)
			skipped = append(skipped, id)
			continue
		}
		if len(content) == 0 {
			return InjectionResult{}, fmt.Errorf("skill %q: embedded asset exists but is empty — build may be corrupt", id)
		}

		path := filepath.Join(skillDir, string(id), "SKILL.md")
		writeResult, writeErr := filemerge.WriteFileAtomic(path, []byte(content), 0o644)
		if writeErr != nil {
			return InjectionResult{}, fmt.Errorf("skill %q: write failed: %w", id, writeErr)
		}

		changed = changed || writeResult.Changed
		paths = append(paths, path)

		referencesChanged, referencePaths, copyErr := copyEmbeddedDir(
			assets.FS,
			filepath.ToSlash(filepath.Join("skills", string(id), "references")),
			filepath.Join(skillDir, string(id), "references"),
		)
		if copyErr != nil {
			return InjectionResult{}, fmt.Errorf("skill %q references: copy failed: %w", id, copyErr)
		}

		changed = changed || referencesChanged
		paths = append(paths, referencePaths...)
	}

	return InjectionResult{Changed: changed, Files: paths, Skipped: skipped}, nil
}

func copyEmbeddedDir(embeddedFS fs.FS, srcDir, dstDir string) (bool, []string, error) {
	entries, err := fs.ReadDir(embeddedFS, srcDir)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return false, nil, nil
		}
		return false, nil, err
	}
	if len(entries) == 0 {
		return false, nil, nil
	}

	changed := false
	paths := make([]string, 0)

	err = fs.WalkDir(embeddedFS, srcDir, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d.IsDir() {
			return nil
		}

		relPath := strings.TrimPrefix(path, srcDir+"/")
		if relPath == path {
			relPath = pathpkg.Base(path)
		}

		content, err := fs.ReadFile(embeddedFS, path)
		if err != nil {
			return err
		}

		outPath := filepath.Join(dstDir, relPath)
		writeResult, err := filemerge.WriteFileAtomic(outPath, content, 0o644)
		if err != nil {
			return err
		}

		changed = changed || writeResult.Changed
		paths = append(paths, outPath)
		return nil
	})
	if err != nil {
		return false, nil, err
	}

	return changed, paths, nil
}

// SkillPathForAgent returns the filesystem path where a skill file would be written.
func SkillPathForAgent(homeDir string, adapter agents.Adapter, id model.SkillID) string {
	skillDir := adapter.SkillsDir(homeDir)
	if skillDir == "" {
		return ""
	}
	return filepath.Join(skillDir, string(id), "SKILL.md")
}
