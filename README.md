# C2 Framework

**Overview**
---
This is Command and Control (C2) Framework is designed by [**Mohamed Masmoudi**](https://github.com/enigma522) and [**Moahemd Malek Gharbi**](https://github.com/Malek-trez) as part of the a summer intership at [**Pwn&Patch**](https://www.pwnandpatch.com/) to facilitate red team operations. 

**Architecture**
---
The framework includes:

1. CLI: Command-line interface for interacting with the implant.
2. C2 Server: The server-side component that controls and communicates with the implant.
3. Implant: A malware that runs on the target machine.
<!-- PROJECT LOGO -->
<br />
<div align="center">
<img src="https://github.com/enigma522/C2-framework/img/basicarch.png" width=300>
  <h3 align="center">C2 - Framework</h3>
  </p>
</div>

**Features**
---
- Secure communication between implant and C2 server.
- Easy-to-use CLI for managing implants.
- Real-time data collection and command execution.
- modulare implemantion that makes it easy to add more featres.
  
**Installation**
---

1. Prerequisites
   
Go (for implant development)

Docker (for C2 server)

Python (for C2 server)

2. Clone the Repository

```shell
git clone https://github.com/enigma522/C2-framework
cd c2-framework
```

3. Setting Up the Implant
   
    (https://github.com/enigma522/C2-framework/tree/main/BFimplant) 
- Navigate to the implant directory:

` cd implant `

-change the ip address and for your c2 server in cmd/main.go

- Build the implant:
- 
`go build -o implant main.go`

4. Setting Up the CLI:
   
 (https://github.com/enigma522/C2-framework/tree/main/Cli)

Navigate to the CLI directory:
`cd cli`
`pip install -r requirements.txt`
`python cli.py`

Install dependencies and build the CLI:

5. Setting Up the C2 Server:
   
   (https://github.com/enigma522/C2-framework/tree/main/C2-Server)

`cd c2-server`

`docker-compose up`


<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to learn, inspire, and create. Any contributions you make are **greatly appreciated**.

If you have a suggestion that would make this better, please fork the repo and create a pull request. You can also simply open an issue with the tag "enhancement".
Don't forget to give the project a star! Thanks again!

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- LICENSE -->
## License

Distributed under the MIT License. See `LICENSE.txt` for more information.

<p align="right">(<a href="#readme-top">back to top</a>)</p>



<!-- CONTACT -->
## Contact

Moahemd Masmoudi ~ enigma - mohamedmasmoudi745@gmail.com  
Mohamed Malek Gharbi ~ Trez13  - melek.gharbi1@gmail.com

<p align="right">(<a href="#readme-top">back to top</a>)</p>

