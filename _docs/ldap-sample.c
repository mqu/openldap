
// search / read example

	for ( e = ldap_first_entry( ld, res ); e != NULL; e = ldap_next_entry( ld, e ) )
	{
			struct berval **v = ldap_get_values_len( ld, e, attr );

			if ( v != NULL ) {
					int n = ldap_count_values_len( v );
					int j;

					values = realloc( values, ( nvalues + n + 1 )*sizeof( char * ) );
					for ( j = 0; j < n; j++ ) {
							values[ nvalues + j ] = strdup( v[ j ]->bv_val );
					}
					values[ nvalues + j ] = NULL;
					nvalues += n;
					ldap_value_free_len( v );
			}
	}

	ldap_msgfree( res );
