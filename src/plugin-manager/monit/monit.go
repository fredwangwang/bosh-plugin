package monit

import (
	"fmt"
	"github.com/pkg/errors"
	"log"
	"os/exec"
	"time"
)

func cmdRetry(numRetry int, name string, arg ...string) error {
	var err error
	for i := 0; i < numRetry; i++ {
		cmd := exec.Command(name, arg...)
		if err = cmd.Run(); err != nil {
			log.Println(errors.Wrap(err, fmt.Sprintf("%s %v failed, retrying", name, arg)))
			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
	return err
}

func Reload() error {
	log.Println("monit reload")
	if err := cmdRetry(3, "/var/vcap/bosh/bin/monit", "reload"); err != nil {
		return errors.Wrap(err, "failed to reload monit config")
	}
	return nil
}

func Start(process string) error {
	log.Println("monit start", process)
	if err := cmdRetry(3, "/var/vcap/bosh/bin/monit", "start", process); err != nil {
		return errors.Wrap(err, "failed to monit start "+process)
	}
	return nil
}

func Stop(process string) error {
	log.Println("monit stop", process)
	if err := cmdRetry(3, "/var/vcap/bosh/bin/monit", "stop", process); err != nil {
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
