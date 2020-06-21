# The Story

So, the one liner is this - in my current project, we have this GitLab CI pipeline for a GitLab repo
which contains our helm charts, and the pipeline runs to push all the charts to
[chartmuseum](https://chartmuseum.com/). The pipeline is too slow for me 🙈 I like to make my life easier
and smooth and fast. This is my attempt to do just that - make it easier for me to push charts faster and to
also fail fast if there are any issues in the pipeline! :)

Also, a lot of the ideas started from here https://kubernetes.slack.com/archives/C0NH30761/p1575737530228200

![1](images/story/1.png "1")

![2](images/story/2.png "2")

I'll keep updating this story with more details about the problem I'm trying to solve. I have written
down the problems and ideas in multiple places. 🙈 😅

These are some of the things I wrote in my notes:

1. Push all charts to a chart repo - fast is key
2. When the chart already exists in chart repo, it should not be pushed and no error should come if existing chart present in chart repo and the chart being pushed are the same. Also, try to not download the chart from chart repo to do this check. Use head request or other ways using the http API. See if that's possible
3. When there are child charts to be pushed. Push child charts first. Then parent charts. This way, you can push both at once. First child, then parent will be pushed later. Draw a graph and walk it and push. For lock files - no issues if versions are exact. But there might be issues if versions are not exact! Like range. Or else lock files can be generated. Think on this. How to solve lock file updation.
4. Current pipeline pushes charts one by one in sequence. We can parallelize it. Like have 10 or so threads (workers) that keep picking up tasks of pushing charts. Given a worker, it looks if there's a task to be done, and then does it. Check if there are frameworks for this. The task is to push a chart. In which order should the charts be pushed is a question that comes when you want to solve 3. Especially let's say there are 10 child charts, and 10 workers, and the remaining charts are parent charts and their parents and so on. The first ten workers spawn up, push the charts, let's say one of them has done it really fast cuz it's a small and simple chart, it takes up a parent chart. Now, if this parent depends on a child chart that some worker is still pushing or going to push, then that's a problem. Parallelizing this may not be easy. One way to do this would be - the tasks will be given in such a way that, there are dependent tasks. So if a task's dependencies are done, then that task is put in the task queue for the workers to pickup. How does that sound? Hmm? Makes sense?
5. Also, when child charts and parent charts are involved, and are updated in a single go, then we need to run repository update every time, before doing anything. Actually, come to think of it, helm dependency update already does that, no?

---

Tasks:
1. Write module code to push one chart. Do TDD.

2. Write module code to read charts from a directory. Any directory with Chart.yaml is a chart.

3. Write module code to validate chart from a directory. Only valid charts must be pushed. If linting fails, it must not be pushed I guess. Think on this!

---

Some things to note -
* helm push plugin has quite some configurations that can be done using multiple
flags. See below the flags it has based on version `0.8.1`

```bash
$ helm push -h
Helm plugin to push chart package to ChartMuseum

Examples:

  $ helm push mychart-0.1.0.tgz chartmuseum       # push .tgz from "helm package"
  $ helm push . chartmuseum                       # package and push chart directory
  $ helm push . --version="7c4d121" chartmuseum   # override version in Chart.yaml
  $ helm push . https://my.chart.repo.com         # push directly to chart repo URL

Usage:
  helm push [flags]

Flags:
      --access-token string             Send token in Authorization header [$HELM_REPO_ACCESS_TOKEN]
      --auth-header string              Alternative header to use for token auth [$HELM_REPO_AUTH_HEADER]
      --ca-file string                  Verify certificates of HTTPS-enabled servers using this CA bundle [$HELM_REPO_CA_FILE]
      --cert-file string                Identify HTTPS client using this SSL certificate file [$HELM_REPO_CERT_FILE]
      --check-helm-version              outputs either "2" or "3" indicating the current Helm major version
      --context-path string             ChartMuseum context path [$HELM_REPO_CONTEXT_PATH]
      --debug                           Enable verbose output
  -d, --dependency-update               update dependencies from "requirements.yaml" to dir "charts/" before packaging
  -f, --force                           Force upload even if chart version exists
  -h, --help                            help for helm
      --home string                     Location of your Helm config. Overrides $HELM_HOME (default "/Users/karuppiahn/.helm")
      --host string                     Address of Tiller. Overrides $HELM_HOST
      --insecure                        Connect to server with an insecure way by skipping certificate verification [$HELM_REPO_INSECURE]
      --key-file string                 Identify HTTPS client using this SSL key file [$HELM_REPO_KEY_FILE]
      --keyring string                  location of a public keyring (default "/Users/karuppiahn/.gnupg/pubring.gpg")
      --kube-context string             Name of the kubeconfig context to use
      --kubeconfig string               Absolute path of the kubeconfig file to be used
  -p, --password string                 Override HTTP basic auth password [$HELM_REPO_PASSWORD]
      --tiller-connection-timeout int   The duration (in seconds) Helm will wait to establish a connection to Tiller (default 300)
      --tiller-namespace string         Namespace of Tiller (default "kube-system")
  -u, --username string                 Override HTTP basic auth username [$HELM_REPO_USERNAME]
  -v, --version string                  Override chart version pre-push
```

