package openldap

/*
# include <ldap.h>

*/
// #cgo CFLAGS: -DLDAP_DEPRECATED=1
// #cgo linux CFLAGS: -DLINUX=1
// #cgo LDFLAGS: -lldap -llber
import "C"

type Ldap struct {
	conn *C.LDAP
}

type LdapMessage struct {
	ldap *Ldap
	// conn *C.LDAP
	msg   *C.LDAPMessage
	errno int
}

type LdapAttribute struct{
	name string
	values []string
}


type LdapEntry struct {
	ldap *Ldap
	// conn  *C.LDAP
	entry *C.LDAPMessage
	errno int
	ber   *C.BerElement

	dn string
	values []LdapAttribute
}

type LdapSearchResult struct{
	ldap *Ldap

	scope int
	filter string
	base string
	attributes []string
	
	entries []LdapEntry
}
