package mongo

type MediaDocument struct {
	AdID   string   `bson:"ad_id"`
	Images []string `bson:"images"`
}
