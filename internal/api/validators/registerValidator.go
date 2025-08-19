package validators

import (
	"fmt"
	"regexp"

	"github.com/go-playground/validator/v10"
	requestDTOs "github.com/tahakara/discogo/internal/api/dtos/requestdto"
	serviceconfigloader "github.com/tahakara/discogo/internal/config/serviceconfiguration"
)

var validate *validator.Validate

func init() {
	validate = validator.New()
	// Register custom version validation: dot-separated numbers, e.g. 1.0.0
	validate.RegisterValidation("version", func(fl validator.FieldLevel) bool {
		version := fl.Field().String()
		// Accepts versions like 1.0, 1.0.0, 2.3.4.5 etc.
		matched, _ := regexp.MatchString(`^\d+(\.\d+)*$`, version)
		return matched
	})

	validate.RegisterValidation("type", func(fl validator.FieldLevel) bool {
		typeStr := fl.Field().String()
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, typeStr)
		if matched {
			return serviceconfigloader.IsValidServiceType(typeStr)
		}
		return false
	})

	validate.RegisterValidation("provider", func(fl validator.FieldLevel) bool {
		providerStr := fl.Field().String()
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, providerStr)
		if matched {
			return serviceconfigloader.IsValidProvider(providerStr)
		}
		return false
	})

	validate.RegisterValidation("alphanumanddashandunderscore", func(fl validator.FieldLevel) bool {
		str := fl.Field().String()
		matched, _ := regexp.MatchString(`^[a-zA-Z0-9_-]+$`, str)
		return matched
	})

}

func validationErrorToText(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s alanı zorunludur.", fe.Field())
	case "min":
		return fmt.Sprintf("%s alanı en az %s karakter olmalıdır.", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s alanı en fazla %s karakter olmalıdır.", fe.Field(), fe.Param())
	case "type":
		return "Geçersiz servis tipi."
	case "version":
		return "Geçersiz versiyon formatı."
	case "provider":
		return "Geçersiz provider."
	case "ip4_addr":
		return fmt.Sprintf("%s alanı geçerli bir IPv4 adresi olmalıdır.", fe.Field())
	case "ip6_addr":
		return fmt.Sprintf("%s alanı geçerli bir IPv6 adresi olmalıdır.", fe.Field())
	// Diğer tag'ler için de ekleyebilirsin
	default:
		return fmt.Sprintf("%s alanı için geçersiz değer.", fe.Field())
	}
}

// ValidateRegisterRequest validates a RegisterRequest instance.
func ValidateRegisterRequest(req *requestDTOs.RegisterRequestDTO) []string {
	// Önce field bazlı validasyonları uygula
	if err := validate.Struct(req); err != nil {

		var errors []string
		for _, err := range err.(validator.ValidationErrors) {
			errors = append(errors, validationErrorToText(err))
		}
		return errors
	}

	// Sonra özel cross-field validasyonunu uygula
	addr4ok := req.Addr4 != "" && req.Port4 > 0
	addr6ok := req.Addr6 != "" && req.Port6 > 0

	if !addr4ok && !addr6ok {
		return []string{"either (Addr4 and Port4) or (Addr6 and Port6) must be provided"}
	}

	return nil
}
