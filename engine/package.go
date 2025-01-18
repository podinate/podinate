package engine

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"

	hcl "github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/hcl/v2/hclparse"
	"github.com/podinate/podinate/engine/helpers"
	"github.com/podinate/podinate/kube_client"
	"github.com/podinate/podinate/tui"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Ref: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
var StandardLabels = map[string]string{
	"app.kubernetes.io/managed-by": "Podinate",
}

// Resource represents a group of related Kubernetes Objects, for example a Pod, or the contents of a YAML manifest
type Resource interface {
	// Get the array of Kubernetes objects that is the desired state of the resource
	GetObjects(context.Context) ([]runtime.Object, error)
	// Get the display name type of the resource
	GetType() ResourceType
	// GetName returns the name of the resource
	GetName() string
}

// Package represents a package to be installed, either the current state or the desired state
type Package struct {
	Name      string
	Namespace string
	Resources []Resource
	Labels    map[string]string
}

// PodinateHCLSpec is the HCL spec for a Podinate block
var podinateHCLSpec = &hcldec.BlockSpec{
	TypeName: "podinate",
	Nested: &hcldec.ObjectSpec{
		"namespace": &hcldec.AttrSpec{
			Name:     "namespace",
			Type:     cty.String,
			Required: true,
		},
		"package": &hcldec.AttrSpec{
			Name:     "package",
			Type:     cty.String,
			Required: true,
		},
	},
}

func Parse(paths []string) (*Package, error) {
	// TODO: Make this support multiple files

	extension := strings.ToLower(filepath.Ext(paths[0])[1:])

	// If the file extension shows it is a Kubernetes yaml or json, parse it as such
	if slices.Contains([]string{"yaml", "yml", "json"}, extension) {
		return ParseYaml(paths[0])
	}

	// Assume a PodFile and parse as such
	return ParsePodfile(paths[0])
}

func ParseYaml(path string) (*Package, error) {

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	//fmt.Println("Content: ", string(content))

	var objects []runtime.Object
	var mf Manifest
	mf.SetName(path)
	mf.SetPath(path)

	decoder := yaml.NewYAMLOrJSONDecoder(f, 1000)
	for {
		var us unstructured.Unstructured
		// yamldec := yaml.NewYAMLReader(f)
		err = decoder.Decode(&us)
		if errors.Is(err, io.EOF) {
			break
			//return nil, err
		}
		if err != nil {
			return nil, err
		}
		//fmt.Println("Object: ", us)

		//var object *runtime.Object
		var object = us.DeepCopyObject()

		objects = append(objects, object)
	}

	mf.SetObjects(objects)

	//fmt.Println("Objects: ", objects)

	// Create a new package
	thePackage := Package{
		Name:      path,
		Namespace: "default",
		Resources: []Resource{&mf},
	}

	return &thePackage, nil
}

