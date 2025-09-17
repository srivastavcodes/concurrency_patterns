package main

import "os/exec"

type IntermediateLevelErr struct {
	error
}

func runJob(id string) error {
	const joinBinPath = "/bad/job/binary"

	isExecutable, err := isGloballyExec(joinBinPath)
	if err != nil {
		return IntermediateLevelErr{
			wrapError(err,
				"cannot run job %q: requisite binaries not available", id,
			),
		}
	} else if !isExecutable {
		return IntermediateLevelErr{
			wrapError(nil,
				"cannot run job %q: requisite binaries not available", id,
			),
		}
	}
	return exec.Command(joinBinPath, "--id="+id).Run()
}
