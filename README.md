Repos
=====

Gizmo for the tool belt of industrious developers to keep track of all your local repos.

**Assumption:** You have multiple (git/vcs) repositories and often work on multiple over the day / time.

**Problem:** Forgetting to push some vs the hassle of checking each repo directory.

**Solution:** Check all local repos *with one command* and see if anything has been forgotten.

In addition, this tool serves as a demo command line application showcasing [CLIF (CLI Framework)](https://github.com/ukautz/clif).

Install
-------

Latest release with binaries for Linux, Mac & Windows can be found [here](https://github.com/ukautz/repos/releases)

### Build

``` bash
$ go get github.com/ukautz/repos
```

Usage
-----

![repos-commands](https://cloud.githubusercontent.com/assets/600604/8886554/b9753476-326c-11e5-9a8f-9d54e74713d6.png)

### Add repos

You can add repos in two ways: explicit adding a single directory or scanning a folder (recursively)

**Scan directory**

![repos-scan](https://cloud.githubusercontent.com/assets/600604/8886538/77b13378-326c-11e5-9c1d-168b8c0a4eb6.png)

**Add a single repo**

![repos-add](https://cloud.githubusercontent.com/assets/600604/8886537/779409f6-326c-11e5-9954-25a629530133.png)

### Check repos

Well, this is the primary function of this tool: Check if any of your repos have local (uncommitted/unpushed) changes.

![repos-check](https://cloud.githubusercontent.com/assets/600604/8886590/4b4ba164-326d-11e5-83ca-8fdd26783795.png)

State
-----

Currently only **Git** is supported. Since I don't use anything else atm... Check out the [Repo interface](common/repo.go) if you feel like contributing.


