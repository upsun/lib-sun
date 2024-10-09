package entity

import (
	"maps"

	"gopkg.in/yaml.v3"
)

// Constants files names of configuration.
const (
	PSH_ROUTE       string = "routes.yaml"        // PSH routes yaml file.
	PSH_SERVICE     string = "services.yaml"      // PSH services yaml file.
	PSH_PLATFORM    string = ".platform.app.yaml" // PSH application yaml file.
	PSH_APPLICATION string = "applications.yaml"  // PSH application yaml file. (variante config)
	UPSUN_CONFIG    string = "config.yaml"        // Upsun configuration yaml file.
)

// type Provider string

// const (
// 	Platform Provider = "platform"
// 	Upsun    Provider = "upsun"
// )

const PSH_PROVIDER = "platform"
const UPS_PROVIDER = "upsun"

// Application arguments.
type CliCommonArgs struct {
	Verbose bool
	Silent  bool
	PathLog string
}

type CliConvertArgs struct {
	ProjectSource      string
	ProjectDestination string
	TypeMount          string
}

type CliCloneArgs struct {
	SrcProvider    string
	SrcProjectID   string
	SrcEnvironment string
	DstProvider    string
	DstProjectID   string
	DstEnvironment string
	DstRegion      string
	DstOrga        string
	NoData         bool
	OnlyData       bool
	OnlyDb         bool
	OnlyMount      bool
	KeepData       string
	NoUsers        bool
	PshRepo        bool
	NoLocal        bool
	Sensitive      bool
}

type CliScalingArgs struct {
	Name            string
	IncludeServices bool
	HostCountMin    int
	HostCountMax    int
	CpuUsageMin     float64
	CpuUsageMax     float64
	MemUsageMin     float64
	MemUsageMax     float64
}

// Project context content.
//
//   - Provider: link to CLI name [platfom, upsun, ibexa...].
//   - ProjectID: ID of project.
//   - Environement: Environement name.
type ProjectGlobal struct {
	Provider     string
	Region       string                     `json:"region"`
	ID           string                     `json:"id"`
	Name         string                     `json:"title"`
	DefaultEnv   string                     `json:"-"`
	OrgEmail     string                     `json:"-"`
	Timezone     string                     `json:"timezone"`
	Services     map[string]EnvService      `json:"-"`
	Mounts       map[string]EnvMount        `json:"-"`
	Variables    map[string]ProjectVariable `json:"-"`
	VariablesEnv map[string]ProjectVariable `json:"-"`
	Users        map[string]ProjectUser     `json:"-"`
	Access       map[string]ProjectAccess   `json:"-"`
	//Domains      map[string]ProjectDomain
	//DeployKey    string
	//Integrations []string
	Description string `json:"description"`
	Repository  struct {
		Url string `json:"url"`
		Ssh string `json:"client_ssh_key"`
	} `json:"repository"`
	DefaultDom string `json:"default_domain"`
}

func MakeProjectContext(provider string, id string, defaultEnv string) ProjectGlobal {
	var result = ProjectGlobal{
		Provider:   provider,
		ID:         id,
		DefaultEnv: defaultEnv,
	}
	result.Variables = make(map[string]ProjectVariable)
	result.VariablesEnv = make(map[string]ProjectVariable)
	result.Users = make(map[string]ProjectUser)
	result.Access = make(map[string]ProjectAccess)
	result.Services = make(map[string]EnvService)
	result.Mounts = make(map[string]EnvMount)
	return result
}

func (dst *ProjectGlobal) CopyProjectBase(src ProjectGlobal) {
	if dst.Name == "" {
		dst.Name = src.Name
	}
	if dst.Description == "" {
		dst.Description = src.Description
	}
	if dst.DefaultEnv == "" {
		dst.DefaultEnv = src.DefaultEnv
	}
	if dst.Timezone == "" {
		dst.Timezone = src.Timezone
	}
	if dst.Region == "" {
		dst.Region = src.Region
	}
	if dst.DefaultDom == "" {
		dst.DefaultDom = src.DefaultDom
	}
}

func (dst *ProjectGlobal) CopyVariables(src ProjectGlobal) {
	maps.Copy(dst.Variables, src.Variables)
	maps.Copy(dst.VariablesEnv, src.VariablesEnv)
}

