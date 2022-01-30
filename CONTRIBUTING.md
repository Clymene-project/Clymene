# Contributing Guide


* [New Contributor Guide](#contributing-guide)
  * [Ways to Contribute](#ways-to-contribute)
  * [Find an Issue](#find-an-issue)
  * [Ask for Help](#ask-for-help)
  * [Pull Request Lifecycle](#pull-request-lifecycle)
  * [Development Environment Setup](#development-environment-setup)
  * [Sign Your Commits](#sign-your-commits)
  * [Pull Request Checklist](#pull-request-checklist)


Welcome! We are glad that you want to contribute to our project! 💖

As you get started, you are in the best position to give us feedback on areas of
our project that we need help with including:

* Problems found during setting up a new developer environment
* Gaps in our Quickstart Guide or documentation
* Bugs in our automation scripts

If anything doesn't make sense, or doesn't work when you run it, please open a
bug report and let us know!

## Ways to Contribute

We welcome many different types of contributions including:

* New features
* Builds, CI/CD
* Bug fixes
* Documentation
* Issue Triage
* Answering questions on Slack/Mailing List
* Web design
* Communications / Social Media / Blog Posts
* Release management


Not everything happens through a GitHub pull request. Please come to our
[discussions](https://github.com/Clymene-project/Clymene/discussions) or contact us(clymeneallen@gmail.com) and let's discuss how we can work
together.

## Find an Issue

We have good first issues for new contributors and help wanted issues suitable
for any contributor. [good first issue](https://github.com/Clymene-project/Clymene/labels/good%20first%20issue) has extra information to
help you make your first contribution. [help wanted](https://github.com/Clymene-project/Clymene/labels/help%20wanted) are issues
suitable for someone who isn't a core maintainer and is good to move onto after
your first pull request.

Once you see an issue that you'd like to work on, please post a comment saying
that you want to work on it. Something like "I want to work on this" is fine.

## Ask for Help

The best way to reach us with a question when contributing is to ask on 

* The original github issue
* The developer mailing list(clymeneallan@gmail.com)
* Our Slack channel(clymene-project.slack.com)


## Development Environment Setup
```bash
go mod vendor
```

Tidy up the `go.mod` and `go.sum` files:

```bash
go mod tidy
go mod vendor
git add go.mod go.sum vendor
git commit
```


### Developer Certificates of Origin (DCOs)
Licensing is important to open source projects. It provides some assurances that
the software will continue to be available based under the terms that the
author(s) desired. We require that contributors sign off on commits submitted to
our project's repositories. The [Developer Certificate of Origin
(DCO)](https://developercertificate.org/) is a way to certify that you wrote and
have the right to contribute the code you are submitting to the project.

You sign-off by adding the following to your commit messages. Your sign-off must
match the git user and email associated with the commit.

    This is my commit message

    Signed-off-by: Your Name <your.name@example.com>

Git has a `-s` command line option to do this automatically:

    git commit -s -m 'This is my commit message'

If you forgot to do this and have not yet pushed your changes to the remote
repository, you can amend your commit with the sign-off by running

    git commit --amend -s 


## Pull Request Checklist

When you submit your pull request, or you push new commits to it, our automated
systems will run some checks on your new code. We require that your pull request
passes these checks, but we also have more criteria than just that before we can
accept and merge it. We recommend that you check the following things locally
before you submit your code:
