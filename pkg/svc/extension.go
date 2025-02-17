package svc

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/fuseml/fuseml-core/gen/extension"
	"github.com/fuseml/fuseml-core/pkg/domain"
	"github.com/fuseml/fuseml-core/pkg/util"
)

// extension registry service implementation.
type extensionRegistrySvc struct {
	logger   *log.Logger
	registry domain.ExtensionRegistry
}

// NewExtensionRegistryService returns the extension registry service implementation.
func NewExtensionRegistryService(logger *log.Logger, registry domain.ExtensionRegistry) extension.Service {
	return &extensionRegistrySvc{logger, registry}
}

func extensionToDomain(extension *extension.Extension) (result *domain.Extension) {
	return &domain.Extension{
		ID:            util.DerefString(extension.ID),
		Product:       util.DerefString(extension.Product),
		Version:       util.DerefString(extension.Version),
		Description:   util.DerefString(extension.Description),
		Zone:          util.DerefString(extension.Zone),
		Configuration: extension.Configuration,
	}
}

func extensionServiceToDomain(service *extension.ExtensionService) (result *domain.ExtensionService) {
	return &domain.ExtensionService{
		ExtensionServiceID: domain.ExtensionServiceID{
			ExtensionID: util.DerefString(service.ExtensionID),
			ID:          util.DerefString(service.ID),
		},
		Resource:      util.DerefString(service.Resource),
		Category:      util.DerefString(service.Category),
		Description:   util.DerefString(service.Description),
		AuthRequired:  util.DerefBool(service.AuthRequired),
		Configuration: service.Configuration,
	}
}

func extensionEndpointURLToDomain(URL *string) string {
	// the endpoint URL might be URL-encoded; attempt to decode it and ignore the error
	decodedURL := util.DerefString(URL)
	decodedURL, _ = url.PathUnescape(decodedURL)
	return decodedURL
}

func extensionEndpointToDomain(endpoint *extension.ExtensionEndpoint) (result *domain.ExtensionEndpoint) {
	return &domain.ExtensionEndpoint{
		ExtensionEndpointID: domain.ExtensionEndpointID{
			ExtensionID: util.DerefString(endpoint.ExtensionID),
			ServiceID:   util.DerefString(endpoint.ServiceID),
			URL:         extensionEndpointURLToDomain(endpoint.URL),
		},
		Type:          domain.ExtensionEndpointType(util.DerefString(endpoint.Type, string(domain.EETExternal))),
		Configuration: endpoint.Configuration,
	}
}

func extensionCredentialsToDomain(credentials *extension.ExtensionCredentials) (result *domain.ExtensionCredentials) {
	return &domain.ExtensionCredentials{
		ExtensionCredentialsID: domain.ExtensionCredentialsID{
			ExtensionID: util.DerefString(credentials.ExtensionID),
			ServiceID:   util.DerefString(credentials.ServiceID),
			ID:          util.DerefString(credentials.ID),
		},
		Scope:         domain.ExtensionCredentialScope(util.DerefString(credentials.Scope, string(domain.ECSGlobal))),
		Default:       util.DerefBool(credentials.Default),
		Projects:      credentials.Projects,
		Users:         credentials.Users,
		Configuration: credentials.Configuration,
	}
}

func extensionRecordToDomain(ext *extension.Extension) (result *domain.ExtensionRecord) {
	result = &domain.ExtensionRecord{
		Extension: *extensionToDomain(ext),
		Services:  make([]*domain.ExtensionServiceRecord, 0),
	}

	for _, service := range ext.Services {
		svcRecord := extensionServiceRecordToDomain(service)
		result.Services = append(result.Services, svcRecord)
	}

	return result
}

func extensionServiceRecordToDomain(service *extension.ExtensionService) (result *domain.ExtensionServiceRecord) {
	result = &domain.ExtensionServiceRecord{
		ExtensionService: *extensionServiceToDomain(service),
		Endpoints:        make([]*domain.ExtensionEndpoint, 0),
		Credentials:      make([]*domain.ExtensionCredentials, 0),
	}

	for _, endpoint := range service.Endpoints {
		result.Endpoints = append(result.Endpoints, extensionEndpointToDomain(endpoint))
	}

	for _, credentials := range service.Credentials {
		result.Credentials = append(result.Credentials, extensionCredentialsToDomain(credentials))
	}

	return result
}

