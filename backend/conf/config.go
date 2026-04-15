// SPDX-FileCopyrightText: 2025 Mercedes-Benz Group AG and Mercedes-Benz AG
//
// SPDX-License-Identifier: Apache-2.0

package conf

import (
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"

	"golang.org/x/text/language"

	"github.com/jinzhu/configor"
)

const (
	DocumentResources = "/resources/document/"
	HeaderFile        = "page-header-%s.html"
	ContentFile       = "content-%s.tmpl.html"
	FooterFile        = "page-footer-%s.html"
	Resources         = "/resources/"
)

type DatabaseType string

const (
	DatabaseArangoDB DatabaseType = "ArangoDB"
	DatabaseCouchDB  DatabaseType = "CouchDB"
	DatabaseMongoDB  DatabaseType = "MongoDB"
)

var ConfigFiles = []string{"./conf/config-local.yml", "./conf/config.yml"}

type Server struct {
	Port string `default:"3333"`
	Tls  bool   `default:"false" env:"SERVER_TLS"`
	// the grace period corresponds to terminationGracePeriodSeconds in the deployment.yaml
	// change the value in the deployment.yaml if you change this value
	TerminationGracePeriodSeconds int    `default:"300"`
	ClientRedirectURL             string `default:"https://localhost:3000" env:"CLIENT_REDIRECT_URL"`
	ApplicationToken              string
	Uploadpath                    string `default:"/uploads/" env:"SERVER_UPLOADPATH"`
	BackupPath                    string `default:"/backups"`
	MaxUploadPerHourPerProject    int    `default:"10"`
	Env                           string `default:"local"`
	ProdEnv                       string `default:"prod"`
	EnableFixDataIntegrity        bool   `default:"false"`
	UploadMaxMb                   int64  `default:"200"`
	SbomValidationEnabled         bool   `default:"true"`
	Projectworkpath               string `default:"/srv/disuko/workdir/projects"`
	DisallowedUploadFilenameChars string `default:"<>\\/*:?|\"'"`
	AllowedOrigins                string
	LocalIp                       string // do not set this value!!!
	DevLog                        bool   `default:"false"`
	DevLogDbQueries               bool   `default:"false"`
	TestsWithoutDocker            bool   `default:"false"`
	E2ETests                      bool   `default:"false" env:"E2E_TESTS"`
	VanillaDisuko                 bool   `default:"false" env:"VANILLA_DISUKO"`
	// base path of the disclosure document html templates
	BasePath   string
	DisukoHost string `required:"true" default:"https://localhost:3333" env:"DISUKO_HOST"`
	// used by disclosure pdf
	SBomComponentsPath             string `required:"true" env:"SBOM_COMPONENTS_PATH"`
	SBomLimits                     SBomLimits
	TermsOfUseCurrentVersion       string `required:"true" env:"TERMS_OF_USE_CURRENT_VERSION"`
	MaxVersions                    int    `default:"10"`
	AutoDeleteSbomsAfterUpload     bool   `default:"false" env:"AUTO_DELETE_SBOMS"`
	ProdAutoDeleteSbomsAfterUpload bool   `default:"false" env:"PROD_AUTO_DELETE_SBOMS"`
	InstanceName                   string `default:"standalone"`
	EnforceFOSSOfficeConfirmation  bool
	FOSSOfficeUserId               string
}

type SBomLimits struct {
	MaxComponents           int `default:"15000"`
	MaxLicensesPerComponent int `default:"50"`
	MaxCopyRightTextSize    int `default:"200000"`
	MaxPURLSize             int `default:"1000"`
}

func (s Server) GetResourceFilePath(fileName string) string {
	return filepath.Join(s.BasePath+Resources, fileName)
}

func (s Server) GetDocumentResourceFilePath(fileName string) string {
	return filepath.Join(s.BasePath, DocumentResources, fileName)
}

func (s Server) GetResourcePath(fileName string) string {
	return s.BasePath + DocumentResources + fileName
}

func (s Server) TemplatePageHeaderPath(template string, lang language.Tag) string {
	return filepath.Join(s.BasePath+DocumentResources, template, fmt.Sprintf(HeaderFile, lang.String()))
}

