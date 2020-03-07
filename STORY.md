# The Story

So, the one liner is this - in my current project, we have this GitLab CI pipeline for a GitLab repo
which contains our helm charts, and the pipeline runs to push all the charts to
[chartmuseum](https://chartmuseum.com/). The pipeline is too slow for me 🙈 I like to make my life easier
and smooth and fast. This is my attempt to do just that - make it easier for me to push charts faster and to
also fail fast if there are any issues in the pipeline! :)

I'll keep updating this story with more details about the problem I'm trying to solve. I have written
down the problems and ideas in multiple places. 🙈 😅

These are some of the things I wrote in my notes:

1. Push all charts to a chart repo - fast is key
2. When the chart already exists in chart repo, it should not be pushed and no error should come if existing chart present in chart repo and the chart being pushed are the same. Also, try to not download the chart from chart repo to do this check. Use head request or other ways using the http API. See if that's possible
3. When there are child charts to be pushed. Push child charts first. Then parent charts. This way, you can push both at once. First child, then parent will be pushed later. Draw a graph and walk it and push. For lock files - no issues if versions are exact. But there might be issues if versions are not exact! Like range. Or else lock files can be generated. Think on this. How to solve lock file updation.
4. Current pipeline pushes charts one by one in sequence. We can parallelize it. Like have 10 or so threads (workers) that keep picking up tasks of pushing charts. Given a worker, it looks if there's a task to be done, and then does it. Check if there are frameworks for this. The task is to push a chart. In which order should the charts be pushed is a question that comes when you want to solve 3. Especially let's say there are 10 child charts, and 10 workers, and the remaining charts are parent charts and their parents and so on. The first ten workers spawn up, push the charts, let's say one of them has done it really fast cuz it's a small and simple chart, it takes up a parent chart. Now, if this parent depends on a child chart that some worker is still pushing or going to push, then that's a problem. Parallelizing this may not be easy. One way to do this would be - the tasks will be given in such a way that, there are dependent tasks. So if a task's dependencies are done, then that task is put in the task queue for the workers to pickup. How does that sound? Hmm? Makes sense?
5. Also, when child charts and parent charts are involved, and are updated in a single go, then we need to run repository update every time, before doing anything. Actually, come to think of it, helm dependency update already does that, no?
