import json
import os
import requests
from concurrent.futures import ThreadPoolExecutor, as_completed


def get_repos_from_org(org_url):
    org_name = org_url.split("/")[-1]
    api_url = f"https://api.github.com/orgs/{org_name}/repos"
    response = requests.get(api_url)
    if response.status_code == 200:
        return [repo["html_url"] for repo in response.json()]
    else:
        print(f"Failed to fetch repositories for organization {org_name}")
        return []


def scan_repo(repo_url):
    print(f"Scanning repository: {repo_url}")
    project_name = repo_url.split("/")[-1]
    git_group = repo_url.split("/")[-2]
    command = f"./bin/git-scanner scan --output results/{git_group}/{project_name}.json --repo-url {repo_url}"
    os.system(command)


def main():
    data = {
        "organizations": [
            # "https://github.com/yahoo",
            # "https://github.com/Agoric",
            # "https://github.com/argoproj",
            # "https://github.com/netlify",
            # "https://github.com/MariaDB",
            # "https://github.com/nextcloud",
            # "https://github.com/slackhq",
            # "https://github.com/newrelic",
            # "https://github.com/Shopify",
            # "https://github.com/datadog",
            "https://github.com/azure",
            "https://github.com/mapbox",
            "https://github.com/cloudflare",
            "https://github.com/netflix",
            "https://github.com/openai",
            "https://github.com/google-deepmind"
            "https://github.com/microsoft",
            
            
            
            # "https://github.com/leather-wallet",
            # "https://github.com/fireblocks",
            # "https://github.com/worldcoin",
            # "https://github.com/okx",
            # "https://github.com/tronprotocol",
            # "https://github.com/coralcube-oss",
            # "https://github.com/magiceden-oss",
            # "https://github.com/USStateDept",
            # "https://github.com/MetaMask",
            # "https://github.com/sorare",
            # "https://github.com/strongdm",
            # "https://github.com/rundeck",
            # "https://github.com/endojs",
            # "https://github.com/moonpay",
            # "https://github.com/Electron",
            # "https://github.com/Nginx",
            # "https://github.com/apache",
            # "https://github.com/curl",
            # "https://github.com/django",
            # "https://github.com/libuv",
            # "https://github.com/nodejs",
            # "https://github.com/openssl",
            # "https://github.com/rack",
            # "https://github.com/rails",
            # "https://github.com/ruby",
            # "https://github.com/rubygems",
            # "https://github.com/rust-lang",
            # "https://github.com/spiffe",
            # "https://github.com/18f",
            # "https://github.com/gsa",
            # "https://github.com/snowplow",
            # "https://github.com/hashgraph",
            # "https://github.com/fastify",
            # "https://github.com/WorldHealthOrganization",
            # "https://github.com/stripe",
            # "https://github.com/DopplerHQ",
            # "https://github.com/thesokrin",
            # "https://github.com/grindrlabs",
            # "https://github.com/OpenMage",
            # "https://github.com/jitsi",
            # "https://github.com/DefectDojo",
            # "https://github.com/Dynatrace",
            # "https://github.com/reddit-archive",
            # "https://github.com/mattermost",
            # "https://github.com/etherspot",
            # "https://github.com/pyca",
            # "https://github.com/innocraft",
            # "https://github.com/matomo-org",
            # "https://github.com/8x8",
            # "https://github.com/smartcontractkit",
            # "https://github.com/pixiv",
            # "https://github.com/arkadiyt",
            # "https://github.com/cometbft",
            # "https://github.com/CosmWasm",
            # "https://github.com/cosmos",
            # "https://github.com/strangelove-ventures",
            # "https://github.com/duckduckgo",
            # "https://github.com/crypto-com",
            # "https://github.com/trycourier",
            # "https://github.com/roundcube",
            # "https://github.com/Valvesoftware",
            # "https://github.com/binary-com",
            # "https://github.com/deriv-com",
            # "https://github.com/fetlife",
            # "https://github.com/smooch",
            # "https://github.com/irccloud",
            # "https://github.com/concrete5",
            # "https://github.com/arkime",
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

    for org_url in data["organizations"]:
        os.makedirs(os.path.join("results", org_url.split("/")[-1]), exist_ok=True)

        repos = get_repos_from_org(org_url)

        for repo_url in repos:
            scan_repo(repo_url)
        print(f"Finished scanning repository: {repo_url}")

        # ensure results/<org_name> directory exists
        # with ThreadPoolExecutor(max_workers=3) as executor:
        #     futures = {
        #         executor.submit(scan_repo, repo_url): repo_url for repo_url in repos
        #     }
        #     for future in as_completed(futures):
        #         repo_url = futures[future]
        #         try:
        #             future.result()
        #             print(f"Finished scanning repository: {repo_url}")
        #         except Exception as e:
        #             print(f"Error scanning repository {repo_url}: {e}")


if __name__ == "__main__":
    main()
