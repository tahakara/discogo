package routes

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	serviceconfigloader "github.com/tahakara/discogo/internal/config/serviceconfiguration"
	"github.com/tahakara/discogo/internal/logger"
	redisclient "github.com/tahakara/discogo/internal/redis"
	redishelper "github.com/tahakara/discogo/internal/redis/helper"
	"github.com/tahakara/discogo/internal/utils"
)

type ServiceInfo struct {
	ServiceID   string `json:"serviceID"`
	ServiceAddr string `json:"serviceAddr"`
}

type DiscoverResponse struct {
	Status        string        `json:"status"`
	Message       string        `json:"message,omitempty"`
	ServiceTypes  []string      `json:"serviceTypes,omitempty"`
	ProviderTypes []string      `json:"providerTypes,omitempty"`
	StatusTypes   []string      `json:"statusTypes,omitempty"`
	Services      []ServiceInfo `json:"services"`
}

// DiscoverHandler handles service discovery requests.
//
// @Summary      Discover services
// @Description  Retrieves a list of services filtered by query parameters such as service type, provider, region, zone, network ID, subnet ID, instance ID, and version. Supports pagination.
// @Tags         DiscoGo
// @Accept       json
// @Produce      json
// @Param        servicetype   query     string  true   "Service type to discover"  Enums(mock,test,perftest,loadgen,staging,dev,debug,mq,eventbus,notify,email,sms,push,inmsg,chat,monitor,log,alert,health,cb,lb,discovery,config,util,helper,migrate,cleanup,archive,maint,other,stream,audio,live,transcode,abr,drm,quality,user,auth,authz,profile,prefs,social,watchlist,history,web,mobile,admin,cdn,assets,img,video,catalog,recommend,search,personal,ingest,metadata,subtitle,thumb,sub,billing,payment,pricing,trial,entitle,revenue,secgw,waf,fraud,audit,encrypt,kms,comply,threat,workflow,scheduler,pipeline,etl,batch,eventproc,orchestrate,3rdapi,partner,social,paygate,cdnint,cloudstor,tracker,gw,rest,graphql,grpc,ws,webhook,ratelimit,db,analyticsdb,cache,file,object,datalake,backup,sync,analytics,rtanalytics,abtest,flags,ml,ds,report,metrics)  // Replace with actual service types
// @Param		 status query     string  false  "Service status"            Enums(healthy,unknown,suspicious,registered,deregistered) // Replace with actual statuses
// @Param        provider      query     string  false  "Service provider"          Enums(provider1,provider2,...) // Replace with actual providers
// @Param        region        query     string  false  "Region"
// @Param        zone          query     string  false  "Zone"
// @Param        networkid     query     string  false  "Network ID"
// @Param        subnetid      query     string  false  "Subnet ID"
// @Param        instanceid    query     string  false  "Instance ID"
// @Param        version       query     string  false  "Service version"
// @Param        pagesize      query     int     false  "Number of results per page (1-10)" minimum(1) maximum(10)
// @Param        pageoffset    query     int     false  "Page offset (>= 0)" minimum(0)
// @Success      200  {object}  DiscoverResponse  "List of discovered services"
// @Failure      400  {object}  DiscoverResponse  "Invalid request parameters"
// @Failure      500  {object}  DiscoverResponse  "Internal server error"
// @Router       /disco/discover [get]
func DiscoverHandler(w http.ResponseWriter, r *http.Request, rclient redisclient.Client) {
	startTime := time.Now()
	serviceType := r.URL.Query().Get("servicetype")
	selectedServiceStatus := r.URL.Query().Get("status") // Optional, default to any status
	provider := r.URL.Query().Get("provider")
	region := r.URL.Query().Get("region")
	zone := r.URL.Query().Get("zone")
	networkID := r.URL.Query().Get("networkid")
	subnetID := r.URL.Query().Get("subnetid")
	instanceID := r.URL.Query().Get("instanceid")
	version := r.URL.Query().Get("version")
	pageSizeStr := r.URL.Query().Get("pagesize")
	pageOffsetStr := r.URL.Query().Get("pageoffset")

	const (
		defaultPageSize = 10
		maxPageSize     = 10
		minPageOffset   = 0
	)

	pageSize := defaultPageSize
	pageOffset := minPageOffset

	if pageSizeStr != "" {
		if n, err := strconv.Atoi(pageSizeStr); err == nil && n > 0 && n <= maxPageSize {
			pageSize = n
		} else {
			utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
				Status:  "error",
				Message: "Invalid 'pagesize' query parameter (must be 1-10)",
			})
			return
		}
	}

	if pageOffsetStr != "" {
		if n, err := strconv.Atoi(pageOffsetStr); err == nil && n >= minPageOffset {
			pageOffset = n
		} else {
			utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
				Status:  "error",
				Message: "Invalid 'pageoffset' query parameter (must be >= 0)",
			})
			return
		}
	}

	if serviceType == "" {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
			Status:  "error",
			Message: "Missing 'servicetype' query parameter",
		})
		return
	}

	if !serviceconfigloader.IsValidServiceType(serviceType) {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
			Status:       "error",
			Message:      "Invalid 'servicetype' query parameter",
			ServiceTypes: serviceconfigloader.GetAllServiceTypes(),
		})
		return
	}

	if provider != "" && !serviceconfigloader.IsValidProvider(provider) {
		utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
			Status:        "error",
			Message:       "Invalid 'provider' query parameter",
			ProviderTypes: serviceconfigloader.GetAllProviders(),
		})
		return
	}

	if selectedServiceStatus != "" {
		if !redishelper.IsValidServiceStatus(selectedServiceStatus) {
			utils.WriteJSONResponse(w, http.StatusBadRequest, DiscoverResponse{
				Status:      "error",
				Message:     "Invalid 'status' query parameter",
				StatusTypes: redishelper.GetAllServiceStatuses(),
			})
			return
		}
	}

	services, err := redishelper.GetServicesFiltered(
		rclient,
		serviceType,
		redishelper.DecideStatus(selectedServiceStatus), // Use StatusAny to match any health status
		provider,
		region,
		zone,
		networkID,
		subnetID,
		instanceID,
		version,
		pageSize,
		pageOffset,
	)
	if err != nil {
		utils.WriteJSONResponse(w, http.StatusInternalServerError, DiscoverResponse{
			Status:  "error",
			Message: "Failed to retrieve services",
		})
		return
	}

	var serviceInfos []ServiceInfo
	for _, service := range services {
		addr := ""
		if service.Addr4 != "" {
			addr = service.Addr4 + ":" + strconv.Itoa(service.Port4)
		} else if service.Addr6 != "" {
			addr = "[" + service.Addr6 + "]:" + strconv.Itoa(service.Port6)
		}
		serviceInfos = append(serviceInfos, ServiceInfo{
			ServiceID:   service.ServiceUUID,
			ServiceAddr: addr,
		})
	}

	logger.Discovery(fmt.Sprintf("Discovered '%s':(%v)", serviceType, len(serviceInfos)), time.Since(startTime))
	utils.WriteJSONResponse(w, http.StatusOK, DiscoverResponse{
		Status:   "success",
		Message:  "Services discovered successfully",
		Services: serviceInfos,
	})
}
