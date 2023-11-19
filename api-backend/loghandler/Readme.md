# loghandler
This package handles a bunch of stuff relating to logs, if something goes wrong, give the requestor a vague error message, and write as much context as possible to this package. It should be used like so,

## Use 
The primary use is the logging built on Uber's [zap](https://github.com/uber-go/zap) library. 