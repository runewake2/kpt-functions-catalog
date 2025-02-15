

// Code generated by "mdtogo"; DO NOT EDIT.
package generated

var ApplySettersShort = `Apply setter values on resource fields. Setters serve as parameters for template-free
setting of field values.`
var ApplySettersLong = `
Setters are a safer alternative to other substitution techniques which do not
have the context of the structured data. Setters may be invoked to modify the
package resources using ` + "`" + `apply-setters` + "`" + ` function to set desired values.

We use ConfigMap to configure the ` + "`" + `apply-setters` + "`" + ` function. The desired setter
values are provided as key-value pairs using ` + "`" + `data` + "`" + ` field where key is the name of the
setter and value is the new desired value for the setter.

  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: my-func-config
  data:
    setter_name1: setter_value1
    setter_name2: setter_value2
`
var ApplySettersExamples = `
Setting scalar values:

Let's start with the input resource in a package

  apiVersion: v1
  kind: Deployment
  metadata:
    name: nginx-deployment # kpt-set: ${image}-deployment
  spec:
    replicas: 1 # kpt-set: ${replicas}

Discover the names of setters in the function config file and declare desired values.

  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: apply-setters-fn-config
  data:
    image: ubuntu
    replicas: "3"

Render the declared values by invoking:

  $ kpt fn run --image gcr.io/kpt-fn/apply-setters:unstable --fn-config ./apply-setters-fn-config

Alternatively, setter values can be passed as key-value pairs in the CLI

  $ kpt fn run --image gcr.io/kpt-fn/apply-setters:unstable -- 'image=ubuntu' 'replicas=3'

Rendered resource looks like the following:

  apiVersion: v1
  kind: Deployment
  metadata:
    name: ubuntu-deployment # kpt-set: ${image}-deployment
  spec:
    replicas: 3 # kpt-set: ${replicas}

Setting array values:

Array values can also be parameterized using setters. Since the values of configMap
in pipeline definition must be of string type, the array values must be wrapped into
string. However, the rendered values in the resources will be array type.

Let's start with the input resource

  apiVersion: v1
  kind: MyKind
  metadata:
    name: foo
  environments: # kpt-set: ${env}
    - dev
    - stage

Declare the desired array values, wrapped into string.

  apiVersion: v1
  kind: ConfigMap
  metadata:
    name: apply-setters-fn-config
  data:
    env: |
      - prod
      - dev

Render the declared values by invoking:

  $ kpt fn run --image gcr.io/kpt-fn/apply-setters:unstable --fn-config ./apply-setters-fn-config

Rendered resource looks like the following:

  apiVersion: v1
  kind: MyKind
  metadata:
    name: foo
  environments: # kpt-set: ${env}
    - prod
    - dev
`
