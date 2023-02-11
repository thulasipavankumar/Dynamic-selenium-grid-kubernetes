package models

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/lithammer/shortuuid/v4"
)

type NamespaceDetails struct {
	Namespace string
	Url       string
	Token     string
}

var namespace NamespaceDetails

func init() {
	err := godotenv.Load("../keys.env")
	_ = err
	namespace = NamespaceDetails{Namespace: os.Getenv("Namespace"), Url: os.Getenv("Url"), Token: os.Getenv("Token")}
}

type template interface {
	Deploy()
	SetValues()
	GetName() string
	Init(c Common)
}
type Deployment struct {
	pod     Pod
	service Service
	ingress Ingress
}
type Common struct {
	App    string
	EnvArr []Env
	Port   int
}
type Details struct {
	PodName     string
	ServiceName string
	IngressName string
}

func (d *Deployment) GetDetails() Details {
	return Details{
		PodName:     d.pod.GetName(),
		ServiceName: d.service.GetName(),
		IngressName: d.ingress.GetName(),
	}
}
func (d *Deployment) LoadRequestedCapabilites(matched Match) error {
	err := d.pod.PopulateImagesFronRequest(matched)
	if err != nil {
		return err
	}
	return nil
}
func (d *Deployment) GetService() Service {
	return d.service
}
func (d *Deployment) GetPod() Pod {
	return d.pod
}
func (d *Deployment) GetIngress() Ingress {
	return d.ingress
}
func (d *Deployment) Deploy() error {
	c := Common{App: "app-" + strings.ToLower(shortuuid.New()), EnvArr: nil, Port: 4444}
	d.Init(c)
	serviceErr := d.deployService()
	podErr := d.deployPod()
	//ingerssErr := d.deployIngress()
	_, _ = serviceErr, podErr
	if serviceErr != nil || podErr != nil {
		return fmt.Errorf("Unable to create deployment")
	}
	return nil
}

func (d *Deployment) Init(c Common) {
	d.pod.Init(c, namespace)
	d.service.Init(c, namespace)
	d.ingress.Init(c, namespace)

	log.Printf("pod: %#v\nservice: %#v\ningress: %#v\n", d.pod, d.service, d.ingress)

}
func (d *Deployment) deployPod() error {
	return d.pod.Deploy()
}
func (d *Deployment) deployService() error {
	return d.service.Deploy()
}
func (d *Deployment) deployIngress() error {
	return d.ingress.Deploy()
}

func (i *Deployment) setValues() {

}
