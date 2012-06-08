package main
/*
 * Author : Marc Quinton / 2012.
 * 
 * ldapsearch command mimics openldap/seach command. Supported options :
 *  - host : ldap[s]://hostname:port/ format,
 *  - user,
 *  - password,
 *  - base
 * 
 *  arguments : filter [attributes]
 *  - filter is an LDAP filter (ex: objectClass=*, cn=*admin*", ...
 *  - attributes is an LDAP attribute list ; can be empty. ex: cn, sn, givenName, mail, ...
 * 
 */  
 
import (
	"fmt"
	"os"
	"errors"
	"flag"
	"github.com/mqu/openldap"
)

type LdapSearchOptions struct {
	host string
	user string
	passwd string
	
	base string
	filter string
	attributes []string
	
	scope int
}

type LdapSearchApp struct {
	ldap *openldap.Ldap
	opts *LdapSearchOptions
}

func NewLdapSearchApp() *LdapSearchApp{
	app := new(LdapSearchApp)
	return app
}

// Show ldapsearch app usage
func (self *LdapSearchApp) Usage(){
	fmt.Printf("usage: %s filter [attribute list]\n", os.Args[0])
	flag.PrintDefaults()
}

// Parse ldapsearch app options using flag package.
func (self *LdapSearchApp) ParseOpts() (*LdapSearchOptions, error){
	var opts LdapSearchOptions

	flag.StringVar(&opts.host,   "host",   "ldap://localhost:389/", "ldap server URL format : ldap[s]://hostname:port/")
	flag.StringVar(&opts.user,   "user",   ""                     , "user for authentification")
	flag.StringVar(&opts.passwd, "passwd", ""                     , "password for authentification")
	flag.StringVar(&opts.base,   "base",   ""                     , "base DN for search")

	flag.Parse()

	if flag.NArg() == 0 {
		self.Usage()
		return nil, errors.New(fmt.Sprintf("ParseOpts() error ; see usage for more information"))
	}

	opts.filter = flag.Arg(0)

	if len(flag.Args()) == 1 {
		opts.attributes = []string{}
	} else {
		opts.attributes = flag.Args()[1:]
	}

	return &opts, nil
}

// Connect and Bind to LDAP server using self.opts
func (self *LdapSearchApp) Connect() (error){

	var err error
	self.ldap, err = openldap.Initialize(self.opts.host)
	
	if err != nil {
		return err
	}
	
	//FIXME: should be an external option
	self.ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)

	err = self.ldap.Bind(self.opts.user, self.opts.passwd)
	if err != nil {
		return err
	}
	
	return nil
}

// Close() disconnect application from Ldap server
func (self *LdapSearchApp) Close() (error){
	return self.ldap.Close()
}

// Search using filter and returning attributes list
func (self *LdapSearchApp) Search() (*openldap.LdapSearchResult, error){

	//FIXME: should be an external option
	scope := openldap.LDAP_SCOPE_SUBTREE

	return self.ldap.SearchAll(
		self.opts.base, 
		scope, 
		self.opts.filter, 
		self.opts.attributes)
}

// Print search result
func (self *LdapSearchApp) Print(res *openldap.LdapSearchResult) (error){
	fmt.Println(res)
	return nil
}

func main() {

	var err error
	app := NewLdapSearchApp()

	app.opts, err = app.ParseOpts()

	if err != nil {
		fmt.Println(err)
		return
	}

	err = app.Connect()

	if err != nil {
		fmt.Println(err)
		return
	}

	result, err := app.Search()

	if(err != nil) {
		fmt.Println("search error:", err)
		return
	}

	app.Print(result)
	app.Close()

}
