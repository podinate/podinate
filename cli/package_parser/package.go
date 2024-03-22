package package_parser

import (
	"errors"
	"fmt"
	"os"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/hashicorp/hcl2/hcldec"
	"github.com/hashicorp/hcl2/hclparse"
	"github.com/johncave/podinate/cli/sdk"
	"github.com/zclconf/go-cty/cty/gocty"
)

// Project represents a package to be installed, either the current state or the desired state
type Package struct {
	Projects []Project
	Pods     []Pod
}

func Parse(path string) (*Package, error) {
	fmt.Println("Parsing file: ", path)
	spec := hcldec.ObjectSpec{
		"pods":     podHCLSpec,
		"projects": projectHCLSpec,
	}
	parser := hclparse.NewParser()
	f, diags := parser.ParseHCLFile(path)
	if diags.HasErrors() {
		WriteDiagnostics(diags, parser)
		return nil, errors.New("Error parsing file")
	}
	val, moreDiags := hcldec.Decode(f.Body, spec, nil)
	diags = append(diags, moreDiags...)
	if diags.HasErrors() {
		WriteDiagnostics(diags, parser)
		return nil, errors.New("Error decoding file")
	}

	//var projects []Project

	var thePackage Package
	// Parse all the projects in the file
	for i, projectIn := range val.GetAttr("projects").AsValueMap() {
		var project Project
		project.ID = i
		err := gocty.FromCtyValue(projectIn, &project)
		if err != nil {
			fmt.Println(err)
		}
		thePackage.Projects = append(thePackage.Projects, project)
	}
	//fmt.Printf("Projects: %+v\n", projects)

	// Parse all the pods in the file
	podvalues := val.GetAttr("pods").AsValueMap()
	//var pods []Pod
	for i, podv := range podvalues {
		var pod Pod
		pod.ID = i
		err := gocty.FromCtyValue(podv, &pod)
		if err != nil {
			fmt.Println(err)
		}
		thePackage.Pods = append(thePackage.Pods, pod)
	}

	return &thePackage, nil
}

// Apply takes a Package and makes it the current state of Podinate
func (p *Package) Apply() error {
	// fmt.Printf("Pods: %+v\n", pods)

	//fmt.Println("Parsed file successfully! Let's start deploying!")
	stackProjects := make(map[string]*sdk.Project)
	for _, project := range p.Projects {
		//fmt.Printf("Deploying project: %s\n", project.Name)

		// Check if project exists, update if so
		existing, sdkerr := sdk.GetProjectByID(project.ID)
		if sdkerr == nil {
			// Project exists - try update it
			_, err := existing.Update(project.ToSDK())
			if err != nil {
				return err
			}
			fmt.Println("Updated project: ", project.Name)

			stackProjects[project.ID] = existing
			continue
		}

		var err error
		stackProjects[project.ID], err = sdk.CreateProject(project.ID, project.Name)
		if err != nil {
			return err
		}
		fmt.Println("Created project: ", project.Name)

	}

	for _, pod := range p.Pods {
		//fmt.Printf("Deploying pod: %s\n", pod.Name)
		theProject, err := sdk.GetProjectByID(pod.ProjectID)
		if err != nil {
			return err
		}
		thePod, err := pod.ToSDK()
		if err != nil {
			return err
		}

		// Check if pod exists, update if so
		existing, sdkerr := theProject.GetPodByID(thePod.ID)
		if sdkerr == nil {
			// Pod exists - try update it

			//fmt.Println("pod exists - updating", sdkerr, sdkerr == nil, existing)
			//os.Exit(2)
			fmt.Println("Updated pod: ", thePod.Name)

			_, err := existing.Update(thePod)
			if err != nil {
				return err
			}
			continue
			//	}
			//fmt.Printf("Creating pod: %+v %s\n", thePod, sdkerr.Error())
		} else if sdkerr.Error() == "404: Pod not found" {
			//fmt.Println("Error getting pod", sdkerr, sdkerr == nil, existing)
			//fmt.Printf("Created pod: %+v\n", thePod)
			_, err = theProject.CreatePod(*thePod)
			if err != nil {
				return err
			}
			fmt.Println("Created pod: ", thePod.Name)
		}

	}

	fmt.Println("Stack deployed!")
	return nil

}

// Delete takes a Package and deletes it from Podinate
func (p *Package) Delete() error {
	// fmt.Printf("Pods: %+v\n", pods)

	fmt.Println("Deleting pods...")

	for _, pod := range p.Pods {
		fmt.Printf("Deleting pod: %s\n", pod.Name)
		theProject, err := sdk.GetProjectByID(pod.ProjectID)
		if err != nil {
			return err
		}
		thePod, err := pod.ToSDK()
		if err != nil {
			return err
		}

		// Check if pod exists, update if so
		existing, sdkerr := theProject.GetPodByID(thePod.ID)
		if sdkerr == nil {
			// Pod exists - try update it
			err := existing.Delete()
			if err != nil {
				return err
			}
			continue
		}

	}

	fmt.Println("Deleting projects...")

	stackProjects := make(map[string]*sdk.Project)
	for _, project := range p.Projects {
		fmt.Printf("Deleting project: %s\n", project.Name)

		// Check if project exists, update if so
		existing, sdkerr := sdk.GetProjectByID(project.ID)
		if sdkerr == nil {
			// Project exists - try update it
			err := existing.Delete()
			if err != nil {
				return err
			}

			stackProjects[project.ID] = existing
			continue
		}

	}

	fmt.Println("Stack deleted!")
	return nil

}

// WriteDiagnostic writes the diagnostics to the given writer
func WriteDiagnostics(diags hcl.Diagnostics, parser *hclparse.Parser) {
	wr := hcl.NewDiagnosticTextWriter(os.Stdout, parser.Files(), 80, true)
	// Handle errors
	wr.WriteDiagnostics(diags)
}
