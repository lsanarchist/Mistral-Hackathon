package model

import "time"

type Target struct {
	Type    string `json:"type"`
	BaseURL string `json:"baseUrl"`
}

type PluginInfo struct {
	Name        string      `json:"name"`
	Version     string      `json:"version"`
	SDKVersion  string      `json:"sdkVersion"`
	Capabilities Capabilities `json:"capabilities"`
}

type Capabilities struct {
	Targets   []string `json:"targets"`
	Profiles  []string `json:"profiles"`
}

type Artifact struct {
	Kind        string `json:"kind"`
	ProfileType string `json:"profileType"`
	Path        string `json:"path"`
	ContentType string `json:"contentType"`
}

type ArtifactBundle struct {
	Metadata  Metadata  `json:"metadata"`
	Target    Target    `json:"target"`
	Artifacts []Artifact `json:"artifacts"`
}

type Metadata struct {
	Timestamp   time.Time `json:"timestamp"`
	DurationSec int       `json:"durationSec"`
	Service     string    `json:"service"`
	Scenario    string    `json:"scenario"`
	GitSha      string    `json:"gitSha"`
}

type ProfileBundle struct {
	Metadata  Metadata  `json:"metadata"`
	Target    Target    `json:"target"`
	Plugin    PluginRef `json:"plugin"`
	Artifacts []Artifact `json:"artifacts"`
}

type PluginRef struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Finding struct {
	Category  string      `json:"category"`
	Title     string      `json:"title"`
	Severity  string      `json:"severity"`
	Score     int         `json:"score"`
	Top       []StackFrame `json:"top"`
	Evidence  Evidence    `json:"evidence"`
}

type StackFrame struct {
	Function string  `json:"function"`
	File     string  `json:"file"`
	Line     int     `json:"line"`
	Cum      float64 `json:"cum"`
	Flat     float64 `json:"flat"`
}

type Evidence struct {
	ArtifactPath string    `json:"artifactPath"`
	ProfileType  string    `json:"profileType"`
	ExtractedAt  time.Time `json:"extractedAt"`
}

type FindingsBundle struct {
	Summary   Summary    `json:"summary"`
	Findings  []Finding  `json:"findings"`
}

type Summary struct {
	TopIssueTags []string `json:"topIssueTags"`
	OverallScore int      `json:"overallScore"`
	Notes       []string `json:"notes"`
}

type CollectRequest struct {
	Target     Target            `json:"target"`
	DurationSec int               `json:"durationSec"`
	Profiles   []string           `json:"profiles"`
	OutDir     string            `json:"outDir"`
	Metadata   map[string]string  `json:"metadata"`
}