package library

/* type Proposal struct {
	ID       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title    string             `json:"title" bson:"title"`
	Content  string             `json:"content" bson:"content"`
	AuthorID string             `json:"author" bson:"author"`
	Date     string             `json:"date" bson:"date"`
	Status   string             `json:"status" bson:"status"`
}

func (p *Proposal) InsertToDB() error {
	proposal := bson.M{
		"title":   p.Title,
		"content": p.Content,
		"author":  p.AuthorID,
		"date":    p.Date,
		"status":  p.Status,
	}
	if _, err := db.Proposals.InsertOne(context.TODO(), proposal); err != nil {
		return err
	}
	if err := db.Proposals.FindOne(context.TODO(), bson.M{"title": p.Title}).Decode(&p); err != nil {
		return err
	}
	return nil
}
*/
