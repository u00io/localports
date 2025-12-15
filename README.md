# LocalPorts

**LocalPorts** is a lightweight desktop utility for inspecting local ports and active network connections on a computer.

It allows you to see:
- which **local TCP/UDP ports** are open,
- which processes are using them,
- where exactly your computer is connected to,
- which services and countries are on the remote side of a connection.

The tool is designed for developers, system engineers, and technical users.

---

## Key Features

- 📡 View **open local ports** (TCP by default)
- 🔗 View **all active network connections**
- 🧩 Map connections to **processes and PIDs**
- 🌐 Detect **remote services by port**
- 🗺️ Detect **country of remote IP addresses**
- 🎛️ Filtering by:
  - protocol (TCP / UDP)
  - connection state (LISTEN / ESTABLISHED / others)
- 🖥️ Minimal, distraction-free user interface
- ⚡ Fast startup and low overhead

---

## Default Behavior

By default, **LocalPorts** displays only:
- **local TCP ports**
- in **LISTEN** state

This makes it easy to immediately see which services are listening on the system and may be accessible from the network.

---

## IP Geolocation

Remote IP addresses are resolved to **country level only** using the **GeoLite2 Country** database.

- Only country information is used
- No city, coordinates, or precise location data
- Geolocation data is provided for informational purposes only


---

## Distribution Model

**LocalPorts is distributed as a single executable file.**

- No administrator privileges required
- No installation required
- No drivers or background services
- No configuration files
- No external dependencies

Simply download and run the executable.

---

## Architecture and Privacy

- Runs entirely **locally**
- Does not send any data over the network
- Does not use cloud services or external APIs
- No telemetry or tracking

---

## Intended Audience

- software developers
- system administrators
- security engineers
- technical professionals
- anyone who wants to understand **what their computer is listening on and connecting to**

---

## Project Status

The project is under active development.  
Features are added based on real-world use cases.

---

## License

This project is released under the **MIT License**.

Third-party data used:
- **GeoLite2 Country** — © MaxMind, Inc.
