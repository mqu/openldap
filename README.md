OpenLDAP
====

this is Openldap binding in GO language.


INSTALL:
-----

	# install openldap library and devel packages
	sudo apt-get install libldap libldap2-dev  # debian/ubuntu.
	sudo urpmi openldap-devel # fedora, RH, ...

	# install go
	go get github.com/mqu/openldap

	# verify you've got it :
	go list | grep openldap

Usage
----

- Look a this [exemple](https://github.com/mqu/openldap/blob/master/_examples/test-openldap.go).
- a more complexe example making  [LDAP search](https://github.com/mqu/openldap/blob/master/_examples/ldapsearch.go) that mimics ldapsearch command, printing out result on console.

Todo :
----
 - support binary values ! Search() for "all attributes" will segfault (panic: runtime error: invalid memory address)
   on binary attributes.
 - thread-safe test
 - complete LDAP:GetOption() and LDAP:SetOption() method : now, they work only for integer values.
 - avoid using deprecated function (see LDAP_DEPRECATED flag and "// DEPRECATED" comments in *.go sources)
 - ...

Doc:
---

- look at _examples/*.go to see how to use this library.
- will come soon, complete documentation in Wiki.


Link :
---

 - goc : http://code.google.com/p/go-wiki/wiki/cgo (how to bind native libraries to GO)
 - Openldap library (and server) : http://www.openldap.org/
