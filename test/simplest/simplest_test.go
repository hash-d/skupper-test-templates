package simplest

import (
	"testing"

	"gotest.tools/assert"

	frame2 "github.com/hash-d/frame2/pkg"
	"github.com/hash-d/frame2/pkg/disruptors"
	"github.com/hash-d/frame2/pkg/environment"
	"github.com/hash-d/frame2/pkg/skupperexecute"
	"github.com/hash-d/frame2/pkg/topology"
)

func TestSimplestTemplate(t *testing.T) {
	r := &frame2.Run{
		T: t,
	}
	defer r.Finalize()

	r.AllowDisruptors([]frame2.Disruptor{
		&disruptors.MixedVersionVan{},
		&disruptors.NoConsole{},
		&disruptors.NoFlowCollector{},
		&disruptors.NoHttp{},
		&disruptors.UpgradeAndFinalize{},
		&disruptors.SkipManifestCheck{},
	})

	env := environment.JustSkupperSimple{
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
	prv, err := env.Topo.Get(topology.Private, 1)
	assert.Assert(t, err)
	pub, err := env.Topo.Get(topology.Public, 1)
	assert.Assert(t, err)

	main := frame2.Phase{
		Runner: r,
		Name:   "Check skupper",
		Doc:    "Here go the steps of the actual test",
		MainSteps: []frame2.Step{
			{
				ValidatorFinal: true,
				Validators: []frame2.Validator{
					&skupperexecute.CliSkupper{
						ClusterContext: pub,
						Args:           []string{"version"},
					},
					&skupperexecute.CliSkupper{
						ClusterContext: pub,
						Args:           []string{"status"},
					},
					&skupperexecute.CliSkupper{
						ClusterContext: pub,
						Args:           []string{"network", "status"},
					},
					&skupperexecute.CliSkupper{
						ClusterContext: prv,
						Args:           []string{"version"},
					},
					&skupperexecute.CliSkupper{
						ClusterContext: prv,
						Args:           []string{"status"},
					},
					&skupperexecute.CliSkupper{
						ClusterContext: prv,
						Args:           []string{"network", "status"},
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
