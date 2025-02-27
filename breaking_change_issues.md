# This document contains issues and tasks which break backward compatibility.

## [ON HOLD] Incoming headers must be case-insensitive

Header names which we get from incoming data are case-sensitive.
This means that it is impossible to know about the presence of the header in incoming data,
because the exact format of the header cannot be known.

#### Solution:

Convert all header names in incoming data to lowercase before searching for headers.

## Changed log levels configuration format

The new format of log level configuration properties was introduced.

### Old approach (deprecated)
For env:
```yaml
LOG_LEVEL=debug
LOG_LEVEL_PACKAGE_DBAAS=warn
```

For yaml:
```yaml
log.level: debug
log.level.package:
  dbaas: warn
```

### New approach

For env:
```yaml
LOGGING_LEVEL_ROOT=debug
LOGGING_LEVEL_DBAAS=warn
```

For yaml:
```yaml
logging.level.root: debug
logging.level:
  dbaas: warn
```

The old approach is deprecated and will be removed according to deprecation policy.