package solver

import (
	"github.com/hbagdi/go-kong/kong"
	"github.com/kong/deck/crud"
	"github.com/kong/deck/diff"
	cruds "github.com/kong/deck/solver/kong"
	drycrud "github.com/kong/deck/solver/kong/dry"
	"github.com/pkg/errors"
)

// Solve generates a diff and walks the graph.
func Solve(doneCh chan struct{}, syncer *diff.Syncer,
	client *kong.Client, dry bool) []error {
	var r *crud.Registry
	var err error
	if dry {
		r, err = buildDryRegistry(client)
	} else {
		r, err = buildRegistry(client)
	}
	if err != nil {
		return append([]error{}, errors.Wrapf(err, "cannot build registry"))
	}

	return syncer.Run(doneCh, 10, func(e diff.Event) (crud.Arg, error) {
		return r.Do(e.Kind, e.Op, e)
	})
}

func buildDryRegistry(client *kong.Client) (*crud.Registry, error) {
	var r crud.Registry
	err := r.Register("service", &drycrud.ServiceCRUD{})
	if err != nil {
		return nil, errors.Wrapf(err, "registering 'service' crud")
	}
	err = r.Register("route", &drycrud.RouteCRUD{})
	if err != nil {
		return nil, errors.Wrapf(err, "registering 'route' crud")
	}
	return &r, nil
}

func buildRegistry(client *kong.Client) (*crud.Registry, error) {
	var r crud.Registry
	var err error
	service, err := cruds.NewServiceCRUD(client)
	if err != nil {
		return nil, errors.Wrap(err, "creating a service CRUD")
	}
	err = r.Register("service", service)
	if err != nil {
		return nil, errors.Wrapf(err, "registering 'service' crud")
	}
	route, err := cruds.NewRouteCRUD(client)
	if err != nil {
		return nil, errors.Wrap(err, "creating a route CRUD")
	}
	err = r.Register("route", route)
	if err != nil {
		return nil, errors.Wrapf(err, "registering 'route' crud")
	}
	return &r, nil
}
