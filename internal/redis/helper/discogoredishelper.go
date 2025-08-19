package redishelper

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	env "github.com/tahakara/discogo/internal/config"
	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	"github.com/tahakara/discogo/internal/utils"
)

type ServiceEntry struct {
	ServiceUUID   string            // (DISCO) unique service identifier
	Name          string            // xyz-service human-readable name
	Type          string            // type of service (shortname, örn: "gw")
	Version       string            // version of the service
	Provider      string            // aws, gcp, azure, etc.
	Region        string            // region of the service, e.g., us-east-1
	Zone          string            // availability zone, e.g., us-east-1a
	Cluster       string            // cluster name, e.g., xyz-cluster
	InstanceID    string            // unique instance identifier
	NetworkID     string            // network identifier vpc-12345, vnet-12345, etc.
	SubnetID      string            // subnet identifier, e.g., subnet-12345
	NetworkDomain string            // e.g., internal, public, dmz
	Tags          map[string]string // key-value pairs for additional metadata
	Addr4         string            // IPv4 address
	Addr6         string            // IPv6 address
	Port4         int               // IPv4 port
	Port6         int               // IPv6 port

	CreatedAt    string            // (DISCO) RFC3339 Unix timestamp of creation
	LastHeardAt  string            // (DISCO) RFC3339 Unix timestamp of last heartbeat
	Status       ServiceStatus     // (DISCO) e.g., healthy, degraded, offline
	HeardCount   int64             // (DISCO) Count of heartbeats received
	ReportCount  int64             // (DISCO | Client) Count of reports received
	LastReportAt string            // (DISCO | Client) RFC3339 Unix timestamp of last report
	Metadata     map[string]string // (DISCO | Client) Additional metadata

	TTL int64 // (DISCO) Time to live in seconds
}

// type Providers string
type ServiceStatus string

// type ServiceType string

const (
	StatusAny          ServiceStatus = "*"
	StatusUnknown      ServiceStatus = "unknown"
	StatusRegistered   ServiceStatus = "registered"
	StatusHealthy      ServiceStatus = "healthy"
	StatusDeregistered ServiceStatus = "deregistered"
	StatusSuspicious   ServiceStatus = "suspicious"
)

const (
	// key 550e8400-e29b-41d4-a716-446655440000:api-gateway:aws:us-east-1:us-east-1a:internal:vpc-12345678:subnet-87654321:i-1234567890abcdef0:v1.2.3
	// <service_uuid>:
	// <name>:
	// <type>:
	// <status>:
	// <provider>:
	// <region>:
	// <zone>:
	// <network_id>:
	// <subnet_id>:
	// <instance_id>:
	// <version>
	ServiceKeyPattern = "%s:%s:%s:%s:%s:%s:%s:%s:%s:%s:%s"
)

const (
	defaultTTL time.Duration = 1 * time.Minute // Default TTL for service entries in minutes
)

func _generateServiceKey(serviceUUID string, serviceName string, serviceType string, status ServiceStatus, provider string, region string, zone string, networkID string, subnetID string, instanceID string, version string) string {
	return fmt.Sprintf(ServiceKeyPattern,
		serviceUUID,
		serviceName,
		serviceType,
		status,
		provider,
		region,
		zone,
		networkID,
		subnetID,
		instanceID,
		version,
	)
}

func _GenerateServiceKey(serviceEntry ServiceEntry) string {
	return _generateServiceKey(
		serviceEntry.ServiceUUID,
		serviceEntry.Name,
		serviceEntry.Type,
		serviceEntry.Status,
		serviceEntry.Provider,
		serviceEntry.Region,
		serviceEntry.Zone,
		serviceEntry.NetworkID,
		serviceEntry.SubnetID,
		serviceEntry.InstanceID,
		serviceEntry.Version,
	)
}

func _GenerateCredentialBasedSearchKey(serviceType string, provider string, region string, zone string, networkID string, subnetID string, instanceID string, version string) string {
	return _generateServiceKey(
		"*",
		"*",
		serviceType,
		"*",
		provider,
		region,
		zone,
		networkID,
		subnetID,
		instanceID,
		version,
	)
}

