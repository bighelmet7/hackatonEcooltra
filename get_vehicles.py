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

class Connection():

	def __init__(self, server, user,
							password, remote_path):
		self.server = SERVER
		self.user = USER
		self.password = PASSWORD
		self.remote_path = REMOTE_PATH
	
	def __enter__(self):
		self.ssh = paramiko.SSHClient() 
		self.ssh.load_host_keys(os.path.expanduser(os.path.join("~", ".ssh", "known_hosts")))
		self.ssh.connect(self.server, username=self.user, password=self.password)
		self.sftp = self.ssh.open_sftp()

		return self
	
	def __exit__(self, *args):
		self.sftp.close()
		self.ssh.close()

	def save_vehicles(self):
		self.sftp.put(LOCAL_PATH, self.remote_path)


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

# def connect_ssh_sftp():
# 	ssh = paramiko.SSHClient() 
# 	ssh.load_host_keys(os.path.expanduser(os.path.join("~", ".ssh", "known_hosts")))
# 	ssh.connect(SERVER, username=USER, password=PASSWORD)
# 	sftp = ssh.open_sftp()

# 	return ssh, sftp

# def close_ssh_sftp_close(sftp, ssh):
# 	sftp.close()
# 	ssh.close()

# def save_vehicles(sftp):
# 	sftp.put(LOCAL_PATH, REMOTE_PATH)

def main():
	vehicles = get_vehicles()
	with open(LOCAL_PATH, "w") as file_vehicles:
		file_vehicles.write(vehicles)

		with Connection(SERVER, USER, PASSWORD, REMOTE_PATH) as connection:
			connection.save_vehicles()
	
	# ssh, sftp = connect_ssh_sftp()
	
	# save_vehicles(sftp)
	# close_ssh_sftp_close(sftp, ssh)

if __name__ == "__main__":
	main()