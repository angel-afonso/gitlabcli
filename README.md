# GITLAB CLI
[![pipeline status](https://gitlab.com/angel-afonso/gitlabcli/badges/master/pipeline.svg)](https://gitlab.com/angel-afonso/gitlabcli/-/commits/master)
[![coverage report](https://gitlab.com/angel-afonso/gitlabcli/badges/master/coverage.svg)](https://gitlab.com/angel-afonso/gitlabcli/-/commits/master)
[![Go Report Card](https://goreportcard.com/badge/gitlab.com/angel-afonso/gitlabcli)](https://goreportcard.com/report/gitlab.com/angel-afonso/gitlabcli)

A Gitlab command line interface

## Work in progress

### Avaible Commands

* ***project***
  * ***list***: repositories shows that the user has access
  * ***view [path]***: If the current directory is a git repository with remote in gitlab, it will show information of that project, if not, it will show information of the project with the given path

* ***mergerequest***: 
    * ***create [path]***: If the current directory is a git repository, create a merge request for the current project, if not, create a merge request for the given path
