package engine

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/sirupsen/logrus"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Ref: https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
var StandardLabels = map[string]string{
	"app.kubernetes.io/managed-by": "Podinate",
}

// Package represents a package to be installed, either the current state or the desired state
type Package struct {
	Name          string
	Namespace     string
	Pods          []Pod
	SharedVolumes []SharedVolume
	Labels        map[string]string
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
	// fmt.Println("Parsing file: ", path)
	spec := hcldec.ObjectSpec{
		"podinate":       podinateHCLSpec,
		"pods":           podHCLSpec,
		"shared_volumes": SharedVolumeHCLSpec,
	}
	parser := hclparse.NewParser()

	//var val cty.Value
	var diags hcl.Diagnostics

	val := cty.ObjectVal(make(map[string]cty.Value))

	fmt.Printf("Val start: %#v\n", val)

	path := paths[0]

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

	namespace := val.GetAttr("podinate").GetAttr("namespace").AsString()
	packageName := val.GetAttr("podinate").GetAttr("package").AsString()
	thePackage.Namespace = namespace
	thePackage.Name = packageName

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
		thePackage.Pods = append(thePackage.Pods, pod)
	}

	// Parse all the shared volumes in the file

	// Commenting out - concentrating on pods for now

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
		thePackage.SharedVolumes = append(thePackage.SharedVolumes, sharedVolume)
	}

	thePackage.Labels = StandardLabels
	thePackage.Labels["app.kuberneretes.io/part-of"] = thePackage.Name

	return &thePackage, nil
}

// Apply takes a Package and makes it the current state of Podinate
func (pkg *Package) Apply(ctx context.Context) error {

	plan, err := pkg.Plan(ctx)
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to plan changes")
		os.Exit(1)
	}

	//fmt.Printf("Plan: %+v\n", plan)

	// Display the plan to the user and ask for confirmation
	// If the user confirms, apply the plan
	// If the user cancels, exit

	err = plan.Display()
	if err != nil {
		logrus.WithContext(ctx).WithFields(logrus.Fields{
			"error": err,
		}).Fatal("Failed to display plan")
		return err
	}

	if !plan.Applied {

		fmt.Print("Are you sure you want to apply these changes? (Y/n)")
		var response string
		fmt.Scanln(&response)
		if strings.ToLower(response) == "y" || response == "" {
			// Apply the plan
			err = plan.Apply(ctx)
			if err != nil {
				logrus.WithContext(ctx).WithFields(logrus.Fields{
					"error": err,
				}).Fatal("Failed to apply changes")
				return err
			}
		}
	}

	// out, err := json.MarshalIndent(plan, "", "  ")
	// if err != nil {
	// 	return err
	// }
	// fmt.Println(string(out))

	// Deploy shared volumes
	// zap.S().Infow("Deploying shared volumes", "shared_volumes", p.SharedVolumes)
	// for _, sharedVolume := range p.SharedVolumes {
	// 	zap.S().Infow("Deploying shared volume", "shared_volume", sharedVolume)
	// 	theProject, err := sdk.GetProjectByID(sharedVolume.ProjectID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	theSharedVolume, err := sharedVolume.ToSDK()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	zap.S().Infow("Getting shared volume", "shared_volume", theSharedVolume)
	// 	existing, sdkerr := theProject.GetSharedVolumeByID(theSharedVolume.ID)
	// 	zap.S().Infow("Got shared volume", "shared_volume", theSharedVolume, "existing", existing, "sdkerr", sdkerr)
	// 	if sdkerr == nil {
	// 		// Shared volume exists - try update it
	// 		zap.S().Infow("Shared volume exists - updating", "shared_volume", theSharedVolume)
	// 		err := existing.Update(theSharedVolume)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		continue
	// 	} else if sdkerr.Code == 404 {
	// 		_, err = theProject.CreateSharedVolume(*theSharedVolume)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		fmt.Println("Created shared volume: ", theSharedVolume.ID)
	// 	} else {
	// 		return sdkerr
	// 	}
	// }

	// for _, pod := range p.Pods {
	// 	//fmt.Printf("Deploying pod: %s\n", pod.Name)
	// 	theProject, err := sdk.GetProjectByID(pod.ProjectID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	thePod, err := pod.ToSDK()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Check if pod exists, update if so
	// 	existing, sdkerr := theProject.GetPodByID(thePod.ID)
	// 	if sdkerr == nil {
	// 		// Pod exists - try update it

	// 		//fmt.Println("pod exists - updating", sdkerr, sdkerr == nil, existing)
	// 		//os.Exit(2)
	// 		fmt.Println("Updated pod: ", thePod.Name)

	// 		err := existing.Update(thePod)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		continue
	// 		//	}
	// 		//fmt.Printf("Creating pod: %+v %s\n", thePod, sdkerr.Error())
	// 	} else if sdkerr.Error() == "404: Pod not found" {
	// 		//fmt.Println("Error getting pod", sdkerr, sdkerr == nil, existing)
	// 		//fmt.Printf("Created pod: %+v\n", thePod)
	// 		_, err = theProject.CreatePod(*thePod)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		fmt.Println("Created pod: ", thePod.Name)
	// 	} else {
	// 		return sdkerr
	// 	}

	// }

	//fmt.Println("Not implemented!")
	return nil

}

// Delete takes a Package and deletes it from Podinate
func (p *Package) Delete() error {
	// fmt.Printf("Pods: %+v\n", pods)

	fmt.Println("Deleting pods...")

	// for _, pod := range p.Pods {
	// 	fmt.Printf("Deleting pod: %s\n", pod.Name)
	// 	theProject, err := sdk.GetProjectByID(pod.ProjectID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	thePod, err := pod.ToSDK()
	// 	zap.S().Debugw("Got pod to SDK", "pod", thePod, "err", err)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Check if pod exists, update if so
	// 	existing, sdkerr := theProject.GetPodByID(thePod.ID)
	// 	zap.S().Debugw("Got pod", "pod", thePod, "existing", existing, "sdkerr", sdkerr)
	// 	if sdkerr == nil {
	// 		// Pod exists - try update it
	// 		err := existing.Delete()
	// 		zap.S().Debugw("Deleted pod", "pod", thePod, "err", err)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		continue
	// 	}

	// }

	// Delete shared volumes
	fmt.Println("Deleting shared volumes...")
	// for _, sharedVolume := range p.SharedVolumes {
	// 	fmt.Printf("Deleting shared volume: %s\n", sharedVolume.ID)
	// 	theProject, err := sdk.GetProjectByID(sharedVolume.ProjectID)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	theSharedVolume, err := sharedVolume.ToSDK()
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// Check if shared volume exists, and delete
	// 	existing, sdkerr := theProject.GetSharedVolumeByID(theSharedVolume.ID)
	// 	zap.S().Debugw("Got shared volume", "shared_volume", theSharedVolume, "existing", existing, "sdkerr", sdkerr)
	// 	if sdkerr == nil {
	// 		// Shared volume exists - try update it
	// 		err := existing.Delete()
	// 		zap.S().Debugw("Deleted shared volume", "shared_volume", theSharedVolume, "err", err)
	// 		if err != nil {
	// 			return err
	// 		}
	// 		continue
	// 	}
	// }

	fmt.Println("Not implemented!")
	return nil

}

// WriteDiagnostic writes the diagnostics to stdout
func WriteDiagnostics(diags hcl.Diagnostics, parser *hclparse.Parser) {
	wr := hcl.NewDiagnosticTextWriter(os.Stdout, parser.Files(), 80, true)
	// Handle errors
	wr.WriteDiagnostics(diags)
}
