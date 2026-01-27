#!/bin/bash
set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

# 1. Vérifications Root
if [ "$EUID" -ne 0 ]; then
  echo -e "${RED}Ce script doit être exécuté en tant que root.${NC}"
  exit 1
fi

echo -e "${YELLOW}⚠️  Vous êtes sur le point de désinstaller Crafteur.${NC}"

# 2. Arrêt du Service
if systemctl is-active --quiet crafteur; then
    echo "Arrêt du service Crafteur..."
    systemctl stop crafteur
fi
if systemctl is-enabled --quiet crafteur; then
    systemctl disable crafteur
fi

# 3. Suppression Systemd
echo "Suppression de la configuration systemd..."
rm -f /etc/systemd/system/crafteur.service
systemctl daemon-reload

# 4. Gestion des Données (Safety Check)
echo -e "${YELLOW}Voulez-vous supprimer définitivement les données (mondes, backups, base de données) ?${NC}"
read -p "Tapez 'DELETE' pour tout supprimer, ou appuyez sur Entrée pour conserver les données : " CONFIRM

if [ "$CONFIRM" == "DELETE" ]; then
    echo -e "${RED}Suppression complète de /opt/crafteur...${NC}"
    rm -rf /opt/crafteur
    
    # Suppression utilisateur
    if id "crafteur" &>/dev/null; then
        userdel crafteur
        echo "Utilisateur 'crafteur' supprimé."
    fi
    
    echo -e "${GREEN}Désinstallation complète terminée.${NC}"
else
    echo "Suppression du binaire uniquement..."
    rm -f /opt/crafteur/crafteur
    echo -e "${GREEN}Désinstallation terminée.${NC}"
    echo -e "${GREEN}Vos données sont conservées dans : /opt/crafteur/data${NC}"
    echo "L'utilisateur système 'crafteur' a été conservé pour garder les permissions."
fi
