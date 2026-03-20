// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package mongo

import (
	"context"
	"crypto/tls"
	"fmt"
	"reflect"
	"strings"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"mercedes-benz.ghe.com/foss/disuko/conf"
	"mercedes-benz.ghe.com/foss/disuko/helper/exception"
	"mercedes-benz.ghe.com/foss/disuko/helper/message"
	"mercedes-benz.ghe.com/foss/disuko/infra/repository/database"
	"mercedes-benz.ghe.com/foss/disuko/logy"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Database struct {
	rs         *logy.RequestSession
	client     *mongo.Client
	base       *mongo.Database
	collection *mongo.Collection
	indexes    [][]string
}

func (db *Database) Init(rs *logy.RequestSession, collectionName string, indexes [][]string) {
	uri := fmt.Sprintf(
		"mongodb://%s:%s@%s:%d/?tls=true",
		conf.Config.Database.User,
		conf.Config.Database.Password,
		conf.Config.Database.Host,
		conf.Config.Database.Port,
	)
	api := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(api).SetTLSConfig(&tls.Config{
		InsecureSkipVerify: conf.Config.Database.InsecureSkipVerify,
	})

	client, err := mongo.Connect(opts)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.DatabaseConnection))

	db.client = client
	db.base = client.Database(conf.Config.Database.DatabaseName)
	db.collection = db.base.Collection(collectionName)

	if conf.Config.Server.E2ETests {
		return
	}

	for _, fields := range indexes {
		var keys bson.D
		for _, f := range fields {
			keys = append(keys, bson.E{
				Key:   strings.ToLower(f),
				Value: 1,
			})
		}
		model := mongo.IndexModel{
			Keys: keys,
		}
		_, err = db.collection.Indexes().CreateOne(context.Background(), model)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbEnsureIndex))
	}
}

func (db *Database) QueryQB(qc *database.QueryConfig, createResult func() interface{}) []interface{} {
	filter, opts := buildQuery(qc)
	cursor, err := db.collection.Find(context.Background(), filter, opts)
	if err != nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbRead, filter.String()), err.Error()+" query:"+filter.String())
	}

	var res []interface{}
	for cursor.Next(context.Background()) {
		tmpRes := createResult()
		err = cursor.Decode(tmpRes)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbUnmarshall))
		res = append(res, tmpRes)
	}
	if err := cursor.Err(); err != nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbReadAll), err.Error()+" query:"+filter.String())
	}
	return res
}

func (db *Database) Save(doc interface{}) (string, string) {
	res, err := db.collection.InsertOne(context.Background(), doc)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbCreate, doc))
	id, ok := res.InsertedID.(string)
	if !ok || id == "" {
		logy.Errorf(nil, "unexpected or empty ID after database insertion!")
		return "", ""
	}
	return id, ""
}

func (db *Database) SaveBulk(docs []interface{}) []database.RevKeyHolder {
	_, err := db.collection.InsertMany(context.Background(), docs)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbCreate, "bulk"))
	return nil
}

func (db *Database) DeleteBulk(docs []interface{}) {
	var ids []string
	for _, doc := range docs {
		s := reflect.ValueOf(doc).Elem().FieldByName("Key")
		if !s.IsValid() {
			continue
		}
		ids = append(ids, s.String())
	}
	filter := bson.D{
		bson.E{
			Key: "_id",
			Value: bson.D{
				bson.E{
					Key:   "$in",
					Value: ids,
				},
			},
		},
	}
	_, err := db.collection.DeleteMany(context.Background(), filter)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbDelete))
}

func (db *Database) Delete(key string) {
	filter := bson.D{bson.E{
		Key:   "_id",
		Value: key,
	}}
	_, err := db.collection.DeleteOne(context.Background(), filter)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbDelete))
}

func (db *Database) Update(key string, oldRev string, doc interface{}) string {
	filter := bson.D{bson.E{Key: "_id", Value: key}}
	_, err := db.collection.ReplaceOne(context.Background(), filter, doc)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbUpdate))
	return ""
}

func (db *Database) GetKeyAttribute() string {
	return "_id"
}

func (db *Database) CreateDatabase() {
}

func (db *Database) DropDatabase() {
	err := db.base.Collection(db.collection.Name()).Drop(context.Background())
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbDrop))
}

func (db *Database) Truncate() {
	name := db.collection.Name()
	db.DropDatabase()
	db.collection = db.base.Collection(name)
}
