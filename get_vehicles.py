import sys
import requests
import json
import os
import paramiko

SERVER = "home556586389.1and1-data.host"
USER = "u75998576-ecooltra"
PASSWORD = "ecooltrahack"

LOCAL_PATH = "./vehicles.json"
REMOTE_PATH = "/vehicles.json"


def get_vehicles():

	url = "https://cooltra.electricfeel.net/integrator/v1/vehicles"

	data = {"system_id":sys.argv[1]}

	headers = {
		'Content-Type': "application/json",
		'Accept': "application/json",
		'Host': "cooltra.electricfeel.net",
		'Authorization': "Bearer 0fb6f9fffe309680c17d6fb7203cded9a39fc5b865f36d0763211e70a9948c58"
	}

	response = requests.get(url, headers=headers, params=data)

	if response:
		return response.text

	return 'Error'

def connect_ssh_sftp():
	ssh = paramiko.SSHClient() 
	ssh.load_host_keys(os.path.expanduser(os.path.join("~", ".ssh", "known_hosts")))
	ssh.connect(SERVER, username=USER, password=PASSWORD)
	sftp = ssh.open_sftp()

	return ssh, sftp

def close_ssh_sftp_close(sftp, ssh):
	sftp.close()
	ssh.close()

def save_vehicles(sftp):
	sftp.put(LOCAL_PATH, REMOTE_PATH)

def main():
	vehicles = get_vehicles()
	with open(LOCAL_PATH, "w") as file_vehicles:
		file_vehicles.write(vehicles)
	
	ssh, sftp = connect_ssh_sftp()
	
	save_vehicles(sftp)
	close_ssh_sftp_close(sftp, ssh)

if __name__ == "__main__":
	main()