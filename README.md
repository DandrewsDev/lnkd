# lnkd

## About
LNKD or "Linked" is a self-hosted url shortener service. It was created as a open source home brew version of many existing services for taking long comlicated URLs, and shriking them down for ease of use.

## Demo
https://lnkds.me

Logged-in user experience can be tested using the demo user.
Username:demo
Password:demo

This site is erased on a periodic basis. So we don't recommend you use it as your primary link shortener.

## Running LNKD
LNKD is designed to be super easy to run using a docker container. If you have a URL, and docker, you can run LNKD. 

Below is a quick and simple command that can be run to get a LNKD container up and running.

> Note: You will need to replace the LNDK_URL value with the domain you'll be using to run the service.
	
    docker run -d \
		--name lnkd \
        -p 8070:8070 \
        -e LNKD_URL=https://lnkds.me/ \
        dandrewsdev/lnkd:latest

### First start up
When LNKD starts for the first time it will create a database. By default, recreating the container will currently wipe any users and links you may have created. However, the database can be persisted using a mounted volume.
Here is an example of a container with a persistent db that can be re-created or updated without starting fresh each time.

     docker run -d \
      	--name lnkd \
          -p 8070:8070 \
          -v ./lnkd/db:/data/lnkd 
          -e LNKD_URL=https://lnkds.me/ \
          dandrewsdev/lnkd:latest

During the first time boot process the admin user is created. A random password is generated when that user is created. This password is printed to the container log. If you started the container in detach mode, you can view the password by running the command below.

    docker logs lnkd
It is recommended that you change this password the first time you log in.

### First login 
Once the container is up and running you can connect to the UI using the IP address and Port 8070 (Or the port you specified). When run locally this would generally look like localhost:8070. However, this service doesn't make much sense unless its run on a server with a public facing IP Address. 

The automatically created user has a username of admin, and that randomly generated password we mentioned above. You can login by clicking the login link in the top right.

## Route/Link Management
Once up and running any user can create a shorted link. The homepage will present with a form for a URL to be shorted. On submit LNKD will create a new random short link and automatically copy it to the clipboard. When anyone hits the resulting URL they will automatically be forwarded to the full length URL entered in the initial form. We'll keep an ongoing count of how many times this link gets used, which can be seen by logged in admin users.

### Custom named links
All logged-in users can create custom named links. The same homepage form will have an additional optional field for the custom name. Custom named links can only be removed or edited by the user that created the named link, or any Admin user. 

### Logged in users
As previously mentioned all logged-in users can create custom named links. Additionally, they can see a list of every link they've created. Both custom named links, and automatically generated short links. Users will also be able to see the count of the number of times their links have been used.
Additionally, logged-in users have permission to edit or remove any of the links they have created.

## User management
Admin users can create, edit, or delete any users using the User Management page. Currently, non admin users can not edit their own information or update their passwords. However, both these features are planned additions. 

## Use with a reverse proxy.
The LNKD Service does not create and use its own SSL certificates. And while it can be used via http, that is not recommended. 
We recommend putting LNKD being a reverse proxy like Nginx Proxy Manager, or Traefik. These will handle getting an SSL cert via Letsencrypt and forward all the need traffic to your LNKD container.
The https://lnkds.me demo site uses Nginx Proxy Manager
