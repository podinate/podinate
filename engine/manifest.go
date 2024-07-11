package engine

import "k8s.io/apimachinery/pkg/runtime"

type Manifest struct {
	// Name is the name of the manifest
	name string
	// Path is the path to the manifest file
	path string
	// Objects is the list of objects in the manifest
	objects []runtime.Object
}

func (m *Manifest) GetPath() string {
	return m.path
}

func (m *Manifest) SetPath(path string) {
	m.path = path
}

func (m *Manifest) SetName(name string) {
	m.name = name
}

func (m *Manifest) GetName() string {
	return m.name
}

func (m *Manifest) GetType() ResourceType {
	return ResourceTypeManifest
}

func (m *Manifest) GetObjects() ([]runtime.Object, error) {
	return m.objects, nil
}

func (m *Manifest) SetObjects(objects []runtime.Object) {
	m.objects = objects
}
