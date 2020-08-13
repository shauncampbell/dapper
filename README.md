# Dapper: lightweight authentication for home labs
Dapper is an easy-to-use LDAP server that can be used in very straightforward situations like a home lab.
* Single binary for easy portability.
* Simple YAML configuration for users.
* Automatically updates when configuration file changed.

Dapper has been born out of my frustration trying to set up and debug OpenLDAP for centralised password management on my home network. I realised quite quickly that I didn't need the full complexity of LDAP and that I only needed a few parts of it.

### Quickstart
Dapper is a work in progress and this isn't designed to be used in a production environment.

1. Download the binary from the [releases](https://github.com/shauncampbell/dapper/releases) page.
2. Write a sample configuration file and save it as dapper.yaml: 
```
users:
 - cn: user
   dn: cn=user,ou=users,dc=home,dc=lab
   uid: user
   description: "home lab user"
   objectClass:
    - "inetOrgPerson"
    - "jellyfinUser"
   mail: user@home.lab
   userPassword: "{SSHA}I8wq1+4gyJVJUtQW96JGcmCL46ADyPnW"
 - cn: root
   uid: root
   dn: cn=root,dc=home,dc=lab
   userPassword: "{SSHA}I8wq1+4gyJVJUtQW96JGcmCL46ADyPnW"
   objectClass: "posixAccount"

```
3. Start the dapper service:
```
./dapper -f dapper.yaml -b dc=home,dc=lab -p 3389
```
4. Test with `ldapsearch`
```
ldapsearch -H ldap://localhost:3389 -x -b 'dc=home,dc=lab' -D 'cn=root,dc=home,dc=lab' -w test
```

#### Output from ldapsearch
```
shauncampbell@Shaun-Campbell dapper % ldapsearch -H ldap://localhost:3389 -x -b 'dc=home,dc=lab' -D 'cn=root,dc=home,dc=lab' -w test
# extended LDIF
#
# LDAPv3
# base <dc=home,dc=lab> with scope subtree
# filter: (objectclass=*)
# requesting: ALL
#

# user, users, home.lab
dn: cn=user,ou=users,dc=home,dc=lab
description: home lab user
objectClass: inetOrgPerson
objectClass: jellyfinUser
mail: user@home.lab
userPassword:: e1NTSEF9STh3cTErNGd5SlZKVXRRVzk2SkdjbUNMNDZBRHlQblc=
cn: user
dn: cn=user,ou=users,dc=home,dc=lab
uid: user

# root, home.lab
dn: cn=root,dc=home,dc=lab
uid: root
dn: cn=root,dc=home,dc=lab
userPassword:: e1NTSEF9STh3cTErNGd5SlZKVXRRVzk2SkdjbUNMNDZBRHlQblc=
objectClass: posixAccount
cn: root

# search result
search: 2
result: 0 Success

# numResponses: 3
# numEntries: 2
```

#### Output from Dapper
```
shauncampbell@Shaun-Campbell dapper % sudo ./dapper server -f dapper.yaml -b dc=home,dc=lab -p 3389
7:21PM INF starting LDAP server on 0.0.0.0:3389
7:21PM INF reloading configuration file 'dapper.yaml'
7:21PM INF adding attribute attribute=description dn=cn=user,ou=users,dc=home,dc=lab value=["home lab user"]
7:21PM INF adding attribute attribute=objectClass dn=cn=user,ou=users,dc=home,dc=lab value=["inetOrgPerson","jellyfinUser"]
7:21PM INF adding attribute attribute=mail dn=cn=user,ou=users,dc=home,dc=lab value=["user@home.lab"]
7:21PM INF adding attribute attribute=userPassword dn=cn=user,ou=users,dc=home,dc=lab value=["{SSHA}I8wq1+4gyJVJUtQW96JGcmCL46ADyPnW"]
7:21PM INF adding attribute attribute=cn dn=cn=user,ou=users,dc=home,dc=lab value=["user"]
7:21PM INF adding attribute attribute=dn dn=cn=user,ou=users,dc=home,dc=lab value=["cn=user,ou=users,dc=home,dc=lab"]
7:21PM INF adding attribute attribute=uid dn=cn=user,ou=users,dc=home,dc=lab value=["user"]
7:21PM INF adding attribute attribute=uid dn=cn=root,dc=home,dc=lab value=["root"]
7:21PM INF adding attribute attribute=dn dn=cn=root,dc=home,dc=lab value=["cn=root,dc=home,dc=lab"]
7:21PM INF adding attribute attribute=userPassword dn=cn=root,dc=home,dc=lab value=["{SSHA}I8wq1+4gyJVJUtQW96JGcmCL46ADyPnW"]
7:21PM INF adding attribute attribute=objectClass dn=cn=root,dc=home,dc=lab value=["posixAccount"]
7:21PM INF adding attribute attribute=cn dn=cn=root,dc=home,dc=lab value=["root"]
7:21PM INF request received bindDN=cn=root,dc=home,dc=lab operation=bind request_ip=127.0.0.1:65223
7:21PM INF bind request was accepted bindDN=cn=root,dc=home,dc=lab operation=bind request_ip=127.0.0.1:65223
7:21PM INF beginning search with query: (objectclass=*) bindDN=cn=root,dc=home,dc=lab operation=search request_ip=127.0.0.1:65223
7:21PM INF dn 'cn=user,ou=users,dc=home,dc=lab' matches search criteria bindDN=cn=root,dc=home,dc=lab operation=search request_ip=127.0.0.1:65223
7:21PM INF dn 'cn=root,dc=home,dc=lab' matches search criteria bindDN=cn=root,dc=home,dc=lab operation=search request_ip=127.0.0.1:65223
7:21PM INF search completed with 2 results bindDN=cn=root,dc=home,dc=lab operation=search request_ip=127.0.0.1:65223
```

### Supported Features
The following features are supported right now:
* LDAP Bind (Simple)
* LDAP Search

### Supported LDAP queries
* Equality matches (e.g. `(field=a*)`)
* Not Matches (e.g. `(!(field=a))`)
* And conditions (e.g. `(&(field=a)(field2=c))`)
* Or conditions (e.g. `(|(field=a)(field=b))`)