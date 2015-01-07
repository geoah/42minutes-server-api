package models

type Episode struct {
	ID            int               `db:"id"`
	Season        int               `json:"season" db:"season"`
	Episode       int               `json:"episode" db:"episode"`
	Number        int               `json:"number" db:"number"`
	TvdbID        int               `json:"tvdb_id" db:"tvdb_id"`
	Title         string            `json:"title" db:"title"`
	Overview      string            `json:"overview" db:"overview"`
	FirstAired    int               `json:"first_aired" db:"first_aired"`
	FirstAiredIso string            `json:"first_aired_iso" db:"first_aired_iso"`
	FirstAiredUtc int               `json:"first_aired_utc" db:"first_aired_utc"`
	URL           string            `json:"url" db:"url"`
	Screen        string            `json:"screen" db:"screen"`
	Images        map[string]string `json:"images" db:"images"`
}