func _GenerateNewServiceValue(serviceEntry ServiceEntry) ([]byte, error) {
	now := time.Now().Format(time.RFC3339)
	serviceEntry.CreatedAt = now
	serviceEntry.LastHeardAt = now
	serviceEntry.Status = StatusRegistered
	serviceEntry.HeardCount = 0
	serviceEntry.ReportCount = 0
	serviceEntry.LastReportAt = now

	return json.Marshal(serviceEntry)
}

func RegisterNewService(client redisclient.Client, serviceEntry ServiceEntry) bool {
	startTime := time.Now()
	NewServiceData, err := _GenerateNewServiceValue(serviceEntry)
	NewServiceKey := _GenerateServiceKey(serviceEntry)

	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal service entry: %v", err), time.Since(startTime))
		return false
	}
	err = client.Set(NewServiceKey, NewServiceData, defaultTTL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to register new service: %v", err), time.Since(startTime))
		return false
	}
	return true
}

func IsServiceExists(client redisclient.Client, entry ServiceEntry) (bool, ServiceEntry) {

	searchKey := _GenerateCredentialBasedSearchKey(entry.Type, entry.Provider, entry.Region, entry.Zone, entry.NetworkID, entry.SubnetID, entry.InstanceID, entry.Version)
	var foundEntry ServiceEntry

	keys, err := client.FindKeys(searchKey)
	if err != nil {
		return false, ServiceEntry{}
	}
	if len(keys) == 0 {
		return false, ServiceEntry{}
	}
	if len(keys) >= 1 {
		byteVal, err := client.Get(keys[0])
		if err != nil {
			return false, ServiceEntry{}
		}

		err = json.Unmarshal(byteVal, &foundEntry)
		if err != nil {
			return false, ServiceEntry{}
		}

		return true, foundEntry
	}
	return true, ServiceEntry{}
}

func IsServiceExistsByUUID(client redisclient.Client, serviceUUID string) (bool, ServiceEntry) {
	searchKey := _generateServiceKey(serviceUUID, "*", "*", "*", "*", "*", "*", "*", "*", "*", "*")
	var foundEntry ServiceEntry

	keys, err := client.FindKeys(searchKey)
	if err != nil {
		return false, ServiceEntry{}
	}
	if len(keys) == 0 {
		return false, ServiceEntry{}
	}
	if len(keys) > 1 {
		return false, ServiceEntry{}
	}
	if len(keys) == 1 {
		byteVal, err := client.Get(keys[0])
		if err != nil {
			return false, ServiceEntry{}
		}
		err = json.Unmarshal(byteVal, &foundEntry)
		if err != nil {
			return false, ServiceEntry{}
		}

		return true, foundEntry
	}
	return false, ServiceEntry{}
}

func UpdateServiceEntry(client redisclient.Client, uuid string) (bool, error) {
	startTime := time.Now()

	// Fetch the existing entry to preserve CreatedAt and HeardCount, etc.
	exists, existingEntry := IsServiceExistsByUUID(client, uuid)
	oldKey := _GenerateServiceKey(existingEntry)
	if exists {
		// Preserve CreatedAt, HeardCount, ReportCount, etc.
		existingEntry.LastHeardAt = utils.GetFormatedCurrentTime()
		existingEntry.HeardCount++
		existingEntry.Status = StatusHealthy // Update status to healthy on heartbeat

		if existingEntry.ReportCount > env.GetReportToleranceCount() {
			existingEntry.Status = StatusSuspicious
			logger.HeartBeat(fmt.Sprintf("Service with UUID %s is marked as suspicious", uuid), time.Since(startTime))
			return false, errors.New("service entry is suspicious")
		}

	} else {
		// logger.Error(fmt.Sprintf("Service with UUID %s does not exist", uuid), time.Since(startTime))
		return false, errors.New("key not found")
	}

	updatedData, err := json.Marshal(existingEntry)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to marshal updated service entry: %v", err), time.Since(startTime))
		return false, errors.New("failed to marshal updated service entry")
	}

	serviceKey := _GenerateServiceKey(existingEntry)
	err = client.Set(serviceKey, updatedData, defaultTTL)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to update service entry: %v", err), time.Since(startTime))
		return false, errors.New("failed to update service entry")
	}

	err = client.Delete(oldKey)
	if err != nil {
		logger.Error(fmt.Sprintf("Failed to delete old service entry: %v", err), time.Since(startTime))
		// return false, errors.New("failed to delete old service entry")
	}

	return true, nil
}

