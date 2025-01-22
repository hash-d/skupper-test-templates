package hello_world

import (
	"github.com/hash-d/frame2/pkg/frames/f2k8s/disruptor"
	"github.com/hash-d/frame2/pkg/frames/f2skupper1"
	disruptor2 "github.com/hash-d/frame2/pkg/frames/f2skupper1/disruptor"
	"github.com/hash-d/frame2/pkg/frames/f2skupper1/f2sk1environment"
	"testing"

	frame2 "github.com/hash-d/frame2/pkg"
	"github.com/hash-d/frame2/pkg/frames/f2k8s"
	"gotest.tools/assert"
)

func TestHelloWorldTemplate(t *testing.T) {
	r := &frame2.Run{
		T: t,
	}
	defer r.Finalize()

	r.AllowDisruptors([]frame2.Disruptor{
		&disruptor.PodSecurityAdmission{},
		&disruptor.PSADeployment{},
		&disruptor2.ConsoleAuth{},
		&disruptor2.ConsoleOnAll{},
		&disruptor2.FlowCollectorOnAll{},
		&disruptor2.UpgradeAndFinalize{},
		&disruptor2.SkipManifestCheck{},
		&disruptor2.EdgeOnPrivate{},
	})

	helloWorldDefault := &f2sk1environment.HelloWorldDefault{
		AutoTearDown: true,
	}

	setup := frame2.Phase{
		Runner: r,
		// Do not name setup phases and steps, if they use AutoTearDown, as
		// that will cause the AutoTearDown to be triggered at the end of
		// the named phase or step.
		//Name:   "Hello World setup",
		Doc: "Deploy Hello World on the default topology",
		Setup: []frame2.Step{
			{
				//Name: "Deploy Hello World",
				Doc: "Deploy Hello World",
				// As an alternative, check HelloWorldN for HelloWorld on
				// an N topology, or just HelloWorld for a fully configurable
				// environment.
				Modify: helloWorldDefault,
			},
		},
	}

	assert.Assert(t, setup.Run())

	topo := helloWorldDefault.GetTopology()
	pub1, err := topo.Get(f2k8s.Public, 1)
	if err != nil {
		t.Fatalf("Failed to get pub1: %v", err)
	}
	prv1, err := topo.Get(f2k8s.Private, 1)
	if err != nil {
		t.Fatalf("Failed to get prv1: %v", err)
	}

	main := frame2.Phase{
		Runner: r,
		Name:   "Replace me",
		Doc:    "Here goes the steps of the actual test",
		MainSteps: []frame2.Step{
			{
				Name: "Expose frontend",
				Modify: &f2skupper1.SkupperExpose{
					Namespace: pub1,
					Type:      "deployment",
					Name:      "frontend",
					Ports:     []int{8080},
				},
				Validators: []frame2.Validator{
					&f2k8s.Curl{
						Namespace: pub1,
						Url:       "frontend:8080",
					},
					&f2k8s.Curl{
						Namespace: prv1,
						Url:       "frontend:8080",
					},
				},
				ValidatorRetry: frame2.RetryOptions{
					Allow:  60,
					Ensure: 2,
				},
				ValidatorFinal: true,
			}, {
				Name: "Expose backend",
				Modify: &f2skupper1.SkupperExpose{
					Namespace: prv1,
					Type:      "deployment",
					Name:      "backend",
					Ports:     []int{8080},
				},
				Validators: []frame2.Validator{
					&f2k8s.Curl{
						Namespace: prv1,
						Url:       "backend:8080/api/hello",
					},
					&f2k8s.Curl{
						Namespace: pub1,
						Url:       "backend:8080/api/hello",
					},
				},
				ValidatorRetry: frame2.RetryOptions{
					Allow:  60,
					Ensure: 2,
				},
				ValidatorFinal: true,
			},
		},
	}

	assert.Assert(t, main.Run())

	// Teardown: for the template, all tear down is automatic.
	// If specific tear downs from the main steps are required,
	// create a new phase and specify them.

}
