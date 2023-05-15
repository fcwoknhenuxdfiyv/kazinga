package config

// Tweak holds adjustments for specific window classes
type Tweak struct {
	Class   string
	NudgeX  float64
	NudgeY  float64
	ShrinkW float64
	ShrinkH float64
}
type Tweaks map[string][]Tweak

// Config is the main configuration
type Config struct {
	TopPanelHeight    float64 `json:",omitempty"`
	BottomPanelHeight float64 `json:",omitempty"`
	EdgeBorderX       float64
	EdgeBorderY       float64
	BorderX           float64
	BorderY           float64
	Tweaks            Tweaks
}
