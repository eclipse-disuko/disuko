// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/eclipse-disuko/disuko/conf"
	"github.com/eclipse-disuko/disuko/domain/department"
	"github.com/eclipse-disuko/disuko/domain/job"
	"github.com/eclipse-disuko/disuko/domain/label"
	"github.com/eclipse-disuko/disuko/domain/license"
	"github.com/eclipse-disuko/disuko/domain/obligation"
	"github.com/eclipse-disuko/disuko/domain/project"
	"github.com/eclipse-disuko/disuko/domain/schema"
	"github.com/eclipse-disuko/disuko/domain/user"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
	"github.com/eclipse-disuko/disuko/infra/repository/database"
	"github.com/eclipse-disuko/disuko/logy"
)

var entityCreatorMap = map[string]func() interface{}{
	"labels": func() interface{} {
		return &label.Label{}
	},
	"spdxSchemas": func() interface{} {
		return &schema.SpdxSchema{}
	},
	"jobs": func() interface{} {
		return &job.Job{}
	},
	"departments": func() interface{} {
		return &department.Department{}
	},
	"licenses": func() interface{} {
		return &license.License{}
	},
	"obligations": func() interface{} {
		return &obligation.Obligation{}
	},
	"projects": func() interface{} {
		return &project.Project{}
	},
	"rules": func() interface{} {
		return &license.PolicyRules{}
	},
	"users": func() interface{} {
		return &user.User{}
	},
}

func (db *dbRepos) seedDb(requestSession *logy.RequestSession) error {
	seedPath := "./conf/dbseeds/defaultdb/"

	if conf.Config.Server.VanillaDisuko {
		seedPath = "./conf/dbseeds/disuko/"
	}
	entries, err := os.ReadDir(seedPath)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Type().IsDir() || !strings.HasSuffix(e.Name(), ".jsonl") {
			continue
		}

		err = db.processSeedFile(requestSession, seedPath, e.Name())
		if err != nil {
			return err
		}
	}

	if err := db.seedI18n(requestSession); err != nil {
		return err
	}

	return nil
}

func (db *dbRepos) seedI18n(requestSession *logy.RequestSession) error {
	if db.i18nLocale.GetLocaleCount(requestSession) > 0 {
		return nil
	}

	const i18nSeedPath = "./conf/dbseeds/i18n/"
	matches, err := filepath.Glob(i18nSeedPath + "*.json")
	if err != nil {
		return fmt.Errorf("i18n seed glob: %w", err)
	}
	if len(matches) == 0 {
		logy.Debugf(requestSession, "no i18n seed files found, skipping")
		return nil
	}

	localeDisplayNames := map[string][2]string{
		"en": {"English", "English"},
		"de": {"German", "Deutsch"},
	}

	for _, filePath := range matches {
		fileName := filepath.Base(filePath)
		parts := strings.Split(strings.TrimSuffix(fileName, ".json"), ".")
		localeCode := parts[len(parts)-1]

		raw, err := os.ReadFile(filePath)
		if err != nil {
			return fmt.Errorf("i18n seed read %s: %w", filePath, err)
		}
		var entries map[string]string
		if err := json.Unmarshal(raw, &entries); err != nil {
			return fmt.Errorf("i18n seed parse %s: %w", filePath, err)
		}

		displayName, nativeName := localeCode, localeCode
		if names, ok := localeDisplayNames[localeCode]; ok {
			displayName, nativeName = names[0], names[1]
		}
		isDefault := strings.EqualFold(localeCode, "en")
		db.i18nLocale.UpsertLocaleMetadata(requestSession, localeCode, displayName, nativeName, isDefault, "portal")
		for k, v := range entries {
			db.i18nLocale.SetTranslation(requestSession, localeCode, k, v, "Seeded from JSON", "SYSTEM")
		}
		logy.Debugf(requestSession, "i18n seed: loaded %d keys for locale %s from %s", len(entries), localeCode, fileName)
	}

	return nil
}