func extensionQueryToDomain(query *extension.ExtensionQuery) (result *domain.ExtensionQuery) {
	result = &domain.ExtensionQuery{
		ExtensionID:        query.ExtensionID,
		Product:            query.Product,
		VersionConstraints: query.Version,
		Zone:               query.Zone,
		// If a zone is supplied, it must be used to do strict zone matching
		StrictZoneMatch: true,
		ServiceID:       query.ServiceID,
		ServiceResource: query.ServiceResource,
		ServiceCategory: query.ServiceCategory,
	}

	return result
}

func extensionToRest(ext *domain.Extension) *extension.Extension {
	return &extension.Extension{
		ID:            util.RefString(ext.ID),
		Product:       util.RefString(ext.Product),
		Version:       util.RefString(ext.Version),
		Description:   util.RefString(ext.Description),
		Zone:          util.RefString(ext.Zone),
		Configuration: ext.Configuration,
		Status: &extension.ExtensionStatus{
			Registered: ext.Registered.Format(time.RFC3339),
			Updated:    ext.Updated.Format(time.RFC3339),
		},
		Services: make([]*extension.ExtensionService, 0),
	}
}

func extensionServiceToRest(service *domain.ExtensionService) *extension.ExtensionService {
	return &extension.ExtensionService{
		ID:            util.RefString(service.ID),
		ExtensionID:   util.RefString(service.ExtensionID),
		Resource:      util.RefString(service.Resource),
		Category:      util.RefString(service.Category),
		Description:   util.RefString(service.Description),
		AuthRequired:  &service.AuthRequired,
		Configuration: service.Configuration,
		Status: &extension.ExtensionServiceStatus{
			Registered: service.Registered.Format(time.RFC3339),
			Updated:    service.Updated.Format(time.RFC3339),
		},
		Endpoints:   make([]*extension.ExtensionEndpoint, 0),
		Credentials: make([]*extension.ExtensionCredentials, 0),
	}
}

func extensionEndpointToRest(endpoint *domain.ExtensionEndpoint) *extension.ExtensionEndpoint {
	return &extension.ExtensionEndpoint{
		URL:           util.RefString(endpoint.URL),
		ExtensionID:   util.RefString(endpoint.ExtensionID),
		ServiceID:     util.RefString(endpoint.ServiceID),
		Type:          util.RefString(string(endpoint.Type)),
		Configuration: endpoint.Configuration,
		Status:        &extension.ExtensionEndpointStatus{},
	}
}

func obfuscateCredentials(config map[string]string) (result map[string]string) {
	result = make(map[string]string)
	for k := range config {
		result[k] = "<hidden>"
	}
	return result
}

func extensionCredentialsToRest(credentials *domain.ExtensionCredentials) *extension.ExtensionCredentials {
	return &extension.ExtensionCredentials{
		ID:            util.RefString(credentials.ID),
		ExtensionID:   util.RefString(credentials.ExtensionID),
		ServiceID:     util.RefString(credentials.ServiceID),
		Scope:         util.RefString(string(credentials.Scope)),
		Default:       &credentials.Default,
		Projects:      credentials.Projects,
		Users:         credentials.Users,
		Configuration: obfuscateCredentials(credentials.Configuration),
		Status: &extension.ExtensionCredentialsStatus{
			Created: credentials.Created.Format(time.RFC3339),
			Updated: credentials.Updated.Format(time.RFC3339),
		},
	}
}

func extensionRecordToRest(extRecord *domain.ExtensionRecord) (result *extension.Extension) {

	result = extensionToRest(&extRecord.Extension)

	for _, svcRecord := range extRecord.Services {
		result.Services = append(result.Services, extensionServiceRecordToRest(svcRecord))
	}

	return result
}

func extensionServiceRecordToRest(svcRecord *domain.ExtensionServiceRecord) (result *extension.ExtensionService) {

	result = extensionServiceToRest(&svcRecord.ExtensionService)

	for _, epRecord := range svcRecord.Endpoints {
		result.Endpoints = append(result.Endpoints, extensionEndpointToRest(epRecord))
	}

	for _, credsRecord := range svcRecord.Credentials {
		result.Credentials = append(result.Credentials, extensionCredentialsToRest(credsRecord))
	}

	return result
}

