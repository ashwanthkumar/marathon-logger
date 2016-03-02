# marathon-logger

Marathon logger is a simple tool that's meant to be deployed on all the Mesos slaves and run as a daemon. It monitors the logs for each app task and pushes it to a logging backend (syslog for now).

## How does it work?
Each `marathon-logger` instance is configured to talk to a Marathon instance. It polls all the running apps periodically and finds the one's that has logging enabled. They are identified using a label `logs.enabled=true`. Then for each app's task it queries the local instance of the Mesos slave to find the working directory for the task. We then tail the files within the working directory again specified using the label `logs.files="logs/access.log"` - if nothing is specified. We'll tail the `stdout` in the working directory. We then push these to local syslog server which can then have any custom routing like pushing it to [Loggly](http://loggly.com) for example.

## Status
This tool is a WIP.

## License
http://www.apache.org/licenses/LICENSE-2.0
