package repository

import (
	"backend-challenge/internal/domain/entities"
	"backend-challenge/internal/domain/repositories"
	"context"
	"log"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(db *mongo.Client, dbName, collectionName string) repositories.ProductRepository {
    return &productRepository{
        collection: db.Database(dbName).Collection(collectionName),
    }
}

func (r *productRepository) GetByID(id string) (*entities.Product, error) {
    var product entities.Product

    objectId, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return nil, err
    }

    filter := bson.M{"_id": objectId}

    err = r.collection.FindOne(context.TODO(), filter).Decode(&product)

    if err != nil {
        return nil, err
    }

    return &product, nil
}

func (r *productRepository) GetPaginated(offset, limit int) ([]*entities.Product, error) {
    var products []*entities.Product

    findOptions := options.Find()
    findOptions.SetSkip(int64(offset))
    findOptions.SetLimit(int64(limit))
    
    cursor, err := r.collection.Find(context.TODO(), bson.M{}, findOptions)

    if err != nil {
        return nil, err
    }
    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var product entities.Product
        err := cursor.Decode(&product)
        if err != nil {
            return nil, err
        }
        products = append(products, &product)
    }
    
    if err := cursor.Err(); err != nil {
        return nil, err
    }

    return products, nil
}

func (r *productRepository) GetAll() ([]*entities.Product, error) {
    var products []*entities.Product
    
    cursor, err := r.collection.Find(context.TODO(), bson.M{})

    if err != nil {
        return nil, err
    }

    defer cursor.Close(context.TODO())

    for cursor.Next(context.TODO()) {
        var elem *entities.Product
        err := cursor.Decode(&elem)
        if err != nil {
            return nil, err
        }

        products = append(products, elem)
    }

    return products, nil
}

func (r *productRepository) Create(product *entities.Product) error {
    product.ID = primitive.NewObjectID()

    _, err := r.collection.InsertOne(context.TODO(), product)
    return err
}

func (r *productRepository) Update(product *entities.Product) error {
    
    typeData := reflect.TypeOf(*product)

    values := reflect.ValueOf(*product)

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


    result, err := r.collection.UpdateOne(context.TODO(), bson.M{"_id": product.ID}, bson.M{"$set": updates})

	if err != nil {
		return err
	}

    if result.ModifiedCount == 0 {
        log.Println("Produt found but no update performed")
        return nil
    }

	return nil
}

func (r *productRepository) Delete(id string) error {
    objectId, err := primitive.ObjectIDFromHex(id)

    if err != nil {
        return err
    }

    _, err = r.collection.DeleteOne(context.TODO(), bson.M{"_id": objectId})
    return err
}