func ParsePodfile(path string) (*Package, error) {
	// fmt.Println("Parsing file: ", path)
	spec := hcldec.ObjectSpec{
		"podinate":       podinateHCLSpec,
		"pods":           podHCLSpec,
		"shared_volumes": SharedVolumeHCLSpec,
	}
	parser := hclparse.NewParser()

	//var val cty.Value
	var diags hcl.Diagnostics

	//path := paths[0]

	//for _, path := range paths {
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		WriteDiagnostics(diags, parser)
		return nil, errors.New("Error parsing file")
	}

	// Needs a lot of work from here down...
	// Need to decode the blocks one by one so we can handle errors better
	// And so we can get values from the blocks to add to others

	val, moreDiags := hcldec.Decode(f.Body, spec, nil)
	diags = append(diags, moreDiags...)
	//fmt.Printf("NEWVAL: %#v\n", val)

	//var err error

	// val, err = stdlib.Merge(val, newval)
	// if err != nil {
	// 	return nil, err
	// }
	//val.Add(newval)

	//}
	if diags.HasErrors() {
		WriteDiagnostics(diags, parser)
		return nil, errors.New("Error decoding file")
	}

	// Create a new package
	var thePackage Package

	//fmt.Printf("WHAT IS THE FUCKING TYPE %#v\n", val.Type())
	//fmt.Printf("Podinate is %#v\n", val.GetAttr("podinate").AsValueMap())
	// key, err := gocty.ToCtyValue("podinate", cty.String)
	// if err != nil {
	// 	logrus.Error(err)
	// }
	// valmap := val.AsValueMap()
	// _, ok := valmap["podinate"]

	pBlock := val.GetAttr("podinate")
	ok := !val.GetAttr("podinate").IsNull()
	if ok {
		logrus.Trace("Podinate block found")
		namespace := pBlock.GetAttr("namespace").AsString()
		packageName := pBlock.GetAttr("package").AsString()

		thePackage.Namespace = namespace
		thePackage.Name = packageName
	} else {
		thePackage.Name = path
		ns, err := kube_client.GetDefaultNamespace()
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"package": thePackage,
				"error":   err,
			})
		}
		thePackage.Namespace = ns
		logrus.Trace("No Podinate block found")
	}

	sharedVolumeValues := val.GetAttr("shared_volumes").AsValueMap()
	for i, sharedVolumeV := range sharedVolumeValues {
		var sharedVolume SharedVolume
		sharedVolume.ID = i
		err := gocty.FromCtyValue(sharedVolumeV, &sharedVolume)
		if err != nil {
			fmt.Println(err)
		}

		if sharedVolume.Namespace == nil {
			sharedVolume.Namespace = &thePackage.Namespace
		}
		thePackage.Resources = append(thePackage.Resources, sharedVolume)
	}
	// Parse all the pods in the file
	podvalues := val.GetAttr("pods").AsValueMap()
	//var pods []Pod
	for i, podv := range podvalues {
		var pod Pod
		pod.ID = i
		err := gocty.FromCtyValue(podv, &pod)
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Error parsing pod %s, issue: %s", i, err))
		}

		if pod.Namespace == nil {
			pod.Namespace = &thePackage.Namespace
		}
		//thePackage.Pods = append(thePackage.Pods, pod)
		thePackage.Resources = append(thePackage.Resources, pod)
	}

	// Parse all the shared volumes in the file

	// Commenting out - concentrating on pods for now

	thePackage.Labels = StandardLabels
	thePackage.Labels["app.kuberneretes.io/part-of"] = thePackage.Name

	return &thePackage, nil
}

// Apply takes a Package and makes it the current state of Podinate
func (pkg *Package) Apply(ctx context.Context, delete bool) error {

	var plan *Plan
	var err error
	if delete {
		plan, err = pkg.PlanDelete(ctx)
	} else {
		plan, err = pkg.Plan(ctx)
	}
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error":     err,
			"package":   pkg,
			"resources": pkg.Resources,
		}).Fatal("Failed to plan changes")
		os.Exit(1)
	}

	//fmt.Printf("Plan: %+v\n", plan)

	// Display the plan to the user and ask for confirmation
	// If the user confirms, apply the plan
	// If the user cancels, exit

	err = plan.Display(ctx)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to display plan")
		return err
	}

	if !plan.Applied {

		fmt.Print("Are you sure you want to apply these changes? (y/N)")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" || strings.ToLower(response) == "yes" {
			// Apply the plan
			err = plan.Apply(ctx)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to apply changes")
				return err
			}
		} else {
			fmt.Println()
			fmt.Println(tui.StyleSuccess.Render("Changes not applied"))
			fmt.Println()
		}
	}

	return nil

}

// Export takes a Package and exports it to Kubernetes YAML
func (pkg *Package) Export(ctx context.Context) (*string, error) {
	var out string
	count := 0

	for _, resource := range pkg.Resources {

		if viper.GetBool("debug") {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"resource": resource,
				"type":     resource.GetType(),
				"name":     resource.GetName(),
			}).Trace("Exporting resource")
		}

		objects, err := resource.GetObjects(ctx)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Failed to get objects")
			return nil, err
		}

		for _, object := range objects {
			if count > 0 {
				out += "---\n"
			}
			printing, err := helpers.ObjectToYAML(object)
			if err != nil {
				return nil, err
			}
			out += printing
			count++

		}

	}

	return &out, nil
}

// WriteDiagnostic writes the diagnostics to stdout
func WriteDiagnostics(diags hcl.Diagnostics, parser *hclparse.Parser) {
	wr := hcl.NewDiagnosticTextWriter(os.Stderr, parser.Files(), 80, true)
	// Handle errors
	wr.WriteDiagnostics(diags)
}
