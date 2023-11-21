package cmd

import (
	"errors"
	"flag"
	"fmt"
	"github.com/dddplayer/hugoverse/internal/application"
	"github.com/dddplayer/hugoverse/pkg/log"
	"os"
)

type buildCmd struct {
	parent       *flag.FlagSet
	cmd          *flag.FlagSet
	hugoProjPath *string
}

func NewBuildCmd(parent *flag.FlagSet) (*buildCmd, error) {
	nCmd := &buildCmd{
		parent: parent,
	}

	nCmd.cmd = flag.NewFlagSet("build", flag.ExitOnError)
	nCmd.hugoProjPath = nCmd.cmd.String("p", "", fmt.Sprintf(
		"[required] target hugo project path \n(e.g. %s)", "path/to/your/hugo/project"))
	err := nCmd.cmd.Parse(parent.Args()[1:])
	if err != nil {
		return nil, err
	}

	return nCmd, nil
}

func (oc *buildCmd) Usage() {
	oc.cmd.Usage()
}

func (oc *buildCmd) Run() error {
	if *oc.hugoProjPath == "" {
		oc.cmd.Usage()
		return errors.New("please specify a target hugo project path")
	}

	_, err := os.Stat(*oc.hugoProjPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("directory %s does not exist", *oc.hugoProjPath)
	}

	if err != nil {
		return err
	}

	l := log.NewStdLogger()
	if err = application.GenerateStaticSite(*oc.hugoProjPath); err != nil {
		l.Fatalf("failed to generate static site: %v", err)
		return err
	}

	return nil
}
