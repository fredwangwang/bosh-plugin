package monit

import (
	"fmt"
	"github.com/pkg/errors"
	"os/exec"
)

func Reload() error {
	cmd := exec.Command("/var/vcap/bosh/bin/monit", "reload")
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to reload monit config")
	}
	return nil
}

func Start(process string) error {
	cmd := exec.Command("/var/vcap/bosh/bin/monit", "start", process)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to monit start "+process)
	}
	return nil
}

func Stop(process string) error {
	cmd := exec.Command("/var/vcap/bosh/bin/monit", "stop", process)
	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "failed to monit stop "+process)
	}
	return nil
}

func Monitrc(process string) string {
	return fmt.Sprintf(monitTmpl, process, process, process, process, process)
}

var monitTmpl = `check process %s
  with pidfile /var/vcap/sys/run/bpm/%s/%s.pid
  start program "/var/vcap/jobs/bpm/bin/bpm start %s"
  stop program "/var/vcap/jobs/bpm/bin/bpm stop %s" with timeout 60 seconds
  group vcap
`