func GetServicesFiltered(rclient redisclient.Client, serviceType string, healthStatus ServiceStatus, provider string, region string, zone string, networkID string, subnetID string, instanceID string, version string, pageSize int, pageOffset int) ([]ServiceEntry, error) {
	startTime := time.Now()

	if serviceType == "" {
		return nil, errors.New("serviceType is required")
	}
	if healthStatus == "" {
		healthStatus = StatusAny
	}
	if provider == "" {
		provider = "*"
	}
	if region == "" {
		region = "*"
	}
	if zone == "" {
		zone = "*"
	}
	if networkID == "" {
		networkID = "*"
	}
	if subnetID == "" {
		subnetID = "*"
	}
	if instanceID == "" {
		instanceID = "*"
	}
	if version == "" {
		version = "*"
	}

	searchKey := _generateServiceKey("*", "*", serviceType, healthStatus, provider, region, zone, networkID, subnetID, instanceID, version)
	keys, err := rclient.FindKeys(searchKey)
	if err != nil {
		return nil, err
	}

	// Apply pagination to keys before fetching data
	start := pageOffset * pageSize
	end := start + pageSize
	if start > len(keys) {
		return []ServiceEntry{}, nil
	}
	if end > len(keys) {
		end = len(keys)
	}
	keys = keys[start:end]

	var services []ServiceEntry
	for _, key := range keys {
		data, err := rclient.Get(key)
		if err != nil {
			continue
		}

		var entry ServiceEntry
		if err := json.Unmarshal(data, &entry); err != nil {
			continue
		}

		services = append(services, entry)
	}

	logger.Info(fmt.Sprintf("Found %d services for type (%s), status (%s), provider (%s), region (%s), zone (%s), networkID (%s), subnetID (%s), instanceID (%s), version (%s)", len(services), serviceType, healthStatus, provider, region, zone, networkID, subnetID, instanceID, version), time.Since(startTime))
	return services, nil
}

func DeregisterServiceEntry(rclient redisclient.Client, serviceUUID string) (bool, error) {

	keys, err := rclient.FindKeys(
		_generateServiceKey(serviceUUID, "*", "*", "*", "*", "*", "*", "*", "*", "*", "*"),
	)
	if err != nil {
		return false, err
	}

	if len(keys) == 0 {
		return true, nil
	}
	if len(keys) > 1 {
		logger.Debug(fmt.Sprintf("Multiple entries found for UUID %s", serviceUUID), 0)
		return true, fmt.Errorf("multiple entries found for UUID %s", serviceUUID)
	}
	if len(keys) == 1 {
		err := rclient.Delete(keys[0])
		if err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

func IsValidServiceStatus(status string) bool {
	switch status {
	case string(StatusHealthy), string(StatusUnknown), string(StatusSuspicious), string(StatusAny), string(StatusRegistered), string(StatusDeregistered):
		return true
	default:
		return false
	}
}

func DecideStatus(status string) ServiceStatus {
	logger.Debug(fmt.Sprintf("Deciding status for: %s", status), 0)
	switch status {
	case string(StatusHealthy):
		return StatusHealthy
	case string(StatusUnknown):
		return StatusUnknown
	case string(StatusSuspicious):
		return StatusSuspicious
	case string(StatusAny):
		return StatusAny
	case string(StatusRegistered):
		return StatusRegistered
	case string(StatusDeregistered):
		return StatusDeregistered
	default:
		return StatusUnknown
	}
}

func GetAllServiceStatuses() []string {
	return []string{
		string(StatusHealthy),
		string(StatusUnknown),
		string(StatusSuspicious),
		string(StatusAny),
		string(StatusRegistered),
		string(StatusDeregistered),
	}
}
