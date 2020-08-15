package cliopt

type Version struct {
	Constraint Constraint
	IsLatest   bool
}

func (v *Version) UnmarshalFlag(data string) error {
	if data == "latest" {
		v.IsLatest = true

		return nil
	}

	return v.Constraint.UnmarshalFlag(data)
}
