package structs

import "time"

// Content represents the content of a file.
type Content struct {
	ID     string `yaml:"-"`                   // used by Search
	Source string `yaml:"-"`                   // path to the file
	HTML   string `yaml:"-" json:",omitempty"` // for Markdown files

	// for everything
	Name        string    `yaml:",omitempty" json:",omitempty"` // name of the file, used in the breadcrumbs
	Title       string    `yaml:",omitempty" json:",omitempty"` // override for the name, used as page title, fallback to Name
	Subtitle    string    `yaml:",omitempty" json:",omitempty"`
	Year        int       `yaml:",omitempty" json:",omitempty"`
	Authors     oneOrMany `yaml:",omitempty" json:",omitempty"`
	Developers  oneOrMany `yaml:",omitempty" json:",omitempty"`
	Description string    `yaml:",omitempty" json:",omitempty"`
	CoverArtist string    `yaml:"cover_artist,omitempty" json:",omitempty"`
	Designer    string    `yaml:",omitempty" json:",omitempty"`

	BasedOn  oneOrMany `yaml:"based_on,omitempty" json:",omitempty"`
	Series   Reference `yaml:",omitempty" json:",omitempty"`
	Previous Reference `yaml:",omitempty" json:",omitempty"` // reference to previous in the series

	// for people
	DOB     string `yaml:",omitempty" json:",omitempty"` // date of birth
	DOD     string `yaml:",omitempty" json:",omitempty"` // date of death
	Contact string `yaml:"contact,omitempty" json:",omitempty"`

	Parent   string    `yaml:",omitempty" json:",omitempty"` // for companies
	Founded  string    `yaml:",omitempty" json:",omitempty"` // for companies
	Founders oneOrMany `yaml:",omitempty" json:",omitempty"` // for companies
	Released string    `yaml:",omitempty" json:",omitempty"` // for games, ...

	// general external links
	Website          string   `yaml:",omitempty" json:",omitempty"`
	Websites         []string `yaml:",omitempty" json:",omitempty"`
	Wikipedia        string   `yaml:",omitempty" json:",omitempty"`
	GoodReads        string   `yaml:",omitempty" json:",omitempty"`
	Bookshop         string   `yaml:",omitempty" json:",omitempty"`
	AnimeNewsNetwork string   `yaml:"anime_news_network,omitempty" json:",omitempty"`
	Twitch           string   `yaml:",omitempty" json:",omitempty"`
	YouTube          string   `yaml:",omitempty" json:",omitempty"`
	Vimeo            string   `yaml:",omitempty" json:",omitempty"`
	IMDB             string   `yaml:",omitempty" json:",omitempty"`
	TMDB             string   `yaml:",omitempty" json:",omitempty"`
	TPDB             string   `yaml:",omitempty" json:",omitempty"`
	Steam            string   `yaml:",omitempty" json:",omitempty"`
	Netflix          string   `yaml:",omitempty" json:",omitempty"`
	Spotify          string   `yaml:",omitempty" json:",omitempty"`
	Soundcloud       string   `yaml:",omitempty" json:",omitempty"`
	Hulu             string   `yaml:",omitempty" json:",omitempty"`
	Max              string   `yaml:",omitempty" json:",omitempty"`
	AdultSwim        string   `yaml:",omitempty" json:",omitempty"`
	AppStore         string   `yaml:"app_store,omitempty" json:",omitempty"`
	Fandom           string   `yaml:",omitempty" json:",omitempty"`
	RottenTomatoes   string   `yaml:"rotten_tomatoes,omitempty" json:",omitempty"`
	Metacritic       string   `yaml:",omitempty" json:",omitempty"`
	Opencritic       string   `yaml:",omitempty" json:",omitempty"`
	Twitter          string   `yaml:",omitempty" json:",omitempty"`
	Mastodon         string   `yaml:",omitempty" json:",omitempty"`
	Reddit           string   `yaml:",omitempty" json:",omitempty"`
	Facebook         string   `yaml:",omitempty" json:",omitempty"`
	Instagram        string   `yaml:",omitempty" json:",omitempty"`
	Threads          string   `yaml:",omitempty" json:",omitempty"`
	LinkedIn         string   `yaml:"linkedin,omitempty" json:",omitempty"`
	TikTok           string   `yaml:",omitempty" json:",omitempty"`
	TelegramChannel  string   `yaml:"telegram_channel,omitempty" json:",omitempty"`
	PlayStation      string   `yaml:"playstation,omitempty" json:",omitempty"`
	XBox             string   `yaml:"xbox,omitempty" json:",omitempty"`
	GOG              string   `yaml:"gog,omitempty" json:",omitempty"`
	X                string   `yaml:",omitempty" json:",omitempty"`
	Discord          string   `yaml:",omitempty" json:",omitempty"`
	Epic             string   `yaml:",omitempty" json:",omitempty"`
	IGN              string   `yaml:"ign,omitempty" json:",omitempty"`
	Amazon           string   `yaml:",omitempty" json:",omitempty"`
	PrimeVideo       string   `yaml:"prime_video,omitempty" json:",omitempty"`
	AppleTV          string   `yaml:"apple_tv,omitempty" json:",omitempty"`
	ApplePodcasts    string   `yaml:"apple_podcasts,omitempty" json:",omitempty"`
	AppleBooks       string   `yaml:"apple_books,omitempty" json:",omitempty"`
	Peacock          string   `yaml:",omitempty" json:",omitempty"`
	GooglePlay       string   `yaml:"google_play,omitempty" json:",omitempty"`
	DisneyPlus       string   `yaml:"disney_plus,omitempty" json:",omitempty"`
	MicrosoftStore   string   `yaml:"microsoft_store,omitempty" json:",omitempty"`
	Nintendo         string   `yaml:",omitempty" json:",omitempty"`
	HumbleBundle     string   `yaml:"humble_bundle,omitempty" json:",omitempty"`
	Row8             string   `yaml:",omitempty" json:",omitempty"`
	Redbox           string   `yaml:",omitempty" json:",omitempty"`
	Vudu             string   `yaml:",omitempty" json:",omitempty"`
	DarkHorse        string   `yaml:",omitempty" json:",omitempty"`

	// for books
	ISBN        string    `yaml:",omitempty" json:",omitempty"`
	ISBN10      string    `yaml:",omitempty" json:",omitempty"`
	ISBN13      string    `yaml:",omitempty" json:",omitempty"`
	OCLC        string    `yaml:",omitempty" json:",omitempty"`
	Publishers  oneOrMany `yaml:",omitempty" json:",omitempty"`
	Publication string    `yaml:",omitempty" json:",omitempty"` // date or year of publication

	// for comics
	Artists      oneOrMany `yaml:",omitempty" json:",omitempty"`
	Colorist     string    `yaml:",omitempty" json:",omitempty"`
	Illustrators oneOrMany `yaml:",omitempty" json:",omitempty"`
	Imprint      string    `yaml:",omitempty" json:",omitempty"`
	UPC          string    `yaml:",omitempty" json:",omitempty"`

	// for movies, games, series, ...
	Genres         []string      `yaml:",omitempty" json:",omitempty"`
	Engine         string        `yaml:",omitempty" json:",omitempty"`
	Trailer        string        `yaml:",omitempty" json:",omitempty"`
	Rating         string        `yaml:",omitempty" json:",omitempty"`
	Length         time.Duration `yaml:",omitempty" json:",omitempty"`
	Creators       oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Writers        oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Editors        oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Directors      oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Cinematography oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Producers      oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Screenplay     oneOrMany     `yaml:",omitempty" json:",omitempty"`
	StoryBy        oneOrMany     `yaml:"story_by,omitempty" json:",omitempty"`
	DialoguesBy    oneOrMany     `yaml:"dialogues_by,omitempty" json:",omitempty"`
	Music          oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Production     oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Distributors   oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Network        string        `yaml:",omitempty" json:",omitempty"`
	Composers      oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Programmers    oneOrMany     `yaml:",omitempty" json:",omitempty"`
	Designers      oneOrMany     `yaml:",omitempty" json:",omitempty"`

	// for podcasts
	Hosts  oneOrMany `yaml:",omitempty" json:",omitempty"`
	Guests oneOrMany `yaml:",omitempty" json:",omitempty"`

	RemakeOf Reference `yaml:"remake_of,omitempty" json:",omitempty"`

	Characters []Character `yaml:",omitempty" json:",omitempty"`

	// for awards
	Categories []Category `yaml:",omitempty" json:",omitempty"`

	// unknown fields are stored in the Extra map
	Extra map[string]interface{} `yaml:",inline" json:",omitempty"`

	References []Reference `yaml:"refs,omitempty" json:",omitempty"`

	Episodes []Episode `yaml:",omitempty" json:",omitempty"` // for series

	// fields populated by the generator
	Image                *Media  `yaml:"-" json:",omitempty"`
	Awards               []Award `yaml:"-" json:",omitempty"`
	EditorsAwards        []Award `yaml:"-" json:",omitempty"`
	WritersAwards        []Award `yaml:"-" json:",omitempty"`
	DirectorsAwards      []Award `yaml:"-" json:",omitempty"`
	CinematographyAwards []Award `yaml:"-" json:",omitempty"`
	MusicAwards          []Award `yaml:"-" json:",omitempty"`
	ScreenplayAwards     []Award `yaml:"-" json:",omitempty"`
}
