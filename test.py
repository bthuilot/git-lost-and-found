import time
import os
import requests
from concurrent.futures import ThreadPoolExecutor, as_completed
import subprocess

def get_repos_from_org(org_url):
    org_name = org_url.split("/")[-1]
    repos = []
    
    url = f"https://api.github.com/orgs/{org_name}/repos"
    params = {'per_page': 100, 'page': 1}
    
    
    while True:
        time.sleep(.3)
        response = requests.get(url, params=params)
        data = response.json()
        if response.status_code != 200 or not data:
            break
        repos.extend([repo['html_url'] for repo in data])
        params['page'] += 1
    
    return repos


# def get_repos_from_org(org_url):
#     org_name = org_url.split("/")[-1]
#     api_url = f"https://api.github.com/orgs/{org_name}/repos"
#     response = requests.get(api_url)
#     if response.status_code == 200:
#         return [repo["html_url"] for repo in response.json()]
#     else:
#         print(f"Failed to fetch repositories for organization {org_name}")
#         return []


def scan_repo(repo_url):
    project_name = repo_url.split("/")[-1]
    git_group = repo_url.split("/")[-2]
    command = [
        "docker",
        "run",
        "-e",
        "LOG_LEVEL=info",
        "-v",
        f"{os.getcwd()}/results:/results",
        "-t",
        "git-scanner:latest",
        "scan",
        "--output",
        f"/results/{git_group}/{project_name}.json",
        "--repo-url",
        repo_url,
    ]
    process = subprocess.Popen(command)
    process.wait()


[]


