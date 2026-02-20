package research

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sean/apollo/api/internal/schema"
)

// AssembleFromDir reads the file-per-lesson directory tree under workDir and
// assembles it into a CurriculumOutput. The tree must contain topic.json at
// the root and a modules/ directory with NN-<slug>/ subdirectories, each
// holding a module.json and one or more NN-<slug>.json lesson files.
//
// The assembled output is validated against the curriculum schema before
// being returned. If any file is missing, malformed, or the assembled
// output fails validation, a descriptive error is returned.
func AssembleFromDir(workDir string) (*CurriculumOutput, error) {
	// Read topic.json.
	topicPath := filepath.Join(workDir, TopicFileName)

	topicData, err := os.ReadFile(topicPath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", TopicFileName, err)
	}

	var topic TopicFile
	if err := json.Unmarshal(topicData, &topic); err != nil {
		return nil, fmt.Errorf("parse %s: %w", topicPath, err)
	}

	// Discover and sort module directories.
	modulesDir := filepath.Join(workDir, ModulesDirName)

	moduleDirs, err := readSortedDirs(modulesDir)
	if err != nil {
		return nil, fmt.Errorf("read %s directory: %w", ModulesDirName, err)
	}

	if len(moduleDirs) == 0 {
		return nil, fmt.Errorf("%s directory is empty: no module directories found", ModulesDirName)
	}

	// Assemble modules.
	modules := make([]ModuleOutput, 0, len(moduleDirs))

	for _, modDirName := range moduleDirs {
		modDirPath := filepath.Join(modulesDir, modDirName)

		mod, err := assembleModule(modDirPath)
		if err != nil {
			return nil, fmt.Errorf("module %s: %w", modDirName, err)
		}

		modules = append(modules, *mod)
	}

	// Build CurriculumOutput from topic + modules.
	curriculum := &CurriculumOutput{
		ID:             topic.ID,
		Title:          topic.Title,
		Description:    topic.Description,
		Difficulty:     topic.Difficulty,
		EstimatedHours: topic.EstimatedHours,
		Tags:           topic.Tags,
		Prerequisites:  topic.Prerequisites,
		RelatedTopics:  topic.RelatedTopics,
		Modules:        modules,
		SourceURLs:     topic.SourceURLs,
		GeneratedAt:    topic.GeneratedAt,
		Version:        topic.Version,
	}

	// Validate assembled output against the curriculum schema.
	assembled, err := json.Marshal(curriculum)
	if err != nil {
		return nil, fmt.Errorf("marshal assembled curriculum: %w", err)
	}

	if err := schema.Validate(assembled); err != nil {
		return nil, fmt.Errorf("assembled curriculum schema validation: %w", err)
	}

	return curriculum, nil
}

// assembleModule reads module.json and all lesson files from a module directory.
func assembleModule(modDirPath string) (*ModuleOutput, error) {
	// Read module.json.
	modFilePath := filepath.Join(modDirPath, ModuleFileBaseName)

	modData, err := os.ReadFile(modFilePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", ModuleFileBaseName, err)
	}

	var modFile ModuleFile
	if err := json.Unmarshal(modData, &modFile); err != nil {
		return nil, fmt.Errorf("parse %s: %w", modFilePath, err)
	}

	// Discover and sort lesson files (*.json excluding module.json).
	lessonFiles, err := readSortedLessonFiles(modDirPath)
	if err != nil {
		return nil, fmt.Errorf("read lesson files: %w", err)
	}

	lessons := make([]LessonOutput, 0, len(lessonFiles))

	for _, lessonFileName := range lessonFiles {
		lessonPath := filepath.Join(modDirPath, lessonFileName)

		lessonData, err := os.ReadFile(lessonPath)
		if err != nil {
			return nil, fmt.Errorf("read lesson %s: %w", lessonFileName, err)
		}

		var lesson LessonOutput
		if err := json.Unmarshal(lessonData, &lesson); err != nil {
			return nil, fmt.Errorf("parse %s: %w", lessonPath, err)
		}

		lessons = append(lessons, lesson)
	}

	return &ModuleOutput{
		ID:                 modFile.ID,
		Title:              modFile.Title,
		Description:        modFile.Description,
		LearningObjectives: modFile.LearningObjectives,
		EstimatedMinutes:   modFile.EstimatedMinutes,
		Order:              modFile.Order,
		Lessons:            lessons,
		Assessment:         modFile.Assessment,
	}, nil
}

// readSortedDirs returns directory names under parentDir sorted by numeric prefix.
func readSortedDirs(parentDir string) ([]string, error) {
	entries, err := os.ReadDir(parentDir)
	if err != nil {
		return nil, err
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() {
			dirs = append(dirs, e.Name())
		}
	}

	sort.Slice(dirs, func(i, j int) bool {
		return numericPrefix(dirs[i]) < numericPrefix(dirs[j])
	})

	return dirs, nil
}

// readSortedLessonFiles returns *.json file names in dir, excluding module.json,
// sorted by numeric prefix.
func readSortedLessonFiles(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := e.Name()
		if name == ModuleFileBaseName {
			continue
		}

		if strings.HasSuffix(name, ".json") {
			files = append(files, name)
		}
	}

	sort.Slice(files, func(i, j int) bool {
		return numericPrefix(files[i]) < numericPrefix(files[j])
	})

	return files, nil
}

// numericPrefix extracts the leading integer from a name like "01-introduction"
// or "02-basics.json". Returns 0 if no numeric prefix is found.
func numericPrefix(name string) int {
	parts := strings.SplitN(name, "-", 2)
	if len(parts) == 0 {
		return 0
	}

	n, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	return n
}
