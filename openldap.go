/*
 * Openldap (2.4.30) binding in GO 
 * 
 * 
 *  link to ldap or ldap_r (for thread-safe binding)
 * 
 */

package openldap

/*

#include <stdlib.h>
#include <ldap.h>

static inline char* to_charptr(const void* s) { return (char*)s; }

*/
// #cgo CFLAGS: -DLDAP_DEPRECATED=1
// #cgo linux CFLAGS: -DLINUX=1
// #cgo LDFLAGS: -lldap -llber
import "C"

import (
	"errors"
	"fmt"
	"unsafe"
	"strings"
)

/* Intialize() open an LDAP connexion ; supported url formats :
 * 
 *   ldap://host:389/
 *   ldaps://secure-host:636/
 * 
 * return values :
 *  - on success : LDAP object, nil
 *  - on error : nil and error with error description.
 */
func Initialize(url string) (*Ldap, error) {
	_url := C.CString(url)
	defer C.free(unsafe.Pointer(_url))

	var ldap *C.LDAP

	// API: int ldap_initialize (LDAP **ldp, LDAP_CONST char *url )
	rv := C.ldap_initialize(&ldap, _url)

	if rv != 0 {
		err := errors.New(fmt.Sprintf("LDAP::Initialize() error (%d) : %s", rv, ErrorToString(int(rv))))
		return nil, err
	}

	return &Ldap{ldap}, nil
}

/* 
 * Bind() is used for LDAP authentifications
 * 
 * if who is empty this is an anonymous bind
 * else this is an authentificated bind
 * 
 * return value : 
 *  - nil on succes,
 *  - error with error description on error.
 *
 */
func (self *Ldap) Bind(who, cred string) error {
	var rv int

	authmethod := C.int(LDAP_AUTH_SIMPLE)

	// API: int ldap_bind_s (LDAP *ld,	LDAP_CONST char *who, LDAP_CONST char *cred, int authmethod );
	if who == "" {
		_who := C.to_charptr(unsafe.Pointer(nil))
		_cred := C.to_charptr(unsafe.Pointer(nil))

		rv = int(C.ldap_bind_s(self.conn, _who, _cred, authmethod))
	} else {
		_who := C.CString(who)
		_cred := C.CString(cred)
		defer C.free(unsafe.Pointer(_who))
		rv = int(C.ldap_bind_s(self.conn, _who, _cred, authmethod))
	}

	if rv == LDAP_OPT_SUCCESS {
		return nil
	}

	self.conn = nil
	return errors.New(fmt.Sprintf("LDAP::Bind() error (%d) : %s", rv, ErrorToString(rv)))
}

/* 
 * close LDAP connexion
 * 
 * return value : 
 *  - nil on succes,
 *  - error with error description on error.
 *
 */
func (self *Ldap) Close() error {

	// API: int ldap_unbind(LDAP *ld)
	rv := C.ldap_unbind(self.conn)

	if rv == LDAP_OPT_SUCCESS {
		return nil
	}

	self.conn = nil
	return errors.New(fmt.Sprintf("LDAP::Close() error (%d) : %s", int(rv), ErrorToString(int(rv))))

}
/* 
 * Unbind() close LDAP connexion
 * 
 * an alias to Ldap::Close()
 *
 */
func (self *Ldap) Unbind() error {
	return self.Close()
}

func (self *Ldap) Search(base string, scope int, filter string, attributes []string) (*LdapMessage, error) {

	var attrsonly int = 0 // false: returns all, true, returns only attributes without values

	_base := C.CString(base)
	defer C.free(unsafe.Pointer(_base))

	_filter := C.CString(filter)
	defer C.free(unsafe.Pointer(_filter))

	// transform []string to C.char** null terminated array (attributes argument)
	_attributes := make([]*C.char, len(attributes)+1) // default set to nil (NULL in C)

	for i, arg := range attributes {
		_attributes[i] = C.CString(arg)
		defer C.free(unsafe.Pointer(_attributes[i]))
	}

	var msg *C.LDAPMessage

	// API: int ldap_search_s (LDAP *ld, char *base, int scope, char *filter, char **attrs, int attrsonly, LdapMessage * ldap_res)
	rv := int(C.ldap_search_s(self.conn, _base, C.int(scope), _filter, &_attributes[0], C.int(attrsonly), &msg))

	if rv == LDAP_OPT_SUCCESS {
		_msg := new(LdapMessage)
		_msg.ldap = self
		_msg.errno = rv
		_msg.msg = msg
		return _msg, nil
	}

	return nil, errors.New(fmt.Sprintf("LDAP::Search() error : %d (%s)", rv, ErrorToString(rv)))
}

