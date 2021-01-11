# vorteil-nginx-unit-demo

The repository is used as a demo platform for GitHub, Ansible, AWS, GCP an Vorteil integration using NGINX Unit as the application server. The NGINX Unit application server serves Go pages.


# What are we trying to achieve
In a couple of simple steps:

1. GitHub actions will offload a pull and build request to Ansible (or Ansible Tower)
2. Ansible (Tower) will pull the repository and build a new Vorteil machine using the configuration file
3. Ansible (Tower) will provision the image (or ask Vorteil to do it) to AWS and GCP
4. Ansible (Tower) will create the machine (or replace) the machine if it currently exists

To change the colour of the background - simply PUT the following to the IP address

```bash

$ cat change-colour.json 
{
	"applications": {
		"helloworld": {
			"type": "external",
			"working_directory": "/www/helloworld",
			"executable": "bin/helloworld",
			"arguments": [
				"--colour=#FF0000"
			]
		}
	},

	"listeners": {
		"*:8888": {
			"pass": "applications/helloworld"
		}
	}
}

$ curl -X PUT --data-binary @change-colour.json http://localhost:8080/config

```
