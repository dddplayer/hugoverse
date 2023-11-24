package entity

import "github.com/dddplayer/hugoverse/pkg/lazy"

type HugoSitesInit struct {
	// Performs late initialization (before render) of the templates.
	Layouts *lazy.Init
}
