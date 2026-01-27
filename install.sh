#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'

# 1. Vérifications
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Ce script doit être exécuté en tant que root.${NC}"
  exit 1
fi

echo -e "${GREEN}Vérification des dépendances...${NC}"
for cmd in tar curl; do
  if ! command -v $cmd &> /dev/null; then
    echo -e "${RED}$cmd est introuvable. Veuillez l'installer.${NC}"
    exit 1
  fi
done

# 2. Configuration Système
echo -e "${GREEN}Configuration de l'utilisateur et des dossiers...${NC}"
if ! id "crafteur" &>/dev/null; then
    useradd -r -s /bin/false crafteur
    echo "Utilisateur 'crafteur' créé."
fi

mkdir -p /opt/crafteur/data
chown -R crafteur:crafteur /opt/crafteur

# 3. Installation Binaire
echo -e "${GREEN}Installation du binaire...${NC}"
# Simulation : on s'attend à ce que le binaire soit dans le dossier courant ou téléchargé
if [ -f "crafteur-bin" ]; then
    cp crafteur-bin /opt/crafteur/crafteur
elif [ -f "backend/crafteur-bin" ]; then
    cp backend/crafteur-bin /opt/crafteur/crafteur
elif [ -f "crafteur" ]; then
    cp crafteur /opt/crafteur/crafteur
else
    echo -e "${RED}Binaire 'crafteur' ou 'crafteur-bin' introuvable dans le dossier courant.${NC}"
    # Ici on pourrait ajouter une logique de téléchargement curl
    exit 1
fi
chmod +x /opt/crafteur/crafteur
chown crafteur:crafteur /opt/crafteur/crafteur

# 4. Service Systemd
echo -e "${GREEN}Configuration du service systemd...${NC}"
cat > /etc/systemd/system/crafteur.service <<EOF
[Unit]
Description=Crafteur Game Server Manager
After=network.target

[Service]
User=crafteur
Group=crafteur
WorkingDirectory=/opt/crafteur
ExecStart=/opt/crafteur/crafteur
Restart=always
Environment=PORT=8080

[Install]
WantedBy=multi-user.target
EOF

systemctl daemon-reload

# 5. Initialisation Admin
echo -e "${GREEN}Initialisation du compte Admin${NC}"
read -p "Entrez le nom d'utilisateur Admin : " ADMIN_USER
read -s -p "Entrez le mot de passe : " ADMIN_PASS
echo ""

if [ -z "$ADMIN_USER" ] || [ -z "$ADMIN_PASS" ]; then
    echo -e "${RED}L'utilisateur et le mot de passe sont requis.${NC}"
    exit 1
fi

echo "Création de l'utilisateur admin..."
sudo -u crafteur /opt/crafteur/crafteur create-user "$ADMIN_USER" "$ADMIN_PASS"

# 6. Démarrage
echo -e "${GREEN}Démarrage du service...${NC}"
systemctl enable crafteur
systemctl start crafteur

IP=$(hostname -I | awk '{print $1}')
echo -e "${GREEN}Installation terminée ! Accédez à http://${IP}:8080${NC}"
