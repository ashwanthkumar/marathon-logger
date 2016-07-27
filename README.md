# marathon-logger

Marathon logger is a simple tool that's meant to be deployed on all the Mesos slaves and run as a daemon. It monitors the logs for each app task and pushes it to a logging backend (syslog for now).

## How does it work?
Each `marathon-logger` instance is configured to talk to a Marathon instance. It polls all the running apps periodically and finds the one's that has logging enabled. They are identified using a label `logs.enabled:true`. Then for each app's task it queries the local instance of the Mesos slave to find the working directory for the task. We then create a rsyslog config file for each task, tailing the files within the working directory again specified using the label `logs.files:"logs/access.log"` - if nothing is specified, we'll tail the `stdout` in the working directory. Once we create/modify the configuration, we reload the syslog daemon. This can also be used to push the logs to any backends like [Loggly](http://loggly.com) for example.

## Usage
```
$ marathon-logger --help
Usage of marathon-logger:
      --app-check-interval duration             Frequency at which we check for new tasks (default 30s)
      --rsyslog-configuration-dir string        Location on the Filesystem where the rsyslog configurations needs to be written (default "/etc/rsyslog.d")
      --rsyslog-restart-cmd string              Restart command for rsyslog backend (default "service rsyslog restart")
      --slave-port int                          Mesos slave port (default 5051)
      --task-max-heart-beat-interval duration   Max heartbeat interval after which the task is considered dead and logger is removed (default 30m0s)
      --uri string                              Marathon URI to connect
```

## App Labels
Apart from the flags that are used while starting up, the functionality can be controlled at an app level using labels in the app specification. The following table explains the properties and it's usage.

| Property | Description | Example |
| --- | --- | --- |
| logs.enabled | Enable or disable log monitoring for the app. Default - `false` | true |
| logs.files | List of files to monitor via the backend. Default - `stdout` | stdout,stderr |

## License
http://www.apache.org/licenses/LICENSE-2.0
