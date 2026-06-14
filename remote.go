package git

import (
	"github.com/go-git/go-git/v5/plumbing/protocol/packp"
	"github.com/go-git/go-git/v5/plumbing/storer"
)

// updateShallow updates the shallow boundaries in the storer based on the
// shallow and unshallow lines received from the server.
func (r *Remote) updateShallow(s storer.Storer, shallow []string, unshallow []string) error {
	current, err := s.Shallow()
	if err != nil {
		return err
	}

	// Apply unshallow: remove commits that are no longer shallow
	for _, u := range unshallow {
		for i, c := range current {
			if c.String() == u {
				current = append(current[:i], current[i+1:]...)
				break
			}
		}
	}

	// Apply shallow: add new shallow commits
	for _, sh := range shallow {
		found := false
		for _, c := range current {
			if c.String() == sh {
				found = true
				break
			}
		}
		if !found {
			current = append(current, plumbing.NewHash(sh))
		}
	}

	return s.SetShallow(current)
}