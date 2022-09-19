# git clone https://github.com/shashn-fortanix/sq-dsm_1.0.1_amd64.deb.git
# dpkg -i sq-dsm_1.0.1_amd64.deb
cd pgp-sq

echo FORTANIX_API_ENDPOINT="https://sdkms.fortanix.com/"
echo FORTANIX_API_KEY="MzRkYmYwNTItZmY5ZC00MTEwLWFkZTUtYWQ3MWRmNzU2YWQ1OmdQM1pSYWZrOFFkcUkyR1FrZG80SUp6ZU9kRGhjTVNqdDlhT3M5bVZ3VURuNXlWNWpyYzh4TEozeVZGUVE1NU5PcmdZVEltRGJSSE93T2ROS0NwT2RR"
sq-dsm --force key generate --dsm-key="signerJ" --userid="shashidhar.naraparaju@fortanix.com"
sq-dsm --force key extract-cert --dsm-key="signerJ" > signerJ.asc
 
find ~/Jenkintegration -type f -print0 | sort -z | xargs -0 sha1sum | sha1sum > hash

sq-dsm sign --dsm-key="signerJ" hash > hash.sign
sq-dsm verify --signer-cert="signer.asc" hash.sign
