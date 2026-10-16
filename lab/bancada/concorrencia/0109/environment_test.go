/**
 * Crystalline Lineage
 * @prompt 00_nucleo/prompts/lab/bancada/concorrencia/0109/environment_test.md
 * @layer lab
 * @updated 2026-10-06
 */
package concurrency

import (
	"reflect"
	"testing"
)

func environmentFixture() (Request, Environment) {
	r := request("package p;import \"example/dep\";func f(){_=make(dep.C,3)}")
	dep := "package dep;type C chan int"
	r.Files = append(r.Files, SourceFile{Path: "dep/type.go", Content: dep, SHA256: Digest(dep)})
	e := Environment{SchemaVersion: "1", Toolchain: "AUTHOR_FIXTURE", GOOS: "linux", GOARCH: "amd64", Packages: []PackageBinding{{ImportPath: "example/dep", Directory: "dep"}, {ImportPath: "example/main", Directory: "."}}}
	return r, e
}
func TestSuppliedSourceEnvironmentAndFailures(t *testing.T) {
	r, e := environmentFixture()
	observed, err := ObserveWithEnvironment(r, e)
	if err != nil || profile(observed)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 1 {
		t.Fatalf("environment %+v %v", observed, err)
	}
	if observed.EnvironmentSHA256 == "" || hasGap(observed, "SUPPLIED_PACKAGE_TYPECHECK_FAILED") {
		t.Fatal("unbound/failed environment")
	}
	for _, o := range observed.Observations {
		if o.Property == "channel_type_identity" && o.Authority != "GO_TYPES_SUPPLIED_ENVIRONMENT" {
			t.Fatal("wrong type authority")
		}
	}
	e.Packages[0], e.Packages[1] = e.Packages[1], e.Packages[0]
	again, err := ObserveWithEnvironment(r, e)
	if err != nil || !reflect.DeepEqual(observed, again) {
		t.Fatal("environment order changed report")
	}
	for _, kind := range []string{"missing", "invalid-dependency", "cycle", "corrupt-export", "missing-export"} {
		t.Run(kind, func(t *testing.T) {
			input, env := environmentFixture()
			switch kind {
			case "missing":
				input.Files = input.Files[:1]
			case "invalid-dependency":
				input.Files[1].Content = "package dep;type C chan missing"
				input.Files[1].SHA256 = Digest(input.Files[1].Content)
			case "cycle":
				input.Files[1].Content = "package dep;import _ \"example/main\";type C chan int"
				input.Files[1].SHA256 = Digest(input.Files[1].Content)
			case "corrupt-export":
				input.Files = input.Files[:1]
				env.Packages = env.Packages[1:]
				env.Exports = []ExportData{{ImportPath: "example/dep", Origin: "AUTHOR_BAD_DATA", Data: []byte("not export data"), SHA256: Digest("not export data")}}
			case "missing-export":
				input.Files[0].Content = "package p;import \"fmt\";type C chan int;func f(){_=make(C,3);fmt.Println()}"
				input.Files[0].SHA256 = Digest(input.Files[0].Content)
			}
			actual, err := ObserveWithEnvironment(input, env)
			if err != nil || profile(actual)["channel_type_identity=RESOLVED_CHANNEL_TYPE"] != 0 || !hasGap(actual, "SUPPLIED_PACKAGE_TYPECHECK_FAILED") {
				t.Fatalf("false credit %+v %v", actual, err)
			}
		})
	}
}
func TestEnvironmentValidationAndNoMutation(t *testing.T) {
	r, e := environmentFixture()
	before := append([]PackageBinding(nil), e.Packages...)
	if _, err := ObserveWithEnvironment(r, e); err != nil || !reflect.DeepEqual(before, e.Packages) {
		t.Fatal("environment mutated")
	}
	for _, mutate := range []func(*Environment){
		func(e *Environment) { e.Packages = append(e.Packages, e.Packages[0]) },
		func(e *Environment) { e.Packages[0].Directory = ".." },
		func(e *Environment) { e.Packages[0].ImportPath = "../escape" },
		func(e *Environment) {
			e.Exports = []ExportData{{ImportPath: "example/dep", Origin: "FIXTURE", Data: []byte("x"), SHA256: Digest("x")}}
		},
		func(e *Environment) {
			e.Exports = []ExportData{{ImportPath: "sync", Origin: "FIXTURE", Data: []byte("x"), SHA256: Digest("different")}}
		},
		func(e *Environment) { e.GOARCH = "arm64" },
		func(e *Environment) { e.Exports = make([]ExportData, 129) },
	} {
		_, candidate := environmentFixture()
		mutate(&candidate)
		if _, err := ObserveWithEnvironment(r, candidate); err == nil {
			t.Fatal("invalid environment accepted")
		}
	}
}
