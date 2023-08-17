# function-runtime-exec
A Composition Function runtime that wraps a binary

## Usage

```sh
$ go run cmd/function-runtime-exec/main.go -d start -- jq '.desired.composite.resource.labels |= {"labelizer.xfn.crossplane.io/processed": "true"} + .'
2023-08-23T09:33:43+02:00       DEBUG   function-runtime-exec   Listening       {"network": "tcp", "address": "0.0.0.0:1234", "command": "jq", "args": [".desired.composite.resource.labels |= {\"labelizer.xfn.crossplane.io/processed\": \"true\"} + ."]}
2023-08-23T09:33:46+02:00       DEBUG   function-runtime-exec   Running {"command": "jq", "args": [".desired.composite.resource.labels |= {\"labelizer.xfn.crossplane.io/processed\": \"true\"} + ."]}
2023-08-23T09:33:46+02:00       DEBUG   function-runtime-exec   Ran     {"command": "jq", "args": [".desired.composite.resource.labels |= {\"labelizer.xfn.crossplane.io/processed\": \"true\"} + ."], "stdout": "{\n  \"desired\": {\n    \"composite\": {\n      \"resource\": {\n        \"something\": \"something\",\n        \"labels\": {\n          \"labelizer.xfn.crossplane.io/processed\": \"true\"\n        }\n      }\n    }\n  }\n}\n", "stderr": ""}
```

```sh
grpcurl -plaintext -d @ localhost:1234 apiextensions.fn.proto.v1beta1.FunctionRunnerService.RunFunction <<EOM                                                                                                                                          130 ↵
{
  "desired": {
    "composite": {
      "resource": {
        "something": "something"
      }
    }
  }
}
EOM
{
  "desired": {
    "composite": {
      "resource": {
          "labels": {
                "labelizer.xfn.crossplane.io/processed": "true"
              },
          "something": "something"
        }
    }
  }
}
```