func errToRest(err error) error {
	switch err.(type) {
	case *domain.ErrExtensionNotFound:
		return extension.MakeNotFound(err)
	case *domain.ErrExtensionServiceNotFound:
		return extension.MakeNotFound(err)
	case *domain.ErrExtensionEndpointNotFound:
		return extension.MakeNotFound(err)
	case *domain.ErrExtensionCredentialsNotFound:
		return extension.MakeNotFound(err)
	case *domain.ErrExtensionExists:
		return extension.MakeConflict(err)
	case *domain.ErrExtensionServiceExists:
		return extension.MakeConflict(err)
	case *domain.ErrExtensionEndpointExists:
		return extension.MakeConflict(err)
	case *domain.ErrExtensionCredentialsExists:
		return extension.MakeConflict(err)
	default:
		return extension.MakeBadRequest(err)
	}
}

// Register an extension with the FuseML extension registry.
func (s *extensionRegistrySvc) RegisterExtension(ctx context.Context, req *extension.Extension) (*extension.Extension, error) {
	s.logger.Print("extension.registerExtension")
	extRecord, err := s.registry.RegisterExtension(ctx, extensionRecordToDomain(req))
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionRecordToRest(extRecord), nil
}

// Retrieve information about an extension.
func (s *extensionRegistrySvc) GetExtension(ctx context.Context, req *extension.GetExtensionPayload) (res *extension.Extension, err error) {
	s.logger.Print("extension.getExtension")
	extRecord, err := s.registry.GetExtension(ctx, req.ID, true)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionRecordToRest(extRecord), nil
}

// List extensions registered in FuseML
func (s *extensionRegistrySvc) ListExtensions(ctx context.Context, query *extension.ExtensionQuery) (res []*extension.Extension, err error) {
	s.logger.Print("extension.listExtensions")
	extRecords, err := s.registry.ListExtensions(ctx, extensionQueryToDomain(query))
	if err != nil {
		return nil, errToRest(err)
	}

	res = make([]*extension.Extension, len(extRecords))
	for i, extRecord := range extRecords {
		res[i] = extensionRecordToRest(extRecord)
	}

	return res, nil
}

// Update an extension registered in FuseML
func (s *extensionRegistrySvc) UpdateExtension(ctx context.Context, req *extension.Extension) (res *extension.Extension, err error) {
	s.logger.Print("extension.updateExtension")
	ext := extensionToDomain(req)
	extRecord, err := s.registry.GetExtension(ctx, ext.ID, false)
	if err != nil {
		return nil, errToRest(err)
	}
	// update only attributes present in the update request
	extUpdate := domain.Extension{
		ID:            ext.ID,
		Product:       util.DerefString(req.Product, extRecord.Product),
		Version:       util.DerefString(req.Version, extRecord.Version),
		Description:   util.DerefString(req.Description, extRecord.Description),
		Zone:          util.DerefString(req.Zone, extRecord.Zone),
		Configuration: extRecord.Configuration,
	}
	if req.Configuration != nil {
		extUpdate.Configuration = req.Configuration
	}

	err = s.registry.UpdateExtension(ctx, &extUpdate)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionToRest(&extUpdate), nil
}

// Delete an extension and its subtree of services, endpoints and credentials
func (s *extensionRegistrySvc) DeleteExtension(ctx context.Context, req *extension.DeleteExtensionPayload) (err error) {
	s.logger.Print("extension.deleteExtension")
	err = s.registry.RemoveExtension(ctx, req.ID)
	if err != nil {
		return errToRest(err)
	}
	return nil
}

// Add a service to an existing extension registered with the FuseML extension
// registry.
func (s *extensionRegistrySvc) AddService(ctx context.Context, service *extension.ExtensionService) (res *extension.ExtensionService, err error) {
	svcRecord, err := s.registry.AddService(ctx, extensionServiceRecordToDomain(service))
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionServiceRecordToRest(svcRecord), nil
}

// Retrieve information about a service belonging to an extension.
func (s *extensionRegistrySvc) GetService(ctx context.Context, req *extension.GetServicePayload) (res *extension.ExtensionService, err error) {
	svcRecord, err := s.registry.GetService(ctx, domain.ExtensionServiceID{
		ExtensionID: req.ExtensionID,
		ID:          req.ID,
	}, true)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionServiceRecordToRest(svcRecord), nil
}

// List all services associated with an extension registered in FuseML
func (s *extensionRegistrySvc) ListServices(ctx context.Context, req *extension.ListServicesPayload) (res []*extension.ExtensionService, err error) {
	extRecord, err := s.registry.GetExtension(ctx, req.ExtensionID, true)
	if err != nil {
		return nil, errToRest(err)
	}
	res = make([]*extension.ExtensionService, len(extRecord.Services))
	for i, svcRecord := range extRecord.Services {
		res[i] = extensionServiceRecordToRest(svcRecord)
	}
	return res, nil
}

