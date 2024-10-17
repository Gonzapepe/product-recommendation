package repository

import (
	"backend-challenge/internal/domain/entities"
	"backend-challenge/internal/domain/repositories"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"context"
	"reflect"
)

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(db *mongo.Client, dbName, collectionName string) repositories.CategoryRepository {
	return &categoryRepository{
		collection: db.Database(dbName).Collection(collectionName),
	}
}

func (r *categoryRepository) GetAll() ([]*entities.Category, error) {
	var categories []*entities.Category

	cursor, err := r.collection.Find(context.TODO(), bson.M{})

	if err != nil {
		return nil, err
	}

	defer cursor.Close(context.TODO())

	for cursor.Next(context.TODO()) {
		var category entities.Category

		err := cursor.Decode(&category)

		if err != nil {
			return nil, err
		}
		
		categories = append(categories, &category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	cursor.Close(context.TODO())

	return categories, nil
}


func (r *categoryRepository) GetByID(id string) (*entities.Category, error) {
	var category entities.Category

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": objectId}

	err = r.collection.FindOne(context.TODO(), filter).Decode(&category)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Create(category *entities.Category) error {
	category.ID = primitive.NewObjectID()

	_, err := r.collection.InsertOne(context.TODO(), category)

	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) Update(category *entities.Category) error {
	typeData := reflect.TypeOf(*category)

	values := reflect.ValueOf(*category)
	
	updates := bson.D{}

	for i := 1; i < typeData.NumField(); i++ {
		field := typeData.Field(i)
		val := values.Field(i)
		
		jsonTag := field.Tag.Get("json")
		tagParts := strings.Split(jsonTag, ",")
		tag := tagParts[0]
		
		if !isZeroType(val) {
			update := bson.E{Key: tag, Value: val.Interface()}

			updates = append(updates, update)
		}
	}

	_, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": category.ID}, bson.M{"$set": updates})

	if err != nil {
		return err
	}

	return nil
}

func (r *categoryRepository) Delete(id string) error {
	_, err := r.collection.DeleteOne(context.TODO(), bson.M{"_id": id})

	if err != nil {
		return err
	}

	return nil
}

// isZeroType checks if the value from the struct is the zero value of its type
func isZeroType(value reflect.Value) bool {
	zero := reflect.Zero(value.Type()).Interface()

	switch value.Kind() {
	case reflect.Slice, reflect.Array, reflect.Chan, reflect.Map:
		return value.Len() == 0
	default:
		return reflect.DeepEqual(zero, value.Interface())
	}
}