package main

import (
	"fmt"
	"os"

	"github.com/GoogleContainerTools/kpt-functions-catalog/functions/go/set-namespace/generated"
	"sigs.k8s.io/kustomize/api/k8sdeps/kunstruct"
	"sigs.k8s.io/kustomize/api/konfig/builtinpluginconsts"
	"sigs.k8s.io/kustomize/api/resmap"
	"sigs.k8s.io/kustomize/api/resource"
	"sigs.k8s.io/kustomize/api/types"
	"sigs.k8s.io/kustomize/kyaml/fn/framework"
	"sigs.k8s.io/kustomize/kyaml/fn/framework/command"
	kyaml "sigs.k8s.io/kustomize/kyaml/yaml"
	"sigs.k8s.io/yaml"
)

const (
	fnConfigGroup      = "fn.kpt.dev"
	fnConfigVersion    = "v1alpha1"
	fnConfigAPIVersion = fnConfigGroup + "/" + fnConfigVersion
	fnConfigKind       = "SetNamespaceConfig"
)

//nolint
func main() {
	resourceList := &framework.ResourceList{}
	resourceList.FunctionConfig = &kyaml.RNode{}
	asp := SetNamespaceProcessor{}
	cmd := command.Build(&asp, command.StandaloneEnabled, false)

	cmd.Short = generated.SetNamespaceShort
	cmd.Long = generated.SetNamespaceLong
	cmd.Example = generated.SetNamespaceExamples
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

type SetNamespaceProcessor struct{}

func (snp *SetNamespaceProcessor) Process(resourceList *framework.ResourceList) error {
	err := run(resourceList)
	if err != nil {
		resourceList.Result = &framework.Result{
			Name: "set-namespace",
			Items: []framework.ResultItem{
				{
					Message:  err.Error(),
					Severity: framework.Error,
				},
			},
		}
		return resourceList.Result
	}
	return nil
}

type transformerConfig struct {
	FieldSpecs types.FsSlice `json:"namespace,omitempty" yaml:"namespace,omitempty"`
}

type setNamespaceFunction struct {
	kyaml.ResourceMeta `json:",inline" yaml:",inline"`
	plugin             `json:",inline" yaml:",inline"`
}

func (f *setNamespaceFunction) Config(rn *kyaml.RNode) error {
	switch {
	case f.validGVK(rn, "v1", "ConfigMap"):
		f.plugin.Namespace = rn.GetDataMap()["namespace"]
	case f.validGVK(rn, fnConfigAPIVersion, fnConfigKind):
		// input config is a CRD
		y, err := rn.String()
		if err != nil {
			return fmt.Errorf("cannot get YAML from RNode: %w", err)
		}
		err = yaml.Unmarshal([]byte(y), &f.plugin)
		if err != nil {
			return fmt.Errorf("failed to unmarshal config %#v: %w", y, err)
		}
	default:
		return fmt.Errorf("function config must be a ConfigMap or %s", fnConfigKind)
	}

	if f.plugin.Namespace == "" {
		return fmt.Errorf("input namespace cannot be empty")
	}
	tc, err := getDefaultConfig()
	if err != nil {
		return err
	}
	// set default field specs
	f.plugin.FieldSpecs = append(f.plugin.FieldSpecs, tc.FieldSpecs...)
	return nil
}

func (f *setNamespaceFunction) Run(items []*kyaml.RNode) ([]*kyaml.RNode, error) {
	resmapFactory := newResMapFactory()
	resMap, err := resmapFactory.NewResMapFromRNodeSlice(items)
	if err != nil {
		return nil, err
	}
	err = f.plugin.Transform(resMap)
	if err != nil {
		return nil, fmt.Errorf("failed to run transformer: %w", err)
	}
	return resMap.ToRNodeSlice()
}

func (f *setNamespaceFunction) validGVK(rn *kyaml.RNode, apiVersion, kind string) bool {
	meta, err := rn.GetMeta()
	if err != nil {
		return false
	}
	if meta.APIVersion != apiVersion || meta.Kind != kind {
		return false
	}
	return true
}

//nolint
func getDefaultConfig() (transformerConfig, error) {
	defaultConfigString := builtinpluginconsts.GetDefaultFieldSpecsAsMap()["namespace"]
	var defaultConfig transformerConfig
	err := yaml.Unmarshal([]byte(defaultConfigString), &defaultConfig)
	return defaultConfig, err
}

//nolint
func newResMapFactory() *resmap.Factory {
	resourceFactory := resource.NewFactory(kunstruct.NewKunstructuredFactoryImpl())
	return resmap.NewFactory(resourceFactory, nil)
}

func run(resourceList *framework.ResourceList) error {
	var fn setNamespaceFunction
	err := fn.Config(resourceList.FunctionConfig)
	if err != nil {
		return fmt.Errorf("failed to configure function: %w", err)
	}

	resourceList.Items, err = fn.Run(resourceList.Items)
	if err != nil {
		return fmt.Errorf("failed to run function: %w", err)
	}
	return nil
}
