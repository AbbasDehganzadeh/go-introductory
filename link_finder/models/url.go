package models

type UrlAnchor struct {
	Id    int    `db:"id" json:"-" yaml:"-"`
	Page  string `db:"page" json:"page" yaml:"page"`
	Url   string `db:"-" json:"url" yaml:"url"`
	UrlId int    `db:"urlid" json:"-" yaml:"-"`
	Saved bool   `db:"-" json:"-" yaml:"-"`
	Href  string `db:"href" json:"href" yaml:"href"`
	Title string `db:"title" json:"title" yaml:"title"`
	Text  string `db:"text" json:"text" yaml:"text"`
}

type UrlMap map[string]URLResponse
type URLResponse struct {
	Id      int    `db:"id" json:"-" yaml:"-"`
	Origin  string `db:"origin" json:"-" yaml:"-"`
	Path    string `db:"path" json:"-" yaml:"-"`
	Status  int    `db:"status" json:"status" yaml:"status"`
	Reason  string `db:"reason" json:"reason" yaml:"reason"`
	CanSeen bool   `db:"seen" json:"canseen" yaml:"canseen"`
	Visited bool   `db:"visit" json:"visited" yaml:"visited"`
	IsLocal bool   `db:"-" json:"local" yaml:"local"`
	Saved   bool   `db:"-" json:"-" yaml:"-"`
	Error   *error `db:"error" json:"error" yaml:"error"`
}
