package structs

// ServiceDefinition is used to JSON decode the Service definitions. For
// documentation on specific fields see NodeService which is better documented.
type ServiceDefinition struct {
	Kind              ServiceKind
	ID                string
	Name              string
	Tags              []string
	Address           string
	Meta              map[string]string
	Port              int
	Check             CheckType
	Checks            CheckTypes
	Token             string
	EnableTagOverride bool
	ProxyDestination  string
}

func (s *ServiceDefinition) NodeService() *NodeService {
	ns := &NodeService{
		Kind:              s.Kind,
		ID:                s.ID,
		Service:           s.Name,
		Tags:              s.Tags,
		Address:           s.Address,
		Meta:              s.Meta,
		Port:              s.Port,
		EnableTagOverride: s.EnableTagOverride,
		ProxyDestination:  s.ProxyDestination,
	}
	if ns.ID == "" && ns.Service != "" {
		ns.ID = ns.Service
	}
	return ns
}

func (s *ServiceDefinition) CheckTypes() (checks CheckTypes, err error) {
	if !s.Check.Empty() {
		err := s.Check.Validate()
		if err != nil {
			return nil, err
		}
		checks = append(checks, &s.Check)
	}
	for _, check := range s.Checks {
		if err := check.Validate(); err != nil {
			return nil, err
		}
		checks = append(checks, check)
	}
	return checks, nil
}
