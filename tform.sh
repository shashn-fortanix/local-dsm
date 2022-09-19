go get
make install
pwd

ls
cd ~/workspace/jenkintegration
pwd
ls
echo "sha1sum is"
# cd ~/terraform-provider-dsm
find ./jenkintegration -type f -print0 | sort -z | xargs -0 sha1sum | sha1sum
