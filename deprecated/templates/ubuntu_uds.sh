#!/bin/bash
{{ $fqdn := printf "%s-%s%d.%s.%s" $.Host.Hostname $.Environment.Prefix $.PodID $.Network.Subdomain $.Competition.Domain -}}
sed -i 's/localhost$/localhost {{ $fqdn }}/;' /etc/hosts
echo '{{ $fqdn }}' | tee /etc/hostname
hostname -F /etc/hostname
service networking reload