func (s Server) TemplatePageFooterPath(template string, lang language.Tag) string {
	return filepath.Join(s.BasePath+DocumentResources, template, fmt.Sprintf(FooterFile, lang.String()))
}

func (s Server) TemplateContentName(lang language.Tag) string {
	return fmt.Sprintf(ContentFile, lang.String())
}

func (s Server) TemplateCSSGlob(template string) string {
	return filepath.Join(s.BasePath+DocumentResources, template, "*.css")
}

func (s Server) TemplateGlob(template string) string {
	return filepath.Join(s.BasePath+DocumentResources, template, "*.tmpl.html")
}

func (s Server) GetSBomComponentsLink(project, version, sbom string) string {
	path := strings.Replace(strings.Replace(strings.Replace(s.SBomComponentsPath, "{version}", version, 1), "{project}", project, 1), "{sbom}", sbom, 1)
	return s.DisukoHost + path
}

func (s Server) GetProjectApprovalLink(projectUuid string, isGroup bool) string {
	projectPath := "/#/dashboard/projects/{project}/approvals"
	groupPath := "/#/dashboard/groups/{project}/approvals"
	var path string
	if isGroup {
		path = strings.Replace(groupPath, "{project}", projectUuid, 1)
	} else {
		path = strings.Replace(projectPath, "{project}", projectUuid, 1)
	}
	return s.DisukoHost + path
}

func (s Server) GetUploadPath() string {
	return s.Uploadpath
}

var Config = struct {
	OAuth2 struct {
		ClientId              string `default:""`
		Secret                string
		RedirectURL           string
		InsecureProvider      string
		Provider              string
		AuthorizationEndpoint string
		LogoutEndpoint        string
		// flag is needed, if a http request from outside should use the bearer
		PreventTokenHiJacking bool   `default:"true"`
		DebugLog              bool   `default:"false"`
		UppercaseUsername     bool   `default:"true"`
		RegexToken            string `default:"^[a-zA-Z0-9_\\-]{10,50}$"`
	}
	InternalUsersAllowList              []string
	DeprovisioningInactiveDaysThreshold int `default:"60"`
	PublicAuth                          struct {
		AccessTTLSeconds  int
		RefreshTTLMinutes int
		SigningKey        string
	}
	Auth struct {
		AccessSecret                string
		RefreshSecret               string
		AccessTokenExpiresInMinutes int `default:"30"`
		RefreshTokenExpiresInHours  int `default:"12"`
	}
	Server   Server
	Database struct {
		Host               string `default:""`
		Scheme             string `default:"http"`
		InsecureSkipVerify bool   `default:"true"`
		Port               int    `default:""`
		User               string `default:""`
		Password           string
		DatabaseName       string       `default:"disuko"`
		MigrateOnly        bool         `default:"false"`
		Type               DatabaseType `default:"CouchDB"`
		ShardReplica       int          `default:"3"`
	}
	SPDXLicense struct {
		LicensesInfoPath      string `default:"https://raw.githubusercontent.com/spdx/license-list-data/master/json/licenses.json"`
		LicenseDetailsDirPath string `default:"https://raw.githubusercontent.com/spdx/license-list-data/master/json/details/"`
	}
	Proxy struct {
		HttpProxy string `default:""`
	}
	S3 struct {
		IsEnabled          bool `default:"false"`
		IsRestApiEnabled   bool `default:"false"`
		AwsAccessKeyId     string
		AwsSecretAccessKey string
		AwsRegion          string `default:""`
		AwsEndPoint        string
		BucketName         string
	}
	Fonts struct {
		Stamp string `default:"Roboto-Regular.ttf"`
	}
	Smtp struct {
		Host   string `default:""`
		Port   string `default:""`
		Sender string `default:""`
		User   string
		Pass   string
	}
	Cache struct {
		Host     string `default:""`
		Port     int    `default:""`
		Password string
		Channel  string `default:"disuko"`
	}
	Connector struct {
		Userrole struct {
			Host   string
			Scheme string
			Port   int
		}
		Application struct {
			Host   string
			Scheme string
			Port   int
		}
		Department struct {
			Host   string
			Scheme string
			Port   int
		}
	}
	BasicAuth struct {
		Pepper string
	}
}{}

