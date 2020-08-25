package outputs

import "github.com/eth0izzle/shhgit/internal/types"

type Publisher interface {
	Publish(m types.Match) error
}
