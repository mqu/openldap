package openldap

/*

#include <stdlib.h>
#include <ldap.h>

static LDAPMessage **to_ldap_mesg(void* msg){return (LDAPMessage **) msg;}
static struct timeval *tv_ptr(void* tv){return (struct timeval *) tv;}
static inline char* to_charptr(const void* s) { return (char*)s; }

*/
// #cgo CFLAGS: -DLDAP_DEPRECATED=1
// #cgo linux CFLAGS: -DLINUX=1
// #cgo LDFLAGS: -lldap -llber
import "C"

import (
	"errors"
	"fmt"
	"syscall"
	"unsafe"
)

// ------------------------------------------ RESULTS methods ---------------------------------------------
/*

	openldap C API : 

    int ldap_count_messages( LDAP *ld, LdapMessage *result )
    LdapMessage *ldap_first_message( LDAP *ld, LdapMessage *result )
    LdapMessage *ldap_next_message ( LDAP *ld, LdapMessage *message )

	int ldap_count_entries( LDAP *ld, LdapMessage *result )
	LdapMessage *ldap_first_entry( LDAP *ld, LdapMessage *result )
	LdapMessage *ldap_next_entry ( LDAP *ld, LdapMessage *entry )

    char *ldap_first_attribute(LDAP *ld, LdapMessage *entry, BerElement **berptr )
    char *ldap_next_attribute (LDAP *ld, LdapMessage *entry, BerElement *ber )

    char **ldap_get_values(LDAP *ld, LdapMessage *entry, char *attr)
    struct berval **ldap_get_values_len(LDAP *ld, LdapMessage *entry,char *attr)

    int ldap_count_values(char **vals)
    int ldap_count_values_len(struct berval **vals)
    void ldap_value_free(char **vals)
    void ldap_value_free_len(struct berval **vals)

*/

func (self *LdapMessage) Count() int {
	// API : int ldap_count_messages(LDAP *ld, LdapMessage *chain )
	// err : (count = -1)
	count := int(C.ldap_count_messages(self.ldap.conn, self.msg))
	if count == -1 {
		panic("LDAP::Count() (ldap_count_messages) error (-1)")
	}
	return count

}

func (self *LdapMessage) FirstMessage() *LdapMessage {

	var msg *C.LDAPMessage
	msg = C.ldap_first_message(self.ldap.conn, self.msg)
	if msg == nil {
		return nil
	}
	_msg := new(LdapMessage)
	_msg.ldap = self.ldap
	_msg.errno = 0
	_msg.msg = msg
	return _msg
}

func (self *LdapMessage) NextMessage() *LdapMessage {
	var msg *C.LDAPMessage
	msg = C.ldap_next_message(self.ldap.conn, self.msg)

	if msg == nil {
		return nil
	}
	_msg := new(LdapMessage)
	_msg.ldap = self.ldap
	_msg.errno = 0
	_msg.msg = msg
	return _msg
}

/* an alias to ldap_count_message() ? */
func (self *LdapEntry) CountEntries() int {
	// API : int ldap_count_messages(LDAP *ld, LdapMessage *chain )
	// err : (count = -1)
	return int(C.ldap_count_entries(self.ldap.conn, self.entry))
}

func (self *LdapMessage) FirstEntry() *LdapEntry {

	var msg *C.LDAPMessage
	// API: LdapMessage *ldap_first_entry( LDAP *ld, LdapMessage *result )
	msg = C.ldap_first_entry(self.ldap.conn, self.msg)
	if msg == nil {
		return nil
	}
	_msg := new(LdapEntry)
	_msg.ldap = self.ldap
	_msg.errno = 0
	_msg.entry = msg
	return _msg
}

func (self *LdapEntry) NextEntry() *LdapEntry {
	var msg *C.LDAPMessage
	// API: LdapMessage *ldap_next_entry ( LDAP *ld, LdapMessage *entry )
	msg = C.ldap_next_entry(self.ldap.conn, self.entry)

	if msg == nil {
		return nil
	}
	_msg := new(LdapEntry)
	_msg.ldap = self.ldap
	_msg.errno = 0
	_msg.entry = msg
	return _msg
}

func (self *LdapEntry) FirstAttribute() (string, error) {

	var ber *C.BerElement

	// API: char *ldap_first_attribute(LDAP *ld, LdapMessage *entry, BerElement **berptr )
	rv := C.ldap_first_attribute(self.ldap.conn, self.entry, &ber)

	if rv == nil {
		// error
		return "", nil
	}
	self.ber = ber
	return C.GoString(rv), nil
}

