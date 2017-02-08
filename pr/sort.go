package pr

import (
	"container/list"
	"fmt"

	"github.com/google/go-github/github"
)

// SortByLandingOrder sorts the given list of PRs into the order which they
// should be landed in.
//
// Given a base and a head, the landing order for a list of PRs in between is:
// the first PR off of base, then the next PR off of the head of that PR, and
// so on.
func SortByLandingOrder(base, head string, prs []*github.PullRequest) error {
	// TODO: how do we handle multiple PRs with the same base?

	// map of head to PR with that head
	branches := make(map[string]*github.PullRequest)
	for _, pr := range prs {
		branches[*pr.Head.Ref] = pr
	}

	// Work our way backwards from the head.
	var order list.List
	for currentHead := head; currentHead != base; {
		pr, ok := branches[currentHead]
		if !ok {
			// TODO: mention what we do have in the error
			return fmt.Errorf("could not find a PR with HEAD %q", currentHead)
		}
		order.PushFront(pr)
		currentHead = *pr.Base.Ref
	}

	i := 0
	for e := order.Front(); e != nil; e = e.Next() {
		prs[i] = e.Value.(*github.PullRequest)
		i++
	}

	return nil
}
