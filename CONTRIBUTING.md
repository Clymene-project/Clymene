# Contributing Guide

This is a template contributing guide for CNCF projects that requires editing
before it is ready to use. Read the markdown comments, `<!-- COMMENT -->`, for
additional guidance. The raw markdown uses `TODO` to identify areas that
require customization.

* [New Contributor Guide](#contributing-guide)
  * [Ways to Contribute](#ways-to-contribute)
  * [Find an Issue](#find-an-issue)
  * [Ask for Help](#ask-for-help)
  * [Pull Request Lifecycle](#pull-request-lifecycle)
  * [Development Environment Setup](#development-environment-setup)
  * [Sign Your Commits](#sign-your-commits)
  * [Pull Request Checklist](#pull-request-checklist)
* **TODO**
<!-- Additional Table of Contents
  At the top list level, Link to other related docs such as a reviewing
  guide, developers guide, etc.
-->

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

<!-- TODO: project maintainers fill in exactly which type of contributions you 
are willing to shepherd through your processes to make sure that contributors 
feel successful with their contributions. Make sure that you provide clear 
information about security concerns (are they handled in the open? Submitted to 
a different list with high priority?) -->

* New features
* Builds, CI/CD
* Bug fixes
* Documentation
* Issue Triage
* Answering questions on Slack/Mailing List
* Web design
* Communications / Social Media / Blog Posts
* Release management

<!-- Think about your project's contribution ladder, and if it makes sense, 
encourage people to review pull requests as a way to contribute as well --> 

Not everything happens through a GitHub pull request. Please come to our
[meetings](TODO) or [contact us](TODO) and let's discuss how we can work
together.

<!-- TODO: project maintainers fill in details about what people should not 
do with contributions. Examples might include don’t change version information 
or update changelogs. -->

### Come to meetings!
Absolutely everyone is welcome to come to any of our meetings. You never need an
invite to join us. In fact, we want you to join us, even if you don’t have
anything you feel like you want to contribute. Just being there is enough!

You can find out more about our meetings [here](TODO). You don’t have to turn on
your video. The first time you come, introducing yourself is more than enough.
Over time, we hope that you feel comfortable voicing your opinions, giving
feedback on others’ ideas, and even sharing your own ideas, and experiences.

## Find an Issue

We have good first issues for new contributors and help wanted issues suitable
for any contributor. [good first issue](TODO) has extra information to
help you make your first contribution. [help wanted](TODO) are issues
suitable for someone who isn't a core maintainer and is good to move onto after
your first pull request.

Sometimes there won’t be any issues with these labels. That’s ok! There is
likely still something for you to work on. If you want to contribute but you
don’t know where to start or can't find a suitable issue, you can **TODO**
<!-- say how people can reach out to you for help finding something to work on -->  

Once you see an issue that you'd like to work on, please post a comment saying
that you want to work on it. Something like "I want to work on this" is fine.

## Ask for Help

The best way to reach us with a question when contributing is to ask on **TODO**
<!-- Replace one of the options below with how a contributor can best 
ask for help on your project when working on a issue --> 

* The original github issue
* The developer mailing list
* Our Slack channel

## Pull Request Lifecycle

**TODO**
<!-- This is an optional section but we encourage you to think about your 
pull request process and help set expectations for both contributors and 
reviewers.

Instead of a fixed template, use these questions below as an exercise to uncover
the unwritten rules and norms your project has for both reviewers and
contributors. Using your answers, write a description of what a
contributor can expect during their pull request.

* When should contributors start to submit a PR - when it’s ready for review or
  as a work-in-progress?
* How do contributors signal that a PR is ready for review or that it’s not
  complete and still a work-in-progress?
* When should the contributor should expect initial review? The follow-up
  reviews?
* When and how should the author ping/bump when the pull request is ready for
  further review or appears stalled?
* How to handle stuck pull requests that you can’t seem to get reviewed?
* How to handle follow-up issues and pull requests?
* What kind of pull requests do you prefer: small scope, incremental value or
  feature complete?
* What should contributors do if they no longer want to follow-through with the
  PR? For example, will maintainers potentially refactor and use the code?
  Will maintainers close a PR if the contributor hasn’t responded in a specific
  timeframe?
* Once a PR is merged, what is the process for it getting into the next release?
* When does a contribution show up “live”?

Here are some examples from other projects:
 
* https://porter.sh/src/CONTRIBUTING.md#the-life-of-a-pull-request

-->

## Development Environment Setup

**TODO**
<!-- Provide enough information so that someone can find your project on 
the weekend and get set up, build the code, test it and submit a pull request 
successfully without having to ask any questions. If there is a one-off tool
they need to install, of common error people run into, or useful script they
should run, document it here. 

Document any necessary tools, for example VS Code and recommended extensions.
You don’t have to document the beginner’s guide to these tools, but how they
are used within the scope of your project.

* How to get the source code
* How to get any dependencies
* How to build the source code
* How to run the project locally
* How to test the source code, unit and "integration" or "end-to-end"
* How to generate and preview the documentation locally
* Links to new user documentation videos and examples to get people started and
  understanding how to use the project

-->

## Sign Your Commits

<!-- TODO: Based on your project, keep either the DCO or CLA section below -->

### DCO
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

### CLA
We require that contributors have signed our Contributor License Agreement (CLA).
<!--Explain the process for how to sign or link to it here -->

## Pull Request Checklist

When you submit your pull request, or you push new commits to it, our automated
systems will run some checks on your new code. We require that your pull request
passes these checks, but we also have more criteria than just that before we can
accept and merge it. We recommend that you check the following things locally
before you submit your code:

**TODO**
<!-- list both the automated and any manual checks performed by reviewers, it
is very helpful when the validations are automated in a script for example in a
Makefile target. Below is an example of a checklist:

* It passes tests: run the following command to run all of the tests locally:
  `make build test lint`
* Impacted code has new or updated tests
* Documentation created/updated
* We use [Azure DevOps, GitHub Actions, CircleCI]  to test all pull
  requests. We require that all tests succeed on a pull request before it is merged.

-->