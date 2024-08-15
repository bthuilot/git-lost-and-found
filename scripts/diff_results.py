"""
Script to diff two trufflehog or gitleaks scan results
and output the new secrets found in the new scan result.

Usage:
    python diff_results.py -b <baseline> -n <new> -o <output> -s <scanner>

run `python diff_results.py --help` for more information.
"""
import argparse
import json
import logging
from io import TextIOWrapper
from typing import Any, Callable, List
from sys import stdout, stderr

logger = logging.getLogger(__name__)
logger.setLevel(logging.DEBUG)

# Initialize parser
parser = argparse.ArgumentParser()

parser.add_argument("-b", "--baseline", help="baseline trufflehog scan JSON", required=True)
parser.add_argument("-n", "--new", help="new trufflehog scan JSON", required=True)
parser.add_argument("-o", "--output", help="output file", required=True)
parser.add_argument("-s", "--scanner", help="results scanner schema", choices=["trufflehog", "gitleaks"], required=True)


def parse_json(fp: str) -> Any:
    """Parse a JSON file and return the data as a Python object."""
    with open(fp) as f:
        return json.load(f)


def get_trufflehog_secret(secret: Any) -> str:
    """Return the secret value from a trufflehog scan result."""
    if type(secret) != dict:
        raise ValueError("Secret must be a dictionary")

    if "Raw" not in secret or "RawV2" not in secret:
        raise ValueError("Secret must have 'Raw' and 'RawV2' keys")

    return secret.get("Raw") + secret.get("RawV2")


def get_gitleaks_secret(secret: Any) -> str:
    """Return the secret value from a gitleaks scan result."""
    if type(secret) != dict:
        raise ValueError("Secret must be a dictionary")

    if "Fingerprint" not in secret:
        raise ValueError("Secret must have 'Fingerprint' key")

    return secret.get("Fingerprint")


def main(baseline_path: str, new_path: Any, output_path: str, raw_secret_f: Callable[[Any], str]):
    baseline, new = parse_json(baseline_path), parse_json(new_path)

    if type(baseline) is not list or type(new) is not list:
        raise ValueError("Both baseline and new must be lists")

    baseline_secrets = set()
    for secret in baseline:
        try:
            secret_val = raw_secret_f(secret)
        except ValueError as ve:
            logger.warning(f"error while extracting baseline secret: {ve}, skipping")
            continue

        if secret_val is not None:
            baseline_secrets.add(secret_val)

    diff = []
    for secret in new:
        try:
            secret_val = raw_secret_f(secret)
        except ValueError as ve:
            logger.warning(f"error while extracting new secret: {ve}, skipping")
            continue
        if secret_val is not None and secret_val not in baseline_secrets:
            diff.append(secret)

    io: TextIOWrapper = stdout if output_path == "-" else open(output_path, "w")
    with io as f:
        json.dump(diff, f, indent=4)


if __name__ == '__main__':
    args = parser.parse_args()

    secret_f = get_trufflehog_secret if args.scanner == "trufflehog" else get_gitleaks_secret

    try:
        main(args.baseline, args.new, args.output, secret_f)
    except Exception as e:
        logger.error("%s", e)
        exit(1)
