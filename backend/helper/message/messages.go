// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package message

import (
	"fmt"
	"sync"
)

type I18N struct {
	Code string `json:"code"`
	Text string `json:"text"`
}

var (
	i18nMap   map[string]I18N
	i18nMutex sync.RWMutex
)

func addI18NKeyValue(key string, text string) {
	i18nMutex.Lock()
	defer i18nMutex.Unlock()
	i18nMap[key] = I18N{
		Code: key,
		Text: text,
	}
}

const (
	GetRemarkUnknownKeyError                  = "GET_REMARK_UNKNOWN_KEY_ERROR"
	UpdateRemarkInvalidStatusError            = "UPDATE_REMARK_INVALID_STATUS_ERROR"
	DefaultBranchMainName                     = "main"
	DefaultBranchDevName                      = "dev"
	DefaultBranchMainDescription              = "The main branch serves as the repository for all SBOM files that are included in official releases, ensuring a stable and referenceable collection of release artifacts."
	DefaultBranchDevDescription               = "The dev branch is a workspace for developers to upload temporary SBOMs for quality testing and validation purposes before they are finalized for production."
	ErrorLabelUsed                            = "ERROR_LABEL_USED"
	ErrorLabelNotExist                        = "ERROR_LABEL_NOT_EXIST"
	ErrorLabelAlreadyExist                    = "ERROR_LABEL_ALREADY_EXIST"
	ErrorLabelTypeCanNotChanged               = "ERROR_LABEL_TYPE_CAN_NOT_CHANGED"
	KeyProjectCreateVersion                   = "CREATE_VERSION"
	ErrorDbDeleteNotAllowed                   = "ERROR_REPOSITORY_DELETE_NOT_ALLOWED"
	ErrorDbUpdate                             = "ERROR_REPOSITORY_UPDATE"
	ErrorDbNotFound                           = "ERROR_DB_NOT_FOUND"
	ErrorDbEnsureIndex                        = "ERROR_DB_ENSURE_INDEX"
	ErrorDbTryUpdateOptimized                 = "ERROR_REPOSITORY_TRY_UPDATE_OPTIMIZED"
	ErrorDbCreate                             = "ERROR_REPOSITORY_SAVE"
	ErrorDbRead                               = "ERROR_REPOSITORY_READ"
	ErrorDbReadAll                            = "ERROR_REPOSITORY_READ_ALL"
	ErrorDbUnmarshall                         = "ERROR_REPOSITORY_UNMARSHALL"
	ErrorDbMarshall                           = "ERROR_REPOSITORY_MARSHALL"
	ErrorDbDelete                             = "ERROR_REPOSITORY_DELETE"
	ErrorDbDrop                               = "ERROR_REPOSITORY_DROP"
	ErrorDbExist                              = "ERROR_REPOSITORY_EXIST"
	ErrorKeyProjectAccess                     = "ERROR_PROJECT_ACCESS_KEY"
	ErrorEncrypting                           = "ERROR_ENCRYPTING_PASSWORD"
	ErrorKeyRequestParamEmpty                 = "ERROR_REQUEST_PARAM_EMPTY"
	ErrorKeyRequestParamNotValid              = "ERROR_REQUEST_PARAM_NOT_VALID"
	ErrorStartingJob                          = "ERROR_STARTING_JOB"
	ErrorNotTimedout                          = "ERROR_NOT_TIMEDOUT"
	ErrorPermissionDeniedDownload             = "ERROR_MISSING_DOWNLOAD_PERMISSION"
	ErrorPermissionDeniedUser                 = "ERROR_USER_PERMISSION_DENIED"
	ErrorRunNotOnProd                         = "ERROR_RUN_NOT_ON_PROD"
	ErrorTokenCreate                          = "ERROR_TOKEN_CREATE"
	ErrorAAR                                  = "AAR"
	ErrorS3ClientInit                         = "ERROR_S3_CLIENT_INIT"
	ErrorS3Disabled                           = "ERROR_S3_DISABLED"
	ErrorS3Put                                = "ERROR_S3_PUT"
	ErrorS3Read                               = "ERROR_S3_READ"
	ErrorS3Copy                               = "ERROR_S3_COPY"
	ErrorS3ReadMetaData                       = "ERROR_S3_READ_META_DATA"
	ErrorS3Delete                             = "ERROR_S3_DELETE"
	DiscoTokenUnauthorized                    = "DISCOTOKEN_UNAUTHORIZED"
	ErrorGroupsParentAlreadyExists            = "PARENT_ALREADY_EXISTS"
	ErrorLicenseCanOnlyDeleteCustom           = "ERROR_LICENSE_CAN_ONLY_DELETE_CUSTOM"
	ErrorLicenseFindByObligationKey           = "ERROR_LICENSE_FIND_BY_OBLIGATION_KEY"
	ErrorContentTypeWrong                     = "ERROR_CONTENT_TYPE_WRONG"
	ErrorJsonDecodingInput                    = "ERROR_JSON_DECODING_INPUT"
	ErrorJsonValidating                       = "ERROR_JSON_DECODING_VALIDATING"
	ErrorDbInconsistency                      = "ERROR_DB_INCONSISTENCY"
	ErrorProjectDeleteLastMember              = "ERROR_PROJECT_DELETE_LAST_MEMBER"
	ErrorProjectMemberAlreadyExist            = "ERROR_PROJECT_MEMBER_ALREADY_EXIST"
	ErrorUnexpectError                        = "ERROR_UNEXPECT_ERROR"
	ErrorUnknownFileType                      = "ERROR_FILE_TYPE"
	ErrorUnexpectPanic                        = "ERROR_UNEXPECT_PANIC"
	ErrorQueryUnescape                        = "ERROR_QUERY_UNESCAPE"
	ParamUuidWrong                            = "PARAM_UUID_WRONG"
	ApprovalTypeWrong                         = "PARAM_TYPE_WRONG"
	LicenseDataMissing                        = "LICENSE_DATA_MISSING"
	UrlInvalid                                = "URL_INVALID"
	ParamVersionWrong                         = "PARAM_VERSION_WRONG"
	FindVersion                               = "FIND_VERSION"
	MissingProject                            = "MISSING_PROJECT"
	MultipartReader                           = "MULTIPART_READER"
	PartReader                                = "PART_READER"
	CopyTempFile                              = "COPY_TEMP_FILE"
	ValidateSchema                            = "VALIDATE_SCHEMA"
	WritingContent                            = "WRITING_CONTENT"
	SpdxFileEmptyOrLarge                      = "SPDX_FILE_EMPTY_OR_LARGE"
	MaxFilesize                               = "MAX_FILESIZE"
	ErrorOsOpenFile                           = "ERROR_OS_OPEN_FILE"
	ErrorReadAllAndClose                      = "ERROR_READALL_AND_CLOSE"
	ErrorClose                                = "ERROR_CLOSE"
	ErrorJsonEncode                           = "ERROR_JSON_ENCODE"
	ErrorJsonMarshalling                      = "ERROR_JSON_MARSHALLING"
	DatabaseConnection                        = "DATABASE_CONNECTION"
	GetDatabase                               = "GET_DATABASE"
	GetCollection                             = "GET_COLLECTION"
	SavingJob                                 = "SAVING_JOB"
	ErrorReadLocalFile                        = "ERROR_READ_LOCAL_FILE"
	LoadPolicyRules                           = "LOAD_POLICY_RULES"
	ErrorUserUpdate                           = "USER_UPDATE"
	ErrorUserIdInUse                          = "ERROR_USER_ID_IN_USE"
	UserManagementDryRunSuccess               = "USER_MANAGEMENT_DRY_RUN_SUCCESS"
	UserManagementDeletionSuccess             = "USER_MANAGEMENT_ALL_DELETED"
	UserManagementEntityDeleted               = "USER_MANAGEMENT_ENTITY_DELETED"
	UserManagementUserNotFound                = "USER_MANAGEMENT_USER_NOT_FOUND"
	ErrorLicenseIdInUse                       = "ERROR_LICENSE_ID_IN_USE"
	ErrorOpenFile                             = "ERROR_OPEN_FILE"
	ErrorFileExistAlready                     = "ERROR_FILE_EXIST_ALREADY"
	ErrorCreateFolder                         = "ERROR_CREATE_FOLDER"
	ErrorCreateFile                           = "ERROR_CREATE_FILE"
	ErrorFileWrite                            = "ERROR_FILE_WRITE"
	ErrorFileDelete                           = "ERROR_FILE_DELETE"
	ErrorFileHashCheck                        = "ERROR_FILE_HASH_CHECK"
	ErrorCreateFormFile                       = "ERROR_CREATE_FORM_FILE"
	ErrorIoCopy                               = "ERROR_IO_COPY"
	ErrorCreateRequestPart                    = "ERROR_CREATE_REQUEST_PART"
	FindingSbomKey                            = "FINDING_SBOM_KEY"
	MissingComponentInfo                      = "MISSING_COMPONENT_INFO"
	DataMissing                               = "DATA_MISSING"
	LicenseUpdate                             = "LICENSE_UPDATE"
	UnmarshallingSpdxContent                  = "UNMARSHALLING_SPDX_CONTENT"
	UnmarshallingCache                        = "UNMARSHALLING_CACHE"
	UnmarshallingSpdxContentLicense           = "UNMARSHALLING_SPDX_CONTENT_LICENSE"
	ErrorCsvGenerationDataMissing             = "ERROR_CSV_GENERATION_DATA_MISSING"
	ErrorCsvGeneration                        = "ERROR_CSV_GENERATION"
	ErrorUserNotFound                         = "ERROR_USER_NOT_FOUND"
	ErrorProjectLastOwnerCanNotChanged        = "ERROR_PROJECT_LAST_OWNER_CAN_NOT_CHANGED"
	ReadFile                                  = "READ_FILE"
	ReadFileCorrupted                         = "READ_FILE_CORRUPTED"
	ObligationKeyMissing                      = "OBLIGATION_KEY_MISSING"
	FetchS3File                               = "FETCH_S3_FILE"
	SbomReadInReqFile                         = "SBOM_READ_IN_REQ_FILE"
	SbomMaxComponents                         = "SBOM_MAX_COMPONENTS"
	SbomMaxLicensesPerComponent               = "SBOM_MAX_LICENSES_PER_COMPONENT"
	SbomMaxCopyrightTextSize                  = "SBOM_MAX_COPYRIGHT_TEXT_SIZE"
	SbomMaxPUrlSite                           = "SBOM_MAX_PURL_SIZE"
	SbomComponentsError                       = "SBOM_COMPONENTS_ERROR"
	SbomCreateFile                            = "SBOM_CREATE_FILE"
	IncorectLabelType                         = "INCORECT_LABEL_TYPE"
	ProjectToken                              = "PROJECT_TOKEN"
	Error                                     = "ERROR"
	ErrorNotAJobType                          = "ERROR_NOT_A_JOB_TYPE"
	ErrorAcquiringLock                        = "ERROR_LOCK_ACQUIRE"
	ErrorValidationNotValidSpdxIdentifier     = "ERROR_VALIDATION_NOT_VALID_SPDX_IDENTIFIER"
	DataExistsLicenseId                       = "DATA_EXISTS"
	DataExistsLicenseName                     = "DATA_EXISTS"
	IdConflict                                = "ID_CONFLICT"
	DataMissingObligation                     = "DATA_MISSING"
	ObligationAlreadyExists                   = "OBLIGATION_ALREADY_EXISTS"
	CreatePolicyListAlreadyExists             = "POLICY_LIST_ALREADY_EXISTS"
	UpdatePolicyListAlreadyExists             = "POLICY_LIST_ALREADY_EXISTS"
	PolicyDeprecated                          = "POLICY_DEPRECATED"
	DeleteProject                             = "DELETE_PROJECT"
	ErrorFileMaxUploadPerHourReached          = "ERROR_FILE_MAX_UPLOAD_PER_HOUR_REACHED"
	UpdateVersion                             = "UPDATE_VERSION"
	VersionNameInUse                          = "VERSION_NAME_IN_USE"
	DeleteVersion                             = "DELETE_VERSION"
	ParamSbomUuidEmpty                        = "PARAM_SBOM_UUID_EMPTY"
	ParamRemarkUuidEmpty                      = "PARAM_REMARK_UUID_EMPTY"
	ViewComponents                            = "VIEW_COMPONENTS"
	SchemaValidationFailed                    = "SCHEMA_VALIDATION_FAILED"
	ErrorAddProjectUser                       = "ERROR_ADD_PROJECT_USER"
	UserDeletion                              = "USER_DELETION"
	UserUpdateNotAuthorized                   = "USER_UPDATE"
	ErrorProjectMemberSearch                  = "ERROR_PROJECT_MEMBER_SEARCH"
	NoSchemaProvided                          = "NO_SCHEMA_PROVIDED"
	FoundExistingSchema                       = "FOUND_EXISTING_SCHEMA"
	UploadSbom                                = "UPLOAD_SBOM"
	ErrorInUse                                = "ERROR_IN_USE"
	ErrorHasChildren                          = "ERROR_HAS_CHILDREN"
	ErrorInUseOnUpdate                        = "ERROR_IN_USE_ON_UPDATE"
	SrMissingCopyrightText                    = "SR_MISSING_COPYRIGHT_TEXT"
	SrMissingCopyrightDescription             = "SR_MISSING_COPYRIGHT_DESCRIPTION"
	SrMissingLicenseText                      = "SR_MISSING_LICENSE_TEXT"
	SrMissingLicenseDescription               = "SR_MISSING_LICENSE_DESCRIPTION"
	SrMalformedCopyrightText                  = "SR_MALFORMED_COPYRIGHT_TEXT"
	SrMalformedCopyrightDescription           = "SR_MALFORMED_COPYRIGHT_DESCRIPTION"
	SrCopyrightLongText                       = "SR_COPYRIGHT_LONG_TEXT"
	SrCopyrightToLongDescription              = "SR_COPYRIGHT_TO_LONG_DESCRIPTION"
	SrAliasingUsed                            = "SR_ALIASING_USED"
	SrAliasingUsedDescription                 = "SR_ALIASING_USED_DESCRIPTION"
	SrLicensesDiff                            = "SR_LICENSES_DIFF"
	SrLicensesDiffDescription                 = "SR_LICENSES_DIFF_DESCRIPTION"
	SrMissingLicenseId                        = "SR_MISSING_LICENSE_ID"
	SrMissingLicenseIdDescription             = "SR_MISSING_LICENSE_ID_DESCRIPTION"
	SrMissingName                             = "SR_MISSING_NAME"
	SrMissingNameDescription                  = "SR_MISSING_NAME_DESCRIPTION"
	SrUnmatchedLicenseUsed                    = "SR_UNMATCHED_LICENSE_USED"
	SrUnmatchedLicenseUsedDescription         = "SR_UNMATCHED_LICENSE_USED_DESCRIPTION"
	SrUnknownLicenseUsed                      = "SR_UNKNOWN_LICENSE_USED"
	SrUnknownLicenseUsedDescription           = "SR_UNKNOWN_LICENSE_USED_DESCRIPTION"
	SrMissingVersion                          = "SR_MISSING_VERSION"
	SrMissingVersionDescription               = "SR_MISSING_VERSION_DESCRIPTION"
	SrContainsTooMuchOr                       = "SR_TOO_MUCH_OR"
	SrContainsComplex                         = "SR_CONTAINS_COMPLEX"
	SrContainsSnippet                         = "SR_CONTAINS_SNIPPET"
	SrContainsWith                            = "SR_CONTAINS_WITH"
	SrContainsBadchars                        = "SR_CONTAINS_BADCHARS"
	SrContainsAnnotations                     = "SR_CONTAINS_ANNOTATIONS"
	SrContainsExternalRefs                    = "SR_CONTAINS_EXTERNALREFS"
	SbomReadFile                              = "SBOM_READ_FILE"
	SrTooMuchOrTitle                          = "SR_TOO_MUCH_OR_TITLE"
	SrContentException                        = "SR_CONTENT_EXCEPTION"
	SrProjectAddressException                 = "SR_PROJECT_ADDRESS_EXCEPTION"
	SrProjectAddressDescription               = "SR_PROJECT_ADDRESS_DESCRIPTION"
	SrContainsNonLatinLetters                 = "SR_CONTAINS_NON_LATIN_LETTERS"
	SrContainsNonLatinLettersDescription      = "SR_CONTAINS_NON_LATIN_LETTERS_DESCRIPTION"
	ParamVersionEmpty                         = "PARAM_VERSION_EMPTY"
	ParamProjectNameEmpty                     = "PARAM_PROJECT_NAME_EMPTY"
	ParamVersionToLong                        = "PARAM_VERSION_TO_LONG"
	ErrorVersionDeleted                       = "VERSION_DELETED"
	ErrorTooFewCharacters                     = "ERROR_TOO_FEW_CHARACTERS"
	ErrorVersionMissing                       = "ERROR_VERSION_MISSING"
	ErrorVersionAlreadyExist                  = "ERROR_VERSION_ALREADY_EXIST"
	ErrorTokenNotFoundForKey                  = "ERROR_TOKEN_NOT_FOUND_FOR_KEY"
	ErrorDocumentIdUsesAlready                = "DOCUMENT_ID_USED_ALREADY"
	ReadVersions                              = "READ_VERSIONS"
	CreateVersion                             = "CREATE_VERSION"
	MaxVersionsReached                        = "MAX_VERSIONS_REACHED"
	NoSbom                                    = "NO_SBOM"
	ReadComponentDetails                      = "READ_COMPONENT_DETAILS"
	DownloadExternalSources                   = "DOWNLOAD_EXTERNAL_SOURCES"
	UploadExternalSources                     = "UPLOAD_EXTERNAL_SOURCES"
	DeleteExternalSources                     = "DELETE_EXTERNAL_SOURCE"
	UpdateExternalSource                      = "UPDATE_EXTERNAL_SOURCE"
	ParamSourceidEmpty                        = "PARAM_SOURCEID_EMPTY"
	ParamSpdxidEmpty                          = "PARAM_SPDXID_EMPTY"
	ParamSearchFragmentEmpty                  = "PARAM_SEARCH_FRAGMENT_EMPTY"
	ComponentNotFound                         = "COMPONENT_NOT_FOUND"
	FindingSpdxId                             = "FINDING_SPDX_ID"
	UnmarshallingLicenseContent               = "UNMARSHALLING_LICENSE_CONTENT"
	ReadVersionScanRemarks                    = "READ_VERSION_SCAN_REMARKS"
	ReadVersionLicenseRemarks                 = "READ_VERSION_LICENSE_REMARKS"
	ReadVersionReviewRemarks                  = "READ_VERSION_REVIEW_REMARKS"
	ReadReviewTemplates                       = "READ_REVIEW_TEMPLATES"
	ViewSbom                                  = "VIEW_SBOM"
	CouldNotWrite                             = "COULD_NOT_WRITE"
	CouldNotRead                              = "COULD_NOT_READ"
	FileNotFound                              = "FILE_NOT_FOUND"
	SpdxFileRead                              = "SPDX_FILE_READ"
	StaticReadFile                            = "STATIC_READ_FILE"
	ErrorFindingVersion                       = "ERROR_FINDING_VERSION"
	ReadCompaniesFile                         = "READ_COMPANIES_FILE"
	SpdxFileNotFound                          = "SPDX_FILE_NOT_FOUND"
	ErrorFindingJob                           = "ERROR_FINDING_JOB"
	JobNotFound                               = "JOB_NOT_FOUND"
	JobWrongType                              = "JOB_WRONG_TYPE"
	RequestApp                                = "REQUEST_APP"
	ChangeWithoutConnector                    = "DENY_APP_CHANGE"
	TokenExpired                              = "TOKEN_EXPIRED"
	TokenExpiryExeedsMax                      = "TOKEN_EXPIRY_EXEEDS_MAX"
	ReadProject                               = "READ_PROJECT"
	ProjectOwnerMissing                       = "PROJECT_OWNER_MISSING"
	CompanyValidationFailed                   = "COMPANY_VALIDATION_FAILED"
	UpdateProject                             = "UPDATE_PROJECT"
	CloneProject                              = "CLONE_PROJECT"
	UpdateProjectGroup                        = "UPDATE_PROJECT_GROUP"
	CreateToken                               = "CREATE_TOKEN"
	RenewToken                                = "RENEW_TOKEN"
	RenewTokenError                           = "RENEW_TOKEN"
	RevokingToken                             = "REVOKING_TOKEN"
	RevokingTokenError                        = "REVOKING_TOKEN"
	ProjectRead                               = "PROJECT_READ"
	FindActiveSchemas                         = "FIND_ACTIVE_SCHEMAS"
	S3Error                                   = "S3_ERROR"
	S3ErrorTextFileDownload                   = "S3_ERROR_TEXT_FILE_DOWNLOAD"
	S3ErrorFileUpload                         = "S3_ERROR_FILE_UPLOAD"
	HttpPathUnescape                          = "HTTP_PATH_UNESCAPE"
	ParseSchema                               = "PARSE_SCHEMA"
	SchemaNotFound                            = "SCHEMA_NOT_FOUND"
	DownloadSbomHistory                       = "DOWNLOAD_SBOM_HISTORY"
	UpdateSbom                                = "UPDATE_SBOM"
	LockSbom                                  = "LOCK_SBOM"
	DeleteSbom                                = "DELETE_SBOM"
	SbomIsInApproval                          = "SBOM_IS_IN_APPROVAL"
	SpdxFileKeyNotSet                         = "SPDX_FILE_KEY_NOT_SET"
	ViewNoticeFile                            = "VIEW_NOTICE_FILE"
	WarnNotExists                             = "WARN_NOT_EXISTS"
	ParamProjectUuidEmpty                     = "PARAM_PROJECT_UUID_EMPTY"
	ParamVersionOldEmpty                      = "PARAM_VERSION_OLD_EMPTY"
	ParamSpdxOldWrong                         = "PARAM_SPDX_OLD_WRONG"
	ParamVersionNewEmpty                      = "PARAM_VERSION_NEW_EMPTY"
	ParamSpdxNewWrong                         = "PARAM_SPDX_NEW_WRONG"
	SpdxCompare                               = "SPDX_COMPARE"
	ViewUsers                                 = "VIEW_USERS"
	CreateTokenAuth                           = "CREATETOKEN"
	IdTokenMissing                            = "TOKEN"
	AccessTokenMissing                        = "TOKEN"
	VerifyError                               = "VERIFY"
	Verify                                    = "VERIFY"
	VerifyClaims                              = "VERIFY_CLAIMS"
	Unauthorized                              = "UNAUTHORIZED"
	TokenCreate                               = "TOKEN_CREATE"
	BasicAuth                                 = "BASIC AUTH"
	Auth                                      = "AUTH"
	UserDisabled                              = "USER_DISABLED"
	ErrorJsonUnmarshall                       = "ERROR_JSON_UNMARSHALL"
	ErrorUnexpectedType                       = "ERROR_UNEXPECTED_TYPE"
	RequiresOwner                             = "DENIED_NOT_OWNER"
	ApprovableSPDXParamEmpty                  = "APPROVABLE_PARAM_EMPTY"
	ImageOperationError                       = "IMAGE_OPERATION"
	Conflict                                  = "CONFLICT"
	BadRequest                                = "BAD_REQUEST"
	ResourceInUse                             = "RESOURCE_IN_USE"
	ProblemBadchars                           = "COMP_PROBLEM_BAD_CHARS"
	ProblemComplex                            = "COMP_PROBLEM_COMPLEX"
	ProblemWith                               = "COMP_PROBLEM_WITH"
	ResourceMissing                           = "MISSING_RESOURCE"
	OnlyInternalUsersOwners                   = "INTERNAL_USERS_ONLY"
	NonOwnerResponsible                       = "OWNER_USERS_ONLY"
	OneResponsibleOnly                        = "ONE_MAIN_CONTACT"
	ErrorSpdxInUse                            = "SPDX_IN_USE"
	ErrorProjectDecoupling                    = "ERROR_PROJECT_DECOUPLING"
	TaskNotFound                              = "TASK_NOT_FOUND"
	DelResponsible                            = "DEL_RESPONSIBLE"
	TaskTypeInternalApprovalinfo              = "TASK_TYPE_INTERNAL_APPROVALINFO"
	TaskTypePlausibilityApprovalinfo          = "TASK_TYPE_PLAUSIBILITY_APPROVALINFO"
	TaskTypeInternalApproval                  = "TASK_TYPE_INTERNAL_APPROVAL"
	TaskTypePlausibilityApproval              = "TASK_TYPE_PLAUSIBILITY_APPROVAL"
	ApprovalStatusPlausibilityPending         = "APPROVAL_STATUS_PLAUSIBILITY_PENDING"
	ApprovalStatusPlausibilityDeclined        = "APPROVAL_STATUS_PLAUSIBILITY_DECLINED"
	ApprovalStatusPlausibilityApproved        = "APPROVAL_STATUS_PLAUSIBILITY_APPROVED"
	ApprovalStatusPlausibilityAborted         = "APPROVAL_STATUS_PLAUSIBILITY_ABORTED"
	ApprovalStatusExternalPending             = "APPROVAL_STATUS_EXTERNAL_PENDING"
	ApprovalStatusExternalDeclined            = "APPROVAL_STATUS_EXTERNAL_DECLINED"
	ApprovalStatusExternalAborted             = "APPROVAL_STATUS_EXTERNAL_ABORTED"
	ApprovalStatusExternalGenerationFailed    = "APPROVAL_STATUS_EXTERNAL_GENERATION_FAILED"
	ApprovalStatusExternalSupplierApproved    = "APPROVAL_STATUS_EXTERNAL_SUPPLIER_APPROVED"
	ApprovalStatusExternalCustomerApproved    = "APPROVAL_STATUS_EXTERNAL_CUSTOMER_APPROVED"
	ApprovalStatusInternalPending             = "APPROVAL_STATUS_INTERNAL_PENDING"
	ApprovalStatusInternalDeclined            = "APPROVAL_STATUS_INTERNAL_DECLINED"
	ApprovalStatusInternalApproved            = "APPROVAL_STATUS_INTERNAL_APPROVED"
	ApprovalStatusInternalAborted             = "APPROVAL_STATUS_INTERNAL_ABORTED"
	ApprovalStatusInternalGenerationFailed    = "APPROVAL_STATUS_INTERNAL_GENERATION_FAILED"
	ApprovalStatusInternalDeveloperApproved   = "APPROVAL_STATUS_INTERNAL_DEVELOPER_APPROVED"
	ApprovalStatusPending                     = "APPROVAL_STATUS_PENDING"
	ApprovalStatusDeclined                    = "APPROVAL_STATUS_DECLINED"
	ApprovalStatusApproved                    = "APPROVAL_STATUS_APPROVED"
	ApprovalStatusAborted                     = "APPROVAL_STATUS_ABORTED"
	ApprovalStatusGenerationFailed            = "APPROVAL_STATUS_GENERATION_FAILED"
	FilterSetNotFound                         = "FILTER_SET_NOT_FOUND"
	UserAddedToTheProject                     = "USER_ADDED_TO_THE_PROJECT"
	ConnectorReqFailed                        = "CONNECTOR_REQUEST_FAILED"
	InvalidExchangeCode                       = "INVALID_EXCHANGE_CODE"
	AuthErrorCode                             = "IAM_AUTH_ERROR"
	StatusReviewUnreviewed                    = "SR_UNREVIEWED"
	StatusReviewUnreviewedDE                  = "SR_UNREVIEWED_DE"
	StatusReviewAudited                       = "SR_AUDITED"
	StatusReviewAuditedDE                     = "SR_AUDITED_DE"
	StatusReviewAcceptable                    = "SR_ACCEPTABLE"
	StatusReviewAcceptableDE                  = "SR_ACCEPTABLE_DE"
	StatusReviewNotAcceptable                 = "SR_NOT_ACCEPTABLE"
	StatusReviewNotAcceptableDE               = "SR_NOT_ACCEPTABLE_DE"
	StatusReviewAcceptableAfterChanges        = "SR_ACCEPTABLE_AFTER_CHANGES"
	StatusReviewAcceptableAfterChangesDE      = "SR_ACCEPTABLE_AFTER_CHANGES_DE"
	ReadLicenseRules                          = "READ_LICENSE_RULES"
	CreateLicenseRule                         = "CREATE_LICENSE_RULES"
	EditLicenseRule                           = "EDIT_LICENSE_RULES"
	ParamLicenseRuleUuidEmpty                 = "PARAM_LICENSE_RULE_UUID_EMPTY"
	UnknownPattern                            = "UNKNOWN_PATTERN"
	ActiveLicenseRuleExists                   = "ACTIVE_LICENSE_RULE_EXISTS"
	DeprecatedProjectError                    = "DEPRECATED_PROJECT_ERROR"
	SpdxFileContentTypeValidation             = "SPDX_FILE_CONTENT_TYPE_VALIDATION"
	CustomIdKeyMalformed                      = "MALFORMED_CUSTOM_ID_KEY"
	SchedulerPubSub                           = "SCHEDULER_PUB_SUB_FAILED"
	ErrorProjectHasActiveApprovalsOrReviews   = "ERROR_PROJECT_HAS_ACTIVE_APPROVALS_OR_REVIEWS"
	ErrorGroupHasActiveChildren               = "ERROR_GROUP_HAS_ACTIVE_CHILDREN"
	SpdxAlreadyLocked                         = "SPDX_ALREADY_LOCKED"
	SpdxNotLocked                             = "SPDX_NOT_LOCKED"
	SpdxRetainedForApprovalOrReview           = "SPDX_RETAINED_FOR_APPROVAL_OR_REVIEW"
	ProjectGroupRequired                      = "PROJECT_GROUP_REQUIRED"
	ErrorProjectHasDummyLabel                 = "ERROR_PROJECT_HAS_DUMMY_LABEL"
	ChoiceDeniedResp                          = "CHOICE_DENIED_RESP"
	ChoiceDeniedMassive                       = "CHOICE_DENIED_MASSIVE"
	PolicyDecisionOperationNotAuthorized      = "POLICY_DECISION_OPERATION_NOT_AUTHORIZED"
	ActivePolicyDecisionExists                = "ACTIVE_POLICY_DECISION_EXISTS"
	InvalidPolicyDecisionData                 = "INVALID_POLICY_DECISION_DATA"
	InvalidPolicyDecisionLicenseData          = "INVALID_POLICY_DECISION_LICENSE_DATA"
	InvalidPolicyDecisionLicenseApprovalState = "INVALID_POLICY_DECISION_LICENSE_APPROVAL_STATE"
	InvalidBulkPolicyDecisionData             = "INVALID_BULK_POLICY_DECISION_DATA"
	PolicyDecisionDeniedNotResponsible        = "POLICY_DECISION_DENIED_NOT_RESPONSIBLE"
	PolicyDecisionDeniedNotFossOfficeUser     = "POLICY_DECISION_DENIED_NOT_FOSS_OFFICE_USER"
	ParamPolicyDecisionUuidEmpty              = "PARAM_POLICY_DECISION_UUID_EMPTY"
	DecisionDeniedComponentVersionNotSet      = "DECISION_DENIED_COMPONENT_VERSION_NOT_SET"
	PolicyDecisionDeniedForbiddenLicense      = "POLICY_DECISION_DENIED_FORBIDDEN_LICENSE"
)

