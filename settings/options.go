package settings


type OptionKeySetting struct {
	LongKey string
	ShortKey string
	Required bool
	Flag bool
	Multiple bool
	Description string
	Example string
}

type OptionNoKeySetting struct {
	Name string
	Required bool
	Description string
	Example string
}

type OptionSettings struct {
	KeySettings []OptionKeySetting
	NoKeySettings []OptionNoKeySetting
}
