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
if [ -z "$EMAIL" ]; then
  echo "Please set the EMAIL environment variable."
  exit 1
fi

ssh $host "PASS='$password' bash -s -- '$EMAIL' '$CF_TOKEN' '$IP_RANGE' '$CLOUDFLARE_API_TOKEN' " < initialize.sh 
