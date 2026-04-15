package runner

import (
	"fmt"
	"log"

	"github.com/patchwork-deploy/internal/config"
	"github.com/patchwork-deploy/internal/patch"
	"github.com/patchwork-deploy/internal/state"
)

// Result holds the outcome of a single patch run.
type Result struct {
	PatchID string
	Applied bool
	Skipped bool
	Err     error
}

// Runner orchestrates patch discovery, state checking, and application.
type Runner struct {
	env     config.Environment
	loader  *patch.Loader
	applier *patch.Applier
	st      *state.State
}

// New creates a Runner for the given environment.
func New(env config.Environment, loader *patch.Loader, applier *patch.Applier, st *state.State) *Runner {
	return &Runner{env: env, loader: loader, applier: applier, st: st}
}

// Run discovers patches and applies any that have not yet been applied.
// It returns a slice of Results and the first fatal error encountered.
func (r *Runner) Run() ([]Result, error) {
	patches, err := r.loader.Discover(r.env.PatchDir)
	if err != nil {
		return nil, fmt.Errorf("discover patches: %w", err)
	}

	var results []Result
	for _, p := range patches {
		if r.st.Has(r.env.Name, p) {
			results = append(results, Result{PatchID: p, Skipped: true})
			continue
		}

		log.Printf("[%s] applying patch %s", r.env.Name, p)
		if err := r.applier.Apply(r.env, p); err != nil {
			results = append(results, Result{PatchID: p, Err: err})
			return results, fmt.Errorf("apply patch %s: %w", p, err)
		}

		results = append(results, Result{PatchID: p, Applied: true})
	}

	return results, nil
}
