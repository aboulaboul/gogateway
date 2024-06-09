<h1>gogateway</h1>
<h2>what</h2>
a personal gateway to handle dispatched servers.
<p>
I'm rusty on Go but I need a gateway to dispatch incoming requests from the web to micro services and http servers that run on different Raspberry Pi.<br>
I'm using dynDNS to route my domain to my box and then I want to handle the dispatching with this gateway project.
</p>
<h2>install / run</h2>
<p>
assuming  : <br>
- we're on Raspberry Pi OS (previously called Raspbian, in other words debian)<br>
- we're not using docker, we keep everything as simple as possible<br>
</p>
<h3>prep</h3>
<p>
sudo apt update<br>
sudo apt upgrade<br>
sudo apt install golang<br>
git clone https://github.com/aboulaboul/gogateway<br>
</p>
<h3>cert ssh</h3>
<p>
sudo apt install certbot <br>
sudo certbot certonly --standalone <br>
as a result : <br>
Certificate is saved at: /etc/letsencrypt/live/<i>your_domain_name</i>/fullchain.pem<br>
Key is saved at:         /etc/letsencrypt/live/<i>your_domain_name</i>/privkey.pem<br>
copying to project directory<br>
sudo cp /etc/letsencrypt/live/<i>your_domain_name</i>/fullchain.pem ./gogateway/server.crt<br>
sudo cp /etc/letsencrypt/live/<i>your_domain_name</i>/privkey.pem ./gogateway/server.key<br>
sudo chmod 644 ./gogateway/server.key
</p>
<h3>routes</h3>
<p>
edit routes : <br>
nano ./gogateway/routes.json
</p>
<h2>run</h2>
<p>
go run ./gogateway/main.go<br>
</p>
