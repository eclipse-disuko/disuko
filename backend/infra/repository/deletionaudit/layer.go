package deletionaudit

import (
	"github.com/eclipse-disuko/disuko/domain/deletionaudit"
	"github.com/eclipse-disuko/disuko/infra/repository/base"
)

const DeletionAuditCollectionName = "deletionAuditLog"

type IDeletionAuditRepository interface {
	base.IBaseRepositoryWithHardDelete[*deletionaudit.DeletionAuditEntry]
}