func init() {
	log.Printf("config::LoadConfig: %v", ConfigFiles)

	err := configor.Load(&Config, ConfigFiles...)
	if err != nil {
		log.Printf("Could not load config file: %v %v", err, ConfigFiles)
	}

	// set local server base path
	setServerBasePath(err)

	checkEnvironmentVariables()

	// set local ip
	Config.Server.LocalIp = GetLocalIP().String()

	logConfiguration()
}

func setServerBasePath(err error) {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Could not extract server base path, err= %s", err)
		return
	}
	Config.Server.BasePath = dir
}

func checkEnvironmentVariables() {
	Config.Database.Host = getEnvVariable("DATABASE_HOST", Config.Database.Host)
	Config.Database.Port = getEnvVariableInt("DATABASE_PORT", Config.Database.Port)
	Config.Database.Scheme = getEnvVariable("DATABASE_SCHEME", Config.Database.Scheme)
	Config.Database.User = getEnvVariable("DATABASE_USER", Config.Database.User)
	Config.Database.Password = getEnvVariable("DATABASE_PASSWORD", Config.Database.Password)
	Config.Database.Type = DatabaseType(getEnvVariable("DATABASE_TYPE", string(Config.Database.Type)))
	Config.Database.ShardReplica = getEnvVariableInt("DATABASE_SHARD_REPLICA", Config.Database.ShardReplica)

	Config.Server.Port = getEnvVariable("SERVER_PORT", Config.Server.Port)
	Config.Server.ApplicationToken = getEnvVariable("APPLICATION_TOKEN", Config.Server.ApplicationToken)
	Config.Server.TestsWithoutDocker = getEnvVariableBoolean("TESTS_WITHOUT_DOCKER", Config.Server.TestsWithoutDocker)
	Config.Server.DevLogDbQueries = getEnvVariableBoolean("DEV_LOG_DB_QUERIES", Config.Server.DevLogDbQueries)
	Config.Server.InstanceName = getEnvVariable("INSTANCE_NAME", Config.Server.InstanceName)

	// set default AllowedOrigins ... this value is not possible to set direct with the default annotation
	if len(Config.Server.AllowedOrigins) == 0 {
		// Config.Server.AllowedOrigins = "*,http://*,https://*,https://localhost:3000"
		Config.Server.AllowedOrigins = "https://localhost:3000"
	}
	Config.Server.AllowedOrigins = getEnvVariable("SERVER_ALLOWED_ORIGINS", Config.Server.AllowedOrigins)

	Config.OAuth2.ClientId = getEnvVariable("OAUTH2_CLIENTID", Config.OAuth2.ClientId)
	Config.OAuth2.Secret = getEnvVariable("OAUTH2_SECRET", Config.OAuth2.Secret)
	Config.OAuth2.Provider = getEnvVariable("OAUTH2_PROVIDER", Config.OAuth2.Provider)
	Config.OAuth2.InsecureProvider = getEnvVariable("OAUTH2_INSECURE_PROVIDER", Config.OAuth2.InsecureProvider)
	Config.OAuth2.RedirectURL = getEnvVariable("OAUTH2_REDIRECTURL", Config.OAuth2.RedirectURL)
	Config.OAuth2.AuthorizationEndpoint = getEnvVariable("OAUTH2_AUTHORIZATION_ENDPOINT", Config.OAuth2.AuthorizationEndpoint)
	Config.OAuth2.LogoutEndpoint = getEnvVariable("OAUTH2_LOGOUT_ENDPOINT", Config.OAuth2.LogoutEndpoint)
	Config.OAuth2.PreventTokenHiJacking = getEnvVariableBoolean("OAUTH2_PREVENT_TOKEN_HIJACKING", Config.OAuth2.PreventTokenHiJacking)
	Config.OAuth2.DebugLog = getEnvVariableBoolean("OAUTH2_DEBUG_LOG", Config.OAuth2.DebugLog)
	Config.OAuth2.UppercaseUsername = getEnvVariableBoolean("OAUTH2_UPPERCASEUSERNAME", Config.OAuth2.UppercaseUsername)
	Config.OAuth2.RegexToken = getEnvVariable("OAUTH2_REGEXTOKEN", Config.OAuth2.RegexToken)

	Config.DeprovisioningInactiveDaysThreshold = getEnvVariableInt("DEPROVISIONING_INACTIVE_THRESHOLD", Config.DeprovisioningInactiveDaysThreshold)

	Config.InternalUsersAllowList = strings.Split(getEnvVariable("INTERNAL_USERS_ALLOW_LIST", ""), ",")

	Config.PublicAuth.AccessTTLSeconds = getEnvVariableInt("PUBLICAUTH_ACCESS_TTL", Config.PublicAuth.AccessTTLSeconds)
	Config.PublicAuth.RefreshTTLMinutes = getEnvVariableInt("PUBLICAUTH_REFRESH_TTL", Config.PublicAuth.RefreshTTLMinutes)
	Config.PublicAuth.SigningKey = getEnvVariable("PUBLICAUTH_SIGNING_KEY", Config.PublicAuth.SigningKey)

	Config.Auth.AccessSecret = getEnvVariable("AUTH_ACCESS_SECRET", Config.Auth.AccessSecret)
	Config.Auth.RefreshSecret = getEnvVariable("AUTH_REFRESH_SECRET", Config.Auth.RefreshSecret)

	Config.S3.IsEnabled = getEnvVariableBoolean("S3_ENABLED", Config.S3.IsEnabled)
	Config.S3.IsRestApiEnabled = getEnvVariableBoolean("S3_REST_API_ENABLED", Config.S3.IsRestApiEnabled)
	Config.S3.AwsRegion = getEnvVariable("AWS_REGION", Config.S3.AwsRegion)
	Config.S3.AwsEndPoint = getEnvVariable("AWS_ENDPOINT", Config.S3.AwsEndPoint)
	Config.S3.AwsAccessKeyId = getEnvVariable("AWS_ACCESS_KEY_ID", Config.S3.AwsAccessKeyId)
	Config.S3.AwsSecretAccessKey = getEnvVariable("AWS_SECRET_ACCESS_KEY", Config.S3.AwsSecretAccessKey)
	Config.S3.BucketName = getEnvVariable("S3_BUCKET_NAME", Config.S3.BucketName)

	Config.Server.MaxUploadPerHourPerProject = getEnvVariableInt("MAX_UPLOAD_PER_HOUR_PER_PROJECT", Config.Server.MaxUploadPerHourPerProject)
	Config.Server.Env = getEnvVariable("ENV", Config.Server.Env)
	Config.Server.ProdEnv = getEnvVariable("PROD_ENV", Config.Server.ProdEnv)
	Config.Server.BackupPath = getEnvVariable("BACKUP_PATH", Config.Server.BackupPath)
	Config.Server.EnableFixDataIntegrity = getEnvVariableBoolean("ENABLE_FIX_DATA_INTEGRITY", Config.Server.EnableFixDataIntegrity)

	Config.Cache.Host = getEnvVariable("CACHE_HOST", Config.Cache.Host)
	Config.Cache.Port = getEnvVariableInt("CACHE_PORT", Config.Cache.Port)
	Config.Cache.Password = getEnvVariable("CACHE_PASSWORD", Config.Cache.Password)

	Config.Database.MigrateOnly = getEnvVariableBoolean("MIGRATE_ONLY", Config.Database.MigrateOnly)
	Config.Server.AutoDeleteSbomsAfterUpload = getEnvVariableBoolean("AUTO_DELETE_SBOMS", Config.Server.AutoDeleteSbomsAfterUpload)
	Config.Server.ProdAutoDeleteSbomsAfterUpload = getEnvVariableBoolean("PROD_AUTO_DELETE_SBOMS", Config.Server.ProdAutoDeleteSbomsAfterUpload)

	Config.Server.ClientRedirectURL = getEnvVariable("CLIENT_REDIRECT_URL", Config.Server.ClientRedirectURL)
	Config.Server.Tls = getEnvVariableBoolean("SERVER_TLS", Config.Server.Tls)
	Config.Server.EnforceFOSSOfficeConfirmation = getEnvVariableBoolean("ENFORCE_FOSS_OFFICE_CONFIRMATION", Config.Server.EnforceFOSSOfficeConfirmation)
	Config.Server.FOSSOfficeUserId = getEnvVariable("FOSS_OFFICE_USER_ID", Config.Server.FOSSOfficeUserId)

	Config.Connector.Userrole.Scheme = getEnvVariable("USERROLE_SCHEME", Config.Connector.Userrole.Scheme)
	Config.Connector.Userrole.Host = getEnvVariable("USERROLE_HOST", Config.Connector.Userrole.Host)
	Config.Connector.Userrole.Port = getEnvVariableInt("USERROLE_PORT", Config.Connector.Userrole.Port)

	Config.Connector.Application.Scheme = getEnvVariable("APPLICATION_SCHEME", Config.Connector.Application.Scheme)
	Config.Connector.Application.Host = getEnvVariable("APPLICATION_HOST", Config.Connector.Application.Host)
	Config.Connector.Application.Port = getEnvVariableInt("APPLICATION_PORT", Config.Connector.Application.Port)

	Config.Connector.Department.Scheme = getEnvVariable("DEPARTMENT_SCHEME", Config.Connector.Department.Scheme)
	Config.Connector.Department.Host = getEnvVariable("DEPARTMENT_HOST", Config.Connector.Department.Host)
	Config.Connector.Department.Port = getEnvVariableInt("DEPARTMENT_PORT", Config.Connector.Department.Port)

	Config.BasicAuth.Pepper = getEnvVariable("BASICAUTH_PEPPER", Config.BasicAuth.Pepper)

	Config.Smtp.Host = getEnvVariable("SMTP_HOST", Config.Smtp.Host)
	Config.Smtp.Port = getEnvVariable("SMTP_PORT", Config.Smtp.Port)
	Config.Smtp.Sender = getEnvVariable("SMTP_SENDER", Config.Smtp.Sender)
	Config.Smtp.User = getEnvVariable("SMTP_USER", Config.Smtp.User)
	Config.Smtp.Pass = getEnvVariable("SMTP_PASS", Config.Smtp.Pass)
}

