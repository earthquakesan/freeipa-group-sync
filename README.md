# freeipa-group-sync

freeipa-group-sync is an application, which will sync memberships of users in the groups listed in IPA_GROUPS_YAML_PATH (env var) yaml file to freeipa instance.

## Configuration

| Environment variable   | Description  | Example   |
|---|---|---|
| IPA_HOST  | freeipa host to connect to  | ipa.example.test |
| IPA_REALM | freeipa realm to connect to  | EXAMPLE.TEST |
| IPA_USERNAME | freeipa username | admin |
| IPA_PASSWORD | freeipa password | Pa33s0000rd |
| IPA_GROUPS_YAML_PATH | absolute path to yaml file used for synchronization | /tmp/groups.yaml |

## How to run

Edit groups.yaml and .envrc then:

```
git clone https://github.com/earthquakesan/freeipa-group-sync

# my modification to goipa is not yet in the upstream repo
# you need to have my version of goipa in the parent dir (see go.mod)
git clone https://github.com/earthquakesan/goipa

cd freeipa-group-sync
go build
set -o allexport
source .envrc
./freeipa-group-sync
```

## Example

Given the following yaml file:

```
groups: 
  - name: group1
    users:
  - name: group2
    users:
  - name: group3
    users:
    - user5
    - user6
```

The application will:
- Ensure that groups 'group1', 'group2', 'group3' exist in freeipa
- Remove all users from the groups, which are not listed in yaml file
- Add all users listed in the yaml file to the specified groups
- Tell you that user does not exist if that's the case

## Development Environment

### Running freeipa server

Start freeipa server:

```
docker run --name freeipa-server-container -ti \
    -p 80:80/tcp -p 443:443/tcp -p 389:389/tcp -p 636:636/tcp -p 88:88/tcp -p 464:464/tcp \
    -p 88:88/udp -p 464:464/udp -p 123:123/udp \
    --sysctl net.ipv6.conf.all.disable_ipv6=0 \
    -v /sys/fs/cgroup:/sys/fs/cgroup:ro \
    -h ipa.example.test --read-only \
    -v /tmp/ipa-data-4.6.6:/data:Z freeipa/freeipa-server:centos-7-4.6.6
```

The installation will ask for configuration params, keep them default.

Configure ipa.example.test to resolve to your localhost (i.e. edit /etc/hosts file).

Note: because of -v /sys/fs/cgroup:/sys/fs/cgroup:ro mount, the command will not work on windows docker (even in WSL2), use e.g. debian virtualbox on windows to run docker image.

### Configuring krb5 client

```
sudo apt-get install krb5-user
```

Edit /etc/krb5.conf:

```
[libdefaults]
        default_realm = EXAMPLE.TEST

# The following krb5.conf variables are only for MIT Kerberos.
        kdc_timesync = 1
        ccache_type = 4
        forwardable = true
        proxiable = true

        fcc-mit-ticketflags = true

[realms]
        EXAMPLE.TEST = {
                kdc = ipa.example.test:88
                admin_server = ipa.example.test:464
        }
```

Get /etc/ipa/ca.crt from the docker container and place it under /etc/ipa/ca.crt on your system (where you run code).


## Notes

Freeipa can be accessed via kerberos. To do that, in freeipa create user with the same name as username on your system, set initial password and test connection:

```
$ kinit
Password for iermilov@EXAMPLE.TEST:
Password expired.  You must change it now.
Enter new password:
Enter it again:
```

Install freeipa client:

```
sudo apt-get install freeipa-client
```

To create keytab file:

```
TODO
```

## TODOs

* Cover with tests
* Enable keytab auth
