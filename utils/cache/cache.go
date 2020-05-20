package cache

import (
	"github.com/cornelk/hashmap"
)

var (
	// store bird server
	ServiceInformerCache    = hashmap.New(1000)

	BirdAccountCache = hashmap.New(1000)

	// store the details of the bird proxy services
	BirdCache            = hashmap.New(1000)
)


// todo fill birdAccountCache and watch delete and update password