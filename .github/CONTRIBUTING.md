## Introduction

### Welcome!

First off, thank you for considering contributing to Orlop.

### Tell them why they should read your guidelines.

Following these guidelines helps to communicate that you respect the time of the developers
managing and developing this open source project. In return, they should reciprocate that
respect in addressing your issue, assessing changes, and helping you finalize your pull requests.

### Explain what kinds of contributions you are looking for.

Orlop is an open source project and we love to receive contributions from our community — you! There are many ways to contribute, from writing tutorials or blog posts, improving the documentation, submitting bug reports and feature requests or writing code which can be incorporated into Elasticsearch itself. 

### Using the issue tracker

First things first: Do NOT report security vulnerabilities in public issues! Please disclose
responsibly by letting the [SwitchBit team](mailto:security@switchbit.com) know upfront.
We will assess the issue as soon as possible on a best-effort basis and will give you an
estimate for when we have a fix and release available for an eventual public disclosure.

The issue tracker is the preferred channel for bug reports, features requests and submitting pull requests, but please respect the following restrictions:
* Please do not use the issue tracker for personal support requests.
* Please do not derail or troll issues. Keep the discussion on topic and respect the opinions of others.

## Ground Rules

## Responsibilities
* Ensure that code that goes into core meets all requirements in this checklist:
  * General
    - [ ] Is this change useful kjust to me, or something that I think will benefit others greatly?
    - [ ] Check for overlap with other PRs
    - [ ] Think carefully about the long-term implications of the change. How will it affect existing projects that are dependent on this?
  * Quality
    - [ ] Is it consistent with all other code?
    - [ ] Have you documented it appropriately?
    - [ ] Is it free of lint?
    - [ ] Do all of the existing tests work? Have you added tests to cover your changes?
    - [ ] Take the time to get things right. PRs almost always require additional improvements to meet the bar for quality. Be very strict about quality. This usually takes several commits on top of the original PR.
* Create issues for any major changes and enhancements that you wish to make. Discuss things transparently and get community feedback.
* Keep feature versions as small as possible, preferably one new feature per version.
* Be welcoming to newcomers and encourage diverse new contributors from all backgrounds.

## Your First Contribution

Unsure where to begin contributing? You can start by looking through these beginner and help-wanted issues:
* Beginner issues - issues which should only require a few lines of code, and a test or two.
* Help wanted issues - issues which should be a bit more involved than beginner issues.

Both issue lists are sorted by total number of comments. While not perfect, number of
comments is a reasonable proxy for impact a given change will have.

Working on your first Pull Request? You can learn how from this *free* series,
[How to Contribute to an Open Source Project on GitHub](https://egghead.io/series/how-to-contribute-to-an-open-source-project-on-github).

At this point, you're ready to make your changes! Feel free to ask for help; everyone is a beginner at first.

If a maintainer asks you to "rebase" your PR, they're saying that a lot of code has changed, and that you need to update your branch so it's easier to merge.

## Getting started

### For something that is bigger than a one or two line fix:

1. Create your own fork of the code
2. Do the changes in your fork
3. If you like the change and think the project could use it:
    * Be sure you have followed the code style for the project.
    * Sign the Contributor License Agreement (CLA).
    * Note the  Code of Conduct.
    * Send a pull request indicating that you have a CLA on file.

### If you have a different process for small or "obvious" fixes, let them know.

Small contributions such as fixing spelling errors, where the content is small enough to not be considered intellectual property, can be submitted by a contributor as a patch, without a CLA.

As a rule of thumb, changes are obvious fixes if they do not introduce any new functionality or creative thinking. As long as the change does not affect functionality, some likely examples include the following:
* Spelling / grammar fixes
* Typo correction, white space and formatting changes
* Comment clean up
* Bug fixes that change default return values or error codes stored in constants
* Adding logging messages or debugging output
* Changes to ‘metadata’ files like go.mod, .gitignore, build scripts, etc.
* Moving source files from one directory or package to another

## How to report a bug

If you find a security vulnerability, do NOT open an issue. Email [security@switchvbit.com](mailto:security@switchvbit.com) instead.

In order to determine whether you are dealing with a security issue, ask yourself these two questions:
* Can I access something that's not mine, or something I shouldn't have access to?
* Can I disable something for other people?

 If the answer to either of those two questions are "yes", then you're probably dealing with a security issue. Note that even if you answer "no" to both questions, you may still be dealing with a security issue, so if you're unsure, just email us at [security@switchbit.com](mailto:security@switchbit.com.

 When filing an issue, make sure to answer these five questions:

1. What version of Go are you using (go version)?
2. What did you do?
3. What did you expect to see?
4. What did you see instead?

## How to suggest a feature or enhancement

If you find yourself wishing for a feature that doesn't exist in Orlop, you are
probably not alone. There are bound to be others out there with similar needs.
Many of the features that Orlop has today have been added because our users saw
the need. Open an issue on our issues list on GitHub which describes the feature
you would like to see, why you need it, and how it should work.

The core team looks at Pull Requests on a regular basis in a weekly triage meeting.

After feedback has been given we expect responses within two weeks. After two weeks
we may close the pull request if it isn't showing any activity.

## Commit message and labeling conventions

We have very precise rules over how our git commit messages should be formatted. This leads to **more
readable messages** that are easy to follow when looking through the **project history**.  Our tooling
also relies on it to properly assign release numbers.

Improperly formatted commit messages will result in your change not being merged.

### Commit Message Format

Each commit message consists of a **header**, a **body** and a **footer**. The header has a special
format that includes a **type**, a **scope** and a **subject**:

```html
<type>(<scope>): <subject>
<BLANK LINE>
<body>
<BLANK LINE>
<footer>
```

> Any line of the commit message cannot be longer 100 characters!
  This allows the message to be easier to read on GitHub as well as in various Git tools.

##### Type
Must be one of the following:

* **feat**: A new feature
* **fix**: A bug fix
* **docs**: Documentation only changes
* **style**: Changes that do not affect the meaning of the code (white-space, formatting, missing
  semi-colons, etc)
* **refactor**: A code change that neither fixes a bug nor adds a feature
* **perf**: A code change that improves performance
* **test**: Adding missing tests
* **chore**: Changes to the build process or auxiliary tools and libraries such as documentation
  generation

##### Scope
The scope could be anything that helps specifying the scope (or feature) that is changing.

Examples
- fix(token): 
- fix(vault): 

##### Subject
The subject contains a succinct description of the change:

* use the imperative, present tense: "change" not "changed" nor "changes"
* don't capitalize first letter
* no dot (.) at the end

##### Body
Just as in the **subject**, use the imperative, present tense: "change" not "changed" nor "changes"
The body should include the motivation for the change and contrast this with previous behavior.

##### Footer
The footer should contain any information about **Breaking Changes** and is also the place to
reference GitHub issues that this commit **Closes**, **Fixes**, or **Relates to**.

> Breaking Changes are intended to be highlighted in the ChangeLog as changes that will require
  community users to modify their code after updating to a version that contains this commit.
