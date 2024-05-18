# This is PoC and test app for eatwo app

Proposed folder structrure

```
/myapp
  /cmd
    /myapp
      main.go
  /api
    (API definitions and protocol files)
  /internal
    /service
      (business logic implementation)
    /repository
      (database access code)
    /handler
      (API endpoint handlers)
  /pkg
    (library code that's ok to use by external applications)
  /scripts
    (scripts to perform various build, install, analysis, etc operations)
  /test
    (additional external test apps and test data)
  README.md
  .gitignore
  go.mod
  go.sum
```
