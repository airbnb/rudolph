# Rudolph
Rudolph is the control server counterpart of [Santa](https://github.com/google/santa), and is used to rapidly deploy configurations to Santa agents.

Rudolph is built in Amazon Web Services, and utilizes exclusively serverless components to reduce operational burden. It is designed to be fast, easy-to-use, low-maintenance, and cost-conscious.


# Deployment

## Step 1) Deploy Rudolph
Start by deploying rudolph ([docs/deploy.md](docs/deploy.md)).


## Step 2) Deploying Santa Agents
Next, deploy and configure your Santa sensors ([docs/configuring-santa.md](docs/configuring-santa.md)).


## Step 3) Deploy Rules
Use the cli to sync rules ([docs/rules.md](docs/rules.md)).

