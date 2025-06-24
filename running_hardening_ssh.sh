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

while true; do
  read -p "Did you copy ssh certifate?: (Y/N)" answer
  case "$answer" in
    yes|YES|y|Y)
      echo "You typed YES. Continuing with the script..."
      break
      ;;
    no|NO|n|N)
      echo "You typed NO."
      exit 1
      break
      ;;
    *)
      echo "Invalid input. Please type yes or no."
      ;;
  esac
done
ssh $host "PASS=$password" 'bash -s' < hardening_ssh.sh 
