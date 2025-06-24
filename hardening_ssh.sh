#!/usr/bin/env bash
echo "🔒 Starting SSH hardening script..."
if [[ -z "$PASS" ]]; then
  echo "Usage: $0 <ssh_password>"
  exit 1
fi
echo "🔧 Creating SSH hardening override..."
OVERRIDE_PATH="/etc/ssh/sshd_config.d/hardening.conf"
rm -f "$OVERRIDE_PATH" 2>/dev/null || true
echo "$PASS" | sudo -S tee "$OVERRIDE_PATH" > /dev/null <<EOF
PermitRootLogin no
PasswordAuthentication no
PubkeyAuthentication yes
EOF

# === 4. Test SSH config ===
echo "🧪 Testing SSH config syntax..."
echo "$PASS" | sudo -S sshd -t
if [[ $? -ne 0 ]]; then
  echo "❌ SSH config test failed. Aborting."
  exit 1
fi

# === 5. Restart SSH service ===
echo "🔄 Restarting SSH service..."
echo "$PASS" | sudo -S systemctl restart ssh

echo "✅ SSH hardened successfully."
