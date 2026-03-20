// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package fossdd

import (
	"time"

	"mercedes-benz.ghe.com/foss/disuko/domain/job"
	"mercedes-benz.ghe.com/foss/disuko/domain/license"
	"mercedes-benz.ghe.com/foss/disuko/domain/project"
	"mercedes-benz.ghe.com/foss/disuko/domain/project/components"
	"mercedes-benz.ghe.com/foss/disuko/helper/s3Helper"
	"mercedes-benz.ghe.com/foss/disuko/logy"
)

type ExternCheckCreator interface {
	CreateProjectSPDXExternCheck(requestSession *logy.RequestSession, currentProject *project.Project, componentInfo components.ComponentInfos, sbomUpload *time.Time, sbomKey string) project.SpdxStatusInformation
}

func (g *gen) createPRSnapshot() {
	snapshot := make(map[string]any)
	for _, subData := range g.data.subDatas {
		if len(subData.rules) == 0 {
			snapshot[subData.pr.Name] = "No data available"
			continue
		}
		var dtos []license.PolicyRulePublicResponseDto
		for _, r := range subData.rules {
			dtos = append(dtos, license.PolicyRulePublicResponseDto{
				Key:         r.Key,
				Name:        r.Name,
				Description: r.Description,
				Type:        license.ALLOW,
				Created:     r.Created,
				Updated:     r.Updated,
				Licenses:    license.Licenses(r.allowed).ToPublicResDtos(),
			})
			dtos = append(dtos, license.PolicyRulePublicResponseDto{
				Key:         r.Key,
				Name:        r.Name,
				Description: r.Description,
				Type:        license.WARN,
				Created:     r.Created,
				Updated:     r.Updated,
				Licenses:    license.Licenses(r.warned).ToPublicResDtos(),
			})
			dtos = append(dtos, license.PolicyRulePublicResponseDto{
				Key:         r.Key,
				Name:        r.Name,
				Description: r.Description,
				Type:        license.DENY,
				Created:     r.Created,
				Updated:     r.Updated,
				Licenses:    license.Licenses(r.denied).ToPublicResDtos(),
			})
		}
		snapshot[subData.pr.Name] = dtos
	}
	outPath := g.tempHelper.GetCompleteFileName("pr-snapshot.json")
	s3Helper.SaveObjectToLocalFileSystem(g.rs, outPath, snapshot)
	logy.Infof(g.rs, "policy rule snapshot created: %s", outPath)
	g.jobLog.AddEntry(job.Info, "policy rule snapshot created: %s", outPath)
}

func (g *gen) createPCSnapshot() {
	snapshot := make(map[string]any)
	for _, subData := range g.data.subDatas {
		if subData.evalRes == nil {
			snapshot[subData.pr.Name] = "No data available"
			continue
		}
		snapshot[subData.pr.Name] = g.externCheckCreator.CreateProjectSPDXExternCheck(g.rs, subData.pr, subData.compInfos, subData.spdx.Uploaded, subData.spdx.Key)
	}
	outPath := g.tempHelper.GetCompleteFileName("pc-snapshot.json")
	s3Helper.SaveObjectToLocalFileSystem(g.rs, outPath, snapshot)
	logy.Infof(g.rs, "policy rule evaluation snapshot created: %s", outPath)
	g.jobLog.AddEntry(job.Info, "policy rule evaluation snapshot created: %s", outPath)
}