func (self *LdapEntry) NextAttribute() (string, error) {

	// API: char *ldap_next_attribute (LDAP *ld, LdapMessage *entry, BerElement *ber )
	rv := C.ldap_next_attribute(self.ldap.conn, self.entry, self.ber)

	if rv == nil {
		// error
		return "", nil
	}
	return C.GoString(rv), nil
}

// private func
func sptr(p uintptr) *C.char {
	return *(**C.char)(unsafe.Pointer(p))
}

// private func used to convert null terminated char*[] to go []string
func cstrings_array(x **C.char) []string {
	var s []string
	for p := uintptr(unsafe.Pointer(x)); sptr(p) != nil; p += unsafe.Sizeof(uintptr(0)) {
		s = append(s, C.GoString(sptr(p)))
	}
	return s
}

/*
 FIXME: does not work for binary attributes
 FIXME:
  If  the  attribute values are binary in nature, and thus not suitable to be returned as an array of char *'s, the ldap_get_values_len() routine can be used instead.  It
  takes the same parameters as ldap_get_values(), but returns a NULL-terminated array of pointers to berval structures, each containing the length of and a pointer  to  a
  value.
 */

func (self *LdapEntry) GetValues(attr string) []string {

	_attr := C.CString(attr)
	defer C.free(unsafe.Pointer(_attr))

	// DEPRECATED
	// API: char **ldap_get_values(LDAP *ld, LdapMessage *entry, char *attr)
	values := cstrings_array(C.ldap_get_values(self.ldap.conn, self.entry, _attr))
	// count := C.ldap_count_values(values)

	return values
}

// ------------------------------------------------ RESULTS -----------------------------------------------
/*
    int ldap_result ( LDAP *ld, int msgid, int all, struct timeval *timeout, LdapMessage **result );
	int ldap_msgfree( LdapMessage *msg );
	int ldap_msgtype( LdapMessage *msg );
	int ldap_msgid  ( LdapMessage *msg );

*/

// Result()
// take care to free LdapMessage result with MsgFree()
//
func (self *Ldap) Result() (*LdapMessage, error) {

	var msgid int = 1
	var all int = 1

	timeout := syscall.Timeval{30, 0} // timeout for 30 seconds by default
	tv := C.tv_ptr(unsafe.Pointer(&timeout))

	var _result *LdapMessage
	result := C.to_ldap_mesg(unsafe.Pointer(_result))

	// API: int ldap_result( LDAP *ld, int msgid, int all, struct timeval *timeout, LdapMessage **result );
	rv := C.ldap_result(self.conn, C.int(msgid), C.int(all), tv, result)

	if rv != LDAP_OPT_SUCCESS {
		return nil, errors.New(fmt.Sprintf("LDAP::Result() error :  %d (%s)", rv, ErrorToString(int(rv))))
	}

	return _result, nil
}

// MsgFree() is used to free LDAP::Result() allocated data
//
// returns -1 on error.
//
func (self *LdapMessage) MsgFree() int{
        if self.msg != nil {
                rv := C.ldap_msgfree(self.msg)
                self.msg = nil
                return int(rv)
        }
        return -1
}


//  ---------------------------------------- DN Methods ---------------------------------------------------
/*

	char *ldap_get_dn( LDAP *ld, LdapMessage *entry)
	int   ldap_str2dn( const char *str, LDAPDN *dn, unsigned flags)
	void  ldap_dnfree( LDAPDN dn)
	int   ldap_dn2str( LDAPDN dn, char **str, unsigned flags)

	char **ldap_explode_dn( const char *dn, int notypes)
	char **ldap_explode_rdn( const char *rdn, int notypes)

	char *ldap_dn2ufn  ( const char * dn )
	char *ldap_dn2dcedn( const char * dn )
	char *ldap_dcedn2dn( const char * dn )
	char *ldap_dn2ad_canonical( const char * dn )

*/

func (self *LdapEntry) GetDn() string {
	// API: char *ldap_get_dn( LDAP *ld, LdapMessage *entry )
	rv := C.ldap_get_dn(self.ldap.conn, self.entry)
	defer C.free(unsafe.Pointer(rv))

	if rv == nil {
		err := self.ldap.Errno()
		panic(fmt.Sprintf("LDAP::GetDn() error %d (%s)", err, ErrorToString(err)))
	}

	val := C.GoString(rv)
	return val
}