Some might not make sense, for example, tiller related stuff will not make sense
for Helm v3. Also, Helm v2 will be gone soon. As in, the support for it - in
terms of bug fixes and security fixes and the feature development for it has
already been stopped. Only Helm v3 will get new features. So yeah. I guess it's
better to mainly support v3 and not have any flags related to tiller and all.

To start with, I'll not bring in these flags. At some point, I need to 😅

---

I can see that there are lot of command level features in the helm push plugin
code, other than the package code.

I was thinking if I should use exec and execute and helm push command. Hmm.

So, I was trying to think what kind of functionality this tool is trying to
provide.

The idea is to be a wrapper tool - to be built on top of helm push and provide
the extra feature of bulk push. Keep it fast. And also support pushing charts
which are interdependent.

So, the code of this tool will contain logic for features like
* Bulk push charts in a concurrent and parallel manner
* Know the dependency among charts and accordingly push the charts in proper
order. This order will be determined after drawing a dependency graph of all
the charts

What I thought was - may be, I can provide the above features and let users
hook into the whole functionality of helm push, by letting them tell what is
the command they want to run for each chart. A user interface like this -

```bash
$ helm push-all all-my-charts --command "helm push {.Chart} chartmuseum-repo" 
```

And the above is kind of the simplest feature. For the simplest feature, without
all the options and complex features, I could also provide

```bash
$ helm push-all all-my-charts chartmuseum-repo
```

If they want to use all the crazy options and all, they could just use the
previous and it will invoke the helm push command. Some things to note when it
comes to invoking commands / executing commands, the environment variables that
the command has access to, and the working directory it's present in. Ideally
it should be the same as what's accessible to helm-push-all and where
helm-push-all is running

---

Decisions
1. Support only Helm v3 at the moment. Don't do extra work to make v2 work. If
it comes for free, all cool. Or else not needed. But Helm v3 must be supported!

---

The first task I'm doing is - read all charts from a directory.

Modules
1. Write module code to read charts from a directory. Any directory with
Chart.yaml is a chart but there's more to it. So it's better to use proper
functions from Helm or Helm push to try to load the chart and if it's a valid
chart it'll be able to load, or else it will not load. 

So, what I'm going to do is, if some file or directory is not a valid chart
or not a chart at all, I'm going to silently ignore it and have warning
messages for those files / directories.

---

Running tests in GitHub Actions

I found this cool repo to help with this!

https://github.com/mvdan/github-actions-golang

It shows how GitHub actions can run a matrix of tests - different golang
versions and different platforms. I'm thinking of running it just for Mac OS
and go v1.14.x for now. I think I could support linux too! :) Let me just run
the tests in linux too. Windows, may be not. Will try to support this later if
needed. Or...okay, let me keep it, for tests. If there are any issues, based on
that I'll think about it. 

Changing this

```yaml
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.13.x, 1.14.x]
        platform: [ubuntu-latest, macos-latest, windows-latest]
```

to this

```yaml
jobs:
  test:
    strategy:
      matrix:
        go-version: [1.14.x]
        platform: [ubuntu-latest, macos-latest, windows-latest]
```
