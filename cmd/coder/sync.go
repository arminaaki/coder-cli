package main

import (
	"path/filepath"
	"strings"

	"github.com/spf13/pflag"

	"go.coder.com/cli"
	"go.coder.com/flog"

	"cdr.dev/coder-cli/internal/sync"
)

type syncCmd struct {
	init bool
}

func (cmd *syncCmd) Spec() cli.CommandSpec {
	return cli.CommandSpec{
		Name:  "sync",
		Usage: "[local directory] [<env name>:<remote directory>]",
		Desc:  "establish a one way directory sync to a remote environment",
	}
}

func (cmd *syncCmd) RegisterFlags(fl *pflag.FlagSet) {
	fl.BoolVarP(&cmd.init, "init", "i", false, "do initial transfer and exit")
}

func (cmd *syncCmd) Run(fl *pflag.FlagSet) {
	var (
		local  = fl.Arg(0)
		remote = fl.Arg(1)
	)
	if local == "" || remote == "" {
		exitUsage(fl)
	}

	entClient := requireAuth()

	remoteTokens := strings.SplitN(remote, ":", 2)
	if len(remoteTokens) != 2 {
		flog.Fatal("remote misformatted")
	}
	var (
		envName   = remoteTokens[0]
		remoteDir = remoteTokens[1]
	)

	env := findEnv(entClient, envName)

	absLocal, err := filepath.Abs(local)
	if err != nil {
		flog.Fatal("make abs path out of %v: %v", local, absLocal)
	}

	s := sync.Syncer{
		Environment: env,
		RemoteDir:   remoteDir,
		LocalDir:    absLocal,
	}

	err = s.Create()
	if err != nil {
		flog.Fatal("%v", err)
	}

	err = s.Monitor()
	if err != nil {
		flog.Fatal("%v", err)
	}
}
