a shared runtime for BOSH


## API
- GET `/health` health check endpoint for the plugin manager, returns HTTP 200
- GET `/api/v1/plugins` list all plugins currently installed  
  ```bash
  curl 
  ```  
- POST `/api/v1/plugins` upload plugins to the plugin manager

reference:
BOSH VM: https://bosh.io/docs/vm-config/
BPM: https://github.com/cloudfoundry/bpm-release/blob/master/docs/config.md
Monit: https://mmonit.com/monit/documentation/monit.html


TODOS:
when updating the deployment and stemcell update, the plugin folder structure will be destroyed
needs to recreate.

