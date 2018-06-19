# RIP
Простой DNS сервер с возможность кодирования IP адреса в имени.

# Кодирование IP
DNS сервер парсит имя, пытается вычленить из него требуемый IP и формирует ответ.
Правила:
```
    <optional-prefix>.<IPv4>.4.<zone>  -> returns A record with <IPv4> address
    <optional-prefix>.<IPv6>.6.<zone>  -> returns AAAA record with <IPv6> address
    <proxy-name>.p.<zone>  -> resolve proxy name and returns it
    <cname>.c.<zone>  -> return CNAME record with <cname>
    <any-name>.<zone>  -> returns default address
```

# Пример
Например, запустим DNS сервер дл зоны `example.com` с дефолтными IP `77.88.55.70` и `2a02:6b8:a::a`:
```
$ rip --zone=example.com --ipv4=77.88.55.70 --ipv6=2a02:6b8:a::a
```
При запросе к нему, мы должны получить следующие ответы:
```
# IPv4
    1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a
    foo.1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a
    bar.foo.1-1-1-1.4.example.com ->  1.1.1.1  && 2a02:6b8:a::a

# IPv6
    2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    ->  2a01:7e01::f03c:91ff:fe3b:c9ba  && 77.88.55.70
    foo.2a01-7e01--f03c-91ff-fe3b-c9ba.6.example.com    -> 2a01:7e01::f03c:91ff:fe3b:c9ba  && 77.88.55.70
    foo.--1.6.example.com   ->  ::1

# Cname
    ya.ru.c.example.com ->  canonical name ya.ru
    google.com.c.example.com ->  canonical name google.com

# Proxy
    ya.ru.p.example.com ->  87.250.250.242 && 2a02:6b8::2:242
    google.com.exampl.ecom  ->  64.233.164.102 and 2a00:1450:4010:c07::64
```
