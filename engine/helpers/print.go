package helpers

import (
	"bytes"
	"context"
	"fmt"
	"os"

	"github.com/podinate/podinate/tui"
	"github.com/sirupsen/logrus"
	"github.com/sters/yaml-diff/yamldiff"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/cli-runtime/pkg/printers"
)

// PrintKubernetesValidationError prints a Kubernetes validation error to the console
// If there is no existing object, pass in nil
func PrintKubernetesValidationError(ctx context.Context, currentObject runtime.Object, newObject runtime.Object, theError error) error {
	fmt.Println(tui.StyleError.Render("The following change was rejected by Kubernetes:"))

	if currentObject == nil {
		err := PrintObject(newObject)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
			return err
		}
	} else {

		err := YamlDiffObjects(ctx, currentObject, newObject)
		if err != nil {
			logrus.WithContext(ctx).WithFields(logrus.Fields{
				"error": err,
			}).Error("A further error occurred when trying to display the change that was rejected by the Kubernetes API.")
			return err
		}

	}

	fmt.Println(tui.StyleError.Render("Reason:"), theError, "\n")
	return nil
}

func YamlDiffObjects(ctx context.Context, oldResource runtime.Object, newResource runtime.Object) error {
	y := printers.YAMLPrinter{}
	b := new(bytes.Buffer)

	// Strip the ManagedFields from the old resource
	// as they're not suitable for comparison
	err := stripManagedFields(oldResource)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error stripping managed fields from oldResource")
		return err
	}

	err = y.PrintObj(oldResource, b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error printing old object")
		return err
	}
	oldstr := b.String()
	b.Truncate(0)

	yamldiff.Load(oldstr)

	// Strip managed fields to make comparison easier
	err = stripManagedFields(newResource)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error stripping managed fields from newResource")
		return err
	}

	err = y.PrintObj(newResource, b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error":        err,
			"old_resource": oldResource,
			"new_resource": newResource,
		}).Error("Error printing new object")
		return err
	}
	newstr := b.String()

	oldYaml, err := yamldiff.Load(oldstr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading old YAML")
		return err
	}

	newYaml, err := yamldiff.Load(newstr)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error loading new YAML")
		return err
	}

	for _, diff := range yamldiff.Do(oldYaml, newYaml) {
		fmt.Println(diff.Dump())
	}

	return nil
}

// stripManagedFields removes the "managedFields" field from the object
func stripManagedFields(resource runtime.Object) error {
	// Strip ManagedFields from the old resource
	o, err := ObjectToUnstructured(resource)
	if err != nil {
		return err
	}
	o.SetManagedFields(nil)
	return runtime.DefaultUnstructuredConverter.FromUnstructured(o.Object, resource)
}

func PrintObject(object runtime.Object) error {
	y := printers.YAMLPrinter{}
	err := stripManagedFields(object)
	if err != nil {
		return err
	}
	err = y.PrintObj(object, os.Stdout)
	fmt.Println("Print", err)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error printing object")
		return err
	}
	return nil
}

func ObjectToYAML(object runtime.Object) (string, error) {
	y := printers.YAMLPrinter{}
	err := stripManagedFields(object)
	if err != nil {
		return "", err
	}
	var b bytes.Buffer
	err = y.PrintObj(object, &b)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error printing object")
		return "", err
	}
	return b.String(), nil
}
