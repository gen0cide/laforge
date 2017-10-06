#!/bin/bash
{{ $fqdn := printf "%s-%s%d.%s.%s" $.Host.Hostname $.Environment.Prefix $.PodID $.Network.Subdomain $.Environment.Domain -}}
sed -i 's/^.*ssh-/ssh-/;' /root/.ssh/authorized_keys
sed -i 's/localhost$/localhost {{ $fqdn }}/;' /etc/hosts
echo '{{ $fqdn }}' | tee /etc/hostname
hostname -F /etc/hostname
service networking reload
touch /tmp/finished-user-data