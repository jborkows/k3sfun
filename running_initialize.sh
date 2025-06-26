if [ -z "$1" ]; then
  echo "Usage: $0 <remote_host> <ssh_password>"
  exit 1
fi
host="$1"
if [ -z "$2" ]; then
  echo "Usage: $0 <remote_host> <ssh_password>"
  exit 1
fi
password=$2

FOO=1
BOO=2
ssh $host "PASS='$password' bash -s -- '$FOO' '$BOO'  " < initialize.sh 
