package openldap

/*
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
)

// FIXME : support all kind of option (int, int*, ...)
func (self *Ldap) SetOption(opt int, val int) error {

	// API: ldap_set_option (LDAP *ld,int option, LDAP_CONST void *invalue));
	rv := C.ldap_set_option(self.conn, C.int(opt), unsafe.Pointer(&val))

	if rv == LDAP_OPT_SUCCESS {
		return nil
	}

	return errors.New(fmt.Sprintf("LDAP::SetOption() error (%d) : %s", int(rv), ErrorToString(int(rv))))
}

// FIXME : support all kind of option (int, int*, ...) should take care of all return type for ldap_get_option
func (self *Ldap) GetOption(opt int) (val int, err error) {

	// API: int ldap_get_option (LDAP *ld,int option, LDAP_CONST void *invalue));
	rv := C.ldap_get_option(self.conn, C.int(opt), unsafe.Pointer(&val))

	if rv == LDAP_OPT_SUCCESS {
		return val, nil
	}

	return 0, errors.New(fmt.Sprintf("LDAP::GetOption() error (%d) : %s", rv, ErrorToString(int(rv))))
}

/*
** WORK IN PROGRESS!
**
** OpenLDAP reentrancy/thread-safeness should be dynamically
** checked using ldap_get_option().
**
** The -lldap implementation is not thread-safe.
**
** The -lldap_r implementation is:
**              LDAP_API_FEATURE_THREAD_SAFE (basic thread safety)
** but also be:
**              LDAP_API_FEATURE_SESSION_THREAD_SAFE
**              LDAP_API_FEATURE_OPERATION_THREAD_SAFE
**
** The preprocessor flag LDAP_API_FEATURE_X_OPENLDAP_THREAD_SAFE
** can be used to determine if -lldap_r is available at compile
** time.  You must define LDAP_THREAD_SAFE if and only if you
** link with -lldap_r.
**
** If you fail to define LDAP_THREAD_SAFE when linking with
** -lldap_r or define LDAP_THREAD_SAFE when linking with -lldap,
** provided header definations and declarations may be incorrect.
**
 */

func (self *Ldap) IsThreadSafe() bool {
	// fmt.Println("IsThreadSafe()")
	// opt, err := self.GetOption(LDAP_API_FEATURE_THREAD_SAFE) ; fmt.Println(opt, err)
	// opt, err = self.GetOption(LDAP_THREAD_SAFE) ; fmt.Println(opt, err)
	// opt, err = self.GetOption(LDAP_API_FEATURE_X_OPENLDAP_THREAD_SAFE) ; fmt.Println(opt, err)

	//FIXME: need to implement LDAP::GetOption(LDAP_OPT_API_FEATURE_INFO)
	return false
}

func ErrorToString(err int) string {

	// API: char * ldap_err2string (int err )
	result := C.GoString(C.to_charptr(unsafe.Pointer(C.ldap_err2string(C.int(err)))))
	return result
}

func (self *Ldap) Errno() int {
	opt, _ := self.GetOption(LDAP_OPT_ERROR_NUMBER)
	return opt
}