// Update a service belonging to an extension registered in FuseML
func (s *extensionRegistrySvc) UpdateService(ctx context.Context, req *extension.ExtensionService) (res *extension.ExtensionService, err error) {
	s.logger.Print("extension.updateService")
	service := extensionServiceToDomain(req)
	svcRecord, err := s.registry.GetService(ctx, domain.ExtensionServiceID{
		ID:          service.ID,
		ExtensionID: service.ExtensionID,
	}, false)
	if err != nil {
		return nil, errToRest(err)
	}

	// update only attributes present in the update request
	svcUpdate := domain.ExtensionService{
		ExtensionServiceID: svcRecord.ExtensionServiceID,
		Resource:           util.DerefString(req.Resource, svcRecord.Resource),
		Category:           util.DerefString(req.Category, svcRecord.Category),
		Description:        util.DerefString(req.Description, svcRecord.Description),
		AuthRequired:       util.DerefBool(req.AuthRequired, svcRecord.AuthRequired),
		Configuration:      svcRecord.Configuration,
	}
	if req.Configuration != nil {
		svcUpdate.Configuration = req.Configuration
	}

	err = s.registry.UpdateService(ctx, &svcUpdate)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionServiceToRest(&svcUpdate), nil
}

// Delete an extension service and its subtree of endpoints and credentials
func (s *extensionRegistrySvc) DeleteService(ctx context.Context, req *extension.DeleteServicePayload) (err error) {
	s.logger.Print("extension.deleteService")
	err = s.registry.RemoveService(ctx, domain.ExtensionServiceID{
		ExtensionID: req.ExtensionID,
		ID:          req.ID,
	})
	if err != nil {
		return errToRest(err)
	}
	return nil
}

// Add an endpoint to an existing extension service registered with the FuseML
// extension registry.
func (s *extensionRegistrySvc) AddEndpoint(ctx context.Context, req *extension.ExtensionEndpoint) (res *extension.ExtensionEndpoint, err error) {
	s.logger.Print("extension.addEndpoint")
	endpoint, err := s.registry.AddEndpoint(ctx, extensionEndpointToDomain(req))
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionEndpointToRest(endpoint), nil
}

// Retrieve information about an endpoint belonging to an extension.
func (s *extensionRegistrySvc) GetEndpoint(ctx context.Context, req *extension.GetEndpointPayload) (res *extension.ExtensionEndpoint, err error) {
	s.logger.Print("extension.getEndpoint")
	endpoint, err := s.registry.GetEndpoint(ctx, domain.ExtensionEndpointID{
		ExtensionID: req.ExtensionID,
		ServiceID:   req.ServiceID,
		URL:         extensionEndpointURLToDomain(&req.URL),
	})
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionEndpointToRest(endpoint), nil
}

// List all endpoints associated with an extension service registered in FuseML
func (s *extensionRegistrySvc) ListEndpoints(ctx context.Context, req *extension.ListEndpointsPayload) (res []*extension.ExtensionEndpoint, err error) {
	s.logger.Print("extension.listEndpoints")
	svcRecord, err := s.registry.GetService(ctx, domain.ExtensionServiceID{
		ExtensionID: req.ExtensionID,
		ID:          req.ServiceID,
	}, true)
	if err != nil {
		return nil, errToRest(err)
	}
	res = make([]*extension.ExtensionEndpoint, len(svcRecord.Endpoints))
	for i, endpoint := range svcRecord.Endpoints {
		res[i] = extensionEndpointToRest(endpoint)
	}
	return res, nil
}

// Update an endpoint belonging to an extension service registered in FuseML
func (s *extensionRegistrySvc) UpdateEndpoint(ctx context.Context, req *extension.ExtensionEndpoint) (res *extension.ExtensionEndpoint, err error) {
	s.logger.Print("extension.updateEndpoint")
	endpoint := extensionEndpointToDomain(req)
	ep, err := s.registry.GetEndpoint(ctx, domain.ExtensionEndpointID{
		URL:         extensionEndpointURLToDomain(&endpoint.URL),
		ExtensionID: endpoint.ExtensionID,
		ServiceID:   endpoint.ServiceID,
	})
	if err != nil {
		return nil, errToRest(err)
	}

	// update only attributes present in the update request
	epUpdate := domain.ExtensionEndpoint{
		ExtensionEndpointID: ep.ExtensionEndpointID,
		Type:                domain.ExtensionEndpointType(util.DerefString(req.Type, string(ep.Type))),
		Configuration:       ep.Configuration,
	}
	if req.Configuration != nil {
		epUpdate.Configuration = req.Configuration
	}

	err = s.registry.UpdateEndpoint(ctx, &epUpdate)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionEndpointToRest(&epUpdate), nil
}

