# netpeek

`netpeek` is a layer 7 sniffer. Like `tcpdump`, but at the application layer. 

It has support for the HTTP protocol currently, so it can reconstruct the HTTP request and response from the packets flowing thru the network interface and display it with latency, packet level stats.

Some flags for filtering the packets based on host and ports:


- `sport`  
Destination host of the packet  

- `dport`  
Destination port of the packet  

- `shost`  
Destination host of the packet  

- `dhost`  
Destination port of the packet  

- `i`  
Network interface to sniff on  

- `protocol`  
One of `http`, `drain`, `dump`  
-- `http`  
will dump the request and response on stdout  
-- `drain`  
will dump the packet metadata on stdout  
-- `dump`   
will dump the packet payload on stdout  

- `v`  
Verbose logging  

- `cui` (experimental)  
Use CUI (character user interface) mode  
