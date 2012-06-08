package main

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
	fmt.Println("# Connect()")

	var err error
	self.ldap, err = openldap.Initialize(self.opts.host)
	
	if err != nil {
		return err
	}
	
	self.ldap.SetOption(openldap.LDAP_OPT_PROTOCOL_VERSION, openldap.LDAP_VERSION3)

	err = self.ldap.Bind(self.opts.user, self.opts.passwd)
	if err != nil {
		return err
	}
	
	fmt.Println("# Connect()", self)
	return nil
}

// Close() disconnect application from Ldap server
func (self *LdapSearchApp) Close() (error){
	fmt.Println("# Close()")
	return nil
}

// Search using filter and returning attributes list
func (self *LdapSearchApp) Search() (*openldap.LdapSearchResult, error){
	fmt.Println("# Search()", self)

	scope := openldap.LDAP_SCOPE_SUBTREE
	return self.ldap.SearchAll(
		self.opts.base, 
		scope, 
		self.opts.filter, 
		self.opts.attributes)
}

// Print search result
func (self *LdapSearchApp) Print(res *openldap.LdapSearchResult) (error){
	fmt.Println("# Print()")
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
