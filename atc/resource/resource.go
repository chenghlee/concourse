package resource

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/concourse/concourse/atc"
	"github.com/concourse/concourse/atc/runtime"
)

const resultCachePropertyName = "concourse:resource-result"

type VersionResult struct {
	Version  atc.Version         `json:"version"`
	Metadata []atc.MetadataField `json:"metadata,omitempty"`
}

type Resource struct {
	Source  atc.Source  `json:"source"`
	Params  atc.Params  `json:"params,omitempty"`
	Version atc.Version `json:"version,omitempty"`
}

func (resource Resource) Signature() ([]byte, error) {
	return json.Marshal(resource)
}

func (resource Resource) Check(ctx context.Context, container runtime.Container, stderr io.Writer) ([]atc.Version, error) {
	spec := runtime.ProcessSpec{
		Path: "/opt/resource/check",
	}

	var versions []atc.Version
	if err := resource.run(ctx, container, spec, stderr, false, &versions); err != nil {
		return nil, err
	}
	return versions, nil
}

func (resource Resource) Get(ctx context.Context, container runtime.Container, stderr io.Writer) (VersionResult, error) {
	var versionResult VersionResult

	properties, err := container.Properties()
	if err != nil {
		return VersionResult{}, err
	}

	if result := properties[resultCachePropertyName]; result != "" {
		if err := json.Unmarshal([]byte(result), &versionResult); err != nil {
			return VersionResult{}, err
		}
		return versionResult, nil
	}

	spec := runtime.ProcessSpec{
		Path: "/opt/resource/in",
		Args: []string{ResourcesDir("get")},
	}

	if err := resource.run(ctx, container, spec, stderr, true, &versionResult); err != nil {
		return VersionResult{}, err
	}

	if err := resource.cacheResult(container, versionResult); err != nil {
		return VersionResult{}, err
	}

	return versionResult, nil
}

func (resource Resource) Put(ctx context.Context, container runtime.Container, stderr io.Writer) (VersionResult, error) {
	var versionResult VersionResult

	properties, err := container.Properties()
	if err != nil {
		return VersionResult{}, err
	}

	if result := properties[resultCachePropertyName]; result != "" {
		if err := json.Unmarshal([]byte(result), &versionResult); err != nil {
			return VersionResult{}, err
		}
		return versionResult, nil
	}

	spec := runtime.ProcessSpec{
		Path: "/opt/resource/out",
		Args: []string{ResourcesDir("put")},
	}

	if err := resource.run(ctx, container, spec, stderr, true, &versionResult); err != nil {
		return VersionResult{}, err
	}
	if versionResult.Version == nil {
		return VersionResult{}, fmt.Errorf("resource script (%s %s) output a null version", spec.Path, strings.Join(spec.Args, " "))
	}

	if err := resource.cacheResult(container, versionResult); err != nil {
		return VersionResult{}, err
	}

	return versionResult, nil
}

func (resource Resource) run(ctx context.Context, container runtime.Container, spec runtime.ProcessSpec, stderr io.Writer, attach bool, output interface{}) error {
	input, err := resource.Signature()
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	io := runtime.ProcessIO{
		Stdin:  bytes.NewBuffer(input),
		Stdout: buf,
		Stderr: stderr,
	}

	var result runtime.ProcessResult
	if attach {
		result, err = container.Attach(ctx, spec, io)
		if err == nil {
			goto finished_execution
		}
		buf.Reset()
	}

	result, err = container.Run(ctx, spec, io)
	if err != nil {
		return err
	}

finished_execution:
	if result.ExitStatus != 0 {
		return runtime.ErrResourceScriptFailed{
			Path:       spec.Path,
			Args:       spec.Args,
			ExitStatus: result.ExitStatus,
		}
	}

	return json.Unmarshal(buf.Bytes(), output)
}

func (resource Resource) cacheResult(container runtime.Container, result VersionResult) error {
	payload, err := json.Marshal(result)
	if err != nil {
		return err
	}

	return container.SetProperty(resultCachePropertyName, string(payload))
}

func ResourcesDir(suffix string) string {
	return filepath.Join("/tmp", "build", suffix)
}
