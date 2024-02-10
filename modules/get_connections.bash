TOTAL=$(netstat -an | grep 'ESTABLISHED' | grep ':3399' | grep 'tcp' | wc -l) # :3399 - Cnc port
echo $(($TOTAL / 2)) #cause tcp4 and tcp6