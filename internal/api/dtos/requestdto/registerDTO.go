package requestdto

type RegisterRequestDTO struct {
	Name          string            `validate:"required,min=3,max=64"`
	Type          string            `validate:"required,type"`
	Version       string            `validate:"required,version"`
	Provider      string            `validate:"required,provider"`
	Region        string            `validate:"required,alphanumanddashandunderscore"`
	Zone          string            `validate:"required,alphanumanddashandunderscore"`
	Cluster       string            `validate:"required,alphanumanddashandunderscore"`
	InstanceID    string            `validate:"required,alphanumanddashandunderscore"`
	NetworkID     string            `validate:"required,alphanumanddashandunderscore"`
	SubnetID      string            `validate:"required,alphanumanddashandunderscore"`
	NetworkDomain string            `validate:"required,alphanumanddashandunderscore"`
	Tags          map[string]string `validate:"omitempty"`
	Addr4         string            `validate:"omitempty,ip4_addr"`
	Port4         int               `validate:"omitempty,min=1,max=65535"`
	Addr6         string            `validate:"omitempty,ip6_addr"`
	Port6         int               `validate:"omitempty,min=1,max=65535"`
}