func (db *dbRepos) processSeedFile(requestSession *logy.RequestSession, path string, filename string) error {
	logy.Debugf(nil, "processing seed file %s", filename)
	f, err := os.Open(path + filename)
	if err != nil {
		return err
	}
	defer f.Close()

	collName := filename[:strings.LastIndex(filename, ".jsonl")]
	entityFn, ok := entityCreatorMap[collName]
	if !ok {
		return nil
		// panic("no entity creator for collection " + collName)
	}
	targetDb := base.NewDatabase()
	targetDb.Init(nil, collName, nil)
	dec := json.NewDecoder(f)
	for dec.More() {
		ent := entityFn()
		if err := dec.Decode(&ent); err != nil {
			return fmt.Errorf("decoding %w", err)
		}
		if collName == "labels" || collName == "spdxSchemas" || collName == "jobs" || collName == "licenses" {
			if err := db.legacyInsertIfNotExists(requestSession, ent, collName, targetDb); err != nil {
				return fmt.Errorf("legacy inserting %w", err)
			}
		} else {
			if err := db.insertIfNotExists(requestSession, ent, targetDb); err != nil {
				return fmt.Errorf("inserting %w", err)
			}
		}
	}

	if collName == "labels" {
		db.label.LoadFromDb(requestSession)
	}

	if conf.Config.Server.VanillaDisuko {
		db.department.LoadFromDb(requestSession)
	}

	return nil
}

func (db *dbRepos) insertIfNotExists(requestSession *logy.RequestSession, ent interface{}, entDb base.IDatabase) error {
	key := reflect.ValueOf(ent).Elem().FieldByName("RootEntity").FieldByName("Key").String()
	if key == "" {
		return errors.New("no key attribute in entity")
	}

	qc := database.New().SetMatcher(
		database.AttributeMatcher(
			entDb.GetKeyAttribute(),
			database.EQ,
			key,
		),
	)
	existing := entDb.QueryQB(qc, func() interface{} {
		var x interface{}
		return &x
	})
	if len(existing) > 0 {
		return nil
	}
	logy.Infof(nil, "inserting seed entity %s", key)
	entDb.Save(ent)

	return nil
}

func (db *dbRepos) legacyInsertIfNotExists(requestSession *logy.RequestSession, ent interface{}, collName string, targetDb base.IDatabase) error {
	if collName == "labels" {
		label := ent.(*label.Label)
		existing := db.label.FindByNameAndType(requestSession, label.Name, label.Type)
		if existing == nil {
			logy.Infof(requestSession, "inserting label %s", label.Name)
			targetDb.Save(ent)
		}
	} else if collName == "spdxSchemas" {
		schema := ent.(*schema.SpdxSchema)
		existingLabel := db.label.FindByNameAndType(requestSession, schema.Label, label.SCHEMA)
		if existingLabel == nil || existingLabel.Type != label.SCHEMA {
			return nil
		}
		existingSchema := db.schema.FindSpdxSchemaByNameAndVersion(requestSession, schema.Name, schema.Version)
		if existingSchema == nil {
			logy.Infof(requestSession, "inserting schema %s", schema.Label)
			schema.Label = existingLabel.Key
			targetDb.Save(schema)
		}
	} else if collName == "jobs" {
		job := ent.(*job.Job)
		existingJob := db.job.FindByTypeAndExecution(requestSession, job.JobType, job.Execution)
		if existingJob == nil {
			logy.Infof(requestSession, "inserting job %s", job.JobType)
			targetDb.Save(ent)
		}
	} else if collName == "licenses" {
		lic := ent.(*license.License)
		existing := db.licenses.FindByIdCaseInsensitive(requestSession, lic.LicenseId)
		if existing == nil {
			logy.Infof(requestSession, "inserting license %s", lic.LicenseId)
			targetDb.Save(ent)
		}
	}
	return nil
}