// ------------------------------------- Ldap* method (object oriented) -------------------------------------------------------------------

func (self *LdapEntry) Append(a LdapAttribute){
	self.values = append(self.values, a)
}

func (self *LdapAttribute) String() string{
	return self.ToText()
}

func (self *LdapAttribute) ToText() string{
	return fmt.Sprintf("%s: [%s]", self.name, strings.Join(self.values, ", "))
}

func (self *LdapAttribute) Name() string{
	return self.name
}

func (self *LdapAttribute) Values() []string{
	return self.values
}

func (self *LdapEntry) Dn() string{
	return self.dn
}

func (self *LdapEntry) Attributes() []LdapAttribute{
	return self.values
}

func (self *LdapEntry) String() string{
	return self.ToText()
}

func (self *LdapEntry) ToText() string{

	txt := fmt.Sprintf("dn: %s\n", self.dn)
	
	for _, a := range self.values{
		txt = txt + fmt.Sprintf("%s\n", a.ToText())
	}

	return txt
}

func (self *LdapSearchResult) Append(e LdapEntry){
	self.entries = append(self.entries, e)
}

func (self *LdapSearchResult) ToText() string{

	txt := fmt.Sprintf("# query : %s\n", self.filter)
	txt = txt + fmt.Sprintf("# num results : %d\n", self.Count())
	txt = txt + fmt.Sprintf("# search : %s\n", self.Filter())
	txt = txt + fmt.Sprintf("# base : %s\n", self.Base())
	txt = txt + fmt.Sprintf("# attributes : [%s]\n", strings.Join(self.Attributes(), ", "))

	for _, e := range self.entries{
		txt = txt + fmt.Sprintf("%s\n", e.ToText())
	}

	return txt
}

func (self *LdapSearchResult) String() string{
	return self.ToText()
}

func (self *LdapSearchResult) Entries() []LdapEntry{
	return self.entries
}

func (self *LdapSearchResult) Count() int{
	return len(self.entries)
}

func (self *LdapSearchResult) Filter() string{
	return self.filter
}

func (self *LdapSearchResult) Base() string{
	return self.base
}

func (self *LdapSearchResult) Scope() int{
	return self.scope
}

func (self *LdapSearchResult) Attributes() []string{
	return self.attributes
}

func (self *Ldap) SearchAll(base string, scope int, filter string, attributes []string) (*LdapSearchResult, error) {

	sr := new(LdapSearchResult)

	sr.ldap   = self
	sr.base   = base
	sr.scope  = scope
	sr.filter = filter
	sr.attributes = attributes

	// Search(base string, scope int, filter string, attributes []string) (*LDAPMessage, error)	
	result, err := self.Search(base, scope, filter, attributes)

	if err != nil {
		fmt.Println(err)
		return sr, err
	}

	// Free LDAP::Result() allocated data
	defer result.MsgFree()

	e := result.FirstEntry()

	for e != nil {
		_e := new(LdapEntry)
		
		_e.dn = e.GetDn()

		attr, _ := e.FirstAttribute()
		for attr != "" {

			_attr := new(LdapAttribute)
			_attr.values = e.GetValues(attr)
			_attr.name = attr

			_e.Append(*_attr)

			attr, _ = e.NextAttribute()

		}

		sr.Append(*_e)

		e = e.NextEntry()
	}
	
	return sr, nil
}
