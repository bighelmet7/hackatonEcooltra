import sys
import requests
import json
import os
import paramiko
import argparse

SERVER = "home556586389.1and1-data.host"
USER = "u75998576-ecooltra"
PASSWORD = "ecooltrahack"

REMOTE_PATH = "/vehicles.json"


class ScrapperException(Exception):
    pass


class Connection(object):

    def __init__(self, server, user, password, remote_path, local_path):
        self.server = SERVER
        self.user = USER
        self.password = PASSWORD
        self.remote_path = REMOTE_PATH
        self.local_path = local_path

    def __enter__(self):
        self.ssh = paramiko.SSHClient()
        self.ssh.load_host_keys(os.path.expanduser(
            os.path.join("~", ".ssh", "known_hosts")))
        self.ssh.connect(self.server, username=self.user,
                         password=self.password)
        self.sftp = self.ssh.open_sftp()

        return self

    def __exit__(self, *args):
        self.sftp.close()
        self.ssh.close()

    def save_vehicles(self):
        self.sftp.put(self.local_path, self.remote_path)


def create_parse_arguments():
    parser = argparse.ArgumentParser(
        description='Usage'
    )
    parser.add_argument('--path', help='Local JSON file path', required=True)
    parser.add_argument('--city', default='barcelona',
                        help='City to get vehicles', required=False)

    return parser


def get_args():
    parser = create_parse_arguments()
    return parser.parse_args()


def get_vehicles(city):

    url = "https://cooltra.electricfeel.net/integrator/v1/vehicles"

    data = {"system_id": city}

    headers = {
        'Content-Type': "application/json",
        'Accept': "application/json",
        'Host': "cooltra.electricfeel.net",
        'Authorization': "Bearer 0fb6f9fffe309680c17d6fb7203cded9a39fc5b865f36d0763211e70a9948c58"
    }

    response = requests.get(url, headers=headers, params=data)

    if response:
        return response.text

    raise ScrapperException('No se ha conseguido una respuesta')


def main():
    args = get_args()

    local_path = args.path[0]
    city = args.city

    vehicles = get_vehicles(city)

    with open(local_path, "w") as file_vehicles:
        file_vehicles.write(vehicles)

        with Connection(SERVER, USER, PASSWORD, REMOTE_PATH, local_path) as connection:
            connection.save_vehicles()


if __name__ == "__main__":
    main()