func rightPad2Len(s string, padStr string, overallLen int) string {
	var padCountInt int
	padCountInt = 1 + ((overallLen - len(padStr)) / len(padStr))
	retStr := s + strings.Repeat(padStr, padCountInt)
	return retStr[:overallLen]
}

func dumpStructToSystemOut(parentTitle string, data interface{}) interface{} {
	if reflect.ValueOf(data).Kind() == reflect.Struct {
		v := reflect.ValueOf(data)
		typeOfS := v.Type()
		for i := 0; i < v.NumField(); i++ {
			if v.Field(i).Kind() == reflect.Struct {
				dumpStructToSystemOut(typeOfS.Field(i).Name, v.Field(i).Interface())
			} else {
				title := rightPad2Len(strings.ToUpper(parentTitle)+"_"+strings.ToUpper(typeOfS.Field(i).Name), " ", 33)
				valueStr := fmt.Sprintf("%v", v.Field(i).Interface())
				titleLower := strings.ToLower(title)
				if len(valueStr) <= 0 {
					valueStr = "[WARNING is empty]"
				} else if strings.Index(titleLower, "token") > -1 ||
					strings.Index(titleLower, "secret") > -1 ||
					strings.Index(titleLower, "clientid") > -1 ||
					strings.Index(titleLower, "key") > -1 ||
					strings.Index(titleLower, "pass") > -1 {
					if strings.Index(titleLower, "preventtokenhijacking") == -1 {
						valueStr = "***"
					}
				}
				fmt.Printf("\n%s = %v", title, valueStr)

			}
		}
	}

	return data
}

func logConfiguration() {
	fmt.Printf("\n[CONFIGURATION] from %v", ConfigFiles)
	dumpStructToSystemOut("", Config)
	fmt.Printf("\n\n")
}

func getEnvVariable(envKey string, defaultValue string) string {
	val, exists := os.LookupEnv(envKey)
	if exists {
		return val
	}
	return defaultValue
}

func getEnvVariableBoolean(envKey string, defaultValue bool) bool {
	val, exists := os.LookupEnv(envKey)
	if exists {
		ret, _ := strconv.ParseBool(val)
		return ret
	}
	return defaultValue
}

func getEnvVariableInt(envKey string, defaultValue int) int {
	val, exists := os.LookupEnv(envKey)
	if exists {
		ret, _ := strconv.ParseInt(val, 10, 64)
		return int(ret)
	}
	return defaultValue
}

func GetLocalIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func IsProdEnv() bool {
	return Config.Server.Env == Config.Server.ProdEnv
}
