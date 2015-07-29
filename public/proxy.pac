 
function FindProxyForURL(url,host)
{
	/*
    if (isInNet(dnsResolve(host), "172.24.0.0", "255.255.0.0")) {
        return "DIRECT";
    }
    
    if (isInNet(dnsResolve(host), "192.168.1.0", "255.255.255.0")) {
        return "DIRECT";
    }
    */
 
    return "SOCKS 172.24.222.54:8889";
}

 
