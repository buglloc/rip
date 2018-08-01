# RIP
A simple DNS server that extracts IP address (or cname) from the requested domain name and sends it back in the response.

# Encoding rules
DNS server parses requested name to extract the requested mode, IP or CNAME by the following rules:
```
    <optional-prefix>.<IPv4>.4.<zone>  -> returns A record with <IPv4> address
    <optional-prefix>.<IPv6>.6.<zone>  -> returns AAAA record with <IPv6> address
    <proxy-name>.p.<zone>  -> resolve proxy name and returns it
    <ip1>.<ip2>.r.<zone>  -> pick random <ip1> or <ip2>
    <ip1>.<ip2>.l.<zone>  -> loop over <ip1> and <ip2>
    <cname>.c.<zone>  -> return CNAME record with <cname>
    <any-name>.<zone>  -> returns default address
    [(<IPv4>.4|<IPv6>.6)...(<IPv4>.4|<IPv6>.6)].m.<zone>  -> returns multiple address according to order and type
```

# IP address format
IP address can be presented in two versions - dash-delimited and base16-form.
For example, these entries `0a000001` and `10-0-0-1` are equivalent and point to `10.0.0.1`
You can also use the built-in converter to encode IP address:
```
$ rip ip2hex fe80::fa94:c2ff:fee5:3cf6 127.0.0.1
fe80000000000000fa94c2fffee53cf6	7f000001
```


# Cname/ProxyName format
`cname` and `proxy ' modes support two name resolution logic - prefixed and dash-delimited:
For eg:
```
    something.victim.com.c.evil.com -> CNAME to something.victim.com
    something.victim-com.c.evil.com -> CNAME to victim.com
```

# Usage
Run NS server for zone `example.com` with default IP `77.88.55.70` and `2a02:6b8: a:: a`:
```
$ rip ns --zone=example.com --ipv4=77.88.55.70 --ipv6=2a02:6b8:a::a
```

When requesting it, we should get the following responses:
```
# IPv4
    1-1-1-1.4.example.com ->  1.1.1.1 
    foo.1-1-1-1.4.example.com ->  1.1.1.1
    bar.foo.1-1-1-1.4.example.com ->  1.1.1.1

# IPv6
    2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    ->  2a01:7e01::f03c:91ff:fe3b:c9ba
    foo.2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    -> 2a01:7e01::f03c:91ff:fe3b:c9ba
    foo.--1.6.example.com   ->  ::1

# Random
    0a000002.0a000001.r.example.com ->  random between 10.0.0.1 and 10.0.0.2

# Loop
    8ba299a7.8ba299a8.l.example.com ->  loop over 139.162.153.167 and 139.162.153.168

# Cname
    ya.ru.c.example.com ->  canonical name ya.ru
    google.com.c.example.com ->  canonical name google.com

# Proxy
    ya.ru.p.example.com ->  87.250.250.242 and 2a02:6b8::2:242
    google.com.p.example.com  ->  64.233.164.102 and 2a00:1450:4010:c07::64

# Multi
    1-1-1-1.4.8ba299a7.4.2a017e0100000000f03c91fffe3bc9ba.6.m.example.com   ->  1.1.1.1, 139.162.153.167, 2a01:7e01::f03c:91ff:fe3b:c9ba
```
