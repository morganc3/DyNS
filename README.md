# DyNS

DyNS is a dynamic DNS server, which allows clients to specify which records should be returned. For example, 
`a-record.169.254.169.254.dyns.in` would return an A Record for 169.254.169.254

# Usage
Currently, port 53 is hardcoded. There is a live example of this on dyns.in

## Supported Records and Formatting
Multiple records can be requested at once. This is useful for DNS rebinding attacks.
For example, a-record.169.254.169.254.a-record.127.0.0.1.cname-record-2.my.cname.dyns.in will return A Records for both IP addresses and a CNAME record for my.cname.

#### A Record
`a-record.<ipv4 address>.dyns.in`
  
#### AAAA Record
`aaaa-record.<ipv6 address>.dyns.in`
*NOTE*: Since colons are not allowed in domain names, they should be replaced with "." in the 
  request. 
  Example: `aaaa-record.2606.2800.220.1.248.1893.25c8.1946.dyns.in`
  
#### CNAME Record
`cname-record-<amount of subdomains in cname>.<cname>.dyns.in`
Example: cname-record-3.my.example.com.dyns.in returns CNAME for my.example.com
  
  
