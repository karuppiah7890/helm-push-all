module github.com/karuppiah7890/helm-push-all

go 1.15

require (
	github.com/chartmuseum/helm-push v0.8.1
	github.com/spf13/cobra v1.0.0
	github.com/stretchr/testify v1.5.1
	helm.sh/helm/v3 v3.2.1
	k8s.io/helm v2.16.9+incompatible
)

replace (
	github.com/docker/docker => github.com/docker/docker v0.0.0-20190731150326-928381b2215c
	helm.sh/helm/v3 => helm.sh/helm/v3 v3.2.1
	k8s.io/helm => k8s.io/helm v2.16.9+incompatible
)
