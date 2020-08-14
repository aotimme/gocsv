module github.com/aotimme/gocsv

go 1.14

replace (
	github.com/aotimme/gocsv/cmd => ./cmd
	github.com/aotimme/gocsv/csv => ./csv
)

require (
	github.com/Masterminds/goutils v1.1.0 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/Masterminds/sprig v2.22.0+incompatible // indirect
	github.com/aotimme/gocsv/cmd v0.0.0
	github.com/aotimme/gocsv/csv v0.0.0
	github.com/google/uuid v1.1.1 // indirect
	github.com/huandu/xstrings v1.3.2 // indirect
	github.com/imdario/mergo v0.3.10 // indirect
	github.com/mitchellh/copystructure v1.0.0 // indirect
	golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect
)