def main():
    data = {
        "organizations": [
            # "https://github.com/slackhq",  # Slack
            # "https://github.com/newrelic",  # New Relic
            # "https://github.com/Shopify",  # Shopify
            # "https://github.com/datadog",  # Datadog
            # "https://github.com/azure",  # Microsoft Azure
            "https://github.com/mapbox",  # Mapbox
            "https://github.com/cloudflare",  # Cloudflare
            "https://github.com/netflix",  # Netflix
            "https://github.com/uber",  # Uber
            "https://github.com/github",  # GitHub
            "https://github.com/coinbase",  # Coinbase
            "https://github.com/spotify",  # Spotify
            "https://github.com/airbnb",  # Airbnb
            "https://github.com/intel",  # Intel
            "https://github.com/stripe",  # Stripe
            "https://github.com/square",  # Square
            "https://github.com/lyft",  # Lyft
            "https://github.com/zoom",  # Zoom
            "https://github.com/okta",  # Okta
            "https://github.com/pinterest",  # Pinterest
            "https://github.com/teslamotors",  # Tesla
            "https://github.com/atlassian",  # Atlassian
            "https://github.com/dropbox",  # Dropbox
            "https://github.com/indeedeng",  # Indeed
            "https://github.com/twilio",  # Twilio
            "https://github.com/opensea",  # OpenSea
            "https://github.com/yelp",  # Yelp
            "https://github.com/bugcrowd",  # Bugcrowd
            "https://github.com/cisco",  # Cisco
            "https://github.com/rapid7",  # Rapid7
            "https://github.com/elastic",  # Elastic
            "https://github.com/sonatype",  # Sonatype
            "https://github.com/launchdarkly",  # LaunchDarkly
            "https://github.com/hackerone",  # HackerOne
            "https://github.com/segmentio",  # Segment
            "https://github.com/coinbase",  # Coinbase
            "https://github.com/fastly",  # Fastly
            "https://github.com/1password",  # 1Password
            "https://github.com/algolia",  # Algolia
            "https://github.com/bugcrowd",  # Bugcrowd
            "https://github.com/basecamp",  # Basecamp
            "https://github.com/kayak",  # KAYAK
            "https://github.com/valvesoftware",  # Valve (Steam)
            "https://github.com/gitlabhq",  # GitLab
            "https://github.com/linkedin",  # LinkedIn
            "https://github.com/expedia",  # Expedia
            "https://github.com/bookingcom",  # Booking.com
            "https://github.com/snapchat",  # Snapchat
            "https://github.com/adobe",  # Adobe
            "https://github.com/facebook",  # Facebook
            "https://github.com/google",  # Google
            "https://github.com/amazon",  # Amazon
            "https://github.com/apple",  # Apple
            "https://github.com/microsoft",  # Microsoft
            "https://github.com/ibm",  # IBM
            "https://github.com/oracle",  # Oracle
            "https://github.com/akamai",  # Akamai
            "https://github.com/duckduckgo",  # DuckDuckGo
            "https://github.com/bitwarden",  # Bitwarden
            "https://github.com/discord",  # Discord
        ],
        "repositories": [
            "https://github.com/circlefin/noble-cctp",
            "https://github.com/circlefin/evm-cctp-contracts",
            "https://github.com/circlefin/solana-cctp-contracts",
            "https://github.com/nimiq/core-rs-albatross",
            "https://github.com/nimiq/core-js",
            "https://github.com/nimiq/core-rs",
            "https://github.com/nimiq/ledger-app-nimiq",
            "https://github.com/Chia-Network/chia-blockchain",
            "https://github.com/Chia-Network/chia-blockchain-gui",
            "https://github.com/Chia-Network/chia_rs",
            "https://github.com/Chia-Network/chiapos",
            "https://github.com/Chia-Network/chiavdf",
            "https://github.com/Chia-Network/clvm_rs",
            "https://github.com/leather-wallet/extension",
            "https://github.com/fireblocks/mpc-lib",
            "https://github.com/worldcoin/world-id-contracts",
            "https://github.com/worldcoin/world-id-state-bridge",
            "https://github.com/tronprotocol/java-tron",
            "https://github.com/coralcube-oss/mmm/releases/latest",
            "https://github.com/magiceden-oss/erc721m/releases/latest",
            "https://github.com/magiceden-oss/open_creator_protocol/releases/latest",
            "https://github.com/MetaMask/snaps/tree/main",
            "https://github.com/MetaMask/snaps-directory",
            "https://github.com/strongdm/strongdm-sdk-go",
            "https://github.com/strongdm/strongdm-sdk-java",
            "https://github.com/strongdm/strongdm-sdk-python",
            "https://github.com/strongdm/strongdm-sdk-ruby",
            "https://github.com/rundeck/rundeck",
            "https://github.com/rundeck/rundeck-cli",
            "https://github.com/Agoric/agoric-sdk/tree/master/packages/ERTP",
            "https://github.com/Agoric/agoric-sdk/tree/master/packages/inter-protocol",
            "https://github.com/Agoric/agoric-sdk/tree/master/packages/zoe",
            "https://github.com/cosmos/ibc-go/tree/main",
            "https://github.com/crypto-com/chain-desktop-wallet",
            "https://github.com/crypto-com/cro-staking",
            "https://github.com/crypto-com/sample-chain-wallet",
            "https://github.com/crypto-com/swap-contracts-core",
            "https://github.com/crypto-com/swap-contracts-periphery",
            "https://github.com/endojs/endo/tree/master/packages/ses",
            "https://github.com/daita/files_fulltextsearch_tesseract",
            "https://github.com/nextcloud/server",
            "https://github.com/nextcloud/collectives",
            "https://github.com/nextcloud/files_confidential",
            "https://github.com/nextcloud/tables",
            "https://github.com/roundcube/roundcubemail",
            "https://github.com/ruby/ruby",
            "https://github.com/bwillis/versioncake",
            "https://github.com/revive-adserver/revive-adserver",
            "https://github.com/arkadiyt/aws_public_ips",
            "https://github.com/arkadiyt/bounty-targets",
            "https://github.com/arkadiyt/ddexport",
            "https://github.com/arkadiyt/free-ft",
            "https://github.com/arkadiyt/protodump",
            "https://github.com/arkadiyt/ssrf_filter",
            "https://github.com/arkadiyt/zoom-redirector",
        ],
    }
    repos_to_scan = []
    
    # load repos from file else get repos from organizations above
    if os.path.exists("results/repos.txt"):
        with open("results/repos.txt", "r") as f:
            repos_to_scan = f.read().splitlines()
            
    if len(repos_to_scan) > 0:
        print(f"Found {len(repos_to_scan)} repositories to scan")
    else:
        for org_url in list(set(data["organizations"])):
            os.makedirs(os.path.join("results", org_url.split("/")[-1]), exist_ok=True)
            repos = get_repos_from_org(org_url)
            print(f"Found {len(repos)} repositories in organization: {org_url}")
            repos_to_scan += repos

        repos_to_scan += data["repositories"]
        
    repos_to_scan = list(set(repos_to_scan))
        
    # save repos to file
    with open("results/repos.txt", "w") as f:
        f.write("\n".join(repos_to_scan))
    
    with ThreadPoolExecutor(max_workers=3) as executor:
        futures = {
            executor.submit(scan_repo, repo_url): repo_url for repo_url in repos_to_scan
        }
        for future in as_completed(futures):
            repo_url = futures[future]
            try:
                future.result()
            except Exception as e:
                print(f"Error scanning repository {repo_url}: {e}")


if __name__ == "__main__":
    main()
