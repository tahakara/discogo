package serviceconfigloader

import (
	"encoding/json"
	"os"
	"sync"
	"time"

	"github.com/tahakara/discogo/internal/logger"
)

type ServiceDef struct {
	Name        string `json:"name"`
	Short       string `json:"short"`
	Description string `json:"description"`
}

type ServiceTypeGroup struct {
	Name     string       `json:"name"`
	Services []ServiceDef `json:"services"`
}

type ServiceTypesConfig map[string]ServiceTypeGroup

var (
	serviceTypes     ServiceTypesConfig
	serviceTypesOnce sync.Once
	loadErr          error
)

func LoadConfig() error {
	if err := loadServiceTypes(); err != nil {
		return err
	}
	return nil
}

// conf.json dosyasından service_types'ı yükler
func loadServiceTypes() error {
	startTime := time.Now()
	serviceTypesOnce.Do(func() {
		f, err := os.Open("conf.json")
		if err != nil {
			loadErr = err
			return
		}
		defer f.Close()

		var raw struct {
			ServiceTypes map[string]ServiceTypeGroup `json:"service_types"`
		}
		if err := json.NewDecoder(f).Decode(&raw); err != nil {
			loadErr = err
			return
		}

		serviceTypes = raw.ServiceTypes
		if b, err := json.Marshal(serviceTypes); err == nil {
			logger.Info(string(b), time.Since(startTime))
		} else {
			logger.Info("Failed to marshal serviceTypes: ", time.Since(startTime))
		}
	})
	return loadErr
}

// shortname ile tam isim ve açıklama döndürür (tüm gruplardaki servislerde arar)
func GetServiceTypeInfo(short string) (string, string, bool) {
	if err := loadServiceTypes(); err != nil {
		return "", "", false
	}
	for _, group := range serviceTypes {
		for _, svc := range group.Services {
			if svc.Short == short {
				return svc.Name, svc.Description, true
			}
		}
	}
	return "", "", false
}

// ServiceType'ın geçerliliğini kontrol eden fonksiyon
func IsValidServiceType(short string) bool {
	for _, group := range serviceTypes {
		for _, svc := range group.Services {
			if svc.Short == short {
				return true
			}
		}
	}
	return false
}

// Tüm grupları ve servisleri döndürür
func GetAllServiceTypes() []string {
	var all []string
	for _, group := range serviceTypes {
		for _, svc := range group.Services {
			all = append(all, svc.Short)
		}
	}
	return all
}

// Tüm servislerin kısa isimlerini döndürür
func GetAllServiceShortNames() []string {
	if err := loadServiceTypes(); err != nil {
		return nil
	}
	var shorts []string
	for _, group := range serviceTypes {
		for _, svc := range group.Services {
			shorts = append(shorts, svc.Short)
		}
	}
	return shorts
}

// Bir grup adı ile o gruptaki tüm servisleri döndürür
func GetServicesByGroup(groupName string) []ServiceDef {
	group, ok := serviceTypes[groupName]
	if !ok {
		return nil
	}
	return group.Services
}

func GetAllProviders() []string {
	var all []string
	for _, provider := range providers {
		all = append(all, provider.Short)
	}
	return all
}

type ProviderDef struct {
	Name  string `json:"name"`
	Short string `json:"short"`
}

var (
	providers     map[string]ProviderDef
	providersOnce sync.Once
	providerErr   error
)

// conf.json dosyasından providers'ı yükler
func loadProviders() error {
	providersOnce.Do(func() {
		f, err := os.Open("conf.json")
		if err != nil {
			providerErr = err
			return
		}
		defer f.Close()

		var raw struct {
			Providers map[string]ProviderDef `json:"providers"`
		}
		if err := json.NewDecoder(f).Decode(&raw); err != nil {
			providerErr = err
			return
		}
		providers = raw.Providers
	})
	return providerErr
}

// Provider'ın geçerliliğini kontrol eden fonksiyon
func IsValidProvider(short string) bool {
	_, ok := providers[short]
	return ok
}

// Ortak yükleme fonksiyonu
func LoadAllConfigs() error {
	var err1, err2 error
	err1 = loadServiceTypes()
	err2 = loadProviders()
	if err1 != nil {
		return err1
	}
	if err2 != nil {
		return err2
	}
	return nil
}
