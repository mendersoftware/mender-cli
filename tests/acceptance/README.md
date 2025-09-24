# mender-cli Acceptance tests

## How to run

1. Clone `mender-server`

   ```bash
   export MENDER_SERVER_DIR=$HOME/mender-server # Optional: Replace with desired location
   git clone git@github.com:mendersoftware/mender-server $HOME/mender-server
   ```

2. Run the tests:

   ```bash
   docker compose run --rm acceptance
   ```

3. (Optional): Clean up resources

   ```bash
   docker compose down -v
   ```
