Surfer
======

*Surfer* is a fast, simple and lightweight web-framework for `Go <http://golang.org>`_ programming language.

Objectives
==========

* Don't try to reinvent the wheel. Instead utilize already established libraries, like: `Gorilla web toolkit <www.gorillatoolkit.org>`_ (routing and session management)
* Do not create yet another config schema. Use yaml!
* The framework should be small yet productive:

  * Do not focus on features, which can be delivered by third parties.
  * Focus how to easily extend Surfer components to utilize third parties libraries, like: DB drivers, DB ORMs, template languages, session storages...


Requirements
============

    go get github.com/gorilla/securecookie
    go get github.com/gorilla/sessions
    go get github.com/gorilla/mux
    go get github.com/kylelemons/go-gypsy/yaml
    go get github.com/scale-it/clog
