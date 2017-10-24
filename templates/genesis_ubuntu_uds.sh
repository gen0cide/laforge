#!/bin/bash
{{ $fqdn := printf "%s.%s" $.Host.Hostname  $.Competition.Domain -}}
sed -i 's/localhost$/localhost {{ $fqdn }}/;' /etc/hosts
echo '{{ $fqdn }}' | tee /etc/hostname
hostname -F /etc/hostname
service networking reload