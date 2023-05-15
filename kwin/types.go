package kwin

type dimensions struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
	W float64 `json:"width"`
	H float64 `json:"height"`
}

// display holds display properties
type display struct {
	Id         int64
	Dimensions dimensions `json:"displaySize"`
}
type displays map[int64]display

// window holds window properties
type window struct {
	Id         string     `json:"id"`
	Class      string     `json:"class"`
	Title      string     `json:"title"`
	Active     bool       `json:"active"`
	Minimised  bool       `json:"minimised"`
	Display    int64      `json:"desktop"`
	Dimensions dimensions `json:"geometry"`
}
type windows []window

type dbusResponse string
