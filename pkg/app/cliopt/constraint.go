package cliopt

import (
	"github.com/Masterminds/semver"
	"github.com/pkg/errors"
)

type Constraint struct {
	*semver.Constraints
	RawValue string
}

func (o *Constraint) UnmarshalFlag(data string) error {
	con, err := semver.NewConstraint(data)
	if err != nil {
		return errors.Wrap(err, "parsing version constraint")
	}

	o.Constraints = con
	o.RawValue = data

	return nil
}
