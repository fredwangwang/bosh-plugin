## What
Plugin manager is a lightweight application (plugin) hosting platform enables running small workload
outside of cloud foundry without creating dedicated vms for each workload.

## API
Note: all the api endpoints under `/api/v1` require authentication

- GET `/health` health check endpoint for plugin manager, returns HTTP 200

- GET `/api/v1/plugins` list all plugins currently installed  
  ```
  $ curl -H "Authorization: $(cf oauth-token) https://host:port/api/v1/plugins"
  [
    {
      "Name": "sample-plugin",
      "Description": "description for sample plugin!",
      "Location": "sample-plugin",
      "Enabled": true,
      "Env": {
        "PORT": "4321"
      },
      "Arg": [
        "thisisanexample"
      ],
      "AdditionalEnv": {
        "METRON_CA_CERT_PATH": "/var/vcap/jobs/sample-plugin/config/metron_ca_cert.pem",
        "METRON_CERT_PATH": "/var/vcap/jobs/sample-plugin/config/metron_cert.pem",
        "METRON_KEY_PATH": "/var/vcap/jobs/sample-plugin/config/metron_cert.key"
      },
      "PendingEnv": {}
    }
  ]
  ```

- POST `/api/v1/plugins` upload the plugin to plugin manager and enable the plugin 
  ```
  $ curl -H "Authorization: $(cf oauth-token)" -F 'file=@example/sample-plugin.zip' https://host:port/api/v1/plugins
  {"message":"plugin uploaded successfully"}
  ```

- GET `/api/v1/plugins/:name` get the details of a plugin

- PATCH `/api/v1/plugins/:name` update configuration (environment variable) for the plugin
  ```
  curl -H "Authorization: $(cf oauth-token)" -X PATCH https://host:port/api/v1/plugins/sample-plugin?GREETING=hi
  "config applied, need to disable/enable the plugin to see the effect"
  ```

- DELETE `/api/v1/plugins/:name` delete the plugin from plugin-manager

- POST `/api/v1/plugins/:name/enable` enable the plugin

- POST `/api/v1/plugins/:name/disable` disable the plugin


## Reference:
BOSH VM: https://bosh.io/docs/vm-config/
BPM: https://github.com/cloudfoundry/bpm-release/blob/master/docs/config.md
Monit: https://mmonit.com/monit/documentation/monit.html