// Delete an extension endpoint
func (s *extensionRegistrySvc) DeleteEndpoint(ctx context.Context, req *extension.DeleteEndpointPayload) (err error) {
	s.logger.Print("extension.deleteEndpoint")
	err = s.registry.RemoveEndpoint(ctx, domain.ExtensionEndpointID{
		ExtensionID: req.ExtensionID,
		ServiceID:   req.ServiceID,
		URL:         extensionEndpointURLToDomain(&req.URL),
	})
	if err != nil {
		return errToRest(err)
	}
	return nil
}

// Add a set of credentials to an existing extension service registered with
// the FuseML extension registry.
func (s *extensionRegistrySvc) AddCredentials(ctx context.Context, req *extension.ExtensionCredentials) (res *extension.ExtensionCredentials, err error) {
	credentials, err := s.registry.AddCredentials(ctx, extensionCredentialsToDomain(req))
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionCredentialsToRest(credentials), nil
}

// Retrieve information about a set of credentials belonging to an extension.
func (s *extensionRegistrySvc) GetCredentials(ctx context.Context, req *extension.GetCredentialsPayload) (res *extension.ExtensionCredentials, err error) {
	s.logger.Print("extension.getCredentials")
	credentials, err := s.registry.GetCredentials(ctx, domain.ExtensionCredentialsID{
		ExtensionID: req.ExtensionID,
		ServiceID:   req.ServiceID,
		ID:          req.ID,
	})
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionCredentialsToRest(credentials), nil
}

// List all credentials associated with an extension service registered in
// FuseML
func (s *extensionRegistrySvc) ListCredentials(ctx context.Context, req *extension.ListCredentialsPayload) (res []*extension.ExtensionCredentials, err error) {
	s.logger.Print("extension.listCredentials")
	svcRecord, err := s.registry.GetService(ctx, domain.ExtensionServiceID{
		ExtensionID: req.ExtensionID,
		ID:          req.ServiceID,
	}, true)
	if err != nil {
		return nil, errToRest(err)
	}
	res = make([]*extension.ExtensionCredentials, len(svcRecord.Credentials))
	for i, credentials := range svcRecord.Credentials {
		res[i] = extensionCredentialsToRest(credentials)
	}
	return res, nil
}

// Update a set of credentials belonging to an extension service registered in
// FuseML
func (s *extensionRegistrySvc) UpdateCredentials(ctx context.Context, req *extension.ExtensionCredentials) (res *extension.ExtensionCredentials, err error) {
	s.logger.Print("extension.updateCredentials")
	credentials := extensionCredentialsToDomain(req)
	cred, err := s.registry.GetCredentials(ctx, domain.ExtensionCredentialsID{
		ID:          credentials.ID,
		ExtensionID: credentials.ExtensionID,
		ServiceID:   credentials.ServiceID,
	})
	if err != nil {
		return nil, errToRest(err)
	}

	// update only attributes present in the update request
	credUpdate := domain.ExtensionCredentials{
		ExtensionCredentialsID: cred.ExtensionCredentialsID,
		Scope:                  domain.ExtensionCredentialScope(util.DerefString(req.Scope, string(cred.Scope))),
		Default:                util.DerefBool(req.Default, cred.Default),
		Projects:               cred.Projects,
		Users:                  cred.Users,
		Configuration:          cred.Configuration,
	}
	if req.Configuration != nil {
		credUpdate.Configuration = req.Configuration
	}
	if req.Projects != nil {
		credUpdate.Projects = req.Projects
	}
	if req.Users != nil {
		credUpdate.Users = req.Users
	}

	err = s.registry.UpdateCredentials(ctx, &credUpdate)
	if err != nil {
		return nil, errToRest(err)
	}
	return extensionCredentialsToRest(&credUpdate), nil
}

// Delete a set of extension credentials
func (s *extensionRegistrySvc) DeleteCredentials(ctx context.Context, req *extension.DeleteCredentialsPayload) (err error) {
	s.logger.Print("extension.deleteCredentials")
	err = s.registry.RemoveCredentials(ctx, domain.ExtensionCredentialsID{
		ExtensionID: req.ExtensionID,
		ServiceID:   req.ServiceID,
		ID:          req.ID,
	})
	if err != nil {
		return errToRest(err)
	}
	return nil
}
