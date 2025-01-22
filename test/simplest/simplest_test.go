package simplest

import (
	"github.com/hash-d/frame2/pkg/frames/f2skupper1"
	"github.com/hash-d/frame2/pkg/frames/f2skupper1/disruptor"
	"github.com/hash-d/frame2/pkg/frames/f2skupper1/f2sk1environment"
	"testing"

	"gotest.tools/assert"

	frame2 "github.com/hash-d/frame2/pkg"
	"github.com/hash-d/frame2/pkg/frames/f2k8s"
)

func TestSimplestTemplate(t *testing.T) {
	r := &frame2.Run{
		T: t,
	}
	defer r.Finalize()

	r.AllowDisruptors([]frame2.Disruptor{
		&disruptor.MixedVersionVan{},
		&disruptor.NoConsole{},
		&disruptor.NoFlowCollector{},
		&disruptor.NoHttp{},
		&disruptor.UpgradeAndFinalize{},
		&disruptor.SkipManifestCheck{},
	})

	env := f2sk1environment.JustSkupperSimple{
		Name:         "simplest",
		AutoTearDown: true,
		Console:      true,
	}

	setup := frame2.Phase{
		Runner: r,
		//Name:   "skupper-simplest",
		Doc: "Simplest Skupper deploy: it's just skupper on the topology, nothing else",
		Setup: []frame2.Step{
			{
				//Name:   "Deploy Skupper",
				Modify: &env,
			},
		},
	}
	assert.Assert(t, setup.Run())

	// Get the namespaces
	prv, err := env.Topo.Get(f2k8s.Private, 1)
	assert.Assert(t, err)
	pub, err := env.Topo.Get(f2k8s.Public, 1)
	assert.Assert(t, err)

	main := frame2.Phase{
		Runner: r,
		Name:   "Check skupper",
		Doc:    "Here go the steps of the actual test",
		MainSteps: []frame2.Step{
			{
				ValidatorFinal: true,
				Validators: []frame2.Validator{
					&f2skupper1.CliSkupper{
						F2Namespace: pub,
						Args:        []string{"version"},
					},
					&f2skupper1.CliSkupper{
						F2Namespace: pub,
						Args:        []string{"status"},
					},
					&f2skupper1.CliSkupper{
						F2Namespace: pub,
						Args:        []string{"network", "status"},
					},
					&f2skupper1.CliSkupper{
						F2Namespace: prv,
						Args:        []string{"version"},
					},
					&f2skupper1.CliSkupper{
						F2Namespace: prv,
						Args:        []string{"status"},
					},
					&f2skupper1.CliSkupper{
						F2Namespace: prv,
						Args:        []string{"network", "status"},
					},
				},
			},
		},
	}
	assert.Assert(t, main.Run())

	// Teardown: for the template, all tear down is automatic.
	// If specific tear downs from the main steps are required,
	// create a new phase and specify them.

}
