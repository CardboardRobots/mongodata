package mongodata

import "go.mongodb.org/mongo-driver/bson/primitive"

func NewObjectId(id string) *primitive.ObjectID {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil
	}
	return &oid
}

func ObjectIdString(oid *primitive.ObjectID) string {
	if oid == nil {
		return ""
	}
	return oid.Hex()
}
