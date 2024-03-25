MONGO_VERSION=6.0
sudo apt-get update
sudo apt-get install gnupg -y
# Changes here!
  sudo mkdir -p /etc/apt/keyrings
  curl -fsSL https://www.mongodb.org/static/pgp/server-${MONGO_VERSION}.asc | sudo gpg --dearmor -o /etc/apt/keyrings/mongodb-${MONGO_VERSION}.gpg
cd /etc/apt/sources.list.d/
sudo touch mongodb-org-${MONGO_VERSION}.list
# And here!
  echo "deb [arch=amd64,arm64 signed-by=/etc/apt/keyrings/mongodb-${MONGO_VERSION}.gpg] https://repo.mongodb.org/apt/ubuntu focal/mongodb-org/${MONGO_VERSION} multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-${MONGO_VERSION}.list
sudo apt-get update
sudo apt-get install -y mongodb-org
sudo systemctl start mongod