func InitI18N() {
	i18nMap = make(map[string]I18N)
	addI18NKeyValue(ErrorUnexpectError, "The error occurs, please try again!")
	addI18NKeyValue(ErrorUnexpectPanic, "The error occurs, please try again!")
	addI18NKeyValue(KeyProjectCreateVersion, "Project version already exists")
	addI18NKeyValue(MaxVersionsReached, "Maximum Number of Project Versions reached!")
	addI18NKeyValue(ErrorDbUpdate, "Error update %s to database")
	addI18NKeyValue(ErrorUnknownFileType, "Error unknown file type %s")
	addI18NKeyValue(ErrorDbDeleteNotAllowed, "Deletion of data from the collection %s is not allowed")
	addI18NKeyValue(ErrorDocumentIdUsesAlready, "Document Id already used in %s")
	addI18NKeyValue(ErrorDbInconsistency, "Database error")
	addI18NKeyValue(ErrorDbTryUpdateOptimized, "Error try updating optimized entity %s to database")
	addI18NKeyValue(ErrorDbCreate, "Error save %s to database")
	addI18NKeyValue(ErrorDbNotFound, "Database element not found")
	addI18NKeyValue(ErrorProjectDeleteLastMember, "the last owner cannot be removed")
	addI18NKeyValue(ErrorProjectMemberAlreadyExist, "error project member already exist: %s")
	addI18NKeyValue(ErrorPermissionDeniedDownload, "Permission denied to download %s")
	addI18NKeyValue(ErrorDbRead, "Error read %s from database")
	addI18NKeyValue(ErrorDbReadAll, "Error read all from database")
	addI18NKeyValue(ErrorDbUnmarshall, "Error unmarshall %s data from database")
	addI18NKeyValue(ErrorDbDelete, "Error delete %s from database")
	addI18NKeyValue(ErrorDbExist, "Error check exist %s from database")
	addI18NKeyValue(ErrorKeyRequestParamEmpty, "Request parameter %s is empty")
	addI18NKeyValue(ErrorKeyRequestParamNotValid, "Request parameter %s is not valid")
	addI18NKeyValue(ErrorPermissionDeniedUser, "Permission denied for user %s")
	addI18NKeyValue(ErrorRunNotOnProd, "This function is not available on prod")
	addI18NKeyValue(ErrorTokenCreate, "Could not create token")
	addI18NKeyValue(ErrorAAR, "%s")
	addI18NKeyValue(ErrorS3ClientInit, "Could not init S3 client")
	addI18NKeyValue(ErrorS3Disabled, "S3 is disabled")
	addI18NKeyValue(ErrorS3Put, "Error S3 on put operation for file: %s")
	addI18NKeyValue(ErrorS3Copy, "Error S3 on copy file '%s' to '%s'")
	addI18NKeyValue(ErrorS3Read, "Error S3 on read operation for file: %s")
	addI18NKeyValue(ErrorS3ReadMetaData, "Error S3 on read meta data for file: %s")
	addI18NKeyValue(ErrorS3Delete, "Error S3 on delete file: %s")
	addI18NKeyValue(ErrorLicenseCanOnlyDeleteCustom, "Can't delete license cause source is not custom for key %s")
	addI18NKeyValue(ErrorLicenseFindByObligationKey, "failed to get licenses by obligation key: %s")
	addI18NKeyValue(ErrorGroupsParentAlreadyExists, "One of the children has already a parent set.")
	addI18NKeyValue(ErrorJsonDecodingInput, "The JSON could not be decoded")
	addI18NKeyValue(ErrorJsonValidating, "The json contains validation errors.")
	addI18NKeyValue(ErrorQueryUnescape, "")
	addI18NKeyValue(ParamUuidWrong, "project uuid parameter is wrong")
	addI18NKeyValue(ApprovalTypeWrong, "approval type parameter is wrong")
	addI18NKeyValue(LicenseDataMissing, "Error getting license from db")
	addI18NKeyValue(UrlInvalid, "URL is not valid")
	addI18NKeyValue(ParamVersionWrong, "Project version parameter is wrong")
	addI18NKeyValue(FindVersion, "The project does not contain the requested version")
	addI18NKeyValue(MissingProject, "Could not retrieve project.")
	addI18NKeyValue(MultipartReader, "Error getting multipart reader")
	addI18NKeyValue(PartReader, "Error reading mutli part")
	addI18NKeyValue(CopyTempFile, "Error copying temp file to final")
	addI18NKeyValue(ValidateSchema, "Schema invalid")
	addI18NKeyValue(WritingContent, "Error writing content to ResponseWriter")
	addI18NKeyValue(SpdxFileEmptyOrLarge, "The content of the SPDX is empty or to large!")
	addI18NKeyValue(MaxFilesize, "Error reading request")
	addI18NKeyValue(ErrorOsOpenFile, "Could not open file: %s")
	addI18NKeyValue(ErrorReadAllAndClose, "")
	addI18NKeyValue(ErrorJsonEncode, "")
	addI18NKeyValue(ErrorJsonMarshalling, "Error marshalling object to json")
	addI18NKeyValue(DatabaseConnection, "failed connect to database disclosure")
	addI18NKeyValue(GetDatabase, "failed to get database: %s")
	addI18NKeyValue(GetCollection, "failed to get the collection %s")
	addI18NKeyValue(SavingJob, "Error saving job")
	addI18NKeyValue(ErrorReadLocalFile, "%s")
	addI18NKeyValue(LoadPolicyRules, "Error while getting policy rules from database")
	addI18NKeyValue(ErrorUserUpdate, "Error during user update")
	addI18NKeyValue(ErrorUserIdInUse, "Error user id (%s) in use!")
	addI18NKeyValue(ErrorLicenseIdInUse, "Error license id (%s) in use!")
	addI18NKeyValue(ErrorOpenFile, "%s")
	addI18NKeyValue(ErrorFileExistAlready, "Error file exist already: %s")
	addI18NKeyValue(ErrorCreateFolder, "failed to create folder %s")
	addI18NKeyValue(ErrorCreateFile, "failed to create file %s")
	addI18NKeyValue(ErrorFileWrite, "failed to write to file %s")
	addI18NKeyValue(ErrorFileDelete, "failed to delete file %s")
	addI18NKeyValue(ErrorFileHashCheck, "corrupted file, sha256 hash do not match; file: %s; hash s3: %s; hash db: %s")
	addI18NKeyValue(ErrorCreateFormFile, "%s")
	addI18NKeyValue(ErrorIoCopy, "%s")
	addI18NKeyValue(ErrorCreateRequestPart, "%s")
	addI18NKeyValue(FindingSbomKey, "Could not find sbom")
	addI18NKeyValue(DataMissing, "Error getting license from db")
	addI18NKeyValue(LicenseUpdate, "License to update not found with id: %s")
	addI18NKeyValue(UnmarshallingSpdxContent, "Error unmarshalling spdx content into object")
	addI18NKeyValue(UnmarshallingSpdxContentLicense, "No license info found for spdx component")
	addI18NKeyValue(ErrorCsvGenerationDataMissing, "error fetching %s from the database")
	addI18NKeyValue(ErrorUserNotFound, "Error user not found %s")
	addI18NKeyValue(ErrorProjectLastOwnerCanNotChanged, "Error the last project owner cannot be changed")
	addI18NKeyValue(ReadFile, "The %s file could not be read from HD")
	addI18NKeyValue(ReadFileCorrupted, "The %s file is corrupted")
	addI18NKeyValue(ObligationKeyMissing, "Key is missing for obligation with id %s")
	addI18NKeyValue(SrProjectAddressException, "Missing Contact Address")
	addI18NKeyValue(SrProjectAddressDescription, "The contact address in the project settings for display in the third party notice is missing. It is required to set a custom address in the project settings, or the project owner needs to make sure that the default contact address fits with the local guidelines for disclosure.")
	addI18NKeyValue(SrContainsNonLatinLetters, "Non-Latin Package URL")
	addI18NKeyValue(SrContainsNonLatinLettersDescription, "The Package URL contains one or more non-Latin characters. It is recommended to check this package to prevent potential supply chain dependency attacks, for example combo-squatting.")
	addI18NKeyValue(FetchS3File, "Error fetching file from S3")
	addI18NKeyValue(SbomReadInReqFile, "The SBOM file could not be read (%v)")
	addI18NKeyValue(SbomCreateFile, "The SBOM file could not be created / saved")
	addI18NKeyValue(IncorectLabelType, "label type is incorrect")
	addI18NKeyValue(ProjectToken, "Project uuid wrong")
	addI18NKeyValue(Error, "")
	addI18NKeyValue(DataExistsLicenseId, "License with licenseId %s already exists")
	addI18NKeyValue(DataExistsLicenseName, "License with name %s already exists")
	addI18NKeyValue(IdConflict, "Aliases contain duplicates")
	addI18NKeyValue(DataMissingObligation, "Obligation name is empty")
	addI18NKeyValue(ObligationAlreadyExists, "Error creating classification. The name is already used")
	addI18NKeyValue(CreatePolicyListAlreadyExists, "Error creating allow/deny list. The name is already used")
	addI18NKeyValue(UpdatePolicyListAlreadyExists, "Error updating the allow/deny list name, the name is already used")
	addI18NKeyValue(DeleteProject, "User is not authorized to delete this project")
	addI18NKeyValue(ErrorFileMaxUploadPerHourReached, "max uploads per hour reached for the project (%s / %s): uploads %s, max. uploads: %s")
	addI18NKeyValue(UpdateVersion, "User is not authorized to update channel for this project")
	addI18NKeyValue(VersionNameInUse, "Channel name already in use: %s")
	addI18NKeyValue(DeleteVersion, "User is not authorized to delete channel for this project")
	addI18NKeyValue(ParamSbomUuidEmpty, "Url parameter 'sbomUuid' missing")
	addI18NKeyValue(ParamRemarkUuidEmpty, "Url parameter 'remarkId' missing")
	addI18NKeyValue(ViewComponents, "User is not authorized to view components on version for this project")
	addI18NKeyValue(SchemaValidationFailed, "Error validating the SPDX file")
	addI18NKeyValue(ErrorAddProjectUser, "User is not authorized to add user for this project")
	addI18NKeyValue(UserDeletion, "User is not authorized to delete user for this project")
	addI18NKeyValue(UserUpdateNotAuthorized, "User is not authorized to update user for this project")
	addI18NKeyValue(ErrorProjectMemberSearch, "User is not authorized to read user for this project")
	addI18NKeyValue(NoSchemaProvided, "Schema file not provided")
	addI18NKeyValue(FoundExistingSchema, "Schema with name and version already exists")
	addI18NKeyValue(UploadSbom, "User is not authorized to upload SBOM for this project")
	addI18NKeyValue(ErrorInUse, "Can not be deleted, it is in use: %s")
	addI18NKeyValue(ErrorHasChildren, "Can not be deleted, please remove project children first")
	addI18NKeyValue(ErrorInUseOnUpdate, "ERROR_IN_USE_ON_UPDATE_KEY")
	addI18NKeyValue(ErrorValidationNotValidSpdxIdentifier, "Spdx identifier is not valid: %s")
	addI18NKeyValue(ErrorNotAJobType, "Error not a job type: %s")
	addI18NKeyValue(SrMissingCopyrightText, "Missing Copyright Information")
	addI18NKeyValue(SrMissingCopyrightDescription, "Copyright information is missing; the warn level varies depending on the risk level of distribution. It is required to ascertain and provide this information in an updated SBOM delivery. Copyright information can typically be found on the license text of a component, within the project’s readme file, or written as code header in the main source code files.")
	addI18NKeyValue(SrMissingLicenseText, "Missing License Text")
	addI18NKeyValue(SrMissingLicenseDescription, "License text is missing. Without the license text of a custom license, you cannot check if the text compares to the custom license text stored in the FOSS License Database for the respective license identifier or alias. It is required to properly identify the custom license text and provide this information in an updated SBOM delivery.")
	addI18NKeyValue(SrMalformedCopyrightText, "Malformed Copyright Information")
	addI18NKeyValue(SrMalformedCopyrightDescription, "Copyright information contains special or unexpected characters. It is recommended to check this specific copyright information. In case it is malformed please check with the development team to correct this information and provide it in an updated SBOM delivery.")
	addI18NKeyValue(SrCopyrightLongText, "Lengthy Copyright Information")
	addI18NKeyValue(SrCopyrightToLongDescription, "Copyright information is unusually lengthy. It is recommended to check this specific copyright information. In case it is malformed this information needs to be corrected and provided in an updated SBOM delivery.")
	addI18NKeyValue(SrAliasingUsed, "Aliased License")
	addI18NKeyValue(SrAliasingUsedDescription, "A license alias was used for this component.")
	addI18NKeyValue(SrLicensesDiff, "Consistency Check")
	addI18NKeyValue(SrLicensesDiffDescription, "License Declared and License Concluded differs. It is required to check the components original source, identify the reason for this difference and confirm this difference is acceptable.")
	addI18NKeyValue(SrMissingLicenseId, "Missing License Identifier")
	addI18NKeyValue(SrMissingLicenseIdDescription, "The license identifier is missing. It is required to properly identify the component license and provide this information in an updated SBOM delivery.")
	addI18NKeyValue(SrMissingName, "Missing Component Name")
	addI18NKeyValue(SrMissingNameDescription, "Component name is missing. Therefore, the component cannot be properly handled. It is required to properly identify the component name and provide this information in an updated SBOM delivery.")
	addI18NKeyValue(SrUnknownLicenseUsed, "Unknown Component License")
	addI18NKeyValue(SrUnmatchedLicenseUsed, "Unmatched Component License")
	addI18NKeyValue(SrUnmatchedLicenseUsedDescription, "License can't be matched with a policy rule")
	addI18NKeyValue(SrContainsTooMuchOr, "The license expression contains an unusually large number of licenses to choose from (OR). It is required to carefully review license expressions of this kind for correctness to avoid legal risk. It is recommended to check the components original source.")
	addI18NKeyValue(SrTooMuchOrTitle, "Extensive OR Expression")
	addI18NKeyValue(SrContainsComplex, "SBOM contains unusual content which cannot be processed properly: Complex logical license expressions; Please contact the FOSS Disclosure Portal tech support to handle this exception. Currently only AND-expressions and OR-expressions are handled, but not combinations thereof.")
	addI18NKeyValue(SrContainsSnippet, "SBOM contains unusual content which cannot be processed properly: Snippet items; please contact the FOSS Disclosure Portal tech support to handle this exception.")
	addI18NKeyValue(SrContainsWith, "SBOM contains unusual content which cannot be processed properly: License expression with license exception; please contact the FOSS Disclosure Portal tech support to handle this exception. Currently only specific combinations of licenses and exceptions are handled with a specific custom license identifier, but not arbitrary combinations of licenses and exceptions.")
	addI18NKeyValue(SrContainsBadchars, "SBOM contains unusual content which cannot be processed properly: License name contains bad characters; please contact the FOSS Disclosure Portal tech support to handle this exception or check with the development team to properly identify the component license and provide this information in an updated SBOM delivery.")
	addI18NKeyValue(SrContainsAnnotations, "SBOM contains unusual content which cannot be processed properly: Annotations; please contact the FOSS Disclosure Portal tech support to handle this exception.")
	addI18NKeyValue(SrContainsExternalRefs, "SBOM contains unusual content which cannot be processed properly: External references; please contact the FOSS Disclosure Portal tech support to handle this exception.")
	addI18NKeyValue(SrContentException, "SBOM Content Exception")
	addI18NKeyValue(SrUnknownLicenseUsedDescription, "License is unknown.")
	addI18NKeyValue(SrMissingVersion, "Missing Component Version")
	addI18NKeyValue(SrMissingVersionDescription, "Component Version is missing. Therefore, the component cannot be exactly identified. It is required to properly identify the component version and provide this information in an updated SBOM delivery.")
	addI18NKeyValue(ErrorLabelAlreadyExist, "Label already exists")
	addI18NKeyValue(ErrorLabelTypeCanNotChanged, "Label type can not changed!")
	addI18NKeyValue(ErrorLabelUsed, "Label is already used")
	addI18NKeyValue(ErrorLabelNotExist, "Label not exist")
	addI18NKeyValue(SbomReadFile, "The SBOM file could not be read from HD")
	addI18NKeyValue(ParamVersionEmpty, "Project channel parameter is missing")
	addI18NKeyValue(ParamProjectNameEmpty, "Project name is missing")
	addI18NKeyValue(ParamVersionToLong, "Project channel parameter is to long")
	addI18NKeyValue(ErrorVersionDeleted, "The channel is marked as deleted!")
	addI18NKeyValue(ErrorTooFewCharacters, "Too few characters!")
	addI18NKeyValue(ErrorVersionMissing, "The channel is missing!")
	addI18NKeyValue(ErrorVersionAlreadyExist, "The channel already exist, %s!")
	addI18NKeyValue(ErrorTokenNotFoundForKey, "no token found for key %s")
	addI18NKeyValue(ReadVersions, "User is not authorized to view versions for this project")
	addI18NKeyValue(CreateVersion, "User is not authorized to create channel for this project")
	addI18NKeyValue(NoSbom, "no sbom was uploaded")
	addI18NKeyValue(ReadComponentDetails, "User is not authorized to view component details on version for this project")
	addI18NKeyValue(DownloadExternalSources, "User is not authorized to read external sources for this project")
	addI18NKeyValue(UploadExternalSources, "User is not authorized to create external sources for this project")
	addI18NKeyValue(DeleteExternalSources, "User is not authorized to delete external sources for this project")
	addI18NKeyValue(UpdateExternalSource, "User is not authorized to update external sources for this project")
	addI18NKeyValue(ParamSourceidEmpty, "Url parameter 'sourceId' missing")
	addI18NKeyValue(ParamSpdxidEmpty, "Url parameter 'spdxId' missing")
	addI18NKeyValue(FindingSpdxId, "Could not find requested spdxId (%s) in content")
	addI18NKeyValue(UnmarshallingLicenseContent, "Error unmarshalling license content into object")
	addI18NKeyValue(ReadVersionScanRemarks, "User is not authorized to perform action on version for this project")
	addI18NKeyValue(ReadVersionLicenseRemarks, "User is not authorized to perform action on version for this project")
	addI18NKeyValue(ReadVersionReviewRemarks, "User is not authorized to perform action on version for this project")
	addI18NKeyValue(ReadReviewTemplates, "User is not authorized to view review templates for this project")
	addI18NKeyValue(ViewSbom, "User is not authorized to view SBOM on version for this project")
	addI18NKeyValue(CouldNotWrite, "Cannot write temporary file")
	addI18NKeyValue(FileNotFound, "Cannot read template")
	addI18NKeyValue(SpdxFileRead, "The content of the SPDX could not be read")
	addI18NKeyValue(StaticReadFile, "Static resource unavailable")
	addI18NKeyValue(ErrorFindingVersion, "Error retrieving channel from database")
	addI18NKeyValue(ReadCompaniesFile, "Could not ready companies file.")
	addI18NKeyValue(SpdxFileNotFound, "Found no spdx file with file key: %s")
	addI18NKeyValue(ErrorFindingJob, "Error during job DB request")
	addI18NKeyValue(ErrorAcquiringLock, "Error acquiring lock")
	addI18NKeyValue(JobNotFound, "The requested job could not be found")
	addI18NKeyValue(RequestApp, "Error requesting application connector. Value is wrong.")
	addI18NKeyValue(ChangeWithoutConnector, "Changing the application is only allowed with a connector set.")
	addI18NKeyValue(TokenExpired, "Token already expired")
	addI18NKeyValue(TokenExpiryExeedsMax, "Token expiry exceeds max expiry")
	addI18NKeyValue(ReadProject, "User is not authorized to view this project")
	addI18NKeyValue(ProjectOwnerMissing, "Error creating project. Project owner is not provided")
	addI18NKeyValue(CompanyValidationFailed, "Error during companies validation")
	addI18NKeyValue(UpdateProject, "User is not authorized to update this project")
	addI18NKeyValue(CloneProject, "User is not authorized to clone this project")
	addI18NKeyValue(CreateToken, "User is not authorized to create token for this project")
	addI18NKeyValue(RenewToken, "User is not authorized to renew token for this project")
	addI18NKeyValue(RenewTokenError, "Error updating token")
	addI18NKeyValue(RevokingToken, "User is not authorized to revoke token for this project")
	addI18NKeyValue(RevokingTokenError, "Error revoking token")
	addI18NKeyValue(ProjectRead, "User is not authorized to read this project")
	addI18NKeyValue(FindActiveSchemas, "User is not authorized to view schema for this project")
	addI18NKeyValue(S3Error, "%s")
	addI18NKeyValue(S3ErrorTextFileDownload, "%s")
	addI18NKeyValue(S3ErrorFileUpload, "%s")
	addI18NKeyValue(HttpPathUnescape, "%s")
	addI18NKeyValue(ParseSchema, "Error parsing schema")
	addI18NKeyValue(SchemaNotFound, "Requested schema not found in DB with id %s")
	addI18NKeyValue(DownloadSbomHistory, "User is not authorized to download SBOM history for this project")
	addI18NKeyValue(UpdateSbom, "User is not authorized to update SBOM on version for this project")
	addI18NKeyValue(LockSbom, "User is not authorized to lock/unlock SBOM on version for this project")
	addI18NKeyValue(DeleteSbom, "User is not authorized to delete SBOM on version for this project")
	addI18NKeyValue(SbomIsInApproval, "Update SBOM under approval process is not allowed")
	addI18NKeyValue(SpdxFileKeyNotSet, "No spdx file key set!")
	addI18NKeyValue(ViewNoticeFile, "User is not authorized to view notice file on version for this project")
	addI18NKeyValue(WarnNotExists, "Info currently notice file not exists")
	addI18NKeyValue(ParamProjectUuidEmpty, "Parameter is missing: %s")
	addI18NKeyValue(ParamVersionOldEmpty, "Parameter is missing: %s")
	addI18NKeyValue(ParamSpdxOldWrong, "Parameter is missing: %s")
	addI18NKeyValue(ParamVersionNewEmpty, "Parameter is missing: %s")
	addI18NKeyValue(ParamSpdxNewWrong, "Parameter is missing: %s")
	addI18NKeyValue(SpdxCompare, "Error during spdx file compare")
	addI18NKeyValue(ViewUsers, "User is not authorized to read users for this project")
	addI18NKeyValue(CreateTokenAuth, "Error creating token from code")
	addI18NKeyValue(IdTokenMissing, "Missing token in oauth answer")
	addI18NKeyValue(AccessTokenMissing, "Missing access token in oauth answer")
	addI18NKeyValue(VerifyError, "Verify error")
	addI18NKeyValue(Verify, "Unable to retrieve user info from CIAM")
	addI18NKeyValue(VerifyClaims, "Verify claim error")
	addI18NKeyValue(Unauthorized, "User unauthorized")
	addI18NKeyValue(TokenCreate, "Failed to create a jwt token")
	addI18NKeyValue(BasicAuth, "Missing Basic")
	addI18NKeyValue(Auth, "Unauthorized")
	addI18NKeyValue(UserDisabled, "User disabled")
	addI18NKeyValue(ErrorJsonUnmarshall, "Could not unmarshall JSON.")
	addI18NKeyValue(ErrorUnexpectedType, "Unexpected type.")
	addI18NKeyValue(RequiresOwner, "You have to be the owner of this project to perform this action.")
	addI18NKeyValue(ApprovableSPDXParamEmpty, "Either both or no keys have to be set.")
	addI18NKeyValue(ResourceInUse, "This resource is currently in use by another instance.")
	addI18NKeyValue(OnlyInternalUsersOwners, "Only internal users can be owners")
	addI18NKeyValue(NonOwnerResponsible, "Only owner users can be flagged as project responsibles")
	addI18NKeyValue(ErrorSpdxInUse, "SPDX_IN_USE_MESSAGE_KEY")
	addI18NKeyValue(ErrorProjectDecoupling, "Can not be decoupled, project(s) was(were) subject to an approval: %s")
	addI18NKeyValue(TaskNotFound, "Task not found")
	addI18NKeyValue(ErrorContentTypeWrong, "Unexpected content type. In case you see this message while using the disclosure-cli, please update to the latest version.")
	addI18NKeyValue(ApprovalStatusPlausibilityPending, "Pending")
	addI18NKeyValue(ApprovalStatusPlausibilityDeclined, "NOT OK")
	addI18NKeyValue(ApprovalStatusPlausibilityApproved, "OK")
	addI18NKeyValue(ApprovalStatusPlausibilityAborted, "Aborted")
	addI18NKeyValue(ApprovalStatusExternalPending, "Pending")
	addI18NKeyValue(ApprovalStatusExternalDeclined, "Declined")
	addI18NKeyValue(ApprovalStatusExternalAborted, "Aborted")
	addI18NKeyValue(ApprovalStatusExternalGenerationFailed, "Generation Failed")
	addI18NKeyValue(ApprovalStatusExternalSupplierApproved, "Developer Approved")
	addI18NKeyValue(ApprovalStatusExternalCustomerApproved, "Owner Approved")
	addI18NKeyValue(ApprovalStatusInternalPending, "Pending")
	addI18NKeyValue(ApprovalStatusInternalDeclined, "Declined")
	addI18NKeyValue(ApprovalStatusInternalApproved, "Approved")
	addI18NKeyValue(ApprovalStatusInternalAborted, "Aborted")
	addI18NKeyValue(ApprovalStatusInternalGenerationFailed, "Generation Failed")
	addI18NKeyValue(ApprovalStatusInternalDeveloperApproved, "Developer Approved")
	addI18NKeyValue(ApprovalStatusPending, "Pending")
	addI18NKeyValue(ApprovalStatusDeclined, "Declined")
	addI18NKeyValue(ApprovalStatusApproved, "Approved")
	addI18NKeyValue(ApprovalStatusAborted, "Aborted")
	addI18NKeyValue(ApprovalStatusGenerationFailed, "Generation Failed")
	addI18NKeyValue(TaskTypeInternalApprovalinfo, "Your own approval request status")
	addI18NKeyValue(TaskTypePlausibilityApprovalinfo, "Your own review request status")
	addI18NKeyValue(TaskTypeInternalApproval, "Approval request")
	addI18NKeyValue(TaskTypePlausibilityApproval, "Review request")
	addI18NKeyValue(FilterSetNotFound, "Filter Set not found")
	addI18NKeyValue(ConnectorReqFailed, "Connector request failure")
	addI18NKeyValue(InvalidExchangeCode, "Invalid auth exchange code")
	addI18NKeyValue(AuthErrorCode, "IAM auth error during login process")
	addI18NKeyValue(StatusReviewUnreviewed, "Open")
	addI18NKeyValue(StatusReviewUnreviewedDE, "Offen")
	addI18NKeyValue(StatusReviewAudited, "Management Approved")
	addI18NKeyValue(StatusReviewAuditedDE, "Management Akzeptiert")
	addI18NKeyValue(StatusReviewAcceptable, "Information")
	addI18NKeyValue(StatusReviewAcceptableDE, "Information")
	addI18NKeyValue(StatusReviewNotAcceptable, "Problem")
	addI18NKeyValue(StatusReviewNotAcceptableDE, "Problem")
	addI18NKeyValue(StatusReviewAcceptableAfterChanges, "Investigation")
	addI18NKeyValue(StatusReviewAcceptableAfterChangesDE, "Untersuchung")
	addI18NKeyValue(ReadLicenseRules, "User is not authorized to view license rules for this project")
	addI18NKeyValue(CreateLicenseRule, "User is not authorized to create license rule for this project")
	addI18NKeyValue(EditLicenseRule, "User is not authorized to edit license rule for this project")
	addI18NKeyValue(ParamLicenseRuleUuidEmpty, "Url parameter 'licenseRuleId' missing")
	addI18NKeyValue(SbomMaxPUrlSite, "One or more PURL entries exceed the maximum length (%d).")
	addI18NKeyValue(SbomMaxCopyrightTextSize, "One or more copyright texts exceed the maximum length (%d).<br>")
	addI18NKeyValue(SbomMaxLicensesPerComponent, "One or more components exceed the maximum number of licenses per component (%d).<br>")
	addI18NKeyValue(SbomMaxComponents, "SBOM exceeds the maximum number of components (%d).")
	addI18NKeyValue(SbomComponentsError, "SBOM validation failed: %s")
	addI18NKeyValue(UnknownPattern, "The requested url is unknown.")
	addI18NKeyValue(SpdxFileContentTypeValidation, "The content type of the SPDX could not be checked.")
	addI18NKeyValue(SpdxAlreadyLocked, "SPDX is already locked.")
	addI18NKeyValue(SpdxNotLocked, "SPDX is not locked.")
	addI18NKeyValue(SpdxRetainedForApprovalOrReview, "SPDX cannot be unlocked because it is retained for approval, review, FOSS DD generation, or status review.")
	addI18NKeyValue(ProjectGroupRequired, "Project group is required for this operation.")
	addI18NKeyValue(ErrorJsonValidatingForm, "Form Data is not valid.")
	addI18NKeyValue(PolicyDecisionOperationNotAuthorized, "User is not authorized to perform operations on policy rule decision for this project")
	addI18NKeyValue(InvalidPolicyDecisionData, "Invalid policy decision data. Allowed decisions are: ALLOW, DENY")
	addI18NKeyValue(InvalidPolicyDecisionLicenseData, "Invalid policy decision data: Invalid license ID")
	addI18NKeyValue(InvalidPolicyDecisionLicenseApprovalState, "Invalid policy decision data: Forbidden license not allowed")
	addI18NKeyValue(InvalidBulkPolicyDecisionData, "Invalid bulk policy decision data.")
	addI18NKeyValue(ParamPolicyDecisionUuidEmpty, "Url parameter 'policyDecisionId' missing")
}

