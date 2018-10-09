package tools

import "os"

const (
	// DefaultPermission is the permission that is applied when creating
	// a file or a folder.
	DefaultPermission os.FileMode = 0755

	// GoExtension is the file extension for the Golang files.
	GoExtension string = ".go"

	// KeyMetadata is the key in the toml file Gopkg.toml to indicate the
	// metadatas for a component.
	KeyMetadata = "metadata"

	// KeyConstraint is the key in the toml file Gopkg.toml to indicate a constraint
	KeyConstraint = "constraint"

	// KeyOverride is the key in the toml file to indicate that the entry overrides
	// a non direct constraint.
	KeyOverride = "override"

	// KeyConstraintName is the field name used in the Gopkg.toml file to
	// mark the constraint name field.
	KeyConstraintName = "name"

	// KeyConstraintRevision is the field name used in the Gopkg.toml file to
	// mark the constraint revision field. The revision is the SHA1 of the commit.
	KeyConstraintRevision = "revision"

	// KeyConstraintBranch is the field name used in the Gopkg.toml file to
	// mark the constraint branch field. This is the git branch name.
	KeyConstraintBranch = "branch"

	// KeyConstraintVersion is the field name used in the Gopkg.toml file to
	// mark the constraint version field. Version must be a semver value.
	KeyConstraintVersion = "version"

	// KeyTypeOfComponent is the field in the metadata to mark the type of the
	// component
	KeyTypeOfComponent = "type"

	// MigrationSeparator is the string used in the migration string to separator source from
	// origin
	MigrationSeparator = "->"
)
