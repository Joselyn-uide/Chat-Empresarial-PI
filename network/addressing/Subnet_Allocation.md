# Subnet Allocation Worksheet

Table 1  
*Subnet allocation worksheet for the distributed enterprise network (Site A and Site B)*

| Order | VLAN / Zone | Segment Name | Subnet / Prefix | First Usable IP | Last Usable IP | Broadcast |
| :---: | :--- | :--- | :---: | :---: | :---: | :---: |
| 1 | VLAN 10 (Site A) | Users (Site A) | `10.10.10.0/24` | `10.10.10.1` | `10.10.10.254` | `10.10.10.255` |
| 2 | VLAN 20 (Site A) | Go and FastAPI Servers (Site A) | `10.10.20.0/28` | `10.10.20.1` | `10.10.20.14` | `10.10.20.15` |
| 3 | VLAN 30 (Site A) | Guests (Site A) | `10.10.30.0/25` | `10.10.30.1` | `10.10.30.126` | `10.10.30.127` |
| 4 | VLAN 99 (Site A) | Administration (Site A) | `10.10.99.0/28` | `10.10.99.1` | `10.10.99.14` | `10.10.99.15` |
| 5 | VLAN 10 (Site B) | Remote Users (Site B) | `10.10.11.0/24` | `10.10.11.1` | `10.10.11.254` | `10.10.11.255` |
| 6 | VLAN 20 (Site B) | Branch Servers (Site B) | `10.10.20.16/28` | `10.10.20.17` | `10.10.20.30` | `10.10.20.31` |
| 7 | VLAN 30 (Site B) | Guests (Site B) | `10.10.31.0/25` | `10.10.31.1` | `10.10.31.126` | `10.10.31.127` |
| 8 | VLAN 99 (Site B) | Administration (Site B) | `10.10.99.16/28` | `10.10.99.17` | `10.10.99.30` | `10.10.99.31` |
| 9 | WAN | Shared Wireless WAN (OSPF Area 0) | `192.168.100.0/24` | `192.168.100.1` | `192.168.100.254` | `192.168.100.255` |

> *Note.* This table details the VLSM subnet allocation across the private block 10.10.0.0/16. It defines the network addresses, prefix lengths, usable IP ranges, and broadcast addresses for each departmental VLAN in Site A and Site B, as well as the point-to-point WAN link connecting Router A and Router B within OSPF Area 0.
