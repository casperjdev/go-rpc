package main

type Config struct {
	Credentials struct {
		ApplicationID string `mapstructure:"application_id" json:"application_id"`
	} `mapstructure:"credentials" json:"credentials"`

	Activity DiscordActivity `mapstructure:"activity" json:"activity"`
}

type DiscordActivity struct {
	Type       byte       `mapstructure:"type" json:"type,omitempty"`
	State      string     `mapstructure:"state" json:"state,omitempty"`
	Details    string     `mapstructure:"details" json:"details,omitempty"`
	Timestamps *Timestamps `mapstructure:"timestamps" json:"timestamps,omitempty"`
	Assets     *Assets    `mapstructure:"assets" json:"assets,omitempty"`
	Party      *Party     `mapstructure:"party" json:"party,omitempty"`
	Secrets    *Secrets   `mapstructure:"secrets" json:"secrets,omitempty"`
	Buttons    []Button   `mapstructure:"buttons" json:"buttons,omitempty"`
	Instance   bool       `mapstructure:"instance" json:"instance,omitempty"`
}

type Timestamps struct {
	Start int64 `mapstructure:"start" json:"start,omitempty"`
	End   int64 `mapstructure:"end" json:"end,omitempty"`
}

type Assets struct {
	LargeImage string `mapstructure:"large_image" json:"large_image,omitempty"`
	LargeText  string `mapstructure:"large_text" json:"large_text,omitempty"`
	SmallImage string `mapstructure:"small_image" json:"small_image,omitempty"`
	SmallText  string `mapstructure:"small_text" json:"small_text,omitempty"`
}

type Party struct {
	ID   string  `mapstructure:"id" json:"id,omitempty"`
	Size []int32 `mapstructure:"size" json:"size,omitempty"`
}

type Secrets struct {
	Match    string `mapstructure:"match" json:"match,omitempty"`
	Join     string `mapstructure:"join" json:"join,omitempty"`
	Spectate string `mapstructure:"spectate" json:"spectate,omitempty"`
}

type Button struct {
	Label string `mapstructure:"label" json:"label"`
	URL   string `mapstructure:"url" json:"url"`
}
