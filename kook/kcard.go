package kookNode

import "encoding/json"

type KCard struct {
	Card kCard
}

type kTheme string

const (
	Primary   kTheme = "primary"
	Success   kTheme = "success"
	Danger    kTheme = "danger"
	Warning   kTheme = "warning"
	Info      kTheme = "info"
	Secondary kTheme = "secondary"
	None      kTheme = "none"
)

type kType0 string // card
type kType1 string // modules
type kType2 string // fields
type kSize string

const (
	Large  kSize = "lg"
	Medium kSize = "md"
	Small  kSize = "sm"
	XSmall kSize = "xs"
)
const (
	Card kType0 = "card"
)
const (
	Header    kType1 = "header"
	Section   kType1 = "section"
	Context   kType1 = "context"
	Divider   kType1 = "divider"
	Countdown kType1 = "countdown"
	Container kType1 = "container"
	File      kType1 = "file"
)
const (
	Plaintext kType2 = "plain-text"
	Image     kType2 = "image"
	Kmarkdown kType2 = "kmarkdown"
)

type KField struct {
	Type    kType2 `json:"type"`
	Content string `json:"content,omitempty"`
	Src     string `json:"src,omitempty"`
}

type KModule struct {
	Type kType1 `json:"type,omitempty"`

	// header, section
	Text KField `json:"text,omitempty"`

	// context, container
	Elements []KField `json:"elements,omitempty"`

	// countdown
	Mode      string `json:"mode,omitempty"`
	StartTime int64  `json:"startTime,omitempty"`
	EndTime   int64  `json:"endTime,omitempty"`

	// file
	Title string `json:"title,omitempty"`
	Src   int    `json:"src,omitempty"`
	Size  int    `json:"size,omitempty"`
}

type kCard struct {
	Type    kType0    `json:"type"`
	Theme   kTheme    `json:"theme"`
	Size    kSize     `json:"size"`
	Modules []KModule `json:"modules"`
}

func (card *KCard) Init() *KCard {
	card.Card.Type = Card
	card.Card.Size = Large
	card.Card.Theme = Primary
	return card
}
func (card *KCard) AddModule(module KModule) {
	card.Card.Modules = append(card.Card.Modules, module)
}

func (card *KCard) AddModule_image(url string) {
	card.Card.Modules = append(card.Card.Modules, KModule{
		Type: "container",
		Elements: []KField{
			{
				Type: "image",
				Src:  url,
			},
		},
	})
}
func (card *KCard) AddModule_markdown(content string) {
	card.Card.Modules = append(card.Card.Modules, KModule{
		Type: "section",
		Text: KField{
			Type:    "kmarkdown",
			Content: content,
		},
	})
}
func (card *KCard) AddModule_header(content string) {
	card.Card.Modules = append(card.Card.Modules, KModule{
		Type: "header",
		Text: KField{
			Type:    "plain-text",
			Content: content,
		},
	})
}
func (card *KCard) AddModule_divider() {
	card.Card.Modules = append(card.Card.Modules, KModule{
		Type: "divider",
	})
}
func (card *KCard) String() string {
	jsons, _ := json.Marshal([]kCard{card.Card})
	return string(jsons)
}