func GetI18N(key string, a ...any) I18N {
	i18nMutex.RLock()
	defer i18nMutex.RUnlock()
	result := I18N{
		Code: key,
		Text: key,
	}
	i18n := i18nMap[key]
	if len(i18n.Text) > 0 {
		result.Code = i18n.Code
		if len(a) > 0 {
			result.Text = fmt.Sprintf(i18n.Text, a...)
		} else {
			result.Text = i18n.Text
		}
	}
	return result
}

const (
	PDDescription              = "Disclosure Document"
	Archive                    = "Approval Archive"
	FileDescriptionPolicyRules = "Policy Rules"
	FileDescriptionPolicyCheck = "SBOM Policy Check Results"
	ApprovalTaskCreate         = "Approval Task created"
	ApprovalTaskUpdate         = "Approval Task update"
	GroupCreated               = "Group created"
	ProjectCreated             = "Project created"
	ProjectUpdated             = "Project updated"
	PolicyLabelsUpdated        = "Policy labels updated"
	VersionCreated             = "Version created"
	VersionUpdated             = "Version updated"
	VersionDeleted             = "Version deleted"
	SourceCodeResourceCreated  = "Version updated: Source code resource created"
	SourceCodeResourceUpdated  = "Version updated: Source code resource updated"
	SourceCodeResourceDeleted  = "Version updated: Source code resource deleted"
	ProjectVersionCreated      = "Project updated: Version created"
	ProjectVersionDeleted      = "Project updated: Version deleted"
	ProjectUserCreated         = "Project updated: User created"
	ProjectUserUpdated         = "Project updated: User updated"
	ProjectUserDeleted         = "Project updated: User deleted"
	ProjectTokenCreated        = "Project updated: Token created"
	ProjectTokenUpdated        = "Project updated: Token renewed"
	ProjectTokenDeleted        = "Project updated: Token revoked"
	FossOfficeCommentApproved  = "FOSS Office confirmation completed"
	FossOfficeCommentInvest    = "Please check and provide clarifications on review remarks and scan remarks"
	LicenseCreated             = "License created"
	LicenseUpdated             = "License updated"
	LicenseDeleted             = "License deleted"
	SpdxDatabaseRefreshed      = "SPDX Database refreshed"
	PolicyRulesCreated         = "Policy Rules created"
	PolicyRulesCopied          = "Policy Rules copied"
	PolicyRulesUpdated         = "Policy Rules updated"
	UserCreated                = "User created"
	UserUpdated                = "User updated"
	InternalApprovalCreated    = "Internal Approval created"
	ExternalApprovalCreated    = "External Approval created"
	ReviewCreated              = "Review created"
	InternalApprovalUpdated    = "Internal Approval updated"
	InternalApprovalAborted    = "Internal Approval aborted"
	ReviewUpdated              = "Review updated"
	ReviewAborted              = "Review aborted"
	ClassificationCreated      = "Classification Created"
	ClassificationUpdated      = "Classification Updated"
	ClassificationDeleted      = "Classification Deleted"
	ExternalApprovalUpdated    = "External approval updated"
	ReviewRemarkCreated        = "Review Remark created"
	ReviewRemarkChanged        = "Review Remark changed"
	ReviewRemarkCommented      = "Review Remark commented"
	OverallReviewUpdated       = "Status Review updated"
	SpdxFileDeleted            = "SBOM file deleted"
	SpdxFileUploaded           = "SBOM file uploaded"
	SpdxFileLocked             = "SBOM file locked"
	SpdxFileUnlocked           = "SBOM file unlocked"
	LicenseRuleCreated         = "License Rule created"
	LicenseRuleUpdated         = "License Rule updated"
	BasicAuthExpiry            = "Expiry date cannot be more than 7 days from now"
	ErrorJsonValidatingForm    = "ERROR_JSON_VALIDATING_FORM"
	PolicyDecisionCreated      = "Policy Decision created"
	PolicyDecisionUpdated      = "Policy Decision updated"
)
