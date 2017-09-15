#!/bin/bash
{{ $fqdn := printf "%s.%s.vdi.%s" $.Host.Hostname $.Environment.Name $.Environment.Domain }}
sed -i 's/^.*ssh-/ssh-/;' /root/.ssh/authorized_keys
sed -i 's/localhost$/localhost {{ $fqdn }}/;' /etc/hosts
echo '{{ $fqdn }}' | tee /etc/hostname
echo -e "{{ $.Environment.PodPassword $.PodID }}\n{{ $.Environment.PodPassword $.PodID }}\n" | passwd
echo ""
hostname -F /etc/hostname
service networking restart