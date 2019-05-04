import sys
import requests
import json

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

def main():
	vehicles = get_vehicles()

if __name__ == "__main__":
	main()