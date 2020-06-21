module github.com/karuppiah7890/helm-push-all

go 1.14

require (
	github.com/chartmuseum/helm-push v0.8.1
	github.com/spf13/cobra v0.0.6
	github.com/stretchr/testify v1.4.0
)

replace github.com/docker/docker => github.com/docker/docker v0.0.0-20190731150326-928381b2215c
