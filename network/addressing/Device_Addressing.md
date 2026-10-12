# Device Addressing Worksheet

Table 2  
*Device addressing and interface assignment worksheet*

| VLAN / Zone | Equipment | IP Assignment Type | Interface | IP Address | Subnet Mask |
| :--- | :--- | :--- | :---: | :---: | :---: |
| **VLAN 10 (Site A)** | PC0 (User device) | DHCP (Dynamic Pool) | Fa0 | `10.10.10.0` | `255.255.255.0` |
| **VLAN 10 (Site A)** | Router A (Gateway) | Last Usable IP | g0/0/0.10 | `10.10.10.254` | `255.255.255.0` |
| **VLAN 20 (Site A)** | Go-Chat-Server | 2nd Usable IP | Fa0 | `10.10.20.2` | `255.255.255.240` |
| **VLAN 20 (Site A)** | FastAPI-Auth-Server | 3rd Usable IP | Fa0 | `10.10.20.3` | `255.255.255.240` |
| **VLAN 20 (Site A)** | Router A (Gateway) | Last Usable IP | g0/0/0.20 | `10.10.20.14` | `255.255.255.240` |
| **VLAN 30 (Site A)** | Laptop 0 (Guest device) | DHCP | Fa0 | `10.10.30.0` | `255.255.255.128` |
| **VLAN 30 (Site A)** | Router A (Gateway) | Assigned Static IP | g0/0/0.30 | `10.10.30.65` | `255.255.255.128` |
| **VLAN 99 (Site A)** | Laptop1 (Admin device) | DHCP | Fa0 | `10.10.99.0` | `255.255.255.240` |
| **VLAN 99 (Site A)** | Router A (Gateway) | Assigned Static IP | g0/0/0.99 | `10.10.99.10` | `255.255.255.240` |
| **VLAN 10 (Site B)** | PC1, PC2 (Remote users) | DHCP (Dynamic Pool) | Fa0 | `10.10.11.0` | `255.255.255.0` |
| **VLAN 10 (Site B)** | Router B (Gateway) | Assigned Static IP | g0/0/0.10 | `10.10.11.100` | `255.255.255.0` |
| **VLAN 20 (Site B)** | Server2 (Branch Node) | 2nd Usable IP | Fa0 | `10.10.20.18` | `255.255.255.240` |
| **VLAN 20 (Site B)** | Router B (Gateway) | Last Usable IP | g0/0/0.20 | `10.10.20.30` | `255.255.255.240` |
| **VLAN 30 (Site B)** | Laptop2 (Guest device) | DHCP | Fa0 | `10.10.31.0` | `255.255.255.128` |
| **VLAN 30 (Site B)** | Router B (Gateway) | Assigned Static IP | g0/0/0.30 | `10.10.31.50` | `255.255.255.128` |
| **VLAN 99 (Site B)** | Laptop3 (Admin Device) | DHCP | Fa0 | `10.10.99.16` | `255.255.255.240` |
| **VLAN 99 (Site B)** | Router B (Gateway) | Assigned Static IP | g0/0/0.99 | `10.10.99.25` | `255.255.255.240` |
| **WAN / OSPF** | Router A (Gateway) | Assigned Static IP | g0/0/1 | `192.168.100.10` | `255.255.255.0` |
| **WAN / OSPF** | Router B (Gateway) | Assigned Static IP | g0/0/1 | `192.168.100.11` | `255.255.255.0` |

> *Note.* This worksheet specifies the exact IP assignment strategy, interface bindings, and subnet masks for every device and router subinterface across both sites.
