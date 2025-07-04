package runtime

type Environment string

const (
	Local Environment = "local"
	Dev   Environment = "dev"
	Sit   Environment = "sit"
	Stg   Environment = "stg"
	Prd   Environment = "prd"
)

func ValidateProfile(profile string) Environment {
	switch profile {
	case "local", "dev", "sit", "stg", "prd":
		return Environment(profile)
	default:
		return "local"
	}
}

type RuntimeCfg struct {
	Microservice string
	Env          Environment
}