func (dst *ProjectGlobal) CopyUsers(src ProjectGlobal) {
	maps.Copy(dst.Users, src.Users)
}

func (dst *ProjectGlobal) CopyAccess(src ProjectGlobal) {
	maps.Copy(dst.Access, src.Access)
}

func (dst *ProjectGlobal) CopyServices(src ProjectGlobal) {
	maps.Copy(dst.Services, src.Services)
}

func (dst *ProjectGlobal) CopyMounts(src ProjectGlobal) {
	maps.Copy(dst.Mounts, src.Mounts)
}

func (dst *ProjectGlobal) Copy(src ProjectGlobal) {
	dst.CopyProjectBase(src) // Base on Src
	dst.CopyVariables(src)   // Sync Variables
	dst.CopyUsers(src)       // Sync Users
	dst.CopyAccess(src)      // Sync Access
	dst.CopyServices(src)    // Sync Services
	dst.CopyMounts(src)      // Sync Mounts
}

type ProjectEnv struct {
	IsSendEmail     bool
	IsCrawlerHidden bool
	IsHttpAccess    bool
	//Access       map[string]string //index is login and value is password
	//Ips          string
	IsActive  bool
	IsPaused  bool
	Variables map[string]ProjectVariable
	//Domains      map[string]ProjectDomain
}

type ProjectVariable struct {
	ID            string `json:"id,omitempty"`
	Name          string `json:"name"`
	Value         string `json:"value,omitempty"`
	IsAtBuild     bool   `json:"visible_build"`
	IsAtDeploy    bool   `json:"visible_runtime"`
	IsJson        bool   `json:"is_json"`
	IsSensitive   bool   `json:"is_sensitive"`
	IsInheritable *bool  `json:"is_inheritable,omitempty"`
	// Attributes    map[string]interface{} `json:"attributes"`
}

type ProjectUser struct {
	Id            string `json:"id,omitempty"`
	Desactivation bool   `json:"-"`
	Namespace     string `json:"namespace"`
	Email         string `json:"email"`
	//DisplayName    string `json:"-"`
	//Role           []string
	//IsProjectAdmin bool
}

type ProjectAccess struct {
	UserId      string   `json:"user_id,omitempty"`
	Permissions []string `json:"permissions"`
	AddAuto     *bool    `json:"auto_add_member,omitempty"`
}

type ProjectDomain struct {
	Domain       string
	EnvLinked    string
	AttachedProd string
}

type EnvService struct {
	Type         string `json:"type"`
	Relationship string `json:"-"`
	Application  string `json:"-"`
	TypeService  string `json:"-"`
	DumpPath     string `json:"-"`
}

type EnvMount struct {
	Path        string `json:"-"`
	Type        string `json:"source"`
	SourcePath  string `json:"source_path"`
	Application string `json:"-"`
	DumpPath    string `json:"-"`
}

type EnvVariable struct {
}

// MetaModel (configuration) of project.
//
//   - Services...
//   - Applications...
//   - Routes...
type MetaConfig struct {
	Applications yaml.Node `yaml:"applications"`
	Services     yaml.Node `yaml:"services"`
	Routes       yaml.Node `yaml:"routes"`
}

type ProvisionApplication struct {
	Name       string `yaml:"-"`
	Mainstream struct {
		URL     string `yaml:"url"`
		Shell   string `yaml:"shell"`
		Version string `yaml:"version"`
	} `yaml:"mainstream"`
	Size struct {
		CPU string `yaml:"cpu"`
		MEM string `yaml:"mem"`
	} `yaml:"size"`
	Files    map[string]string `yaml:"files"`
	Services map[string]string `yaml:"services"`
	Mounts   map[string]string `yaml:"mounts"`
	Hook     string            `yaml:"hook"`
	HookPost string            `yaml:"post-hook"`
}

type ProvisionGlobal struct {
	Name string `yaml:"name"`
	// AppName     string                         `yaml:"app_name"`
	Description  string                          `yaml:"description"`
	Applications map[string]ProvisionApplication `yaml:"applications"`
	Size         struct {
		Plan string `yaml:"plan"`
	} `yaml:"size"`
	Variables map[string]string `yaml:"variables"`
	Users     map[string]string `yaml:"users"`
}
