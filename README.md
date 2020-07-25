# pkustomize

[PKS] & [kustomize] in harmony in a [concourse] pipeline

[PKS]: //docs.pivotal.io/pks
[kustomize]: //kustomize.io/
[concourse]: //concourse-ci.org/

## Proposed directory layout

The [example pipeline][pipe] has expects a certain directory structure:

[pipe]: pipeline/example.yaml


<pre>
<span style="font-weight:bold;color:blue;">.</span>
├── <span style="font-weight:bold;color:blue;">clusters</span>
│   ├── <span style="font-weight:bold;color:blue;">cluster-1</span>
│   │   ├── kustomization.yaml
│   │   ├── <span style="font-weight:bold;color:blue;">ns1</span>
│   │   │   ├── kustomization.yaml
│   │   │   ├── namespace.yaml
│   │   │   └── vcap-services.yaml
│   │   ├── <span style="font-weight:bold;color:blue;">ns2</span>
│   │   │   ├── kustomization.yaml
│   │   │   └── namespace.yaml
│   │   └── <span style="font-weight:bold;color:blue;">ns3</span>
│   │       ├── kustomization.yaml
│   │       └── namespace.yaml
│   └── <span style="font-weight:bold;color:blue;">cluster-2</span>
│       ├── kustomization.yaml
│       └── pks.yaml
├── <span style="font-weight:bold;color:blue;">kustomize</span>
│   └── <span style="font-weight:bold;color:blue;">plugin</span>
│       └── <span style="font-weight:bold;color:blue;">generators.hoegaarden.github.com</span>
│           └── <span style="font-weight:bold;color:blue;">v1alpha1</span>
│               └── <span style="font-weight:bold;color:blue;">vcapservices</span>
│                   ├── main.go
│                   └── <span style="font-weight:bold;color:green;">VcapServices</span>
├── <span style="font-weight:bold;color:blue;">pipeline</span>
│   └── example.yaml
├── README.md
└── <span style="font-weight:bold;color:blue;">shared</span>
    ├── <span style="font-weight:bold;color:blue;">ci-deployer</span>
    │   ├── ci-deployer.yaml
    │   └── kustomization.yaml
    ├── <span style="font-weight:bold;color:blue;">groups</span>
    │   ├── admin.yaml
    │   ├── exec.yaml
    │   ├── kustomization.yaml
    │   └── read-only.yaml
    └── <span style="font-weight:bold;color:blue;">monitoring</span>
        ├── kustomization.yaml
        ├── namespace.yaml
        ├── prom-config.yaml
        ├── prom-deploy.yaml
        └── prom-rbac.yaml

16 directories, 25 files
</pre>


All the stuff in this repo are just examples, nothing is to be considered
"production ready, nothing is to be considered "production ready".

### all clusters in `./clusters/...`

The top level config files for clusters reside at `./clusters/${clusterName}/` and are:

#### `pks.yaml`

Example:
```yaml
---
nrOfNodes: 4                           # default: 1
externalHostname: some.host.name.tld   # default: ${clusterName}.local
plan: large                            # small
```

If the file does not exist or does not hold one of the expected configs, the
default value is used. The cluster name itself is the name of the directory
itself.

If a cluster does not exist yet, it will be created with all the above settings.

If a cluster does exist, it will be scaled (the number of the cluster's worker
nodes that is) in case the expected number of workers from `pks.yaml` does not
match the currently deployed number of worker nodes. For now no other settings
for a cluster are changed once it has been created.

Clusters that exist but have no configuration in `./clusters/` will no be
touched, i.e. they won't be deleted.

#### `kustomize.yaml`

The main [kustomize] configuration for the whole cluster. Everything that needs
to be created in or deployed to the cluster needs needs to be created/managed
with [kustomize].

In the example the clusters have some cluster specific configuration in
`./clusters/${clusterName}/...` and use some shared configuration (shared
"bases") from `./shared/...`.

### `./shared/...` bases

Things that are shared across multiple clusters in that environment can be
pulled in as a base from `./shared/...`.

### plugins in `./kustomize/plugin/...`

`XDG_CONFIG_HOME` will be set to `./kustomize`. This allows us to place plugins in
there to be picked up by kustomize.

There is one example plugin `generators.hoegaarden.github.com/v1alpha1/VcapServices`:

It mimics a usecase I saw at a customer, where they had to pull in some data
from a Cloud Foundry Foundation and make that available to the cluster in a
kubernetes secret. This example plugin shows how this could be done. The
example plugin does not really talk to any external system, but it should be
easy enough to imagine how this could be implemented.

The example plugin is implemented as a [exec plugin][ep] which just happens to
be running go via `go run`. It could have also been implemented as a [go
plugin][gp].

[gp]: //kubernetes-sigs.github.io/kustomize/guides/plugins/#go-plugins
[ep]: //kubernetes-sigs.github.io/kustomize/guides/plugins/#exec-plugins

## Image

The image the example pipeline uses is [hhoerl/ops] from [hoegaarden/img/ops], but any image can be used that brings:
- `bash`
- `kubectl`
- `kustomize` in a recent version that has support for plugins
- `golang` for `go run`ning the plugin

[hhoerl/ops] runs as a non-root users. To have access to concourse resource
directories there are some hacks in place where we do a `sudo chown ...` on
some directories. If you use an image that is running as `root` or does not
need write access to resource or cache directories this is probably not needed.

[hhoerl/ops]: //hub.docker.com/repository/docker/hhoerl/ops
[hoegaarden/img/ops]: //github.com/hoegaarden/img/tree/master/ops
