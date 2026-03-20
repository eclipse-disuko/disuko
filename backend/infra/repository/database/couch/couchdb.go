// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package couch

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/helper/exception"
	"github.com/eclipse-disuko/disuko/helper/message"
	"github.com/eclipse-disuko/disuko/helper/stopwatch"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
	couchdbDriver "github.com/leesper/couchdb-golang"
)

var printThreshold = time.Millisecond * 400

const (
	DBKey = "_id"
	DBRev = "_rev"
)

type DatabaseCouch struct {
	rs           *logy.RequestSession
	client       *couchdbDriver.Server
	couchDb      *couchdbDriver.Database
	indexes      [][]string
	DatabaseName string
}

func (db *DatabaseCouch) Init(rs *logy.RequestSession, collectionName string, indexes [][]string) {
	db.DatabaseName = strings.ToLower(collectionName)
	db.rs = rs
	// TODO: changes this after migration tests
	connString := fmt.Sprintf("%s://%s:%s@%s:%s",
		conf.Config.Database.Scheme,
		conf.Config.Database.User,
		conf.Config.Database.Password,
		conf.Config.Database.Host,
		strconv.Itoa(conf.Config.Database.Port),
	)
	var err error
	db.client, err = couchdbDriver.NewServer(connString)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.DatabaseConnection))

	db.couchDb, err = db.client.Get(strings.ToLower(db.DatabaseName))
	if err != nil {
		db.CreateDatabase()
	}

	for _, i := range indexes {
		db.AddIndex(i)
	}
	db.indexes = indexes
}

func (db *DatabaseCouch) query(query string, createResult func() interface{}) []interface{} {
	sw := stopwatch.StopWatch{}
	sw.Start()
	result, errCB := db.couchDb.QueryJSON(query)
	sw.Stop()
	if sw.DiffTime > printThreshold {
		logy.Infof(db.rs, "Query: %s Time: %s", query, sw.DiffTime)
	}
	if errCB != nil {
		exception.ThrowExceptionServerMessage(message.GetI18N(message.ErrorDbRead, query), errCB.Error()+" query:"+query)
	}
	sw.Start()

	res := make([]interface{}, 0)
	for _, v := range result {
		v["_key"] = v["_id"]
		delete(v, "_id")
		tmp := createResult()

		j, errMarshal := json.Marshal(v)
		exception.HandleErrorServerMessage(errMarshal, message.GetI18N(message.ErrorDbCreate, j))
		errMarshal = json.Unmarshal(j, &tmp)
		exception.HandleErrorServerMessage(errMarshal, message.GetI18N(message.ErrorDbUnmarshall))

		res = append(res, tmp)
	}
	sw.Stop()
	// logy.Infof(db.rs, "Query Marshal: %s Time: %s Length: %d", query, sw.DiffTime, len(res))
	return res
}

func (db *DatabaseCouch) QueryQB(qc *database.QueryConfig, createResult func() interface{}) []interface{} {
	q := BuildQuery(qc)
	return db.query(q, createResult)
}

func (db *DatabaseCouch) Save(doc interface{}) (string, string) {
	var jsonMap map[string]interface{}
	j, _ := json.Marshal(doc)
	err := json.Unmarshal(j, &jsonMap)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonUnmarshall))

	jsonMap["_id"] = jsonMap["_key"]
	delete(jsonMap, "_key")
	delete(jsonMap, "_rev")

	id, rev, err := db.couchDb.Save(jsonMap, url.Values{})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbCreate, doc))
	return id, rev
}

func (db *DatabaseCouch) SaveBulk(docs []interface{}) []database.RevKeyHolder {
	result := make([]database.RevKeyHolder, 0)
	bulkArray := make([]map[string]interface{}, 0)
	for _, entity := range docs {
		var jsonMap map[string]interface{}
		j, _ := json.Marshal(entity)
		err := json.Unmarshal(j, &jsonMap)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonUnmarshall))
		if jsonMap["_key"] != "" {
			jsonMap["_id"] = jsonMap["_key"]
		} else {
			delete(jsonMap, "_id")
		}
		delete(jsonMap, "_key")
		delete(jsonMap, "_rev")
		bulkArray = append(bulkArray, jsonMap)
	}
	res, err := db.couchDb.Update(bulkArray, nil)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbUpdate, "bulk"))
	for _, entity := range res {
		result = append(result, database.RevKeyHolder{Key: entity.ID, Rev: entity.Rev})
	}
	return result
}

func (db *DatabaseCouch) Update(key string, oldRev string, doc interface{}) string {
	var jsonMap map[string]interface{}
	j, _ := json.Marshal(doc)
	err := json.Unmarshal(j, &jsonMap)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonUnmarshall))

	jsonMap["_id"] = jsonMap["_key"]
	delete(jsonMap, "_key")
	err = db.couchDb.Set(key, jsonMap, map[string]string{
		"w": strconv.Itoa(conf.Config.Database.ShardReplica),
	})
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbUpdate, key))
	return jsonMap["_rev"].(string)
}

func (db *DatabaseCouch) Delete(key string) {
	err := db.couchDb.Delete(key)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbDelete, key))
}

func (db *DatabaseCouch) DeleteBulk(docs []interface{}) {
	bulkArray := make([]map[string]interface{}, 0)
	for _, entity := range docs {
		var jsonMap map[string]interface{}
		j, _ := json.Marshal(entity)
		err := json.Unmarshal(j, &jsonMap)
		exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorJsonUnmarshall))

		if jsonMap["_key"] != "" {
			jsonMap["_id"] = jsonMap["_key"]
		}
		delete(jsonMap, "_key")
		jsonMap["_deleted"] = true

		bulkArray = append(bulkArray, jsonMap)
	}
	_, err := db.couchDb.Update(bulkArray, nil)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.ErrorDbDelete, "bulk"))
}

// Warning: This does not apply to possible replication instances
func (db *DatabaseCouch) Truncate() {
	db.DropDatabase()
	db.CreateDatabase()
	for _, i := range db.indexes {
		db.AddIndex(i)
	}
}

func (db *DatabaseCouch) DropDatabase() {
	err := db.client.Delete(db.DatabaseName)
	exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
}

func (db *DatabaseCouch) CreateDatabase() {
	var err error
	db.couchDb, err = db.client.Create(strings.ToLower(db.DatabaseName))
	exception.HandleErrorServerMessage(err, message.GetI18N(message.Error))
}

func (db *DatabaseCouch) GetKeyAttribute() string {
	return DBKey
}

func (db *DatabaseCouch) AddIndex(names []string) {
	db.couchDb.PutIndex(names, "", "")
}
