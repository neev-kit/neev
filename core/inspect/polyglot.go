package inspect

import (
	"os"
	"path/filepath"
	"strings"
)

// LanguageDetector defines the interface for language-specific code analysis
type LanguageDetector interface {
	// Detect returns true if this detector can handle the given file
	Detect(filePath string) bool
	
	// ExtractEndpoints finds HTTP endpoints/handlers in the file
	ExtractEndpoints(filePath string, content []byte) ([]Endpoint, error)
	
	// ExtractFunctions finds function/method signatures in the file
	ExtractFunctions(filePath string, content []byte) ([]FunctionSignature, error)
	
	// Language returns the language name
	Language() Language
}

// FunctionSignature represents a parsed function signature
type FunctionSignature struct {
	Name       string
	Parameters []ParameterSpec
	Returns    []ReturnSpec
	File       string
	Line       int
	Visibility string // public, private, protected
}

// PolyglotAnalyzer manages multi-language code analysis
type PolyglotAnalyzer struct {
	detectors []LanguageDetector
}

// NewPolyglotAnalyzer creates a new analyzer with registered detectors
func NewPolyglotAnalyzer() *PolyglotAnalyzer {
	return &PolyglotAnalyzer{
		detectors: []LanguageDetector{
			// Detectors will be registered here
		},
	}
}

// RegisterDetector adds a language detector to the analyzer
func (pa *PolyglotAnalyzer) RegisterDetector(detector LanguageDetector) {
	pa.detectors = append(pa.detectors, detector)
}

// DetectLanguages scans a directory and returns language statistics
func (pa *PolyglotAnalyzer) DetectLanguages(rootDir string, ignoreDirs map[string]bool) (map[string]int, error) {
	languages := make(map[string]int)
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			// Skip ignored directories
			if ignoreDirs[info.Name()] || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Check if any detector can handle this file
		for _, detector := range pa.detectors {
			if detector.Detect(path) {
				langName := string(detector.Language())
				languages[langName]++
				break
			}
		}
		
		return nil
	})
	
	return languages, err
}

// ExtractAllEndpoints finds all HTTP endpoints in a directory
func (pa *PolyglotAnalyzer) ExtractAllEndpoints(rootDir string, ignoreDirs map[string]bool) ([]Endpoint, error) {
	var allEndpoints []Endpoint
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			if ignoreDirs[info.Name()] || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Check if any detector can handle this file
		for _, detector := range pa.detectors {
			if detector.Detect(path) {
				content, err := os.ReadFile(path)
				if err != nil {
					continue // Skip files we can't read
				}
				
				endpoints, err := detector.ExtractEndpoints(path, content)
				if err != nil {
					continue // Skip files with parsing errors
				}
				
				// Add language info to endpoints
				for i := range endpoints {
					endpoints[i].Language = string(detector.Language())
					endpoints[i].File = path
				}
				
				allEndpoints = append(allEndpoints, endpoints...)
				break
			}
		}
		
		return nil
	})
	
	return allEndpoints, err
}

// ExtractAllFunctions finds all function signatures in a directory
func (pa *PolyglotAnalyzer) ExtractAllFunctions(rootDir string, ignoreDirs map[string]bool) ([]FunctionSignature, error) {
	var allFunctions []FunctionSignature
	
	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories
		if info.IsDir() {
			if ignoreDirs[info.Name()] || strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		
		// Check if any detector can handle this file
		for _, detector := range pa.detectors {
			if detector.Detect(path) {
				content, err := os.ReadFile(path)
				if err != nil {
					continue // Skip files we can't read
				}
				
				functions, err := detector.ExtractFunctions(path, content)
				if err != nil {
					continue // Skip files with parsing errors
				}
				
				allFunctions = append(allFunctions, functions...)
				break
			}
		}
		
		return nil
	})
	
	return allFunctions, err
}

// DetectLanguageByExtension is a helper to detect language from file extension
func DetectLanguageByExtension(filePath string) Language {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".go":
		return LangGo
	case ".py":
		return LangPython
	case ".js":
		return LangJavaScript
	case ".ts":
		return LangTypeScript
	case ".java":
		return LangJava
	case ".cs":
		return LangCSharp
	case ".rb":
		return LangRuby
	default:
		return ""
	}
}
